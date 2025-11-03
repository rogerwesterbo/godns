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
build: ## Build the manager binary.
	go build ./...

.PHONY: build-dns
build-dns: ## Build the godns DNS server binary
	@echo "Building godns..."
	@go build -o $(LOCALBIN)/godns ./cmd/godns
	@echo "godns built at $(LOCALBIN)/godns"

.PHONY: build-api
build-api: ## Build the godnsapi HTTP API server binary
	@echo "Building godnsapi..."
	@go build -o $(LOCALBIN)/godnsapi ./cmd/godnsapi
	@echo "godnsapi built at $(LOCALBIN)/godnsapi"

.PHONY: build-cli
build-cli: ## Build the godnscli tool
	@echo "Building godnscli..."
	@go build -o $(LOCALBIN)/godnscli ./cmd/godnscli
	@echo "godnscli built at $(LOCALBIN)/godnscli"

.PHONY: build-all
build-all: build-dns build-api build-cli ## Build all binaries

.PHONY: swagger
swagger: ## Generate Swagger documentation
	@echo "$(CYAN)Generating Swagger documentation...$(RESET)"
	@if command -v swag >/dev/null 2>&1; then \
		swag init -g cmd/godnsapi/main.go -o docs --parseDependency --parseInternal; \
		echo "$(GREEN)✓ Swagger docs generated in docs/$(RESET)"; \
	else \
		echo "$(YELLOW)swag not found. Installing...$(RESET)"; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init -g cmd/godnsapi/main.go -o docs --parseDependency --parseInternal; \
		echo "$(GREEN)✓ Swagger docs generated in docs/$(RESET)"; \
	fi

.PHONY: generate-swagger
generate-swagger: ## Generate Swagger documentation (alias for swagger)
	@$(MAKE) swagger

##@ Docker
.PHONY: docker-build
docker-build: swagger ## Build docker image (generates swagger docs first)
	@echo "$(CYAN)Building Docker image...$(RESET)"
	docker build -t ghcr.io/rogerwesterbo/godns:latest \
		--build-arg VERSION=$(shell git describe --tags --always --dirty) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		.
	@echo "$(GREEN)✓ Docker image built successfully$(RESET)"

.PHONY: docker-build-multiarch
docker-build-multiarch: swagger ## Build multi-arch docker image (requires buildx, generates swagger docs first)
	@echo "$(CYAN)Building multi-arch Docker image...$(RESET)"
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t ghcr.io/rogerwesterbo/godns:latest \
		--build-arg VERSION=$(shell git describe --tags --always --dirty) \
		--build-arg BUILD_TIME=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ") \
		--build-arg GIT_COMMIT=$(shell git rev-parse HEAD) \
		.
	@echo "$(GREEN)✓ Multi-arch Docker image built successfully$(RESET)"

.PHONY: docker-push
docker-push: ## Push docker image to registry
	@echo "$(CYAN)Pushing Docker image...$(RESET)"
	docker push ghcr.io/rogerwesterbo/godns:latest
	@echo "$(GREEN)✓ Docker image pushed successfully$(RESET)"

.PHONY: docker-run
docker-run: ## Run docker container locally
	docker run --rm -p 53:53/tcp -p 53:53/udp -p 8080:8080 -p 8082:8082 ghcr.io/rogerwesterbo/godns:latest

.PHONY: release
release: swagger docker-build docker-push ## Build and push docker image (full release workflow)
	@echo "$(GREEN)$(BOLD)✓ Release complete!$(RESET)"
	@echo "$(CYAN)Image: ghcr.io/rogerwesterbo/godns:latest$(RESET)"
	@echo "$(CYAN)Version: $(shell git describe --tags --always --dirty)$(RESET)"

.PHONY: release-multiarch
release-multiarch: swagger docker-build-multiarch docker-push ## Build and push multi-arch docker image
	@echo "$(GREEN)$(BOLD)✓ Multi-arch release complete!$(RESET)"
	@echo "$(CYAN)Image: ghcr.io/rogerwesterbo/godns:latest$(RESET)"
	@echo "$(CYAN)Platforms: linux/amd64, linux/arm64$(RESET)"
	@echo "$(CYAN)Version: $(shell git describe --tags --always --dirty)$(RESET)"

##@ Code sanity

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint: golangci-lint ## Run go vet against code.
	$(GOLANGCI_LINT) run --timeout 5m ./...

##@ Tests
.PHONY: test
test: ## Run unit tests.
	go test -v ./... -coverprofile coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: bench
bench: ## Run benchmarks (override with BENCH=<regex>, PKG=<package pattern>, COUNT=<n>)
	@bench_regex=$${BENCH:-.}; \
	pkg_pattern=$${PKG:-./...}; \
	count=$${COUNT:-1}; \
	echo "Running benchmarks: regex=$${bench_regex} packages=$${pkg_pattern} count=$${count}"; \
	go test -run=^$$ -bench=$${bench_regex} -benchmem -count=$${count} $${pkg_pattern}

.PHONY: bench-profile
bench-profile: ## Run benchmarks with CPU & memory profiles (outputs bench.cpu, bench.mem)
	@bench_regex=$${BENCH:-.}; \
	pkg_pattern=$${PKG:-./pkg/loggers/vlog}; \
	echo "Profiling benchmarks: regex=$${bench_regex} packages=$${pkg_pattern}"; \
	go test -run=^$$ -bench=$${bench_regex} -cpuprofile bench.cpu -memprofile bench.mem -benchmem $${pkg_pattern}

deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod verify
	@go mod tidy
	@echo "Dependencies updated!"

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "Dependencies updated!"

##@ Tools

.PHONY: golangci-lint
golangci-lint: $(LOCALBIN) ## Download golangci-lint locally if necessary.
	$(call go-install-tool,$(GOLANGCI_LINT),github.com/golangci/golangci-lint/v2/cmd/golangci-lint,$(GOLANGCI_LINT_VERSION))

.PHONY: install-security-scanner
install-security-scanner: $(GOSEC) ## Install gosec security scanner locally (static analysis for security issues)
$(GOSEC): $(LOCALBIN)
	@set -e; echo "Attempting to install gosec $(GOSEC_VERSION)"; \
	if ! GOBIN=$(LOCALBIN) go install github.com/securego/gosec/v2/cmd/gosec@$(GOSEC_VERSION) 2>/dev/null; then \
		echo "Primary install failed, attempting install from @main (compatibility fallback)"; \
		if ! GOBIN=$(LOCALBIN) go install github.com/securego/gosec/v2/cmd/gosec@main; then \
			echo "gosec installation failed for versions $(GOSEC_VERSION) and @main"; \
			exit 1; \
		fi; \
	fi; \
	echo "gosec installed at $(GOSEC)"; \
	chmod +x $(GOSEC)

##@ Security
.PHONY: go-security-scan
go-security-scan: install-security-scanner ## Run gosec security scan (fails on findings)
	$(GOSEC) ./...
# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@[ -f "$(1)-$(3)" ] || { \
set -e; \
package=$(2)@$(3) ;\
echo "Downloading $${package}" ;\
rm -f $(1) || true ;\
GOTOOLCHAIN=$(GO_TOOLCHAIN) GOBIN=$(LOCALBIN) go install $${package} ;\
mv $(1) $(1)-$(3) ;\
} ;\
ln -sf $(1)-$(3) $(1)
endef