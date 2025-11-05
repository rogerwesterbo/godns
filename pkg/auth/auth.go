package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
)

// TokenCache represents cached OAuth2 tokens
type TokenCache struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// DeviceAuthResponse represents the response from device authorization endpoint
type DeviceAuthResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error,omitempty"`
}

// GetTokenCachePath returns the path to the token cache file
func GetTokenCachePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, ".godns")
	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return filepath.Join(cacheDir, "token.json"), nil
}

// LoadTokenCache loads cached tokens from disk
func LoadTokenCache() (*TokenCache, error) {
	cachePath, err := GetTokenCachePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(cachePath) // #nosec G304 -- cachePath is constructed from user's home directory, not user input
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cached token
		}
		return nil, fmt.Errorf("failed to read token cache: %w", err)
	}

	var cache TokenCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse token cache: %w", err)
	}

	return &cache, nil
}

// SaveTokenCache saves tokens to disk
func SaveTokenCache(cache *TokenCache) error {
	cachePath, err := GetTokenCachePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal token cache: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write token cache: %w", err)
	}

	return nil
}

// ClearTokenCache removes cached tokens
func ClearTokenCache() error {
	cachePath, err := GetTokenCachePath()
	if err != nil {
		return err
	}

	if err := os.Remove(cachePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove token cache: %w", err)
	}

	return nil
}

// IsTokenValid checks if the cached token is still valid
func (tc *TokenCache) IsTokenValid() bool {
	if tc == nil || tc.AccessToken == "" {
		return false
	}

	// Consider token valid if it expires more than 1 minute from now
	return time.Now().Add(1 * time.Minute).Before(tc.ExpiresAt)
}

// GetValidToken returns a valid access token, refreshing if needed
func GetValidToken(ctx context.Context) (string, error) {
	cache, err := LoadTokenCache()
	if err != nil {
		return "", err
	}

	// If we have a valid token, return it
	if cache != nil && cache.IsTokenValid() {
		return cache.AccessToken, nil
	}

	// If we have a refresh token, try to refresh
	if cache != nil && cache.RefreshToken != "" {
		token, err := RefreshToken(ctx, cache.RefreshToken)
		if err == nil {
			return token, nil
		}
		// If refresh failed, continue to device flow
		fmt.Println("Token refresh failed, please login again")
	}

	return "", fmt.Errorf("no valid token found, please run 'godnscli login' first")
}

// generateCodeVerifier creates a random code verifier for PKCE (43-128 characters)
func generateCodeVerifier() (string, error) {
	// Generate 32 random bytes (will be 43 chars when base64url encoded)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	// Base64url encode without padding
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// generateCodeChallenge creates a code challenge from the verifier using S256 method
func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

// DeviceCodeLogin performs OAuth2 device code flow authentication with PKCE
func DeviceCodeLogin(ctx context.Context) error {
	keycloakURL := viper.GetString(consts.KEYCLOAK_URL)
	realm := viper.GetString(consts.KEYCLOAK_REALM)
	clientID := viper.GetString(consts.KEYCLOAK_CLI_CLIENT_ID)

	if keycloakURL == "" || realm == "" || clientID == "" {
		return fmt.Errorf("keycloak configuration incomplete: set KEYCLOAK_URL, KEYCLOAK_REALM, and KEYCLOAK_CLI_CLIENT_ID")
	}

	deviceAuthURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth/device", keycloakURL, realm)
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, realm)

	// Generate PKCE code verifier and challenge
	codeVerifier, err := generateCodeVerifier()
	if err != nil {
		return fmt.Errorf("failed to generate PKCE verifier: %w", err)
	}
	codeChallenge := generateCodeChallenge(codeVerifier)

	// Step 1: Request device code with PKCE
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("code_challenge", codeChallenge)
	data.Set("code_challenge_method", "S256")

	resp, err := http.PostForm(deviceAuthURL, data) // #nosec G107 -- URL is constructed from trusted Keycloak configuration, not user input
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("device authorization failed (status %d): %s", resp.StatusCode, string(body))
	}

	var deviceResp DeviceAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&deviceResp); err != nil {
		return fmt.Errorf("failed to parse device auth response: %w", err)
	}

	// Step 2: Display instructions to user
	fmt.Println("\n🔐 GoDNS Authentication")
	fmt.Println("========================")
	fmt.Printf("\nPlease visit: %s\n", deviceResp.VerificationURI)
	fmt.Printf("And enter code: %s\n", deviceResp.UserCode)

	if deviceResp.VerificationURIComplete != "" {
		fmt.Printf("\nOr visit this URL directly:\n%s\n", deviceResp.VerificationURIComplete)
	}

	fmt.Println("\nWaiting for authentication...")

	// Step 3: Poll for token
	interval := time.Duration(deviceResp.Interval) * time.Second
	if interval == 0 {
		interval = 5 * time.Second
	}

	timeout := time.Duration(deviceResp.ExpiresIn) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("authentication timeout")
		case <-ticker.C:
			token, err := pollForToken(tokenURL, clientID, deviceResp.DeviceCode, codeVerifier)
			if err != nil {
				// authorization_pending is expected, continue polling
				if strings.Contains(err.Error(), "authorization_pending") {
					continue
				}
				// slow_down means we need to increase the polling interval
				if strings.Contains(err.Error(), "slow_down") {
					interval = interval + (2 * time.Second)
					ticker.Reset(interval)
					fmt.Printf("Adjusting polling interval to %v...\n", interval)
					continue
				}
				return err
			}

			// Successfully got token!
			fmt.Println("\n✅ Authentication successful!")
			return SaveTokenCache(token)
		}
	}
}

// pollForToken attempts to exchange device code for access token with PKCE
func pollForToken(tokenURL, clientID, deviceCode, codeVerifier string) (*TokenCache, error) {
	data := url.Values{}
	data.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	data.Set("client_id", clientID)
	data.Set("device_code", deviceCode)
	data.Set("code_verifier", codeVerifier)

	resp, err := http.PostForm(tokenURL, data) // #nosec G107 -- URL is constructed from trusted Keycloak configuration, not user input
	if err != nil {
		return nil, fmt.Errorf("failed to poll for token: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("token error: %s", tokenResp.Error)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response")
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &TokenCache{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    tokenResp.TokenType,
	}, nil
}

// RefreshToken refreshes an access token using refresh token
func RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	keycloakURL := viper.GetString(consts.KEYCLOAK_URL)
	realm := viper.GetString(consts.KEYCLOAK_REALM)
	clientID := viper.GetString(consts.KEYCLOAK_CLI_CLIENT_ID)

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, realm)

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", clientID)
	data.Set("refresh_token", refreshToken)

	resp, err := http.PostForm(tokenURL, data) // #nosec G107 -- URL is constructed from trusted Keycloak configuration, not user input
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token refresh failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse refresh response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in refresh response")
	}

	// Save refreshed token
	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	cache := &TokenCache{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    tokenResp.TokenType,
	}

	if err := SaveTokenCache(cache); err != nil {
		return "", fmt.Errorf("failed to save refreshed token: %w", err)
	}

	return cache.AccessToken, nil
}

// Logout removes cached credentials
func Logout() error {
	if err := ClearTokenCache(); err != nil {
		return err
	}
	fmt.Println("✅ Logged out successfully")
	return nil
}
