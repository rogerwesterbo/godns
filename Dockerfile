# Multi-stage build for secure, minimal Go binary
# Stage 1: Build the Go binary
FROM golang:1.25.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary with security hardening flags
# - CGO_ENABLED=0: Static binary, no C dependencies
# - -trimpath: Remove file system paths from binary
# - -ldflags: Strip debug info and set version info
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -trimpath \
    -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT}" \
    -o godns \
    ./cmd/godns

# Stage 2: Create minimal runtime image using Google Distroless
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/godns /godns

# Use nonroot user (UID 65532)
USER nonroot:nonroot

# Expose DNS ports
EXPOSE 53/tcp 53/udp

# Health check endpoint (if your app supports it)
EXPOSE 8080/tcp

# Set entrypoint
ENTRYPOINT ["/godns"]
CMD ["--help"]
