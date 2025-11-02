package dnsserver

import (
	"os"
	"os/signal"

	"github.com/miekg/dns"
	"github.com/rogerwesterbo/godns/internal/dnsserver/handlers"
	"github.com/rogerwesterbo/godns/internal/healthserver"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// Server represents the DNS server with both UDP and TCP listeners
type Server struct {
	udpServer    *dns.Server
	tcpServer    *dns.Server
	healthServer *healthserver.Server
	dnsHandler   *handlers.DNSHandler
}

// New creates a new DNS server instance
func New(addr, livenessProbePort, readinessProbePort string, dnsHandler *handlers.DNSHandler) *Server {
	return &Server{
		udpServer:    &dns.Server{Addr: addr, Net: "udp"},
		tcpServer:    &dns.Server{Addr: addr, Net: "tcp"},
		healthServer: healthserver.New(livenessProbePort, readinessProbePort),
		dnsHandler:   dnsHandler,
	}
}

// Start begins listening on both UDP and TCP
func (s *Server) Start() error {
	// Start health check servers
	if err := s.healthServer.Start(); err != nil {
		return err
	}

	// Register handler
	dns.HandleFunc(".", s.dnsHandler.HandleDNS)

	// Start servers
	errCh := make(chan error, 2)
	go func() { errCh <- s.udpServer.ListenAndServe() }()
	go func() { errCh <- s.tcpServer.ListenAndServe() }()
	vlog.Infof("DNS server listening on %s (udp/tcp)", s.udpServer.Addr)

	// Mark service as ready once DNS servers are started
	s.healthServer.SetReady(true)

	// Graceful shutdown on Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	select {
	case <-c:
		vlog.Info("shutting down...")
		s.healthServer.SetReady(false)
		_ = s.udpServer.Shutdown()
		_ = s.tcpServer.Shutdown()
		_ = s.healthServer.Shutdown()
		return nil
	case err := <-errCh:
		s.healthServer.SetReady(false)
		_ = s.healthServer.Shutdown()
		return err
	}
}

// Shutdown gracefully shuts down both UDP and TCP servers
func (s *Server) Shutdown() error {
	s.healthServer.SetReady(false)
	if err := s.udpServer.Shutdown(); err != nil {
		return err
	}
	if err := s.tcpServer.Shutdown(); err != nil {
		return err
	}
	return s.healthServer.Shutdown()
}
