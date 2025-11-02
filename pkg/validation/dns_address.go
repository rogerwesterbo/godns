package validation

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// ValidateDNSAddress validates a DNS server address
// Valid formats:
//   - :53 (listen on all interfaces with port)
//   - 0.0.0.0:53 (explicit all interfaces)
//   - 192.168.1.1:53 (specific IP with port)
//   - [::]:53 (IPv6 all interfaces)
//   - [2001:db8::1]:53 (specific IPv6 with port)
//
// Invalid formats:
//   - 53 (port only without colon)
//   - 192.168.1.1 (IP without port)
func ValidateDNSAddress(address string) error {
	if address == "" {
		return fmt.Errorf("DNS address cannot be empty")
	}

	// Check if address starts with : (e.g., :53)
	if strings.HasPrefix(address, ":") {
		// Extract port and validate
		portStr := address[1:]
		if portStr == "" {
			return fmt.Errorf("DNS address '%s' is missing port number", address)
		}
		return validatePort(portStr)
	}

	// Must contain a colon for host:port format
	if !strings.Contains(address, ":") {
		return fmt.Errorf("DNS address '%s' must include a port (e.g., :53 or 0.0.0.0:53)", address)
	}

	// Use net.SplitHostPort to validate the format
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid DNS address '%s': %w", address, err)
	}

	// Validate port
	if err := validatePort(portStr); err != nil {
		return err
	}

	// Validate host/IP if provided
	if host != "" {
		// Try to parse as IP address
		ip := net.ParseIP(host)
		if ip == nil {
			// Not a valid IP, could be hostname (we'll allow it)
			// But check it's not obviously invalid
			if strings.TrimSpace(host) == "" {
				return fmt.Errorf("DNS address '%s' has empty host", address)
			}
		}
	}

	return nil
}

func validatePort(portStr string) error {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port '%s': must be a number", portStr)
	}

	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid port %d: must be between 1 and 65535", port)
	}

	return nil
}
