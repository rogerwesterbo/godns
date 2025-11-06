package v1metricsservice

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// MetricsService manages Prometheus metrics for DNS server
type MetricsService struct {
	// DNS query metrics
	QueryTotal    *prometheus.CounterVec
	QueryDuration *prometheus.HistogramVec
	QueryErrors   *prometheus.CounterVec

	// Cache metrics
	CacheHits      prometheus.Counter
	CacheMisses    prometheus.Counter
	CacheSize      prometheus.Gauge
	CacheEvictions prometheus.Counter

	// Rate limiting metrics
	RateLimitedQueries prometheus.Counter
	ActiveRateLimiters prometheus.Gauge

	// Load balancer metrics
	BackendTotal    prometheus.Gauge
	BackendHealthy  prometheus.Gauge
	BackendRequests *prometheus.CounterVec

	// Health check metrics
	HealthCheckTotal   prometheus.Gauge
	HealthCheckSuccess *prometheus.CounterVec
	HealthCheckFailure *prometheus.CounterVec
	HealthCheckLatency *prometheus.HistogramVec

	// Upstream metrics
	UpstreamQueries  prometheus.Counter
	UpstreamErrors   prometheus.Counter
	UpstreamDuration prometheus.Histogram

	registry *prometheus.Registry
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	ms := &MetricsService{
		registry: prometheus.NewRegistry(),
	}

	// DNS query metrics
	ms.QueryTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godns_query_total",
			Help: "Total number of DNS queries",
		},
		[]string{"type", "rcode"},
	)

	ms.QueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "godns_query_duration_seconds",
			Help:    "DNS query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"type"},
	)

	ms.QueryErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godns_query_errors_total",
			Help: "Total number of DNS query errors",
		},
		[]string{"type", "error"},
	)

	// Cache metrics
	ms.CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_cache_hits_total",
			Help: "Total number of cache hits",
		},
	)

	ms.CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_cache_misses_total",
			Help: "Total number of cache misses",
		},
	)

	ms.CacheSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "godns_cache_size",
			Help: "Current number of entries in cache",
		},
	)

	ms.CacheEvictions = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_cache_evictions_total",
			Help: "Total number of cache evictions",
		},
	)

	// Rate limiting metrics
	ms.RateLimitedQueries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_rate_limited_queries_total",
			Help: "Total number of rate limited queries",
		},
	)

	ms.ActiveRateLimiters = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "godns_active_rate_limiters",
			Help: "Current number of active rate limiters (unique IPs)",
		},
	)

	// Load balancer metrics
	ms.BackendTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "godns_backend_total",
			Help: "Total number of backends",
		},
	)

	ms.BackendHealthy = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "godns_backend_healthy",
			Help: "Number of healthy backends",
		},
	)

	ms.BackendRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godns_backend_requests_total",
			Help: "Total number of requests to each backend",
		},
		[]string{"backend", "status"},
	)

	// Health check metrics
	ms.HealthCheckTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "godns_health_check_total",
			Help: "Total number of health checks configured",
		},
	)

	ms.HealthCheckSuccess = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godns_health_check_success_total",
			Help: "Total number of successful health checks",
		},
		[]string{"target", "type"},
	)

	ms.HealthCheckFailure = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "godns_health_check_failure_total",
			Help: "Total number of failed health checks",
		},
		[]string{"target", "type"},
	)

	ms.HealthCheckLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "godns_health_check_latency_seconds",
			Help:    "Health check latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"target", "type"},
	)

	// Upstream metrics
	ms.UpstreamQueries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_upstream_queries_total",
			Help: "Total number of queries forwarded to upstream",
		},
	)

	ms.UpstreamErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "godns_upstream_errors_total",
			Help: "Total number of upstream query errors",
		},
	)

	ms.UpstreamDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "godns_upstream_duration_seconds",
			Help:    "Upstream query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	// Register all metrics
	ms.registerMetrics()

	return ms
}

// registerMetrics registers all metrics with the registry
func (ms *MetricsService) registerMetrics() {
	// DNS metrics
	ms.registry.MustRegister(ms.QueryTotal)
	ms.registry.MustRegister(ms.QueryDuration)
	ms.registry.MustRegister(ms.QueryErrors)

	// Cache metrics
	ms.registry.MustRegister(ms.CacheHits)
	ms.registry.MustRegister(ms.CacheMisses)
	ms.registry.MustRegister(ms.CacheSize)
	ms.registry.MustRegister(ms.CacheEvictions)

	// Rate limiting metrics
	ms.registry.MustRegister(ms.RateLimitedQueries)
	ms.registry.MustRegister(ms.ActiveRateLimiters)

	// Load balancer metrics
	ms.registry.MustRegister(ms.BackendTotal)
	ms.registry.MustRegister(ms.BackendHealthy)
	ms.registry.MustRegister(ms.BackendRequests)

	// Health check metrics
	ms.registry.MustRegister(ms.HealthCheckTotal)
	ms.registry.MustRegister(ms.HealthCheckSuccess)
	ms.registry.MustRegister(ms.HealthCheckFailure)
	ms.registry.MustRegister(ms.HealthCheckLatency)

	// Upstream metrics
	ms.registry.MustRegister(ms.UpstreamQueries)
	ms.registry.MustRegister(ms.UpstreamErrors)
	ms.registry.MustRegister(ms.UpstreamDuration)

	vlog.Info("Metrics registered successfully")
}

// Handler returns the HTTP handler for the metrics endpoint
func (ms *MetricsService) Handler() http.Handler {
	return promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})
}

// RecordQuery records a DNS query
func (ms *MetricsService) RecordQuery(queryType string, rcode string, duration float64) {
	ms.QueryTotal.WithLabelValues(queryType, rcode).Inc()
	ms.QueryDuration.WithLabelValues(queryType).Observe(duration)
}

// RecordQueryError records a DNS query error
func (ms *MetricsService) RecordQueryError(queryType string, errorType string) {
	ms.QueryErrors.WithLabelValues(queryType, errorType).Inc()
}

// RecordCacheHit records a cache hit
func (ms *MetricsService) RecordCacheHit() {
	ms.CacheHits.Inc()
}

// RecordCacheMiss records a cache miss
func (ms *MetricsService) RecordCacheMiss() {
	ms.CacheMisses.Inc()
}

// UpdateCacheSize updates the cache size gauge
func (ms *MetricsService) UpdateCacheSize(size int) {
	ms.CacheSize.Set(float64(size))
}

// RecordCacheEviction records a cache eviction
func (ms *MetricsService) RecordCacheEviction() {
	ms.CacheEvictions.Inc()
}

// RecordRateLimited records a rate limited query
func (ms *MetricsService) RecordRateLimited() {
	ms.RateLimitedQueries.Inc()
}

// UpdateActiveRateLimiters updates the active rate limiters gauge
func (ms *MetricsService) UpdateActiveRateLimiters(count int) {
	ms.ActiveRateLimiters.Set(float64(count))
}

// UpdateBackendStats updates backend statistics
func (ms *MetricsService) UpdateBackendStats(total int, healthy int) {
	ms.BackendTotal.Set(float64(total))
	ms.BackendHealthy.Set(float64(healthy))
}

// RecordBackendRequest records a request to a backend
func (ms *MetricsService) RecordBackendRequest(backend string, status string) {
	ms.BackendRequests.WithLabelValues(backend, status).Inc()
}

// UpdateHealthCheckTotal updates the total health checks gauge
func (ms *MetricsService) UpdateHealthCheckTotal(count int) {
	ms.HealthCheckTotal.Set(float64(count))
}

// RecordHealthCheckSuccess records a successful health check
func (ms *MetricsService) RecordHealthCheckSuccess(target string, checkType string) {
	ms.HealthCheckSuccess.WithLabelValues(target, checkType).Inc()
}

// RecordHealthCheckFailure records a failed health check
func (ms *MetricsService) RecordHealthCheckFailure(target string, checkType string) {
	ms.HealthCheckFailure.WithLabelValues(target, checkType).Inc()
}

// RecordHealthCheckLatency records health check latency
func (ms *MetricsService) RecordHealthCheckLatency(target string, checkType string, duration float64) {
	ms.HealthCheckLatency.WithLabelValues(target, checkType).Observe(duration)
}

// RecordUpstreamQuery records an upstream query
func (ms *MetricsService) RecordUpstreamQuery(duration float64) {
	ms.UpstreamQueries.Inc()
	ms.UpstreamDuration.Observe(duration)
}

// RecordUpstreamError records an upstream error
func (ms *MetricsService) RecordUpstreamError() {
	ms.UpstreamErrors.Inc()
}
