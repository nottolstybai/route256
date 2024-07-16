package entity

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	RequestTotal    *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}
