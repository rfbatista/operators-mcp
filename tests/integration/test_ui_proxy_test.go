package integration

import (
	"context"
	"net/http"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/mcp"
	"operators-mcp/internal/ui"
)

func TestUIProxy_DevModeServesFromVite(t *testing.T) {
	// Start a stub HTTP server simulating Vite dev server on 5173.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<!DOCTYPE html><html><head><title>Designer</title></head><body>Dev</body></html>`))
	})
	srv := &http.Server{Addr: ":5173", Handler: mux}
	go srv.ListenAndServe()
	t.Cleanup(func() { srv.Close() })

	ctx := context.Background()
	cfg := mcp.Config{DevMode: true}
	server := mcp.NewServer(cfg)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerProxyHandler(ui.DefaultDevServerURL))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	res, err := session.ReadResource(ctx, &sdkmcp.ReadResourceParams{URI: ui.DesignerURI})
	if err != nil {
		t.Fatalf("ReadResource: %v", err)
	}
	if len(res.Contents) == 0 {
		t.Fatal("expected content")
	}
	if !strings.Contains(res.Contents[0].Text, "Dev") || !strings.Contains(res.Contents[0].Text, "<html") {
		t.Errorf("expected HTML from proxy, got: %s", res.Contents[0].Text[:min(120, len(res.Contents[0].Text))])
	}
}
