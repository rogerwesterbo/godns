package handlers

import (
	"context"
	"net"
	"net/netip"
	"time"

	"github.com/miekg/dns"

	"github.com/rogerwesterbo/godns/internal/services/v1allowedlans"
	"github.com/rogerwesterbo/godns/internal/services/v1dnsservice"
	"github.com/rogerwesterbo/godns/internal/services/v1upstream"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

type DNSHandler struct {
	dnsService         *v1dnsservice.DNSService
	allowedLANsService *v1allowedlans.AllowedLANsService
	upstreamService    *v1upstream.UpstreamService
}

// NewDNSHandler creates a new DNS handler
func NewDNSHandler(
	dnsService *v1dnsservice.DNSService,
	allowedLANsService *v1allowedlans.AllowedLANsService,
	upstreamService *v1upstream.UpstreamService,
) *DNSHandler {
	return &DNSHandler{
		dnsService:         dnsService,
		allowedLANsService: allowedLANsService,
		upstreamService:    upstreamService,
	}
}

// HandleDNS is the main DNS request handler
func (h *DNSHandler) HandleDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	// Peer IP (for recursion policy)
	peerAddr := w.RemoteAddr()
	var srcIP netip.Addr
	if udp, ok := peerAddr.(*net.UDPAddr); ok {
		srcIP, _ = netip.AddrFromSlice(udp.IP)
	} else if tcp, ok := peerAddr.(*net.TCPAddr); ok {
		srcIP, _ = netip.AddrFromSlice(tcp.IP)
	}

	// Use context with timeout for all operations
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Answer each question
	for _, q := range r.Question {
		name := dns.Fqdn(q.Name)
		qtype := q.Qtype

		vlog.Debugf("DNS query from %v: %s (type %d)", srcIP, name, qtype)

		// Check if we have this zone in our dynamic storage
		_, hasZone := h.dnsService.HasZone(ctx, name)

		if hasZone {
			// We have this zone - lookup record from Valkey
			records, err := h.dnsService.LookupRecord(ctx, name, qtype)
			if err != nil {
				vlog.Warnf("failed to lookup record %s: %v", name, err)
			} else if len(records) > 0 {
				m.Answer = append(m.Answer, records...)
			}
			// If we have the zone but no records, return NXDOMAIN for this zone
			continue
		}

		// Not in our zone: optionally forward if allowed
		if h.allowedLANsService.IsAllowed(srcIP) {
			resp, err := h.upstreamService.Forward(ctx, r)
			if err == nil && resp != nil {
				_ = w.WriteMsg(resp)
				return
			}
			vlog.Warnf("failed to forward query to upstream: %v", err)
		}
		// If not allowed or forward failed: NXDOMAIN
		m.Rcode = dns.RcodeNameError
	}

	_ = w.WriteMsg(m)
}
