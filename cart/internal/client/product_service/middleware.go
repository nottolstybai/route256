package product_service

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/entity"
	"route256.ozon.ru/project/cart/pkg/logger"
	"time"
)

const httpStatusEnhanceUrCalm = 420

type RetryMiddleware struct {
	roundTripper http.RoundTripper
	maxRetries   int
}

func NewRetryMiddleware(roundTripper http.RoundTripper, maxRetries int) *RetryMiddleware {
	return &RetryMiddleware{
		roundTripper: NewMetricsMiddleware(roundTripper),
		maxRetries:   maxRetries,
	}
}

func (rm *RetryMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return resp, err
	}

	for i := 0; i < rm.maxRetries; i++ {
		r.Body = io.NopCloser(bytes.NewReader(reqBody))

		resp, err = rm.roundTripper.RoundTrip(r)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != httpStatusEnhanceUrCalm && resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}

		logger.Warn("product_service client error", zap.Int("statusCode", resp.StatusCode))
		time.Sleep(time.Second)
	}
	return resp, err
}

type MetricsMiddleware struct {
	roundTripper http.RoundTripper
	metric       entity.Metric
}

func NewMetricsMiddleware(roundTripper http.RoundTripper) *MetricsMiddleware {
	return &MetricsMiddleware{
		roundTripper: roundTripper,
		metric: entity.Metric{
			RequestTotal: promauto.NewCounterVec(
				prometheus.CounterOpts{
					Name: "product_service_requests_total",
					Help: "Tracks the number of HTTP requests to product service.",
				}, []string{"method", "path"},
			),
			RequestDuration: promauto.NewHistogramVec(
				prometheus.HistogramOpts{
					Name: "product_service_request_duration_seconds",
					Help: "Tracks the latencies for HTTP requests to product service.",
				},
				[]string{"method", "path", "code"},
			),
		}}
}

func (m MetricsMiddleware) RoundTrip(r *http.Request) (*http.Response, error) {
	m.metric.RequestTotal.WithLabelValues(r.Method, r.URL.Path).Inc()

	start := time.Now()

	resp, err := m.roundTripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	m.metric.RequestDuration.WithLabelValues(r.Method, r.URL.Path, resp.Status).Observe(time.Since(start).Seconds())

	return resp, nil
}
