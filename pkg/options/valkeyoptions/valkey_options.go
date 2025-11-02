package valkeyoptions

import (
	"fmt"
	"time"
)

type ValkeyOptions struct {
	Host              string
	Port              string
	Username          string
	APIToken          string
	TimeoutSec        int
	MaxRetries        int
	InitialRetryDelay time.Duration
}

func DefaultValkeyOptions(host string, port string, username string, token string) *ValkeyOptions {
	return &ValkeyOptions{
		Host:              host,
		Port:              port,
		Username:          username,
		APIToken:          token,
		TimeoutSec:        30,
		MaxRetries:        3,
		InitialRetryDelay: 100 * time.Millisecond,
	}
}

func (o *ValkeyOptions) ApplyOptions(opts ...func(*ValkeyOptions)) {
	for _, opt := range opts {
		opt(o)
	}
}

func WithValkeyHost(host string) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.Host = host
	}
}

func WithValkeyPort(port string) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.Port = port
	}
}

func WithValkeyUsername(username string) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.Username = username
	}
}

func WithValkeyAPIToken(token string) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.APIToken = token
	}
}

func WithValkeyTimeoutSec(timeout int) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.TimeoutSec = timeout
	}
}

func WithValkeyMaxRetries(maxRetries int) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.MaxRetries = maxRetries
	}
}

func WithValkeyInitialRetryDelay(delay time.Duration) func(*ValkeyOptions) {
	return func(o *ValkeyOptions) {
		o.InitialRetryDelay = delay
	}
}

// GetAddress returns the combined host:port address
func (o *ValkeyOptions) GetAddress() string {
	return fmt.Sprintf("%s:%s", o.Host, o.Port)
}

func (o *ValkeyOptions) Validate() error {
	if o.Host == "" {
		return fmt.Errorf("valkey host cannot be empty")
	}
	if o.Port == "" {
		return fmt.Errorf("valkey port cannot be empty")
	}
	// Username is optional (for backward compatibility with password-only auth)
	// APIToken can be empty for local development without auth
	if o.TimeoutSec <= 0 {
		return fmt.Errorf("valkey timeout must be positive")
	}
	if o.MaxRetries < 0 {
		return fmt.Errorf("valkey max retries cannot be negative")
	}
	if o.InitialRetryDelay < 0 {
		return fmt.Errorf("valkey initial retry delay cannot be negative")
	}
	return nil
}
