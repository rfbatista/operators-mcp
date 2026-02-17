package ui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
)

const DesignerURI = "ui://designer"

// Code for structured errors (FR-007).
const (
	CodeUIAssetsMissing      = "UI_ASSETS_MISSING"
	CodeDevServerUnreachable = "DEV_SERVER_UNREACHABLE"
)

// NewDesignerResourceHandler returns an MCP resource handler for ui://designer.
// When devMode is false, it serves from the embedded Dist FS.
// When devMode is true, the caller should use a proxy handler (see proxy.go); this handler still serves from embed if proxy is not used.
func NewDesignerResourceHandler(devMode bool, embedFS fs.FS) sdkmcp.ResourceHandler {
	return func(ctx context.Context, req *sdkmcp.ReadResourceRequest) (*sdkmcp.ReadResourceResult, error) {
		if req.Params.URI != DesignerURI {
			return nil, sdkmcp.ResourceNotFoundError(req.Params.URI)
		}
		if devMode {
			slog.Warn("designer resource requested in dev mode but embed handler used", "code", CodeDevServerUnreachable)
			return nil, structuredError(CodeDevServerUnreachable, "dev mode is enabled but designer proxy is not available")
		}
		html, err := readEmbeddedIndex(embedFS)
		if err != nil {
			slog.Warn("designer UI assets missing", "code", CodeUIAssetsMissing, "error", err)
			return nil, structuredError(CodeUIAssetsMissing, fmt.Sprintf("designer UI assets not available: %v", err))
		}
		return &sdkmcp.ReadResourceResult{
			Contents: []*sdkmcp.ResourceContents{{
				URI:      DesignerURI,
				MIMEType: "text/html",
				Text:     html,
			}},
		}, nil
	}
}

func readEmbeddedIndex(embedFS fs.FS) (string, error) {
	if embedFS == nil {
		return "", errors.New("embedded FS is nil")
	}
	// Dist embeds "static" so entries are "static/index.html"
	data, err := fs.ReadFile(embedFS, "static/index.html")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func structuredError(code, message string) error {
	return &jsonrpc.Error{
		Code:    -32002, // use SDK's CodeResourceNotFound for compatibility; data carries our code
		Message: message,
		Data:    json.RawMessage([]byte(fmt.Sprintf(`{"code":%q,"message":%q}`, code, message))),
	}
}
