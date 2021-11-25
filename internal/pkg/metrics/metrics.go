package metrics

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Hits    *prometheus.CounterVec
	Timings *prometheus.HistogramVec
}

func NewMetrics(r *mux.Router) *Metrics {
	var metrics Metrics
	metrics.Hits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hits",
		}, []string{"status", "path", "method"},
	)

	metrics.Timings = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "timings",
		}, []string{"status", "path", "method"},
	)

	prometheus.MustRegister(metrics.Hits, metrics.Timings)
	r.Handle("/metrics", promhttp.Handler())
	return &metrics
}
