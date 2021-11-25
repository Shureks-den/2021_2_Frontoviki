package delivery

import (
	"net/http"
	"strconv"
	"time"
	"yula/internal/pkg/metrics"

	"github.com/urfave/negroni"
)

type MetricsMiddleware struct {
	metric metrics.Metrics
}

func NewMetricsMiddleware(metric *metrics.Metrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metric: *metric,
	}
}

func (mm *MetricsMiddleware) ScanMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		nrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(nrw, r)
		if r.URL.Path != "/metrics" {
			path := r.URL.Path
			mm.metric.Hits.WithLabelValues(strconv.Itoa(nrw.Status()), path, r.Method).Inc()
			mm.metric.Timings.WithLabelValues(strconv.Itoa(nrw.Status()), path, r.Method).Observe(float64(time.Since(start).Seconds()))
		}
	})
}
