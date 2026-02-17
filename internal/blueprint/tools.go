package blueprint

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

// ListMatchingPathsIn is the input for list_matching_paths.
type ListMatchingPathsIn struct {
	Pattern string `json:"pattern" jsonschema:"required"`
	Root    string `json:"root,omitempty"`
}

// ListMatchingPathsOut is the output for list_matching_paths.
type ListMatchingPathsOut struct {
	Paths []string `json:"paths"`
}

// ListTreeIn is the input for list_tree.
type ListTreeIn struct {
	Root  string `json:"root,omitempty"`
	Depth int    `json:"depth,omitempty"`
}

// ListTreeOut is the output for list_tree. Tree is typed as any to avoid
// recursive JSON schema inference (TreeNode.Children references TreeNode).
type ListTreeOut struct {
	Tree any `json:"tree"`
}

// ListZonesOut is the output for list_zones.
type ListZonesOut struct {
	Zones []*Zone `json:"zones"`
}

// GetZoneIn is the input for get_zone.
type GetZoneIn struct {
	ZoneID string `json:"zone_id" jsonschema:"required"`
}

// GetZoneOut is the output for get_zone.
type GetZoneOut struct {
	Zone *Zone `json:"zone"`
}

// CreateZoneIn is the input for create_zone.
type CreateZoneIn struct {
	Name          string   `json:"name" jsonschema:"required"`
	Pattern       string   `json:"pattern,omitempty"`
	Purpose       string   `json:"purpose,omitempty"`
	Constraints   []string `json:"constraints,omitempty"`
	AssignedAgent string   `json:"assigned_agent,omitempty"`
}

// CreateZoneOut is the output for create_zone.
type CreateZoneOut struct {
	Zone *Zone `json:"zone"`
}

// UpdateZoneIn is the input for update_zone.
type UpdateZoneIn struct {
	ZoneID        string   `json:"zone_id" jsonschema:"required"`
	Name          string   `json:"name,omitempty"`
	Pattern       string   `json:"pattern,omitempty"`
	Purpose       string   `json:"purpose,omitempty"`
	Constraints   []string `json:"constraints,omitempty"`
	AssignedAgent string   `json:"assigned_agent,omitempty"`
}

// UpdateZoneOut is the output for update_zone.
type UpdateZoneOut struct {
	Zone *Zone `json:"zone"`
}

// AssignPathToZoneIn is the input for assign_path_to_zone.
type AssignPathToZoneIn struct {
	ZoneID string `json:"zone_id" jsonschema:"required"`
	Path   string `json:"path" jsonschema:"required"`
}

// AssignPathToZoneOut is the output for assign_path_to_zone.
type AssignPathToZoneOut struct {
	Zone *Zone `json:"zone"`
}

// RegisterTools registers all blueprint MCP tools on the server.
func RegisterTools(server *sdkmcp.Server, root string, store *Store) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_matching_paths",
		Description: "Return paths under project root that match the given regex pattern.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in ListMatchingPathsIn) (*sdkmcp.CallToolResult, ListMatchingPathsOut, error) {
		r := root
		if in.Root != "" {
			r = in.Root
		}
		paths, err := ListMatchingPaths(r, in.Pattern)
		if err != nil {
			return nil, ListMatchingPathsOut{}, err
		}
		return nil, ListMatchingPathsOut{Paths: paths}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_tree",
		Description: "Return the project's folder structure as a hierarchical tree.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in ListTreeIn) (*sdkmcp.CallToolResult, ListTreeOut, error) {
		r := root
		if in.Root != "" {
			r = in.Root
		}
		tree, err := ListTree(r)
		if err != nil {
			return nil, ListTreeOut{}, err
		}
		return nil, ListTreeOut{Tree: tree}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_zones",
		Description: "Return all zones.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in struct{}) (*sdkmcp.CallToolResult, ListZonesOut, error) {
		return nil, ListZonesOut{Zones: store.List()}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_zone",
		Description: "Return one zone by id.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in GetZoneIn) (*sdkmcp.CallToolResult, GetZoneOut, error) {
		z := store.Get(in.ZoneID)
		if z == nil {
			return nil, GetZoneOut{}, &StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
		}
		return nil, GetZoneOut{Zone: z}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "create_zone",
		Description: "Create a zone with optional metadata and pattern.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in CreateZoneIn) (*sdkmcp.CallToolResult, CreateZoneOut, error) {
		z, err := store.Create(in.Name, in.Pattern, in.Purpose, in.Constraints, in.AssignedAgent)
		if err != nil {
			return nil, CreateZoneOut{}, err
		}
		return nil, CreateZoneOut{Zone: z}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "update_zone",
		Description: "Update zone name, pattern, purpose, constraints, assigned_agent.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in UpdateZoneIn) (*sdkmcp.CallToolResult, UpdateZoneOut, error) {
		z, err := store.Update(in.ZoneID, in.Name, in.Pattern, in.Purpose, in.Constraints, in.AssignedAgent)
		if err != nil {
			return nil, UpdateZoneOut{}, err
		}
		return nil, UpdateZoneOut{Zone: z}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "assign_path_to_zone",
		Description: "Add a path to a zone's explicit path set.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in AssignPathToZoneIn) (*sdkmcp.CallToolResult, AssignPathToZoneOut, error) {
		z, err := store.AssignPath(in.ZoneID, NormalizePath(in.Path))
		if err != nil {
			return nil, AssignPathToZoneOut{}, err
		}
		return nil, AssignPathToZoneOut{Zone: z}, nil
	})
}
