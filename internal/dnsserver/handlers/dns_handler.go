package handlers

import (
	"context"
	"net"
	"net/netip"
	"time"

	"github.com/miekg/dns"
	"github.com/spf13/viper"

	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1cacheservice"
	"github.com/rogerwesterbo/godns/internal/services/v1dnsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1healthcheckservice"
	"github.com/rogerwesterbo/godns/internal/services/v1loadbalancerservice"
	"github.com/rogerwesterbo/godns/internal/services/v1metricsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1querylogservice"
	"github.com/rogerwesterbo/godns/internal/services/v1ratelimitservice"
	"github.com/rogerwesterbo/godns/internal/services/v1upstream"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

type DNSHandler struct {
	dnsService         *v1dnsservice.DNSService
	allowedLANsService *v1allowedlans.AllowedLANsService
	upstreamService    *v1upstream.UpstreamService
	cacheService       *v1cacheservice.DNSCache
	rateLimiter        *v1ratelimitservice.RateLimiter
	loadBalancer       *v1loadbalancerservice.LoadBalancer
	healthCheck        *v1healthcheckservice.HealthCheckService
	queryLog           *v1querylogservice.QueryLogService
	metrics            *v1metricsservice.MetricsService
}

// NewDNSHandler creates a new DNS handler with all optional services
func NewDNSHandler(
	dnsService *v1dnsservice.DNSService,
	allowedLANsService *v1allowedlans.AllowedLANsService,
	upstreamService *v1upstream.UpstreamService,
	cacheService *v1cacheservice.DNSCache,
	rateLimiter *v1ratelimitservice.RateLimiter,
	loadBalancer *v1loadbalancerservice.LoadBalancer,
	healthCheck *v1healthcheckservice.HealthCheckService,
	queryLog *v1querylogservice.QueryLogService,
	metrics *v1metricsservice.MetricsService,
) *DNSHandler {
	return &DNSHandler{
		dnsService:         dnsService,
		allowedLANsService: allowedLANsService,
		upstreamService:    upstreamService,
		cacheService:       cacheService,
		rateLimiter:        rateLimiter,
		loadBalancer:       loadBalancer,
		healthCheck:        healthCheck,
		queryLog:           queryLog,
		metrics:            metrics,
	}
}

// HandleDNS handles incoming DNS queries with caching, rate limiting, load balancing, logging, and metrics
func (h *DNSHandler) HandleDNS(w dns.ResponseWriter, r *dns.Msg) {
	startTime := time.Now()

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	// Extract peer IP
	peerAddr := w.RemoteAddr()
	var srcIP netip.Addr
	if udp, ok := peerAddr.(*net.UDPAddr); ok {
		srcIP, _ = netip.AddrFromSlice(udp.IP)
	} else if tcp, ok := peerAddr.(*net.TCPAddr); ok {
		srcIP, _ = netip.AddrFromSlice(tcp.IP)
	}

	// Use context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Track if query was blocked/rate-limited
	wasBlocked := false
	wasUpstream := false

	// Track cache hit status
	cacheHit := false

	// Extract query details for logging
	var question dns.Question
	if len(r.Question) > 0 {
		question = r.Question[0]
	}

	// Defer query logging at the end
	defer func() {
		latency := time.Since(startTime)

		// Log the query if query logging is enabled
		if h.queryLog != nil {
			h.queryLog.LogQuery(ctx, srcIP, question, m, latency, cacheHit, wasUpstream, wasBlocked)
		}

		// Record metrics if metrics service is enabled
		if h.metrics != nil {
			qtypeStr := dns.TypeToString[question.Qtype]
			rcodeStr := dns.RcodeToString[m.Rcode]
			h.metrics.RecordQuery(qtypeStr, rcodeStr, latency.Seconds())

			if cacheHit {
				h.metrics.RecordCacheHit()
			} else {
				h.metrics.RecordCacheMiss()
			}
		}
	}()

	// 1. Rate Limiting Check (first to prevent abuse)
	if h.rateLimiter != nil && srcIP.IsValid() {
		if !h.rateLimiter.Allow(ctx, srcIP) {
			vlog.Debugf("Rate limit exceeded for %s", srcIP.String())
			wasBlocked = true
			m.Rcode = dns.RcodeRefused
			if err := w.WriteMsg(m); err != nil {
				vlog.Warnf("failed to write rate limit response: %v", err)
			}
			if h.metrics != nil {
				h.metrics.RecordRateLimited()
			}
			return
		}
	}

	// Answer each question
	for _, q := range r.Question {
		name := dns.Fqdn(q.Name)
		qtype := q.Qtype

		vlog.Debugf("DNS query from %v: %s (type %d)", srcIP, name, qtype)

		// 2. Cache Lookup
		if h.cacheService != nil {
			cacheKey := name + ":" + dns.TypeToString[qtype]
			cachedMsg, found := h.cacheService.Get(ctx, cacheKey)
			if found && cachedMsg != nil {
				vlog.Debugf("Cache hit for %s (type %d)", name, qtype)
				cacheHit = true
				// Set reply from cache
				m = cachedMsg
				m.SetReply(r)
				if err := w.WriteMsg(m); err != nil {
					vlog.Warnf("failed to write cached response: %v", err)
				}
				return
			}
		}

		// 3. Check if we have this zone in our dynamic storage
		_, hasZone := h.dnsService.HasZone(ctx, name)
		vlog.Debugf("HasZone check for %s: %v", name, hasZone)

		if hasZone {
			// We have this zone - lookup record from Valkey
			records, err := h.dnsService.LookupRecord(ctx, name, qtype)
			if err != nil {
				vlog.Warnf("failed to lookup record %s: %v", name, err)
			} else if len(records) > 0 {
				vlog.Debugf("Found %d records for %s", len(records), name)

				// 4. Load Balancing - if multiple records and load balancer enabled
				if h.loadBalancer != nil && len(records) > 1 {
					// Use load balancer to select best record
					recordTypeStr := dns.TypeToString[qtype]
					selectedRecord, found := h.loadBalancer.GetBackend(ctx, name, recordTypeStr)
					if found {
						selectedValue := selectedRecord.GetRData()
						vlog.Debugf("Load balancer selected backend: %s", selectedValue)
						// Find the matching DNS record and return it
						for _, rec := range records {
							if aRec, ok := rec.(*dns.A); ok && aRec.A.String() == selectedValue {
								m.Answer = append(m.Answer, rec)
								break
							} else if aaaaRec, ok := rec.(*dns.AAAA); ok && aaaaRec.AAAA.String() == selectedValue {
								m.Answer = append(m.Answer, rec)
								break
							}
						}
					} else {
						// No healthy backend via load balancer, return all records
						m.Answer = append(m.Answer, records...)
					}
				} else {
					// No load balancing needed, return all records
					m.Answer = append(m.Answer, records...)
				}

				// 5. Cache the successful response
				if h.cacheService != nil {
					cacheKey := name + ":" + dns.TypeToString[qtype]
					h.cacheService.Set(ctx, cacheKey, m)
				}
			} else {
				vlog.Debugf("Zone exists but no records found for %s (type %d)", name, qtype)
			}
			continue
		}

		// Not in our zone: optionally forward if allowed
		var isAllowed bool
		if viper.GetBool(consts.DNS_ENABLE_ALLOWED_LANS_CHECK) {
			isAllowed = h.allowedLANsService.IsAllowed(srcIP)
			vlog.Debugf("IsAllowed check for %v: %v (check enabled)", srcIP, isAllowed)
		} else {
			isAllowed = true
			vlog.Debugf("IsAllowed check bypassed (check disabled), allowing %v", srcIP)
		}

		if isAllowed {
			vlog.Debugf("Forwarding query for %s to upstream", name)
			resp, err := h.upstreamService.Forward(ctx, r)
			if err == nil && resp != nil {
				vlog.Debugf("Upstream responded successfully for %s", name)
				wasUpstream = true

				// Cache the upstream response
				if h.cacheService != nil && resp.Rcode == dns.RcodeSuccess {
					cacheKey := name + ":" + dns.TypeToString[qtype]
					h.cacheService.Set(ctx, cacheKey, resp)
				}

				// Record upstream metrics
				if h.metrics != nil {
					h.metrics.RecordUpstreamQuery(time.Since(startTime).Seconds())
				}

				if err := w.WriteMsg(resp); err != nil {
					vlog.Warnf("failed to write upstream response: %v", err)
				}
				return
			}
			vlog.Warnf("failed to forward query to upstream: %v", err)
			if h.metrics != nil {
				h.metrics.RecordUpstreamError()
			}
		}

		// If not allowed or forward failed: NXDOMAIN
		vlog.Debugf("Setting NXDOMAIN for %s", name)
		m.Rcode = dns.RcodeNameError
	}

	vlog.Debugf("Sending final response with %d answers, rcode=%d", len(m.Answer), m.Rcode)
	if err := w.WriteMsg(m); err != nil {
		vlog.Warnf("failed to write DNS response: %v", err)
	} else {
		vlog.Debugf("Successfully wrote DNS response to %v", w.RemoteAddr())
	}
}
