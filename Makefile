# Makefile for the project
# inspired by kubebuilder.io

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

# Basic colors
BLACK=\033[0;30m
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[0;37m

# Text formatting
BOLD=\033[1m
UNDERLINE=\033[4m
RESET=\033[0m

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

GOLANGCI_LINT = $(LOCALBIN)/golangci-lint
GOSEC ?= $(LOCALBIN)/gosec

# Use the Go toolchain version declared in go.mod when building tools
GO_VERSION := $(shell awk '/^go /{print $$2}' go.mod)
GO_TOOLCHAIN := go$(GO_VERSION)
GOSEC_VERSION ?= latest
GOLANGCI_LINT_VERSION ?= latest

##@ Help
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Build
.PHONY: build
build: swagger ## Build the manager binary.
	@printf "$(CYAN)Building all packages...$(RESET)\n"
	go build ./...
	@printf "$(GREEN)✓ Build complete$(RESET)\n"

.PHONY: build-dns
build-dns: ## Build the godns DNS server binary
	@printf "$(CYAN)Building godns...$(RESET)\n"
	@go build -o $(LOCALBIN)/godns ./cmd/godns
	@printf "$(GREEN)✓ godns built at $(BOLD)$(LOCALBIN)/godns$(RESET)\n"

.PHONY: build-api
build-api: ## Build the godnsapi HTTP API server binary
	@printf "$(CYAN)Building godnsapi...$(RESET)\n"
	@go build -o $(LOCALBIN)/godnsapi ./cmd/godnsapi
	@printf "$(GREEN)✓ godnsapi built at $(BOLD)$(LOCALBIN)/godnsapi$(RESET)\n"

.PHONY: build-cli
build-cli: ## Build the godnscli tool
	@printf "$(CYAN)Building godnscli...$(RESET)\n"
	@go build -o $(LOCALBIN)/godnscli ./cmd/godnscli
	@printf "$(GREEN)✓ godnscli built at $(BOLD)$(LOCALBIN)/godnscli$(RESET)\n"

.PHONY: build-all
build-all: build-dns build-api build-cli ## Build all binaries
	@printf "$(GREEN)$(BOLD)✓ All binaries built successfully!$(RESET)\n"

.PHONY: clean
clean: ## Clean build artifacts and binaries
	@printf "$(YELLOW)Cleaning build artifacts...$(RESET)\n"
	@rm -rf $(LOCALBIN)/godns $(LOCALBIN)/godnsapi $(LOCALBIN)/godnscli
	@rm -f coverage.out coverage.html bench.cpu bench.mem
	@printf "$(GREEN)✓ Clean complete$(RESET)\n"

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@printf "$(CYAN)Generating Swagger documentation...$(RESET)\n"
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g internal/httpserver/swagger.go -o internal/httpserver/swaggerdocs --parseDependency --parseInternal; \
		printf "$(GREEN)✓ Swagger docs generated in internal/httpserver/swaggerdocs/$(RESET)\n"; \
	else \
		printf "$(YELLOW)swag not found. Installing...$(RESET)\n"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g internal/httpserver/swagger.go -o internal/httpserver/swaggerdocs --parseDependency --parseInternal; \
		printf "$(GREEN)✓ Swagger docs generated in internal/httpserver/swaggerdocs/$(RESET)\n"; \
	fi

.PHONY: generate-swagger
generate-swagger: ## Generate Swagger documentation (alias for swagger)
	@$(MAKE) swagger

##@ Docker
.PHONY: docker-build
docker-build: swagger ## Build docker image (generates swagger docs first)
	@printf "$(CYAN)Building Docker image...$(RESET)\n"
	docker build -t ghcr.io/rogerwesterbo/godns:latest \
		--build-arg VERSION=$(shell git describe --tags --always --dirty) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		.
	@printf "$(GREEN)✓ Docker image built successfully$(RESET)\n"

.PHONY: docker-build-multiarch
docker-build-multiarch: swagger ## Build multi-arch docker image (requires buildx, generates swagger docs first)
	@printf "$(CYAN)Building multi-arch Docker image...$(RESET)\n"
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t ghcr.io/rogerwesterbo/godns:latest \
		--build-arg VERSION=$(shell git describe --tags --always --dirty) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		.
	@printf "$(GREEN)✓ Multi-arch Docker image built successfully$(RESET)\n"

.PHONY: docker-push
docker-push: ## Push docker image to registry
	@printf "$(CYAN)Pushing Docker image...$(RESET)\n"
	docker push ghcr.io/rogerwesterbo/godns:latest
	@printf "$(GREEN)✓ Docker image pushed successfully$(RESET)\n"

.PHONY: docker-run
docker-run: ## Run docker container locally
	docker run --rm -p 53:53/tcp -p 53:53/udp -p 14000:14000 -p 14001:14001 -p 14002:14002 ghcr.io/rogerwesterbo/godns:latest

.PHONY: release
release: swagger docker-build docker-push ## Build and push docker image (full release workflow)
	@printf "$(GREEN)$(BOLD)✓ Release complete!$(RESET)\n"
	@printf "$(CYAN)Image: ghcr.io/rogerwesterbo/godns:latest$(RESET)\n"
	@printf "$(CYAN)Version: $(shell git describe --tags --always --dirty)$(RESET)\n"

.PHONY: release-multiarch
release-multiarch: swagger docker-build-multiarch docker-push ## Build and push multi-arch docker image
	@printf "$(GREEN)$(BOLD)✓ Multi-arch release complete!$(RESET)\n"
	@printf "$(CYAN)Image: ghcr.io/rogerwesterbo/godns:latest$(RESET)\n"
	@printf "$(CYAN)Platforms: linux/amd64, linux/arm64$(RESET)\n"
	@printf "$(CYAN)Version: $(shell git describe --tags --always --dirty)$(RESET)\n"

##@ Code sanity

.PHONY: fmt
fmt: ## Run go fmt against code.
	@printf "$(CYAN)Running go fmt...$(RESET)\n"
	@go fmt ./...
	@printf "$(GREEN)✓ Code formatted$(RESET)\n"

.PHONY: vet
vet: ## Run go vet against code.
	@printf "$(CYAN)Running go vet...$(RESET)\n"
	@go vet ./...
	@printf "$(GREEN)✓ Vet complete$(RESET)\n"

.PHONY: lint
lint: golangci-lint ## Run go vet against code.
	@printf "$(CYAN)Running golangci-lint...$(RESET)\n"
	@$(GOLANGCI_LINT) run --timeout 5m ./...
	@printf "$(GREEN)✓ Lint complete$(RESET)\n"

##@ Tests
.PHONY: test
test: ## Run unit tests.
	@printf "$(CYAN)Running unit tests...$(RESET)\n"
	@go test -v ./... -coverprofile coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@printf "$(GREEN)✓ Tests complete - coverage report: $(BOLD)coverage.html$(RESET)\n"

.PHONY: bench
bench: ## Run benchmarks (override with BENCH=<regex>, PKG=<package pattern>, COUNT=<n>)
	@bench_regex=$${BENCH:-.}; \
	pkg_pattern=$${PKG:-./...}; \
	count=$${COUNT:-1}; \
	printf "$(CYAN)Running benchmarks: $(RESET)regex=$${bench_regex} packages=$${pkg_pattern} count=$${count}\n"; \
	go test -run=^$$ -bench=$${bench_regex} -benchmem -count=$${count} $${pkg_pattern}; \
	printf "$(GREEN)✓ Benchmarks complete$(RESET)\n"

.PHONY: bench-profile
bench-profile: ## Run benchmarks with CPU & memory profiles (outputs bench.cpu, bench.mem)
	@bench_regex=$${BENCH:-.}; \
	pkg_pattern=$${PKG:-./pkg/loggers/vlog}; \
	printf "$(CYAN)Profiling benchmarks: $(RESET)regex=$${bench_regex} packages=$${pkg_pattern}\n"; \
	go test -run=^$$ -bench=$${bench_regex} -cpuprofile bench.cpu -memprofile bench.mem -benchmem $${pkg_pattern}; \
	printf "$(GREEN)✓ Profiling complete: $(BOLD)bench.cpu, bench.mem$(RESET)\n"

deps: ## Download and verify dependencies
	@printf "$(CYAN)Downloading dependencies...$(RESET)\n"
	@go mod download
	@go mod verify
	@go mod tidy
	@printf "$(GREEN)✓ Dependencies updated!$(RESET)\n"

update-deps: ## Update dependencies
	@printf "$(CYAN)Updating dependencies...$(RESET)\n"
	@go get -u ./...
	@go mod tidy
	@printf "$(GREEN)✓ Dependencies updated!$(RESET)\n"

##@ Tools

.PHONY: golangci-lint
golangci-lint: $(LOCALBIN) ## Download golangci-lint locally if necessary.
	@$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: install-security-scanner
install-security-scanner: $(GOSEC) ## Install gosec security scanner locally (static analysis for security issues)
$(GOSEC): $(LOCALBIN)
	@set -e; printf "$(CYAN)Installing gosec $(GOSEC_VERSION)...$(RESET)\n"; \
	if ! GOBIN=$(LOCALBIN) go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION) 2>/dev/null; then \
		printf "$(YELLOW)Primary install failed, attempting fallback to @main...$(RESET)\n"; \
		if ! GOBIN=$(LOCALBIN) go install github.com/securego/gosec/v2/cmd/gosec@main; then \
			printf "$(RED)✗ gosec installation failed$(RESET)\n"; \
			exit 1; \
		fi; \
	fi; \
	printf "$(GREEN)✓ gosec installed at $(BOLD)$(GOSEC)$(RESET)\n"; \
	chmod +x $(GOSEC)

##@ Security
.PHONY: go-security-scan
go-security-scan: install-security-scanner ## Run gosec security scan (fails on findings)
	@printf "$(CYAN)Running gosec security scan...$(RESET)\n"
	@$(GOSEC) ./...
	@printf "$(GREEN)✓ Security scan complete$(RESET)\n"
# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
printf "$(CYAN)Downloading $${package}...$(RESET)\n" ;\
rm -f $(1) || true ;\
GOTOOLCHAIN=$(GO_TOOLCHAIN) GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
printf "$(GREEN)✓ Installed $(BOLD)$(1)-$(3)$(RESET)\n" ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef