package integration

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/adapter/in/mcp"
	"operators-mcp/internal/adapter/in/ui"
	"operators-mcp/internal/adapter/out/filesystem"
	"operators-mcp/internal/adapter/out/persistence/memory"
	"operators-mcp/internal/application/blueprint"
)

// TestBlueprintTools_WithDesignerResource runs server with both ui://designer
// and blueprint tools; invokes list_tree and list_zones to verify integration.
func TestBlueprintTools_WithDesignerResource(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "cmd"), 0755)
	_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n"), 0644)

	projectStore := memory.NewProjectStore()
	zoneStore := memory.NewStore()
	pathMatcher := filesystem.NewMatcher()
	treeLister := filesystem.NewLister()
	svc := blueprint.NewService(projectStore, zoneStore, pathMatcher, treeLister, root)
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	mcp.RegisterTools(server, svc)
	// Also register designer resource so server has both
	server.AddResource(&sdkmcp.Resource{URI: ui.DesignerURI, Name: "Designer", MIMEType: "text/html"},
		ui.NewDesignerResourceHandler(false, nil))

	t1, t2 := sdkmcp.NewInMemoryTransports()
	if _, err := server.Connect(context.Background(), t1, nil); err != nil {
		t.Fatalf("server connect: %v", err)
	}
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "client", Version: "0.0.1"}, nil)
	session, err := client.Connect(context.Background(), t2, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer session.Close()

	ctx := context.Background()

	// Call list_tree
	res, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_tree",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("list_tree: %v", err)
	}
	if res.IsError {
		t.Fatalf("list_tree error: %v", res.Content)
	}
	if len(res.Content) == 0 {
		t.Fatal("expected list_tree content")
	}
	text := ""
	if tc, ok := res.Content[0].(*sdkmcp.TextContent); ok {
		text = tc.Text
	}
	var treeOut struct {
		Tree struct {
			Name string `json:"name"`
		} `json:"tree"`
	}
	if err := json.Unmarshal([]byte(text), &treeOut); err != nil {
		t.Fatalf("unmarshal list_tree: %v", err)
	}
	if treeOut.Tree.Name != "." {
		t.Errorf("list_tree root name: got %q", treeOut.Tree.Name)
	}

	// Create a project so we can list zones
	p, err := svc.CreateProject("testproj", root)
	if err != nil {
		t.Fatalf("CreateProject: %v", err)
	}
	// Call list_zones
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_zones",
		Arguments: map[string]any{"project_id": p.ID},
	})
	if err != nil {
		t.Fatalf("list_zones: %v", err)
	}
	if res.IsError {
		t.Fatalf("list_zones error: %v", res.Content)
	}
}
