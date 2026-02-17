package integration

import (
	"context"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/mcp"
	"operators-mcp/internal/ui"
)

// TestUIEmbed_ServerServesDesignerFromEmbed verifies that when the server runs in production
// mode with embedded UI (internal/ui/static populated from web/dist), requesting ui://designer
// returns HTML. Populate static before running: cp -r web/dist/* internal/ui/static/
func TestUIEmbed_ServerServesDesignerFromEmbed(t *testing.T) {
	ctx := context.Background()
	cfg := mcp.Config{DevMode: false}
	server := mcp.NewServer(cfg)
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerResourceHandler(false, ui.Dist))

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
		t.Fatalf("ReadResource: %v (ensure internal/ui/static is populated: cp -r web/dist/* internal/ui/static/)", err)
	}
	if len(res.Contents) == 0 {
		t.Fatal("expected at least one content")
	}
	c := res.Contents[0]
	if c.MIMEType != "text/html" {
		t.Errorf("MIMEType = %q, want text/html", c.MIMEType)
	}
	if !strings.Contains(c.Text, "<html") && !strings.Contains(c.Text, "<!DOCTYPE") {
		t.Errorf("expected HTML content, got: %s", c.Text[:min(100, len(c.Text))])
	}
}
