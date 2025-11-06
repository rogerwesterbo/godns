package v1querylogservice

import (
	"context"
	"encoding/json"
	"net/netip"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// QueryLog represents a single DNS query log entry
type QueryLog struct {
	Timestamp    time.Time `json:"timestamp"`
	ClientIP     string    `json:"client_ip"`
	QueryName    string    `json:"query_name"`
	QueryType    string    `json:"query_type"`
	ResponseCode string    `json:"response_code"`
	AnswerCount  int       `json:"answer_count"`
	Latency      int64     `json:"latency_ms"`
	CacheHit     bool      `json:"cache_hit"`
	Upstream     bool      `json:"upstream"`
	Blocked      bool      `json:"blocked"`
}

// QueryLogService manages DNS query logging
type QueryLogService struct {
	enabled       atomic.Bool
	logToConsole  bool
	buffer        []*QueryLog
	bufferSize    int
	mu            sync.RWMutex
	flushInterval time.Duration
	stopCh        chan struct{}
	valkeyClient  valkeyinterface.ValkeyInterface

	// Statistics
	totalQueries   atomic.Uint64
	cachedQueries  atomic.Uint64
	blockedQueries atomic.Uint64
}

// NewQueryLogService creates a new query log service
func NewQueryLogService(bufferSize int, flushInterval time.Duration, valkeyClient valkeyinterface.ValkeyInterface) *QueryLogService {
	qls := &QueryLogService{
		buffer:        make([]*QueryLog, 0, bufferSize),
		bufferSize:    bufferSize,
		flushInterval: flushInterval,
		logToConsole:  true, // Default to console logging
		stopCh:        make(chan struct{}),
		valkeyClient:  valkeyClient,
	}

	qls.enabled.Store(true)

	// Load persisted stats from Valkey
	if valkeyClient != nil {
		qls.loadStatsFromValkey(context.Background())
	}

	// Start periodic flush and stats persistence
	go qls.periodicFlush()
	if valkeyClient != nil {
		go qls.periodicStatsPersistence()
	}

	return qls
}

// LogQuery logs a DNS query
func (qls *QueryLogService) LogQuery(ctx context.Context,
	clientIP netip.Addr,
	question dns.Question,
	response *dns.Msg,
	latency time.Duration,
	cacheHit bool,
	upstream bool,
	blocked bool) {

	if !qls.enabled.Load() {
		return
	}

	// Update statistics
	qls.totalQueries.Add(1)
	if cacheHit {
		qls.cachedQueries.Add(1)
	}
	if blocked {
		qls.blockedQueries.Add(1)
	}

	responseCode := "NOERROR"
	answerCount := 0

	if response != nil {
		responseCode = dns.RcodeToString[response.Rcode]
		answerCount = len(response.Answer)
	}

	log := &QueryLog{
		Timestamp:    time.Now(),
		ClientIP:     clientIP.String(),
		QueryName:    question.Name,
		QueryType:    dns.TypeToString[question.Qtype],
		ResponseCode: responseCode,
		AnswerCount:  answerCount,
		Latency:      latency.Milliseconds(),
		CacheHit:     cacheHit,
		Upstream:     upstream,
		Blocked:      blocked,
	}

	// Log to console if enabled
	if qls.logToConsole {
		qls.logToConsoleFunc(log)
	}

	// Add to buffer
	qls.mu.Lock()
	qls.buffer = append(qls.buffer, log)
	shouldFlush := len(qls.buffer) >= qls.bufferSize
	qls.mu.Unlock()

	// Flush if buffer is full
	if shouldFlush {
		qls.Flush()
	}
}

// Flush writes buffered logs
func (qls *QueryLogService) Flush() {
	qls.mu.Lock()
	if len(qls.buffer) == 0 {
		qls.mu.Unlock()
		return
	}

	// Get buffer and create new one
	buffer := qls.buffer
	qls.buffer = make([]*QueryLog, 0, qls.bufferSize)
	qls.mu.Unlock()

	// In a production system, you'd write to a file or send to a log aggregator
	// For now, we'll just log the count
	vlog.Debugf("Flushed %d query logs", len(buffer))
}

// GetRecentQueries returns the most recent queries from the buffer
func (qls *QueryLogService) GetRecentQueries(ctx context.Context, count int) []*QueryLog {
	qls.mu.RLock()
	defer qls.mu.RUnlock()

	if count > len(qls.buffer) {
		count = len(qls.buffer)
	}

	// Return the last N entries
	start := len(qls.buffer) - count
	if start < 0 {
		start = 0
	}

	result := make([]*QueryLog, count)
	copy(result, qls.buffer[start:])

	return result
}

// Stats returns query logging statistics
func (qls *QueryLogService) Stats() map[string]interface{} {
	total := qls.totalQueries.Load()
	cached := qls.cachedQueries.Load()
	blocked := qls.blockedQueries.Load()

	cacheHitRate := 0.0
	if total > 0 {
		cacheHitRate = float64(cached) / float64(total) * 100
	}

	blockedRate := 0.0
	if total > 0 {
		blockedRate = float64(blocked) / float64(total) * 100
	}

	qls.mu.RLock()
	bufferSize := len(qls.buffer)
	qls.mu.RUnlock()

	return map[string]interface{}{
		"total_queries":   total,
		"cached_queries":  cached,
		"blocked_queries": blocked,
		"cache_hit_rate":  cacheHitRate,
		"blocked_rate":    blockedRate,
		"buffer_size":     bufferSize,
		"enabled":         qls.enabled.Load(),
	}
}

// Enable enables query logging
func (qls *QueryLogService) Enable() {
	qls.enabled.Store(true)
	vlog.Info("Query logging enabled")
}

// Disable disables query logging
func (qls *QueryLogService) Disable() {
	qls.enabled.Store(false)
	vlog.Info("Query logging disabled")
}

// SetLogToConsole enables/disables console logging
func (qls *QueryLogService) SetLogToConsole(enabled bool) {
	qls.logToConsole = enabled
}

// Stop stops the query log service
func (qls *QueryLogService) Stop() {
	close(qls.stopCh)
	qls.Flush()

	// Persist final stats before stopping
	if qls.valkeyClient != nil {
		qls.persistStatsToValkey(context.Background())
	}

	vlog.Info("Query log service stopped")
}

// periodicFlush flushes the buffer periodically
func (qls *QueryLogService) periodicFlush() {
	ticker := time.NewTicker(qls.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			qls.Flush()
		case <-qls.stopCh:
			return
		}
	}
}

// periodicStatsPersistence persists stats to Valkey periodically
func (qls *QueryLogService) periodicStatsPersistence() {
	ticker := time.NewTicker(30 * time.Second) // Persist every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			qls.persistStatsToValkey(context.Background())
		case <-qls.stopCh:
			return
		}
	}
}

// persistStatsToValkey saves query statistics to Valkey
func (qls *QueryLogService) persistStatsToValkey(ctx context.Context) {
	if qls.valkeyClient == nil {
		return
	}

	stats := map[string]string{
		"dns:stats:total_queries":   strconv.FormatUint(qls.totalQueries.Load(), 10),
		"dns:stats:cached_queries":  strconv.FormatUint(qls.cachedQueries.Load(), 10),
		"dns:stats:blocked_queries": strconv.FormatUint(qls.blockedQueries.Load(), 10),
	}

	for key, value := range stats {
		if err := qls.valkeyClient.SetData(ctx, key, value); err != nil {
			vlog.Warnf("Failed to persist query stat %s: %v", key, err)
		}
	}
}

// loadStatsFromValkey loads query statistics from Valkey
func (qls *QueryLogService) loadStatsFromValkey(ctx context.Context) {
	if qls.valkeyClient == nil {
		return
	}

	// Load total queries
	if val, err := qls.valkeyClient.GetData(ctx, "dns:stats:total_queries"); err == nil && val != "" {
		if count, err := strconv.ParseUint(val, 10, 64); err == nil {
			qls.totalQueries.Store(count)
			vlog.Infof("Loaded total_queries from Valkey: %d", count)
		}
	}

	// Load cached queries
	if val, err := qls.valkeyClient.GetData(ctx, "dns:stats:cached_queries"); err == nil && val != "" {
		if count, err := strconv.ParseUint(val, 10, 64); err == nil {
			qls.cachedQueries.Store(count)
			vlog.Infof("Loaded cached_queries from Valkey: %d", count)
		}
	}

	// Load blocked queries
	if val, err := qls.valkeyClient.GetData(ctx, "dns:stats:blocked_queries"); err == nil && val != "" {
		if count, err := strconv.ParseUint(val, 10, 64); err == nil {
			qls.blockedQueries.Store(count)
			vlog.Infof("Loaded blocked_queries from Valkey: %d", count)
		}
	}
}

// logToConsoleFunc logs to console in a readable format
func (qls *QueryLogService) logToConsoleFunc(log *QueryLog) {
	// Create JSON for structured logging
	jsonBytes, err := json.Marshal(log)
	if err != nil {
		vlog.Warnf("Failed to marshal query log: %v", err)
		return
	}

	// Log as JSON for easy parsing
	vlog.Infof("DNS_QUERY: %s", string(jsonBytes))
}
