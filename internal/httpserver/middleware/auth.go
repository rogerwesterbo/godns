package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// Context keys for user information
	userEmailKey contextKey = "user_email"
	userNameKey  contextKey = "user_name"
	usernameKey  contextKey = "username"
	userIDKey    contextKey = "user_id"
)

// AuthMiddleware validates JWT tokens from Keycloak
type AuthMiddleware struct {
	verifier *oidc.IDTokenVerifier
	enabled  bool
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware() (*AuthMiddleware, error) {
	// Check if authentication is enabled
	authEnabled := viper.GetBool(consts.AUTH_ENABLED)
	if !authEnabled {
		vlog.Info("Authentication is disabled")
		return &AuthMiddleware{enabled: false}, nil
	}

	// Get Keycloak configuration
	keycloakURL := viper.GetString(consts.KEYCLOAK_URL)
	realm := viper.GetString(consts.KEYCLOAK_REALM)
	clientID := viper.GetString(consts.KEYCLOAK_API_CLIENT_ID)

	if keycloakURL == "" || realm == "" || clientID == "" {
		return nil, fmt.Errorf("keycloak configuration incomplete: URL=%s, realm=%s, clientID=%s",
			keycloakURL, realm, clientID)
	}

	issuerURL := fmt.Sprintf("%s/realms/%s", keycloakURL, realm)

	vlog.Infof("Initializing OIDC authentication middleware with issuer: %s", issuerURL)

	// Create OIDC provider
	provider, err := oidc.NewProvider(context.Background(), issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	// Create token verifier
	verifier := provider.Verifier(&oidc.Config{
		ClientID: clientID,
	})

	vlog.Info("Authentication middleware initialized successfully")

	return &AuthMiddleware{
		verifier: verifier,
		enabled:  true,
	}, nil
}

// Authenticate is the middleware handler that validates JWT tokens
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If authentication is disabled, pass through
		if !am.enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			vlog.Warn("No Authorization header provided")
			http.Error(w, `{"error": "missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Check for Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			vlog.Warn("Invalid Authorization header format")
			http.Error(w, `{"error": "invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Verify the token
		token, err := am.verifier.Verify(r.Context(), tokenString)
		if err != nil {
			vlog.Warnf("Token verification failed: %v", err)
			http.Error(w, `{"error": "invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Extract claims
		var claims struct {
			Email         string   `json:"email"`
			EmailVerified bool     `json:"email_verified"`
			Name          string   `json:"name"`
			PreferredUser string   `json:"preferred_username"`
			RealmRoles    []string `json:"realm_access.roles"`
		}

		if err := token.Claims(&claims); err != nil {
			vlog.Warnf("Failed to parse token claims: %v", err)
			http.Error(w, `{"error": "invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		// Add user information to request context
		ctx := context.WithValue(r.Context(), userEmailKey, claims.Email)
		ctx = context.WithValue(ctx, userNameKey, claims.Name)
		ctx = context.WithValue(ctx, usernameKey, claims.PreferredUser)
		ctx = context.WithValue(ctx, userIDKey, token.Subject)

		// Log successful authentication
		vlog.Debugf("Authenticated user: %s (%s)", claims.PreferredUser, claims.Email)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthenticateFunc wraps an http.HandlerFunc with authentication
func (am *AuthMiddleware) AuthenticateFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		am.Authenticate(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}

// GetUserFromContext extracts user information from the request context
func GetUserFromContext(ctx context.Context) (userID, username, email string) {
	if val := ctx.Value(userIDKey); val != nil {
		userID, _ = val.(string)
	}
	if val := ctx.Value(usernameKey); val != nil {
		username, _ = val.(string)
	}
	if val := ctx.Value(userEmailKey); val != nil {
		email, _ = val.(string)
	}
	return
}
