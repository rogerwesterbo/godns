package v1adminhandler

import (
	"context"
	"fmt"
	"math"
	"net/http"

	"github.com/rogerwesterbo/godns/internal/httpserver/helpers"
	"github.com/rogerwesterbo/godns/internal/services/v1cacheservice"
	"github.com/rogerwesterbo/godns/internal/services/v1healthcheckservice"
	"github.com/rogerwesterbo/godns/internal/services/v1loadbalancerservice"
	"github.com/rogerwesterbo/godns/internal/services/v1querylogservice"
	"github.com/rogerwesterbo/godns/internal/services/v1ratelimitservice"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// AdminHandler handles administrative API endpoints for DNS features
type AdminHandler struct {
	cacheService *v1cacheservice.DNSCache
	rateLimiter  *v1ratelimitservice.RateLimiter
	loadBalancer *v1loadbalancerservice.LoadBalancer
	healthCheck  *v1healthcheckservice.HealthCheckService
	queryLog     *v1querylogservice.QueryLogService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(
	cacheService *v1cacheservice.DNSCache,
	rateLimiter *v1ratelimitservice.RateLimiter,
	loadBalancer *v1loadbalancerservice.LoadBalancer,
	healthCheck *v1healthcheckservice.HealthCheckService,
	queryLog *v1querylogservice.QueryLogService,
) *AdminHandler {
	return &AdminHandler{
		cacheService: cacheService,
		rateLimiter:  rateLimiter,
		loadBalancer: loadBalancer,
		healthCheck:  healthCheck,
		queryLog:     queryLog,
	}
}

// CacheStats represents cache statistics
type CacheStats struct {
	Enabled     bool   `json:"enabled"`
	CurrentSize int    `json:"current_size"`
	MaxSize     int    `json:"max_size"`
	TTLMinutes  int    `json:"ttl_minutes"`
	HitRate     string `json:"hit_rate,omitempty"`
}

// CacheStatsDetailed represents detailed cache statistics for the /cache/stats endpoint
type CacheStatsDetailed struct {
	Enabled   bool    `json:"enabled"`
	Size      int     `json:"size"`
	Capacity  int     `json:"capacity"`
	Hits      uint64  `json:"hits"`
	Misses    uint64  `json:"misses"`
	HitRate   float64 `json:"hit_rate"`
	Evictions uint64  `json:"evictions"`
}

// RateLimiterStats represents rate limiter statistics
type RateLimiterStats struct {
	Enabled        bool  `json:"enabled"`
	QPS            int   `json:"qps"`
	Burst          int   `json:"burst"`
	ActiveLimiters int   `json:"active_limiters"`
	TotalBlocked   int64 `json:"total_blocked"`
}

// LoadBalancerStats represents load balancer statistics
type LoadBalancerStats struct {
	Enabled         bool                      `json:"enabled"`
	Strategy        string                    `json:"strategy"`
	BackendGroups   int                       `json:"backend_groups"`
	TotalBackends   int                       `json:"total_backends"`
	HealthyBackends int                       `json:"healthy_backends"`
	Backends        []LoadBalancerBackendInfo `json:"backends,omitempty"`
}

// LoadBalancerBackendInfo represents backend information
type LoadBalancerBackendInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Weight  int    `json:"weight"`
	Healthy bool   `json:"healthy"`
	Enabled bool   `json:"enabled"`
}

// HealthCheckStats represents health check statistics
type HealthCheckStats struct {
	Enabled      bool                    `json:"enabled"`
	TotalTargets int                     `json:"total_targets"`
	Interval     int                     `json:"interval_seconds"`
	Timeout      int                     `json:"timeout_seconds"`
	Results      []HealthCheckResultInfo `json:"results,omitempty"`
}

// HealthCheckResultInfo represents health check result
type HealthCheckResultInfo struct {
	Target       string `json:"target"`
	Healthy      bool   `json:"healthy"`
	LastCheck    string `json:"last_check"`
	LastSuccess  string `json:"last_success,omitempty"`
	LastFailure  string `json:"last_failure,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// QueryLogStats represents query log statistics
type QueryLogStats struct {
	Enabled        bool    `json:"enabled"`
	TotalQueries   int64   `json:"total_queries"`
	CachedQueries  int64   `json:"cached_queries"`
	BlockedQueries int64   `json:"blocked_queries"`
	CacheHitRate   float64 `json:"cache_hit_rate"`
	BlockedRate    float64 `json:"blocked_rate"`
}

// SystemStats represents overall system statistics
type SystemStats struct {
	Cache        CacheStats        `json:"cache"`
	RateLimiter  RateLimiterStats  `json:"rate_limiter"`
	LoadBalancer LoadBalancerStats `json:"load_balancer"`
	HealthCheck  HealthCheckStats  `json:"health_check"`
	QueryLog     QueryLogStats     `json:"query_log"`
}

// GetSystemStats returns overall system statistics
// @Summary Get system statistics
// @Description Get comprehensive statistics about DNS caching, rate limiting, load balancing, health checks, and query logging
// @Tags Admin
// @Produce json
// @Success 200 {object} SystemStats
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/stats [get]
func (h *AdminHandler) GetSystemStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats := SystemStats{
		Cache:        h.getCacheStats(ctx),
		RateLimiter:  h.getRateLimiterStats(ctx),
		LoadBalancer: h.getLoadBalancerStats(ctx, false), // without backends list
		HealthCheck:  h.getHealthCheckStats(ctx, false),  // without results list
		QueryLog:     h.getQueryLogStats(ctx),
	}

	helpers.RespondJSON(w, stats)
}

// GetCacheStats returns cache statistics
// @Summary Get cache statistics
// @Description Get detailed statistics about DNS response caching
// @Tags Admin
// @Produce json
// @Success 200 {object} CacheStatsDetailed
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/cache/stats [get]
func (h *AdminHandler) GetCacheStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := h.getCacheStatsDetailed(ctx)
	helpers.RespondJSON(w, stats)
}

// ClearCache clears the DNS cache
// @Summary Clear DNS cache
// @Description Clear all entries from the DNS response cache
// @Tags Admin
// @Success 204 "Cache cleared successfully"
// @Failure 503 {object} map[string]string "Cache service not available"
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/cache/clear [post]
func (h *AdminHandler) ClearCache(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if h.cacheService == nil {
		http.Error(w, `{"error": "Cache service not enabled"}`, http.StatusServiceUnavailable)
		return
	}

	h.cacheService.Clear(ctx)
	vlog.Info("DNS cache cleared via API")
	w.WriteHeader(http.StatusNoContent)
}

// GetLoadBalancerStats returns load balancer statistics
// @Summary Get load balancer statistics
// @Description Get detailed statistics about load balancing including backend health status
// @Tags Admin
// @Produce json
// @Success 200 {object} LoadBalancerStats
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/loadbalancer/stats [get]
func (h *AdminHandler) GetLoadBalancerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := h.getLoadBalancerStats(ctx, true) // with backends list
	helpers.RespondJSON(w, stats)
}

// GetHealthCheckStats returns health check statistics
// @Summary Get health check statistics
// @Description Get detailed health check results for all monitored backends
// @Tags Admin
// @Produce json
// @Success 200 {object} HealthCheckStats
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/healthcheck/stats [get]
func (h *AdminHandler) GetHealthCheckStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := h.getHealthCheckStats(ctx, true) // with results list
	helpers.RespondJSON(w, stats)
}

// GetQueryLogStats returns query log statistics
// @Summary Get query log statistics
// @Description Get statistics about DNS queries including cache hit rate and blocked queries
// @Tags Admin
// @Produce json
// @Success 200 {object} QueryLogStats
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/querylog/stats [get]
func (h *AdminHandler) GetQueryLogStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := h.getQueryLogStats(ctx)
	helpers.RespondJSON(w, stats)
}

// GetRateLimiterStats returns rate limiter statistics
// @Summary Get rate limiter statistics
// @Description Get statistics about rate limiting including active limiters
// @Tags Admin
// @Produce json
// @Success 200 {object} RateLimiterStats
// @Security BearerAuth
// @Security OAuth2Password
// @Router /api/v1/admin/ratelimiter/stats [get]
func (h *AdminHandler) GetRateLimiterStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stats := h.getRateLimiterStats(ctx)
	helpers.RespondJSON(w, stats)
}

// Helper methods to gather stats from each service

func (h *AdminHandler) getCacheStats(ctx context.Context) CacheStats {
	stats := CacheStats{
		Enabled: h.cacheService != nil,
	}

	if h.cacheService != nil {
		hits, misses, size, hitRate := h.cacheService.Stats()
		stats.CurrentSize = size
		stats.MaxSize = 10000 // TODO: Get from config
		stats.TTLMinutes = 5  // TODO: Get from config

		// Format hit rate as string with percentage
		hitRateStr := "0.00%"
		if hitRate > 0 {
			hitRateStr = fmt.Sprintf("%.2f%%", hitRate)
		}
		stats.HitRate = hitRateStr

		// Store raw numbers for detailed stats endpoint
		vlog.Debugf("Cache stats - hits: %d, misses: %d, size: %d, hit_rate: %.2f%%",
			hits, misses, size, hitRate)
	}

	return stats
}

func (h *AdminHandler) getRateLimiterStats(ctx context.Context) RateLimiterStats {
	stats := RateLimiterStats{
		Enabled: h.rateLimiter != nil,
	}

	if h.rateLimiter != nil {
		activeIPs, rateLimit, burstSize := h.rateLimiter.Stats()
		stats.QPS = int(rateLimit)
		stats.Burst = burstSize
		stats.ActiveLimiters = activeIPs
		stats.TotalBlocked = 0 // TODO: Add blocked counter to rate limiter service
	}

	return stats
}

func (h *AdminHandler) getLoadBalancerStats(ctx context.Context, includeBackends bool) LoadBalancerStats {
	stats := LoadBalancerStats{
		Enabled: h.loadBalancer != nil,
	}

	if h.loadBalancer != nil {
		lbStats := h.loadBalancer.Stats()
		stats.Strategy = lbStats["strategy"].(string)
		stats.BackendGroups = lbStats["groups"].(int)
		stats.TotalBackends = lbStats["total_backends"].(int)
		stats.HealthyBackends = lbStats["healthy_backends"].(int)

		// TODO: Add backend details if includeBackends is true
	}

	return stats
}

func (h *AdminHandler) getHealthCheckStats(ctx context.Context, includeResults bool) HealthCheckStats {
	stats := HealthCheckStats{
		Enabled: h.healthCheck != nil,
	}

	// TODO: Add GetStats() method to health check service
	if h.healthCheck != nil {
		stats.TotalTargets = 0 // Placeholder
		stats.Interval = 30    // Default from config
		stats.Timeout = 5      // Default from config
	}

	return stats
}

func (h *AdminHandler) getQueryLogStats(ctx context.Context) QueryLogStats {
	stats := QueryLogStats{
		Enabled: h.queryLog != nil,
	}

	if h.queryLog != nil {
		logStats := h.queryLog.Stats()
		if total, ok := logStats["total_queries"].(uint64); ok {
			if total > math.MaxInt64 {
				stats.TotalQueries = math.MaxInt64
			} else {
				stats.TotalQueries = int64(total)
			}
		}
		if cached, ok := logStats["cached_queries"].(uint64); ok {
			if cached > math.MaxInt64 {
				stats.CachedQueries = math.MaxInt64
			} else {
				stats.CachedQueries = int64(cached)
			}
		}
		if blocked, ok := logStats["blocked_queries"].(uint64); ok {
			if blocked > math.MaxInt64 {
				stats.BlockedQueries = math.MaxInt64
			} else {
				stats.BlockedQueries = int64(blocked)
			}
		}
		// Convert from percentage (0-100) to decimal (0-1) for frontend consistency
		stats.CacheHitRate = logStats["cache_hit_rate"].(float64) / 100.0
		stats.BlockedRate = logStats["blocked_rate"].(float64) / 100.0
	}

	return stats
}

func (h *AdminHandler) getCacheStatsDetailed(ctx context.Context) CacheStatsDetailed {
	stats := CacheStatsDetailed{
		Enabled: h.cacheService != nil,
	}

	if h.cacheService != nil {
		hits, misses, size, hitRate := h.cacheService.Stats()
		stats.Size = size
		stats.Capacity = 10000 // TODO: Get from config
		stats.Hits = hits
		stats.Misses = misses
		stats.HitRate = hitRate / 100.0 // Convert from percentage (0-100) to decimal (0-1)
		stats.Evictions = 0             // TODO: Track evictions in cache service

		vlog.Debugf("getCacheStatsDetailed - raw hitRate: %.2f%%, converted to decimal: %.4f", hitRate, stats.HitRate)
	}

	return stats
}
