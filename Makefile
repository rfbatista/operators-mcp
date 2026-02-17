# operators-mcp — MCP server + Architecture Designer UI
# See specs/001-go-react-ui-bridge/quickstart.md for usage.

.PHONY: help build run dev-server test clean web-install web-build web-dev copy-ui deps

BINARY   := bin/server
WEB_DIR  := web
UI_STATIC := internal/ui/static

# Default target
help:
	@echo "operators-mcp — targets:"
	@echo "  make build       — build web UI, copy to embed dir, build Go binary ($(BINARY))"
	@echo "  make run         — run production server (requires 'make build' first)"
	@echo "  make dev-server  — run Go server with --dev (run 'make web-dev' in another terminal for hot-reload)"
	@echo "  make web-dev     — start Vite dev server (port 5173)"
	@echo "  make test        — run Go tests (contract + integration)"
	@echo "  make clean       — remove bin/, web/dist/, $(UI_STATIC)/"
	@echo "  make deps        — install Go deps + web npm deps"
	@echo ""
	@echo "Production: make deps && make build && make run"
	@echo "Development: make web-dev (terminal 1), make dev-server (terminal 2)"

# Install Go and web dependencies
deps:
	go mod download
	$(MAKE) web-install

# Web: install npm deps
web-install:
	cd $(WEB_DIR) && npm install

# Web: production build (outputs to web/dist/)
web-build:
	cd $(WEB_DIR) && npm run build

# Web: start Vite dev server (hot-reload on port 5173)
web-dev:
	cd $(WEB_DIR) && npm run dev

# Copy web/dist into internal/ui/static for Go embed (required before 'go build')
copy-ui:
	@mkdir -p $(UI_STATIC)
	@cp -r $(WEB_DIR)/dist/* $(UI_STATIC)/
	@echo "Copied $(WEB_DIR)/dist/ -> $(UI_STATIC)/"

# Build production binary: web build + copy + go build
build: web-build copy-ui
	go build -o $(BINARY) ./cmd/server
	@echo "Built $(BINARY) (production, UI embedded)"

# Run production server
run: build
	./$(BINARY)

# Run Go server in dev mode (proxies ui://designer to Vite; start Vite with 'make web-dev' first)
dev-server:
	go run ./cmd/server --dev

# Run all Go tests
test:
	go test ./...

# Run tests with race detector
test-race:
	go test -race ./...

# Remove build artifacts
clean:
	rm -rf bin
	rm -rf $(WEB_DIR)/dist
	rm -rf $(UI_STATIC)
	@echo "Cleaned bin/, $(WEB_DIR)/dist/, $(UI_STATIC)/"
