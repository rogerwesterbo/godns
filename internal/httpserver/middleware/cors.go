package middleware

import (
	"net/http"
	"strings"

	"github.com/rogerwesterbo/godns/pkg/consts"
	"github.com/spf13/viper"
	"github.com/vitistack/common/pkg/loggers/vlog"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS) headers
type CORSMiddleware struct {
	allowedOrigins map[string]bool
}

// NewCORSMiddleware creates a new CORS middleware with allowed origins from config
func NewCORSMiddleware() *CORSMiddleware {
	// Get allowed origins from config (comma-separated list)
	allowedOriginsStr := viper.GetString(consts.HTTP_API_CORS_ALLOWED_ORIGINS)

	allowedOrigins := make(map[string]bool)
	if allowedOriginsStr != "" {
		origins := strings.Split(allowedOriginsStr, ",")
		for _, origin := range origins {
			trimmedOrigin := strings.TrimSpace(origin)
			if trimmedOrigin != "" {
				allowedOrigins[trimmedOrigin] = true
				vlog.Infof("CORS: Allowing origin: %s", trimmedOrigin)
			}
		}
	}

	if len(allowedOrigins) == 0 {
		vlog.Warn("No CORS allowed origins configured. CORS will be restrictive.")
	}

	return &CORSMiddleware{
		allowedOrigins: allowedOrigins,
	}
}

// Handler wraps an http.Handler with CORS support
func (c *CORSMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if origin != "" && c.allowedOrigins[origin] {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours
		}

		// Handle preflight OPTIONS request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// HandlerFunc wraps an http.HandlerFunc with CORS support
func (c *CORSMiddleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Handler(http.HandlerFunc(next)).ServeHTTP(w, r)
	}
}
