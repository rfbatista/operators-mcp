package ui

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// DefaultDevServerURL is the default Vite dev server URL for --dev mode.
const DefaultDevServerURL = "http://localhost:5173"

// NewDesignerProxyHandler returns an MCP resource handler that proxies ui://designer
// to the given dev server URL. On connection failure, returns a structured error
// (DEV_SERVER_UNREACHABLE + message).
func NewDesignerProxyHandler(devServerURL string) sdkmcp.ResourceHandler {
	return func(ctx context.Context, req *sdkmcp.ReadResourceRequest) (*sdkmcp.ReadResourceResult, error) {
		if req.Params.URI != DesignerURI {
			return nil, sdkmcp.ResourceNotFoundError(req.Params.URI)
		}
		html, mime, err := fetchDevServer(ctx, devServerURL)
		if err != nil {
			slog.Warn("dev server unreachable", "code", CodeDevServerUnreachable, "url", devServerURL, "error", err)
			return nil, structuredError(CodeDevServerUnreachable, fmt.Sprintf("dev server unreachable: %v", err))
		}
		return &sdkmcp.ReadResourceResult{
			Contents: []*sdkmcp.ResourceContents{{
				URI:      DesignerURI,
				MIMEType: mime,
				Text:     html,
			}},
		}, nil
	}
}

func fetchDevServer(ctx context.Context, baseURL string) (html string, mime string, err error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/", nil)
	if err != nil {
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	mime = resp.Header.Get("Content-Type")
	if mime == "" {
		mime = "text/html"
	}
	return string(body), mime, nil
}
