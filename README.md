# operators-mcp

MCP server for the AI-Architecture Orchestrator (Stateful MCP Server). Serves the Architecture Designer UI at resource `ui://designer`.

**Architecture**: The codebase follows [hexagonal (ports & adapters)](docs/ARCHITECTURE.md): `internal/domain`, `internal/application/blueprint`, and `internal/adapter` (in/out).

**Quick start (Makefile)**:
- `make help` — list all targets
- `make deps && make build && make run` — production
- `make web-dev` (terminal 1) + `make dev-server` (terminal 2) — development with hot-reload
- `make run` or `./bin/server` — MCP server on :8081, HTTP server on :8080 (UI: /, API: /api)
- `make air` — run with [Air](https://github.com/air-verse/air): auto-reload on Go or React changes, UI at http://localhost:8080 (install: `go install github.com/air-verse/air@latest`)

**Details**: [specs/001-go-react-ui-bridge/quickstart.md](specs/001-go-react-ui-bridge/quickstart.md)

**Pattern Management UI** (002): The designer includes a source tree with zone highlighting, regex playground, zone metadata form (Purpose, Constraints, Assigned Agent), and “Assign to Zone” from the tree. See [specs/002-pattern-management-ui/quickstart.md](specs/002-pattern-management-ui/quickstart.md). The UI calls backend tools via the MCP host bridge (`window.__callTool__` when provided); without it, mock data is used so the layout is still usable.

---

## Two servers (same process)

The binary starts **two separate servers**:

| Server   | Default address   | Purpose                    |
|----------|-------------------|----------------------------|
| **MCP**  | http://localhost:8081 | Tools, resources (IDE connects here) |
| **HTTP** | http://localhost:8080 | UI at `/`, API at `/api`   |

```bash
make build && make run
```

Startup logs show both ports, e.g.:
- `MCP server listening on :8081 (use this URL in your IDE)`
- `HTTP server listening on :8080 — UI: /, API: /api`

---

## Setting up the MCP server in your IDE

MCP runs on its **own port**. Start the server, then point your IDE at the MCP URL.

1. **Start the server**: `make run` (or `./bin/server`). MCP is on port **8081** by default.
2. In your IDE MCP config, set **url** to `http://localhost:8081` (or the host/port you set with `-mcp.addr`).

**Cursor** (`~/.cursor/mcp.json` or project `.cursor/mcp.json`):
```json
{
  "mcpServers": {
    "operators-mcp": {
      "url": "http://localhost:8081"
    }
  }
}
```

**Other IDEs**: use the MCP “URL” or “HTTP/SSE” option and set it to `http://localhost:8081`. Restart or reload MCP after changing config.

---

### Optional server flags

- `-mcp.addr <addr>` — MCP server listen address (default: `:8081`).
- `-http.addr <addr>` — HTTP server listen address (default: `:8080`).
- `-db <path>` — SQLite DB path (default: `data.db`). Use `:memory:` for in-memory.
- `-dev` — Proxy `ui://designer` to Vite; run `make web-dev` separately.
