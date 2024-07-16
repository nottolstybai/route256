package mw

import (
	"errors"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func HandleMetrics(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reg := prometheus.WrapRegistererWith(prometheus.Labels{
			"path": r.URL.Path,
		}, prometheus.DefaultRegisterer)

		requestsTotal := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Tracks the number of HTTP requests.",
			}, []string{"method", "code"},
		)
		requestDuration := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "Tracks the latencies for HTTP requests.",
			},
			[]string{"method", "code"},
		)

		reqNumber := TryRegistering(reg, requestsTotal)
		reqDuration := TryRegistering(reg, requestDuration)

		// Wraps the provided http.Handler to observe the request result with the provided metrics.
		base := promhttp.InstrumentHandlerCounter(
			reqNumber.(*prometheus.CounterVec),
			promhttp.InstrumentHandlerDuration(
				reqDuration.(*prometheus.HistogramVec),
				handler,
			),
		)
		base(w, r)
	}
}

func TryRegistering(registry prometheus.Registerer, metric prometheus.Collector) prometheus.Collector {
	if err := registry.Register(metric); err != nil {
		are := &prometheus.AlreadyRegisteredError{}
		if errors.As(err, are) {
			switch are.ExistingCollector.(type) {
			case *prometheus.HistogramVec, *prometheus.CounterVec:
				return are.ExistingCollector
			default:
				panic("unknown metric type")
			}
		} else {
			// Something else went wrong!
			panic(err)
		}
	}
	return metric
}
