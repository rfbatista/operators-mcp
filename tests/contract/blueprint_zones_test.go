package contract

import (
	"context"
	"encoding/json"
	"testing"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/blueprint"
)

func TestZones_ListCreateGetUpdateAssignPath(t *testing.T) {
	root := t.TempDir()
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

	// list_zones empty
	res, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_zones",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("list_zones: %v", err)
	}
	if res.IsError {
		t.Fatalf("list_zones error: %v", res.Content)
	}
	var listOut struct {
		Zones []map[string]any `json:"zones"`
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &listOut); err != nil {
		t.Fatalf("unmarshal list_zones: %v", err)
	}
	if len(listOut.Zones) != 0 {
		t.Errorf("expected 0 zones, got %d", len(listOut.Zones))
	}

	// create_zone
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name: "create_zone",
		Arguments: map[string]any{
			"name":   "backend",
			"pattern": "cmd/.*",
			"purpose": "Server and CLI",
		},
	})
	if err != nil {
		t.Fatalf("create_zone: %v", err)
	}
	if res.IsError {
		t.Fatalf("create_zone error: %v", res.Content)
	}
	var createOut struct {
		Zone struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Pattern string `json:"pattern"`
		} `json:"zone"`
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &createOut); err != nil {
		t.Fatalf("unmarshal create_zone: %v", err)
	}
	if createOut.Zone.ID == "" {
		t.Error("expected zone id")
	}
	if createOut.Zone.Name != "backend" {
		t.Errorf("expected name backend, got %q", createOut.Zone.Name)
	}
	zoneID := createOut.Zone.ID

	// get_zone
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "get_zone",
		Arguments: map[string]any{"zone_id": zoneID},
	})
	if err != nil {
		t.Fatalf("get_zone: %v", err)
	}
	if res.IsError {
		t.Fatalf("get_zone error: %v", res.Content)
	}
	var getOut struct {
		Zone struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
		} `json:"zone"`
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &getOut); err != nil {
		t.Fatalf("unmarshal get_zone: %v", err)
	}
	if getOut.Zone.Name != "backend" {
		t.Errorf("get_zone name: got %q", getOut.Zone.Name)
	}

	// update_zone
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name: "update_zone",
		Arguments: map[string]any{
			"zone_id": zoneID,
			"name":    "backend-updated",
		},
	})
	if err != nil {
		t.Fatalf("update_zone: %v", err)
	}
	if res.IsError {
		t.Fatalf("update_zone error: %v", res.Content)
	}

	// assign_path_to_zone
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "assign_path_to_zone",
		Arguments: map[string]any{"zone_id": zoneID, "path": "internal/blueprint"},
	})
	if err != nil {
		t.Fatalf("assign_path_to_zone: %v", err)
	}
	if res.IsError {
		t.Fatalf("assign_path_to_zone error: %v", res.Content)
	}
	var assignOut struct {
		Zone struct {
			ExplicitPaths []string `json:"explicit_paths"`
		} `json:"zone"`
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &assignOut); err != nil {
		t.Fatalf("unmarshal assign: %v", err)
	}
	if len(assignOut.Zone.ExplicitPaths) != 1 || assignOut.Zone.ExplicitPaths[0] != "internal/blueprint" {
		t.Errorf("expected explicit_paths [internal/blueprint], got %v", assignOut.Zone.ExplicitPaths)
	}

	// list_zones now has one
	res, err = session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "list_zones",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("list_zones: %v", err)
	}
	if err := json.Unmarshal([]byte(contentText(res.Content[0])), &listOut); err != nil {
		t.Fatalf("unmarshal list_zones: %v", err)
	}
	if len(listOut.Zones) != 1 {
		t.Errorf("expected 1 zone, got %d", len(listOut.Zones))
	}
}

func TestGetZone_NotFound_StructuredError(t *testing.T) {
	root := t.TempDir()
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

	res, err := session.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "get_zone",
		Arguments: map[string]any{"zone_id": "nonexistent"},
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	if !res.IsError {
		t.Fatal("expected tool error for nonexistent zone")
	}
}
