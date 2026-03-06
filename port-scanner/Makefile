# Port Scanner Makefile
# =====================

# Verzija (moze se override-ati sa: make VERSION=1.1.0)
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parametri
BINARY_NAME := port-scanner
GO := go
GOFLAGS := -trimpath
LDFLAGS := -s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GitCommit=$(GIT_COMMIT)'

# Instalacijski direktorij
PREFIX ?= /usr/local
BINDIR := $(PREFIX)/bin

# Boje za output
GREEN := \033[0;32m
YELLOW := \033[1;33m
RED := \033[0;31m
NC := \033[0m

.PHONY: all build install uninstall clean test help version update

# Default target
all: build

# Kompajliranje
build:
	@echo "$(GREEN)[*] Building $(BINARY_NAME) v$(VERSION)...$(NC)"
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "$(GREEN)[+] Build complete: ./$(BINARY_NAME)$(NC)"

# Kompajliranje za produkciju (optimizirano)
build-prod:
	@echo "$(GREEN)[*] Building $(BINARY_NAME) v$(VERSION) (production)...$(NC)"
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .
	@echo "$(GREEN)[+] Production build complete: ./$(BINARY_NAME)$(NC)"

# Instalacija u sistem
install: build
	@echo "$(GREEN)[*] Installing to $(BINDIR)/$(BINARY_NAME)...$(NC)"
	@sudo install -d $(BINDIR)
	@sudo install -m 755 $(BINARY_NAME) $(BINDIR)/$(BINARY_NAME)
	@echo "$(GREEN)[+] Installed! Run with: $(BINARY_NAME)$(NC)"

# Deinstalacija
uninstall:
	@echo "$(YELLOW)[*] Removing $(BINDIR)/$(BINARY_NAME)...$(NC)"
	@sudo rm -f $(BINDIR)/$(BINARY_NAME)
	@echo "$(GREEN)[+] Uninstalled$(NC)"

# Ciscenje build artefakata
clean:
	@echo "$(YELLOW)[*] Cleaning...$(NC)"
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_NAME)-*
	@$(GO) clean
	@echo "$(GREEN)[+] Clean complete$(NC)"

# Pokretanje testova
test:
	@echo "$(GREEN)[*] Running tests...$(NC)"
	$(GO) test -v ./...

# Update - pull, rebuild, reinstall
update:
	@echo "$(GREEN)[*] Updating $(BINARY_NAME)...$(NC)"
	@git pull origin main
	@$(MAKE) install
	@echo "$(GREEN)[+] Update complete! New version:$(NC)"
	@$(BINARY_NAME) version

# Prikazi verziju
version:
	@echo "$(BINARY_NAME) v$(VERSION)"
	@echo "Build: $(BUILD_TIME)"
	@echo "Commit: $(GIT_COMMIT)"

# Cross-compile za razne platforme
build-all: clean
	@echo "$(GREEN)[*] Building for multiple platforms...$(NC)"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-arm64 .
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 $(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows-amd64.exe .
	@echo "$(GREEN)[+] All builds complete$(NC)"

# Pomoc
help:
	@echo ""
	@echo "$(GREEN)Port Scanner - Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Usage:$(NC)"
	@echo "  make              - Build binary"
	@echo "  make build        - Build binary"
	@echo "  make build-prod   - Build optimized production binary"
	@echo "  make install      - Build and install to $(BINDIR)"
	@echo "  make uninstall    - Remove from $(BINDIR)"
	@echo "  make update       - Git pull + rebuild + reinstall"
	@echo "  make clean        - Remove build artifacts"
	@echo "  make test         - Run tests"
	@echo "  make version      - Show version info"
	@echo "  make build-all    - Cross-compile for all platforms"
	@echo "  make help         - Show this help"
	@echo ""
	@echo "$(YELLOW)Options:$(NC)"
	@echo "  VERSION=x.y.z     - Override version number"
	@echo "  PREFIX=/path      - Override install prefix (default: /usr/local)"
	@echo ""
	@echo "$(YELLOW)Examples:$(NC)"
	@echo "  make VERSION=1.1.0 build"
	@echo "  make PREFIX=/usr install"
	@echo ""
