# GoDNS Enhanced Features Guide

This guide covers the advanced features added to GoDNS for production deployments: caching, rate limiting, load balancing, health checks, query logging, and metrics.

## Table of Contents

1. [DNS Response Caching](#dns-response-caching)
2. [Rate Limiting](#rate-limiting)
3. [Load Balancing](#load-balancing)
4. [Health Checks](#health-checks)
5. [Query Logging](#query-logging)
6. [Prometheus Metrics](#prometheus-metrics)
7. [Configuration Reference](#configuration-reference)
8. [Testing Examples](#testing-examples)

---

## DNS Response Caching

### Overview

DNS caching reduces latency and load on upstream servers by storing successful DNS responses in memory with configurable TTL (Time To Live).

### Features

- **LRU Eviction**: Automatically evicts least recently used entries when cache is full
- **TTL-based Expiration**: Cached entries expire after configured time
- **Background Cleanup**: Periodic cleanup of expired entries
- **Thread-safe**: Concurrent access from multiple goroutines

### Configuration

```bash
# Enable DNS caching
DNS_ENABLE_CACHE=true

# Maximum number of cached DNS responses (default: 10000)
DNS_CACHE_SIZE=10000

# Cache entry TTL in minutes (default: 5)
DNS_CACHE_TTL_MINUTES=5
```

### How It Works

1. **Cache Key**: `<domain>:<query_type>` (e.g., `www.example.lan.:A`)
2. **Lookup Flow**:
   - Check cache first
   - If hit: return cached response immediately
   - If miss: lookup from Valkey or upstream, then cache the result

### Metrics

- `godns_cache_hits_total`: Number of cache hits
- `godns_cache_misses_total`: Number of cache misses
- `godns_cache_size`: Current number of cached entries
- `godns_cache_evictions_total`: Number of cache evictions

---

## Rate Limiting

### Overview

Protects against DNS amplification attacks and abuse by limiting queries per IP address using the token bucket algorithm.

### Features

- **Per-IP Tracking**: Each client IP has its own rate limit
- **Token Bucket Algorithm**: Allows bursts while maintaining average rate
- **Automatic Cleanup**: Removes inactive limiters after 10 minutes
- **Graceful Refusal**: Returns DNS `REFUSED` response when limit exceeded

### Configuration

```bash
# Enable rate limiting
DNS_ENABLE_RATE_LIMIT=true

# Queries per second per IP (default: 100)
DNS_RATE_LIMIT_QPS=100

# Burst size (default: 200)
DNS_RATE_LIMIT_BURST=200
```

### How It Works

1. **Token Bucket**: Each IP gets a bucket with `BURST` tokens
2. **Refill Rate**: Tokens refill at `QPS` rate
3. **Query Processing**:
   - If tokens available: process query (consume 1 token)
   - If no tokens: return `REFUSED` response

### Example Rate Calculation

With `QPS=100` and `BURST=200`:

- Client can send 200 queries instantly
- Then limited to 100 queries/second
- Unused capacity accumulates up to 200 tokens

### Metrics

- `godns_rate_limited_total`: Number of rate-limited queries
- `godns_rate_limiters_active`: Number of active rate limiters

---

## Load Balancing

### Overview

Distributes DNS queries across multiple backend servers to improve availability and performance.

### Supported Strategies

#### 1. Round Robin (default)

Distributes requests evenly in circular order.

```bash
DNS_LOAD_BALANCER_STRATEGY=round-robin
```

#### 2. Weighted Round Robin

Distributes based on backend weights (higher weight = more traffic).

```bash
DNS_LOAD_BALANCER_STRATEGY=weighted-round-robin
```

#### 3. Least Connections

Routes to backend with fewest active connections.

```bash
DNS_LOAD_BALANCER_STRATEGY=least-connections
```

#### 4. Random

Randomly selects a healthy backend.

```bash
DNS_LOAD_BALANCER_STRATEGY=random
```

### Configuration

```bash
# Enable load balancing
DNS_ENABLE_LOAD_BALANCER=true

# Strategy (round-robin, weighted-round-robin, least-connections, random)
DNS_LOAD_BALANCER_STRATEGY=round-robin
```

### Adding Backends via API

When you add multiple DNS records for the same name and type, they're automatically added to the load balancer:

```bash
# Example: Add 3 web servers for load balancing
curl -X POST http://localhost:8080/api/v1/zones/example.lan./records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "www.example.lan.",
    "type": "A",
    "value": "192.168.1.10",
    "ttl": 300
  }'

curl -X POST http://localhost:8080/api/v1/zones/example.lan./records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "www.example.lan.",
    "type": "A",
    "value": "192.168.1.11",
    "ttl": 300
  }'

curl -X POST http://localhost:8080/api/v1/zones/example.lan./records \
  -H "Content-Type: application/json" \
  -d '{
    "name": "www.example.lan.",
    "type": "A",
    "value": "192.168.1.12",
    "ttl": 300
  }'
```

### Metrics

- `godns_backends_total`: Total number of backends
- `godns_backends_healthy`: Number of healthy backends
- `godns_backend_requests_total`: Requests per backend

---

## Health Checks

### Overview

Continuously monitors backend server health and removes unhealthy backends from rotation.

### Supported Check Types

#### 1. TCP Health Checks

Attempts to establish TCP connection.

#### 2. HTTP Health Checks

Performs HTTP GET request, expects 2xx status code.

#### 3. HTTPS Health Checks

Performs HTTPS GET request with TLS verification.

### Configuration

```bash
# Enable health checks
DNS_ENABLE_HEALTH_CHECK=true

# Check interval in seconds (default: 30)
DNS_HEALTH_CHECK_INTERVAL=30

# Check timeout in seconds (default: 5)
DNS_HEALTH_CHECK_TIMEOUT=5
```

### How It Works

1. **Periodic Checks**: Runs health check every `INTERVAL` seconds
2. **Health Status**: Updates backend health status
3. **Load Balancer Integration**: Unhealthy backends excluded from selection
4. **Auto-Recovery**: When backend recovers, automatically re-added to pool

### Example Health Check Setup

```bash
# For HTTP backends, health checks automatically use HTTP
# The health check will GET http://<backend_ip>/
```

### Metrics

- `godns_health_checks_total`: Total number of health checks
- `godns_health_check_success_total`: Successful health checks
- `godns_health_check_failures_total`: Failed health checks
- `godns_health_check_duration_seconds`: Health check latency histogram

---

## Query Logging

### Overview

Logs all DNS queries with detailed information for troubleshooting and auditing.

### Features

- **Buffered Logging**: Batches log entries for performance
- **JSON Format**: Structured logs for easy parsing
- **Statistics**: Tracks cache hit rate, blocked queries
- **Configurable Output**: Console and/or file output

### Configuration

```bash
# Enable query logging
DNS_ENABLE_QUERY_LOG=true

# Log to console
DNS_QUERY_LOG_TO_CONSOLE=true

# Buffer size (default: 1000)
DNS_QUERY_LOG_BUFFER_SIZE=1000

# Flush interval in seconds (default: 60)
DNS_QUERY_LOG_FLUSH_INTERVAL=60
```

### Log Entry Format

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "client_ip": "192.168.1.100",
  "query_name": "www.example.lan.",
  "query_type": "A",
  "response_code": "NOERROR",
  "answer_count": 1,
  "latency_ms": 5,
  "cache_hit": false,
  "upstream": false,
  "blocked": false
}
```

### Fields Explanation

- `timestamp`: Query time (ISO 8601)
- `client_ip`: Client IP address
- `query_name`: Requested domain name
- `query_type`: DNS record type (A, AAAA, MX, etc.)
- `response_code`: DNS response code (NOERROR, NXDOMAIN, REFUSED, etc.)
- `answer_count`: Number of answer records
- `latency_ms`: Query processing time in milliseconds
- `cache_hit`: Whether response came from cache
- `upstream`: Whether query was forwarded to upstream
- `blocked`: Whether query was rate-limited

---

## Prometheus Metrics

### Overview

Exposes comprehensive metrics in Prometheus format for monitoring and alerting.

### Configuration

```bash
# Enable metrics
DNS_ENABLE_METRICS=true

# Metrics HTTP port (default: :9090)
DNS_METRICS_PORT=:9090
```

### Accessing Metrics

```bash
curl http://localhost:9090/metrics
```

### Available Metrics

#### Query Metrics

- `godns_queries_total{type,rcode}`: Total queries by type and response code
- `godns_query_duration_seconds{type,rcode}`: Query latency histogram

#### Cache Metrics

- `godns_cache_hits_total`: Cache hits
- `godns_cache_misses_total`: Cache misses
- `godns_cache_size`: Current cache size
- `godns_cache_evictions_total`: Cache evictions

#### Rate Limiting Metrics

- `godns_rate_limited_total`: Rate-limited queries
- `godns_rate_limiters_active`: Active rate limiters

#### Load Balancing Metrics

- `godns_backends_total`: Total backends
- `godns_backends_healthy`: Healthy backends
- `godns_backend_requests_total{backend,status}`: Backend request counts

#### Health Check Metrics

- `godns_health_checks_total`: Total health checks
- `godns_health_check_success_total{target,type}`: Successful checks
- `godns_health_check_failures_total{target,type}`: Failed checks
- `godns_health_check_duration_seconds{target,type}`: Check latency

#### Upstream Metrics

- `godns_upstream_queries_total`: Upstream queries
- `godns_upstream_errors_total`: Upstream errors
- `godns_upstream_duration_seconds`: Upstream query latency

### Sample Prometheus Queries

```promql
# Cache hit rate (percentage)
100 * rate(godns_cache_hits_total[5m]) /
  (rate(godns_cache_hits_total[5m]) + rate(godns_cache_misses_total[5m]))

# Queries per second by type
sum(rate(godns_queries_total[1m])) by (type)

# 95th percentile query latency
histogram_quantile(0.95, rate(godns_query_duration_seconds_bucket[5m]))

# Rate limiting effectiveness
rate(godns_rate_limited_total[5m])

# Backend health percentage
100 * godns_backends_healthy / godns_backends_total
```

---

## Configuration Reference

### Complete Environment Variables

```bash
#########################################
# DNS Caching
#########################################
DNS_ENABLE_CACHE=true
DNS_CACHE_SIZE=10000
DNS_CACHE_TTL_MINUTES=5

#########################################
# Rate Limiting
#########################################
DNS_ENABLE_RATE_LIMIT=true
DNS_RATE_LIMIT_QPS=100
DNS_RATE_LIMIT_BURST=200

#########################################
# Load Balancing
#########################################
DNS_ENABLE_LOAD_BALANCER=true
DNS_LOAD_BALANCER_STRATEGY=round-robin

#########################################
# Health Checks
#########################################
DNS_ENABLE_HEALTH_CHECK=true
DNS_HEALTH_CHECK_INTERVAL=30
DNS_HEALTH_CHECK_TIMEOUT=5

#########################################
# Query Logging
#########################################
DNS_ENABLE_QUERY_LOG=true
DNS_QUERY_LOG_TO_CONSOLE=true
DNS_QUERY_LOG_BUFFER_SIZE=1000
DNS_QUERY_LOG_FLUSH_INTERVAL=60

#########################################
# Metrics
#########################################
DNS_ENABLE_METRICS=true
DNS_METRICS_PORT=:9090
```

---

## Testing Examples

### 1. Test DNS Caching

```bash
# First query (cache miss)
time dig @localhost -p 53 www.example.lan A

# Second query (cache hit - should be faster)
time dig @localhost -p 53 www.example.lan A

# Check cache metrics
curl -s http://localhost:9090/metrics | grep godns_cache
```

### 2. Test Rate Limiting

```bash
# Send 300 rapid queries (exceeds 200 burst)
for i in {1..300}; do
  dig @localhost -p 53 www.example.lan A +short &
done
wait

# Check rate limiting metrics
curl -s http://localhost:9090/metrics | grep godns_rate_limited_total
```

### 3. Test Load Balancing

```bash
# Query multiple times and observe different IPs returned
for i in {1..10}; do
  echo "Query $i:"
  dig @localhost -p 53 www.example.lan A +short
  sleep 1
done

# Check backend distribution
curl -s http://localhost:9090/metrics | grep godns_backend_requests_total
```

### 4. Monitor Query Logs

```bash
# Enable query logging in terminal
# Observe JSON logs for each query

# Example log output:
# {"timestamp":"2024-01-15T10:30:45.123Z","client_ip":"127.0.0.1",...}
```

### 5. Test Health Checks

```bash
# Stop one backend server
# Watch health check detect failure

# Check health metrics
curl -s http://localhost:9090/metrics | grep godns_health_check_failures_total

# Restart backend and watch recovery
curl -s http://localhost:9090/metrics | grep godns_backends_healthy
```

### 6. Grafana Dashboard

Example Prometheus queries for Grafana dashboards:

```yaml
# Panel 1: Query Rate
- expr: sum(rate(godns_queries_total[1m])) by (type)
  title: "Queries per Second by Type"

# Panel 2: Cache Performance
- expr: |
    100 * rate(godns_cache_hits_total[5m]) / 
    (rate(godns_cache_hits_total[5m]) + rate(godns_cache_misses_total[5m]))
  title: "Cache Hit Rate (%)"

# Panel 3: Latency
- expr: |
    histogram_quantile(0.95, 
      rate(godns_query_duration_seconds_bucket[5m]))
  title: "95th Percentile Latency"

# Panel 4: Backend Health
- expr: godns_backends_healthy / godns_backends_total
  title: "Backend Health Ratio"
```

---

## Performance Tips

### 1. Optimize Cache Size

- Monitor `godns_cache_evictions_total`
- If evictions are high, increase `DNS_CACHE_SIZE`
- Balance memory usage vs. hit rate

### 2. Tune Rate Limits

- Set `QPS` based on expected legitimate traffic
- Use `BURST` to allow short spikes
- Monitor `godns_rate_limited_total` for false positives

### 3. Choose Load Balancer Strategy

- **Round Robin**: Best for identical backends
- **Weighted**: For heterogeneous backend capacity
- **Least Connections**: For long-lived connections
- **Random**: Simplest, good for stateless workloads

### 4. Health Check Tuning

- Lower `INTERVAL` for faster failure detection
- Higher `TIMEOUT` for slow backends
- Balance responsiveness vs. overhead

### 5. Query Log Performance

- Increase `BUFFER_SIZE` for high query rates
- Adjust `FLUSH_INTERVAL` based on disk I/O
- Consider disabling console logging in production

---

## Troubleshooting

### High Cache Miss Rate

- Check `DNS_CACHE_TTL_MINUTES` (may be too low)
- Verify cache size is adequate for query patterns
- Look for query diversity in logs

### Excessive Rate Limiting

- Review `godns_rate_limited_total` metric
- Check if legitimate clients are being blocked
- Increase `DNS_RATE_LIMIT_QPS` or `BURST`

### Backend Health Failures

- Check `godns_health_check_failures_total`
- Verify backend server connectivity
- Review `DNS_HEALTH_CHECK_TIMEOUT` setting

### High Query Latency

- Check `godns_query_duration_seconds` percentiles
- Review cache hit rate
- Monitor upstream query latency
- Check Valkey connection performance

---

## Security Considerations

1. **Rate Limiting**: Essential for public-facing deployments
2. **Metrics Endpoint**: Consider firewall rules for `:9090`
3. **Query Logs**: May contain sensitive information, secure accordingly
4. **Cache Poisoning**: TTL prevents stale data
5. **Health Checks**: Use HTTPS for sensitive backends

---

## Further Reading

- [Prometheus Metrics Best Practices](https://prometheus.io/docs/practices/naming/)
- [DNS Amplification Attack Prevention](https://www.cloudflare.com/learning/ddos/dns-amplification-ddos-attack/)
- [Load Balancing Algorithms Comparison](https://kemptechnologies.com/load-balancer/load-balancing-algorithms-techniques/)
- [Token Bucket Algorithm](https://en.wikipedia.org/wiki/Token_bucket)
