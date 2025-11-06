package v1cacheservice

import (
	"context"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// CacheEntry represents a cached DNS response
type CacheEntry struct {
	Response  *dns.Msg
	ExpiresAt time.Time
}

// DNSCache implements a thread-safe DNS response cache with TTL and LRU eviction
type DNSCache struct {
	entries    map[string]*CacheEntry
	mu         sync.RWMutex
	maxSize    int
	ttl        time.Duration
	accessList *accessList // For LRU eviction
	hits       uint64
	misses     uint64
}

// accessList implements a simple LRU tracking
type accessList struct {
	keys []string
	mu   sync.Mutex
}

// NewDNSCache creates a new DNS cache with the specified max size and default TTL
func NewDNSCache(maxSize int, defaultTTL time.Duration) *DNSCache {
	cache := &DNSCache{
		entries:    make(map[string]*CacheEntry),
		maxSize:    maxSize,
		ttl:        defaultTTL,
		accessList: &accessList{keys: make([]string, 0, maxSize)},
	}

	// Start background cleanup goroutine
	go cache.cleanupExpired()

	return cache
}

// Get retrieves a DNS response from the cache
func (c *DNSCache) Get(ctx context.Context, key string) (*dns.Msg, bool) {
	c.mu.RLock()
	entry, exists := c.entries[key]
	c.mu.RUnlock()

	if !exists {
		c.misses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		c.mu.Lock()
		delete(c.entries, key)
		c.mu.Unlock()
		c.misses++
		return nil, false
	}

	// Update access list for LRU
	c.accessList.touch(key)

	c.hits++
	// Return a copy to avoid modifications
	return entry.Response.Copy(), true
}

// Set stores a DNS response in the cache
func (c *DNSCache) Set(ctx context.Context, key string, response *dns.Msg) {
	if response == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict
	if len(c.entries) >= c.maxSize {
		c.evictLRU()
	}

	// Determine TTL from the response or use default
	ttl := c.ttl
	if len(response.Answer) > 0 {
		// Use the minimum TTL from all records
		minTTL := uint32(c.ttl.Seconds())
		for _, rr := range response.Answer {
			if rr.Header().Ttl < minTTL {
				minTTL = rr.Header().Ttl
			}
		}
		if minTTL > 0 {
			ttl = time.Duration(minTTL) * time.Second
		}
	}

	c.entries[key] = &CacheEntry{
		Response:  response.Copy(),
		ExpiresAt: time.Now().Add(ttl),
	}

	c.accessList.add(key)
}

// Delete removes an entry from the cache
func (c *DNSCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
	c.accessList.remove(key)
}

// Clear removes all entries from the cache
func (c *DNSCache) Clear(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*CacheEntry)
	c.accessList.clear()
	vlog.Info("DNS cache cleared")
}

// Stats returns cache statistics
func (c *DNSCache) Stats() (hits, misses uint64, size int, hitRate float64) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate = 0.0
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}

	return c.hits, c.misses, len(c.entries), hitRate
}

// evictLRU removes the least recently used entry
func (c *DNSCache) evictLRU() {
	key := c.accessList.removeLRU()
	if key != "" {
		delete(c.entries, key)
		vlog.Debugf("Evicted cache entry: %s", key)
	}
}

// cleanupExpired removes expired entries periodically
func (c *DNSCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		removed := 0

		for key, entry := range c.entries {
			if now.After(entry.ExpiresAt) {
				delete(c.entries, key)
				c.accessList.remove(key)
				removed++
			}
		}

		if removed > 0 {
			vlog.Debugf("Cleaned up %d expired cache entries", removed)
		}
		c.mu.Unlock()
	}
}

// MakeCacheKey creates a cache key from a DNS question
func MakeCacheKey(q dns.Question) string {
	return dns.Fqdn(q.Name) + ":" + dns.TypeToString[q.Qtype]
}

// accessList methods

func (al *accessList) touch(key string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Remove key if it exists
	for i, k := range al.keys {
		if k == key {
			al.keys = append(al.keys[:i], al.keys[i+1:]...)
			break
		}
	}

	// Add to end (most recently used)
	al.keys = append(al.keys, key)
}

func (al *accessList) add(key string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	// Check if already exists
	for _, k := range al.keys {
		if k == key {
			return
		}
	}

	al.keys = append(al.keys, key)
}

func (al *accessList) remove(key string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	for i, k := range al.keys {
		if k == key {
			al.keys = append(al.keys[:i], al.keys[i+1:]...)
			return
		}
	}
}

func (al *accessList) removeLRU() string {
	al.mu.Lock()
	defer al.mu.Unlock()

	if len(al.keys) == 0 {
		return ""
	}

	// First item is least recently used
	key := al.keys[0]
	al.keys = al.keys[1:]
	return key
}

func (al *accessList) clear() {
	al.mu.Lock()
	defer al.mu.Unlock()

	al.keys = make([]string, 0)
}
