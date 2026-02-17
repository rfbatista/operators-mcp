package mcp

import (
	"context"
	"log/slog"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// Config holds options for the MCP server.
type Config struct {
	DevMode bool
}

// NewServer returns an MCP server instance with stdio transport.
// Resource handlers (e.g. ui://designer) are registered by the caller before Run.
func NewServer(cfg Config) *sdkmcp.Server {
	opts := &sdkmcp.ServerOptions{}
	if slog.Default().Enabled(context.Background(), slog.LevelDebug) {
		opts.Logger = slog.Default()
	}
	s := sdkmcp.NewServer(&sdkmcp.Implementation{
		Name:    "operators-mcp",
		Version: "0.0.1",
	}, opts)
	_ = cfg // DevMode used by cmd/server when registering ui resource
	return s
}
