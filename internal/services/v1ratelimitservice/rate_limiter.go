package v1ratelimitservice

import (
	"context"
	"net/netip"
	"sync"
	"time"

	"github.com/vitistack/common/pkg/loggers/vlog"
	"golang.org/x/time/rate"
)

// RateLimiter implements per-IP rate limiting using token bucket algorithm
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit // queries per second
	burst    int        // burst size
	cleanup  time.Duration
}

// NewRateLimiter creates a new rate limiter
// ratePerSecond: maximum queries per second per IP
// burst: maximum burst size (allows short bursts above the rate)
func NewRateLimiter(ratePerSecond int, burst int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(ratePerSecond),
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	// Start cleanup goroutine to remove old limiters
	go rl.cleanupLimiters()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ctx context.Context, ip netip.Addr) bool {
	if !ip.IsValid() {
		// If IP is invalid, allow but log warning
		vlog.Warn("Invalid IP address for rate limiting")
		return true
	}

	ipStr := ip.String()

	rl.mu.RLock()
	limiter, exists := rl.limiters[ipStr]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check after acquiring write lock
		limiter, exists = rl.limiters[ipStr]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ipStr] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter.Allow()
}

// AllowN checks if N requests from the given IP should be allowed
func (rl *RateLimiter) AllowN(ctx context.Context, ip netip.Addr, n int) bool {
	if !ip.IsValid() {
		vlog.Warn("Invalid IP address for rate limiting")
		return true
	}

	ipStr := ip.String()

	rl.mu.RLock()
	limiter, exists := rl.limiters[ipStr]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter, exists = rl.limiters[ipStr]
		if !exists {
			limiter = rate.NewLimiter(rl.rate, rl.burst)
			rl.limiters[ipStr] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter.AllowN(time.Now(), n)
}

// Reset resets the rate limit for a specific IP
func (rl *RateLimiter) Reset(ctx context.Context, ip netip.Addr) {
	if !ip.IsValid() {
		return
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	delete(rl.limiters, ip.String())
}

// Clear removes all rate limiters
func (rl *RateLimiter) Clear(ctx context.Context) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.limiters = make(map[string]*rate.Limiter)
	vlog.Info("Rate limiters cleared")
}

// Stats returns rate limiter statistics
func (rl *RateLimiter) Stats() (activeIPs int, rateLimit float64, burstSize int) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return len(rl.limiters), float64(rl.rate), rl.burst
}

// cleanupLimiters removes limiters that haven't been used recently
func (rl *RateLimiter) cleanupLimiters() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()

		// In a production system, you'd track last access time
		// For now, we'll just clear if we have too many entries
		if len(rl.limiters) > 10000 {
			// Clear half of the entries (simple approach)
			newLimiters := make(map[string]*rate.Limiter)
			count := 0
			for k, v := range rl.limiters {
				if count < 5000 {
					newLimiters[k] = v
					count++
				}
			}
			rl.limiters = newLimiters
			vlog.Infof("Cleaned up rate limiters, kept %d entries", len(rl.limiters))
		}

		rl.mu.Unlock()
	}
}

// UpdateLimits updates the rate and burst settings
func (rl *RateLimiter) UpdateLimits(ratePerSecond int, burst int) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.rate = rate.Limit(ratePerSecond)
	rl.burst = burst

	// Update all existing limiters
	for _, limiter := range rl.limiters {
		limiter.SetLimit(rl.rate)
		limiter.SetBurst(rl.burst)
	}

	vlog.Infof("Updated rate limits: %d queries/sec, burst: %d", ratePerSecond, burst)
}
