package v1loadbalancerservice

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/internal/models"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// LoadBalancerStrategy defines the load balancing algorithm
type LoadBalancerStrategy int

const (
	// RoundRobin distributes requests evenly across all backends
	RoundRobin LoadBalancerStrategy = iota
	// WeightedRoundRobin distributes based on weights
	WeightedRoundRobin
	// LeastConnections selects backend with fewest active connections
	LeastConnections
	// Random selects a random backend
	Random
)

// Backend represents a single backend server
type Backend struct {
	Record      models.DNSRecord
	Weight      int          // For weighted load balancing
	Healthy     bool         // Health check status
	Enabled     bool         // Can be manually disabled
	Connections atomic.Int32 // Active connections (for least connections)
}

// LoadBalancer manages multiple backends for a DNS record
type LoadBalancer struct {
	backends map[string]*BackendGroup // key: name:type (e.g., "api.example.lan.:A")
	mu       sync.RWMutex
	strategy LoadBalancerStrategy
}

// BackendGroup manages a group of backends for the same DNS name/type
type BackendGroup struct {
	Backends []*Backend
	Counter  atomic.Uint64 // For round-robin
	Strategy LoadBalancerStrategy
	mu       sync.RWMutex
}

// NewLoadBalancer creates a new load balancer
func NewLoadBalancer(strategy LoadBalancerStrategy) *LoadBalancer {
	return &LoadBalancer{
		backends: make(map[string]*BackendGroup),
		strategy: strategy,
	}
}

// AddBackend adds a backend to the load balancer
func (lb *LoadBalancer) AddBackend(ctx context.Context, record models.DNSRecord, weight int) {
	key := makeKey(record.Name, record.Type)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	group, exists := lb.backends[key]
	if !exists {
		group = &BackendGroup{
			Backends: make([]*Backend, 0),
			Strategy: lb.strategy,
		}
		lb.backends[key] = group
	}

	backend := &Backend{
		Record:  record,
		Weight:  weight,
		Healthy: true, // Assume healthy initially
		Enabled: true,
	}

	group.Backends = append(group.Backends, backend)
	vlog.Debugf("Added backend: %s -> %s (weight: %d)", record.Name, record.Value, weight)
}

// GetBackend returns the next backend according to the load balancing strategy
func (lb *LoadBalancer) GetBackend(ctx context.Context, name, recordType string) (*models.DNSRecord, bool) {
	key := makeKey(name, recordType)

	lb.mu.RLock()
	group, exists := lb.backends[key]
	lb.mu.RUnlock()

	if !exists || len(group.Backends) == 0 {
		return nil, false
	}

	return group.Next()
}

// GetAllHealthyBackends returns all healthy backends for a name/type
func (lb *LoadBalancer) GetAllHealthyBackends(ctx context.Context, name, recordType string) []models.DNSRecord {
	key := makeKey(name, recordType)

	lb.mu.RLock()
	group, exists := lb.backends[key]
	lb.mu.RUnlock()

	if !exists {
		return nil
	}

	return group.GetHealthy()
}

// SetBackendHealth updates the health status of a backend
func (lb *LoadBalancer) SetBackendHealth(ctx context.Context, name, recordType, value string, healthy bool) {
	key := makeKey(name, recordType)

	lb.mu.RLock()
	group, exists := lb.backends[key]
	lb.mu.RUnlock()

	if !exists {
		return
	}

	group.SetHealth(value, healthy)
}

// RemoveBackend removes a backend from the load balancer
func (lb *LoadBalancer) RemoveBackend(ctx context.Context, name, recordType, value string) {
	key := makeKey(name, recordType)

	lb.mu.Lock()
	defer lb.mu.Unlock()

	group, exists := lb.backends[key]
	if !exists {
		return
	}

	group.Remove(value)

	// Remove group if empty
	if len(group.Backends) == 0 {
		delete(lb.backends, key)
	}
}

// Stats returns load balancer statistics
func (lb *LoadBalancer) Stats() map[string]interface{} {
	lb.mu.RLock()
	defer lb.mu.RUnlock()

	totalBackends := 0
	healthyBackends := 0
	groups := len(lb.backends)

	for _, group := range lb.backends {
		group.mu.RLock()
		totalBackends += len(group.Backends)
		for _, backend := range group.Backends {
			if backend.Healthy && backend.Enabled {
				healthyBackends++
			}
		}
		group.mu.RUnlock()
	}

	return map[string]interface{}{
		"groups":           groups,
		"total_backends":   totalBackends,
		"healthy_backends": healthyBackends,
		"strategy":         lb.strategy.String(),
	}
}

// BackendGroup methods

// Next returns the next backend according to the strategy
func (bg *BackendGroup) Next() (*models.DNSRecord, bool) {
	bg.mu.RLock()
	defer bg.mu.RUnlock()

	// Filter healthy and enabled backends
	available := make([]*Backend, 0)
	for _, backend := range bg.Backends {
		if backend.Healthy && backend.Enabled {
			available = append(available, backend)
		}
	}

	if len(available) == 0 {
		return nil, false
	}

	var selected *Backend

	switch bg.Strategy {
	case RoundRobin:
		selected = bg.roundRobin(available)
	case WeightedRoundRobin:
		selected = bg.weightedRoundRobin(available)
	case LeastConnections:
		selected = bg.leastConnections(available)
	case Random:
		selected = bg.random(available)
	default:
		selected = bg.roundRobin(available)
	}

	if selected != nil {
		record := selected.Record
		return &record, true
	}

	return nil, false
}

// GetHealthy returns all healthy backends
func (bg *BackendGroup) GetHealthy() []models.DNSRecord {
	bg.mu.RLock()
	defer bg.mu.RUnlock()

	result := make([]models.DNSRecord, 0)
	for _, backend := range bg.Backends {
		if backend.Healthy && backend.Enabled {
			result = append(result, backend.Record)
		}
	}

	return result
}

// SetHealth updates health status
func (bg *BackendGroup) SetHealth(value string, healthy bool) {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	for _, backend := range bg.Backends {
		if backend.Record.Value == value {
			backend.Healthy = healthy
			status := "unhealthy"
			if healthy {
				status = "healthy"
			}
			vlog.Debugf("Backend %s marked as %s", value, status)
		}
	}
}

// Remove removes a backend by value
func (bg *BackendGroup) Remove(value string) {
	bg.mu.Lock()
	defer bg.mu.Unlock()

	newBackends := make([]*Backend, 0)
	for _, backend := range bg.Backends {
		if backend.Record.Value != value {
			newBackends = append(newBackends, backend)
		}
	}
	bg.Backends = newBackends
}

// Load balancing algorithms

func (bg *BackendGroup) roundRobin(available []*Backend) *Backend {
	if len(available) == 0 {
		return nil
	}

	idx := bg.Counter.Add(1) % uint64(len(available))
	return available[idx]
}

func (bg *BackendGroup) weightedRoundRobin(available []*Backend) *Backend {
	if len(available) == 0 {
		return nil
	}

	// Calculate total weight as uint64 to avoid overflow issues
	var totalWeight uint64
	for _, backend := range available {
		weight := backend.Weight
		if weight <= 0 {
			weight = 1
		}
		// #nosec G115 -- Safe conversion: weight is validated to be >= 1
		totalWeight += uint64(weight)
	}

	if totalWeight == 0 {
		return bg.roundRobin(available)
	}

	// Select based on weight
	counterVal := bg.Counter.Add(1)
	index := counterVal % totalWeight
	var currentWeight uint64

	for _, backend := range available {
		weight := backend.Weight
		if weight <= 0 {
			weight = 1
		}
		// #nosec G115 -- Safe conversion: weight is validated to be >= 1
		currentWeight += uint64(weight)
		if index < currentWeight {
			return backend
		}
	}

	return available[0]
}

func (bg *BackendGroup) leastConnections(available []*Backend) *Backend {
	if len(available) == 0 {
		return nil
	}

	minConns := int32(1<<31 - 1) // Max int32
	var selected *Backend

	for _, backend := range available {
		conns := backend.Connections.Load()
		if conns < minConns {
			minConns = conns
			selected = backend
		}
	}

	if selected != nil {
		selected.Connections.Add(1)
	}

	return selected
}

func (bg *BackendGroup) random(available []*Backend) *Backend {
	if len(available) == 0 {
		return nil
	}

	// Use counter as pseudo-random
	idx := bg.Counter.Add(1) % uint64(len(available))
	return available[idx]
}

// Helper functions

func makeKey(name, recordType string) string {
	return dns.Fqdn(name) + ":" + recordType
}

// String returns the strategy name
func (s LoadBalancerStrategy) String() string {
	switch s {
	case RoundRobin:
		return "round-robin"
	case WeightedRoundRobin:
		return "weighted-round-robin"
	case LeastConnections:
		return "least-connections"
	case Random:
		return "random"
	default:
		return "unknown"
	}
}
