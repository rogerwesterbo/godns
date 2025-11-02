package v1valkeyclient

import (
	"context"
	"fmt"
	"time"

	"github.com/rogerwesterbo/godns/pkg/interfaces/valkeyinterface"
	"github.com/rogerwesterbo/godns/pkg/options/valkeyoptions"
	"github.com/valkey-io/valkey-go"
)

var _ valkeyinterface.ValkeyInterface = (*V1ValkeyClient)(nil)

type V1ValkeyClient struct {
	client            valkey.Client
	timeout           time.Duration
	maxRetries        int
	initialRetryDelay time.Duration
}

// NewV1ValkeyClient creates a new Valkey client using the provided options
func NewV1ValkeyClient(opts *valkeyoptions.ValkeyOptions) (*V1ValkeyClient, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid valkey options: %w", err)
	}

	clientOpts := valkey.ClientOption{
		InitAddress: []string{opts.GetAddress()},
	}

	// Add authentication if credentials are provided
	if opts.Username != "" {
		clientOpts.Username = opts.Username
	}
	if opts.APIToken != "" {
		clientOpts.Password = opts.APIToken
	}

	client, err := valkey.NewClient(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	return &V1ValkeyClient{
		client:            client,
		timeout:           time.Duration(opts.TimeoutSec) * time.Second,
		maxRetries:        opts.MaxRetries,
		initialRetryDelay: opts.InitialRetryDelay,
	}, nil
}

// retry executes a function with exponential backoff retry logic
func (c *V1ValkeyClient) retry(ctx context.Context, operation string, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Check if context is cancelled before attempting
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("context cancelled before attempt %d: %w", attempt+1, err)
		}

		// Execute the operation
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		// If this was the last attempt, don't sleep
		if attempt == c.maxRetries {
			break
		}

		// Calculate backoff delay with exponential increase
		// delay = initialDelay * 2^attempt
		// Cap the attempt to prevent overflow (max 30 to avoid 2^31+ overflow)
		cappedAttempt := attempt
		if cappedAttempt > 30 {
			cappedAttempt = 30
		}
		delay := c.initialRetryDelay * time.Duration(1<<cappedAttempt)

		// Wait before retrying, but respect context cancellation
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled after attempt %d: %w", attempt+1, ctx.Err())
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return fmt.Errorf("operation %s failed after %d attempts: %w", operation, c.maxRetries+1, lastErr)
}

// GetData retrieves data by key
func (c *V1ValkeyClient) GetData(ctx context.Context, key string) (string, error) {
	// Create a timeout context if the parent context doesn't have a deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	var value string
	err := c.retry(ctx, "GetData", func() error {
		cmd := c.client.B().Get().Key(key).Build()
		resp := c.client.Do(ctx, cmd)

		if err := resp.Error(); err != nil {
			if valkey.IsValkeyNil(err) {
				return fmt.Errorf("key not found: %s", key)
			}
			return fmt.Errorf("failed to get data: %w", err)
		}

		v, err := resp.ToString()
		if err != nil {
			return fmt.Errorf("failed to convert response to string: %w", err)
		}
		value = v
		return nil
	})

	return value, err
}

// SetData sets data for key
func (c *V1ValkeyClient) SetData(ctx context.Context, key string, data string) error {
	// Create a timeout context if the parent context doesn't have a deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	return c.retry(ctx, "SetData", func() error {
		cmd := c.client.B().Set().Key(key).Value(data).Build()
		resp := c.client.Do(ctx, cmd)

		if err := resp.Error(); err != nil {
			return fmt.Errorf("failed to set data: %w", err)
		}
		return nil
	})
}

// DeleteData deletes data for key
func (c *V1ValkeyClient) DeleteData(ctx context.Context, key string) error {
	// Create a timeout context if the parent context doesn't have a deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	return c.retry(ctx, "DeleteData", func() error {
		cmd := c.client.B().Del().Key(key).Build()
		resp := c.client.Do(ctx, cmd)

		if err := resp.Error(); err != nil {
			return fmt.Errorf("failed to delete data: %w", err)
		}
		return nil
	})
}

// ListKeys lists all keys (WARNING: use with caution in production)
func (c *V1ValkeyClient) ListKeys(ctx context.Context) ([]string, error) {
	// Create a timeout context if the parent context doesn't have a deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	var keys []string
	err := c.retry(ctx, "ListKeys", func() error {
		// Using SCAN instead of KEYS for better performance
		keys = []string{} // Reset keys on each retry
		cursor := uint64(0)

		for {
			cmd := c.client.B().Scan().Cursor(cursor).Match("*").Count(100).Build()
			resp := c.client.Do(ctx, cmd)

			if err := resp.Error(); err != nil {
				return fmt.Errorf("failed to scan keys: %w", err)
			}

			scanResp, err := resp.AsScanEntry()
			if err != nil {
				return fmt.Errorf("failed to parse scan response: %w", err)
			}

			keys = append(keys, scanResp.Elements...)

			if scanResp.Cursor == 0 {
				break
			}
			cursor = scanResp.Cursor
		}
		return nil
	})

	return keys, err
}

// Close closes the Valkey client connection
func (c *V1ValkeyClient) Close() {
	c.client.Close()
}
