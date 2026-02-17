package contract

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/blueprint"
)

func TestListTree_ValidRoot_ReturnsTree(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "cmd", "server"), 0755)
	_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n"), 0644)

	store := blueprint.NewStore()
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	blueprint.RegisterTools(server, root, store)

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
	res, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_tree",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool returned error: %v", res.Content)
	}
	if len(res.Content) == 0 {
		t.Fatal("expected content")
	}
	var out struct {
		Tree struct {
			Path     string `json:"path"`
			Name     string `json:"name"`
			IsDir    bool   `json:"is_dir"`
			Children []any  `json:"children"`
		} `json:"tree"`
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Tree.Name != "." {
		t.Errorf("expected root name \".\", got %q", out.Tree.Name)
	}
	if !out.Tree.IsDir {
		t.Error("expected root to be dir")
	}
	if len(out.Tree.Children) == 0 {
		t.Error("expected at least one child (cmd or go.mod)")
	}
}

func TestListTree_UnreadableRoot_StructuredError(t *testing.T) {
	store := blueprint.NewStore()
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	blueprint.RegisterTools(server, "/nonexistent/path/12345", store)

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
	res, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_tree",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected tool to return error for unreadable root")
	}
}
