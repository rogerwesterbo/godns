package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/rogerwesterbo/godns/pkg/auth"
)

// AuthenticatedClient is an HTTP client that automatically adds authentication headers
type AuthenticatedClient struct {
	httpClient *http.Client
}

// NewAuthenticatedClient creates a new authenticated HTTP client
func NewAuthenticatedClient() *AuthenticatedClient {
	return &AuthenticatedClient{
		httpClient: &http.Client{},
	}
}

// Get performs an authenticated GET request
func (c *AuthenticatedClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.Do(req)
}

// Post performs an authenticated POST request
func (c *AuthenticatedClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// Put performs an authenticated PUT request
func (c *AuthenticatedClient) Put(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	return c.Do(req)
}

// Delete performs an authenticated DELETE request
func (c *AuthenticatedClient) Delete(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return c.Do(req)
}

// Do executes an HTTP request with authentication
func (c *AuthenticatedClient) Do(req *http.Request) (*http.Response, error) {
	// Get valid access token
	token, err := auth.GetValidToken(req.Context())
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w\n\nPlease run 'godnscli login' to authenticate", err)
	}

	// Add authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Execute request
	return c.httpClient.Do(req)
}
