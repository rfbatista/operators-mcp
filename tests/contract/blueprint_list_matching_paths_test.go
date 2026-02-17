package contract

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/adapter/in/mcp"
	"operators-mcp/internal/adapter/out/filesystem"
	"operators-mcp/internal/adapter/out/persistence/memory"
	"operators-mcp/internal/application/blueprint"
)

func TestListMatchingPaths_ValidPattern_ReturnsPaths(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "cmd", "server"), 0755)
	_ = os.MkdirAll(filepath.Join(root, "internal", "mcp"), 0755)
	_ = os.WriteFile(filepath.Join(root, "go.mod"), []byte("module test\n"), 0644)

	projectStore := memory.NewProjectStore()
	zoneStore := memory.NewStore()
	pathMatcher := filesystem.NewMatcher()
	treeLister := filesystem.NewLister()
	svc := blueprint.NewService(projectStore, zoneStore, pathMatcher, treeLister, root)
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	mcp.RegisterTools(server, svc)

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
		Name:      "list_matching_paths",
		Arguments: map[string]any{"pattern": "cmd"},
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
	text := contentText(res.Content[0])
	var out struct {
		Paths []string `json:"paths"`
	}
	if err := json.Unmarshal([]byte(text), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	found := false
	for _, p := range out.Paths {
		if p == "cmd" || p == "cmd/server" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected paths to contain cmd or cmd/server, got %v", out.Paths)
	}
}

func TestListMatchingPaths_InvalidPattern_StructuredError(t *testing.T) {
	root := t.TempDir()
	projectStore := memory.NewProjectStore()
	zoneStore := memory.NewStore()
	pathMatcher := filesystem.NewMatcher()
	treeLister := filesystem.NewLister()
	svc := blueprint.NewService(projectStore, zoneStore, pathMatcher, treeLister, root)
	server := sdkmcp.NewServer(&sdkmcp.Implementation{Name: "test", Version: "0.0.1"}, nil)
	mcp.RegisterTools(server, svc)

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
		Name:      "list_matching_paths",
		Arguments: map[string]any{"pattern": "["},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected tool to return error for invalid regex")
	}
	if len(res.Content) == 0 {
		t.Fatal("expected error content")
	}
	text := contentText(res.Content[0])
	// SDK may serialize error as JSON or as Error(); ensure we see INVALID_PATTERN.
	var errOut struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if _ = json.Unmarshal([]byte(text), &errOut); errOut.Code != "" {
		if errOut.Code != "INVALID_PATTERN" {
			t.Errorf("expected code INVALID_PATTERN, got %q", errOut.Code)
		}
	} else if !strings.Contains(text, "INVALID_PATTERN") {
		t.Errorf("expected error content to contain INVALID_PATTERN, got %q", text)
	}
}

func contentText(c sdkmcp.Content) string {
	if tc, ok := c.(*sdkmcp.TextContent); ok {
		return tc.Text
	}
	return ""
}
