# operators-mcp

MCP server for the AI-Architecture Orchestrator (Stateful MCP Server). Serves the Architecture Designer UI at resource `ui://designer`.

**Quick start (Makefile)**:
- `make help` — list all targets
- `make deps && make build && make run` — production
- `make web-dev` (terminal 1) + `make dev-server` (terminal 2) — development with hot-reload

**Details**: [specs/001-go-react-ui-bridge/quickstart.md](specs/001-go-react-ui-bridge/quickstart.md)

**Pattern Management UI** (002): The designer includes a source tree with zone highlighting, regex playground, zone metadata form (Purpose, Constraints, Assigned Agent), and “Assign to Zone” from the tree. See [specs/002-pattern-management-ui/quickstart.md](specs/002-pattern-management-ui/quickstart.md). The UI calls backend tools via the MCP host bridge (`window.__callTool__` when provided); without it, mock data is used so the layout is still usable.
