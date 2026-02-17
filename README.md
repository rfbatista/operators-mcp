# operators-mcp

MCP server for the AI-Architecture Orchestrator (Stateful MCP Server). Serves the Architecture Designer UI at resource `ui://designer`.

**Architecture**: The codebase follows [hexagonal (ports & adapters)](docs/ARCHITECTURE.md): `internal/domain`, `internal/application/blueprint`, and `internal/adapter` (in/out).

**Quick start (Makefile)**:
- `make help` — list all targets
- `make deps && make build && make run` — production
- `make web-dev` (terminal 1) + `make dev-server` (terminal 2) — development with hot-reload
- `make build && ./bin/server -http` — serve the React UI at http://localhost:8080 (use `-http.addr :3000` for another port)
- `make air` — run with [Air](https://github.com/air-verse/air): auto-reload on Go or React changes, UI at http://localhost:8080 (install: `go install github.com/air-verse/air@latest`)

**Details**: [specs/001-go-react-ui-bridge/quickstart.md](specs/001-go-react-ui-bridge/quickstart.md)

**Pattern Management UI** (002): The designer includes a source tree with zone highlighting, regex playground, zone metadata form (Purpose, Constraints, Assigned Agent), and “Assign to Zone” from the tree. See [specs/002-pattern-management-ui/quickstart.md](specs/002-pattern-management-ui/quickstart.md). The UI calls backend tools via the MCP host bridge (`window.__callTool__` when provided); without it, mock data is used so the layout is still usable.
