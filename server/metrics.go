package server

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metrics   *metricsCollector
	metricsMu sync.Once
)

type metricsCollector struct {
	reqCount    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
	cacheOps    *prometheus.CounterVec
}

func getMetrics() *metricsCollector {
	metricsMu.Do(func() {
		metrics = &metricsCollector{
			reqCount: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "image_text_http_requests_total",
					Help: "Total number of HTTP requests",
				},
				[]string{"method", "path", "status"},
			),
			reqDuration: prometheus.NewHistogramVec(
				prometheus.HistogramOpts{
					Name:    "image_text_http_request_duration_seconds",
					Help:    "HTTP request duration in seconds",
					Buckets: prometheus.DefBuckets,
				},
				[]string{"method", "path"},
			),
			cacheOps: prometheus.NewCounterVec(
				prometheus.CounterOpts{
					Name: "image_text_cache_ops_total",
					Help: "Total number of image cache operations",
				},
				[]string{"type"},
			),
		}
		prometheus.MustRegister(metrics.reqCount, metrics.reqDuration, metrics.cacheOps)
	})
	return metrics
}

func (m *metricsCollector) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		begin := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		status := strconv.Itoa(ww.Status())
		m.reqCount.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		m.reqDuration.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(begin).Seconds())
	})
}

func (m *metricsCollector) cacheHit() {
	m.cacheOps.WithLabelValues("hit").Inc()
}

func (m *metricsCollector) cacheMiss() {
	m.cacheOps.WithLabelValues("miss").Inc()
}

func prometheusHandler() http.HandlerFunc {
	return promhttp.Handler().ServeHTTP
}
