package db

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"route256.ozon.ru/project/loms/internal/entity"
	"time"
)

type requestType string

var (
	find   requestType = "select"
	insert requestType = "insert"
	update requestType = "update"
	delete requestType = "delete"
)

var metric = entity.Metric{
	RequestTotal: promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Tracks the number of requests to db.",
		}, []string{"type"},
	),
	RequestDuration: promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_request_duration_seconds",
			Help: "Tracks the latencies for requests to db.",
		},
		[]string{"type", "err"},
	),
}

func measureMetrics(requestType requestType, start time.Time, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	metric.RequestDuration.WithLabelValues(string(requestType), errMsg).Observe(time.Since(start).Seconds())
	metric.RequestTotal.WithLabelValues(string(requestType)).Inc()
}
