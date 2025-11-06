package v1healthcheckservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/vitistack/common/pkg/loggers/vlog"
)

// HealthCheckType defines the type of health check
type HealthCheckType int

const (
	// TCP checks if a TCP connection can be established
	TCP HealthCheckType = iota
	// HTTP checks if an HTTP endpoint returns 200 OK
	HTTP
	// HTTPS checks if an HTTPS endpoint returns 200 OK
	HTTPS
	// ICMP checks if the host responds to ping (requires privileges)
	ICMP
)

// HealthCheck represents a health check configuration
type HealthCheck struct {
	Target   string          // IP address or hostname
	Port     int             // Port to check
	Type     HealthCheckType // Type of check
	Interval time.Duration   // How often to check
	Timeout  time.Duration   // Timeout for each check
	Path     string          // Path for HTTP(S) checks
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Target    string
	Healthy   bool
	LastCheck time.Time
	Message   string
	Latency   time.Duration
}

// HealthCheckService manages health checks for backends
type HealthCheckService struct {
	checks  map[string]*HealthCheck
	results map[string]*HealthCheckResult
	mu      sync.RWMutex
	stopCh  chan struct{}
	wg      sync.WaitGroup
}

// NewHealthCheckService creates a new health check service
func NewHealthCheckService() *HealthCheckService {
	return &HealthCheckService{
		checks:  make(map[string]*HealthCheck),
		results: make(map[string]*HealthCheckResult),
		stopCh:  make(chan struct{}),
	}
}

// AddCheck adds a health check for a target
func (hcs *HealthCheckService) AddCheck(ctx context.Context, target string, check HealthCheck) {
	hcs.mu.Lock()
	defer hcs.mu.Unlock()

	// Use target as key
	hcs.checks[target] = &check

	// Initialize result
	hcs.results[target] = &HealthCheckResult{
		Target:    target,
		Healthy:   true, // Assume healthy until first check
		LastCheck: time.Now(),
		Message:   "Waiting for first health check",
	}

	// Start health check goroutine
	hcs.wg.Add(1)
	go hcs.runHealthCheck(target, &check)

	vlog.Infof("Added health check for %s (type: %s, interval: %s)", target, check.Type.String(), check.Interval)
}

// RemoveCheck removes a health check for a target
func (hcs *HealthCheckService) RemoveCheck(ctx context.Context, target string) {
	hcs.mu.Lock()
	defer hcs.mu.Unlock()

	delete(hcs.checks, target)
	delete(hcs.results, target)

	vlog.Infof("Removed health check for %s", target)
}

// IsHealthy returns the health status of a target
func (hcs *HealthCheckService) IsHealthy(ctx context.Context, target string) bool {
	hcs.mu.RLock()
	defer hcs.mu.RUnlock()

	result, exists := hcs.results[target]
	if !exists {
		return true // If no health check configured, assume healthy
	}

	return result.Healthy
}

// GetResult returns the health check result for a target
func (hcs *HealthCheckService) GetResult(ctx context.Context, target string) (*HealthCheckResult, bool) {
	hcs.mu.RLock()
	defer hcs.mu.RUnlock()

	result, exists := hcs.results[target]
	if !exists {
		return nil, false
	}

	// Return a copy
	resultCopy := *result
	return &resultCopy, true
}

// GetAllResults returns all health check results
func (hcs *HealthCheckService) GetAllResults(ctx context.Context) map[string]*HealthCheckResult {
	hcs.mu.RLock()
	defer hcs.mu.RUnlock()

	results := make(map[string]*HealthCheckResult)
	for k, v := range hcs.results {
		resultCopy := *v
		results[k] = &resultCopy
	}

	return results
}

// Stop stops all health checks
func (hcs *HealthCheckService) Stop() {
	close(hcs.stopCh)
	hcs.wg.Wait()
	vlog.Info("Health check service stopped")
}

// runHealthCheck runs a health check periodically
func (hcs *HealthCheckService) runHealthCheck(target string, check *HealthCheck) {
	defer hcs.wg.Done()

	ticker := time.NewTicker(check.Interval)
	defer ticker.Stop()

	// Run first check immediately
	hcs.performCheck(target, check)

	for {
		select {
		case <-ticker.C:
			hcs.performCheck(target, check)
		case <-hcs.stopCh:
			return
		}
	}
}

// performCheck performs a single health check
func (hcs *HealthCheckService) performCheck(target string, check *HealthCheck) {
	start := time.Now()
	var healthy bool
	var message string

	switch check.Type {
	case TCP:
		healthy, message = hcs.checkTCP(check.Target, check.Port, check.Timeout)
	case HTTP:
		healthy, message = hcs.checkHTTP(check.Target, check.Port, check.Path, check.Timeout, false)
	case HTTPS:
		healthy, message = hcs.checkHTTP(check.Target, check.Port, check.Path, check.Timeout, true)
	case ICMP:
		healthy, message = hcs.checkICMP(check.Target, check.Timeout)
	default:
		healthy = false
		message = "Unknown health check type"
	}

	latency := time.Since(start)

	// Update result
	hcs.mu.Lock()
	hcs.results[target] = &HealthCheckResult{
		Target:    target,
		Healthy:   healthy,
		LastCheck: time.Now(),
		Message:   message,
		Latency:   latency,
	}
	hcs.mu.Unlock()

	if !healthy {
		vlog.Warnf("Health check failed for %s: %s (latency: %s)", target, message, latency)
	} else {
		vlog.Debugf("Health check passed for %s (latency: %s)", target, latency)
	}
}

// checkTCP performs a TCP health check
func (hcs *HealthCheckService) checkTCP(host string, port int, timeout time.Duration) (bool, string) {
	// Use net.JoinHostPort to properly handle IPv6 addresses
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false, fmt.Sprintf("TCP connection failed: %v", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			vlog.Debugf("Error closing TCP connection to %s: %v", address, closeErr)
		}
	}()

	return true, "TCP connection successful"
}

// checkHTTP performs an HTTP/HTTPS health check
func (hcs *HealthCheckService) checkHTTP(host string, port int, path string, timeout time.Duration, https bool) (bool, string) {
	scheme := "http"
	if https {
		scheme = "https"
	}

	if path == "" {
		path = "/"
	}

	url := fmt.Sprintf("%s://%s:%d%s", scheme, host, port, path)

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, fmt.Sprintf("HTTP request failed: %v", err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			vlog.Debugf("Error closing HTTP response body for %s: %v", url, closeErr)
		}
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, fmt.Sprintf("HTTP %d OK", resp.StatusCode)
	}

	return false, fmt.Sprintf("HTTP %d", resp.StatusCode)
}

// checkICMP performs an ICMP ping health check
func (hcs *HealthCheckService) checkICMP(host string, timeout time.Duration) (bool, string) {
	// ICMP requires raw sockets which need elevated privileges
	// For now, we'll use a TCP check as fallback
	// In production, you'd use a library like github.com/go-ping/ping
	return hcs.checkTCP(host, 80, timeout)
}

// Stats returns health check statistics
func (hcs *HealthCheckService) Stats() map[string]interface{} {
	hcs.mu.RLock()
	defer hcs.mu.RUnlock()

	totalChecks := len(hcs.checks)
	healthyCount := 0

	for _, result := range hcs.results {
		if result.Healthy {
			healthyCount++
		}
	}

	return map[string]interface{}{
		"total_checks":    totalChecks,
		"healthy_count":   healthyCount,
		"unhealthy_count": totalChecks - healthyCount,
	}
}

// String returns the health check type name
func (t HealthCheckType) String() string {
	switch t {
	case TCP:
		return "TCP"
	case HTTP:
		return "HTTP"
	case HTTPS:
		return "HTTPS"
	case ICMP:
		return "ICMP"
	default:
		return "Unknown"
	}
}
