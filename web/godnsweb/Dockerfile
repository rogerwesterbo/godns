# Multi-stage build for secure, minimal web application
# Stage 1: Build the React application
FROM node:24-alpine AS builder

# Set working directory
WORKDIR /build

# Copy package files
COPY package.json package-lock.json* ./

# Install dependencies
RUN npm ci --only=production=false

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Stage 2: Build a minimal static file server in Go
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS server-builder

WORKDIR /build

# Create a simple static file server
RUN cat > server.go <<'EOF'
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	root := "/app"
	
	// Custom handler for SPA routing
	fs := http.FileServer(http.Dir(root))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(root, r.URL.Path)
		
		// Check if file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// If not a file, check if it's an API call or needs index.html
			if !strings.HasPrefix(r.URL.Path, "/api/") && !strings.Contains(r.URL.Path, ".") {
				http.ServeFile(w, r, filepath.Join(root, "index.html"))
				return
			}
		}
		
		// Add security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		fs.ServeHTTP(w, r)
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
EOF

ARG TARGETOS
ARG TARGETARCH

# Build the static server binary
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -trimpath \
    -ldflags="-s -w" \
    -o server \
    server.go

# Stage 3: Create minimal runtime image using Google Distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the static file server
COPY --from=server-builder /build/server /server

# Copy built application from builder
COPY --from=builder --chown=65532:65532 /build/dist /app

# Use nonroot user (UID 65532)
USER 65532:65532

# Expose HTTP port
EXPOSE 8080

# Set environment
ENV PORT=8080

# Set entrypoint
ENTRYPOINT ["/server"]
