package mcp

import (
	"context"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"operators-mcp/internal/application/blueprint"
	"operators-mcp/internal/domain"
)

// ListMatchingPathsIn is the input for list_matching_paths.
type ListMatchingPathsIn struct {
	Pattern   string `json:"pattern" jsonschema:"required"`
	Root      string `json:"root,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// ListMatchingPathsOut is the output for list_matching_paths.
type ListMatchingPathsOut struct {
	Paths []string `json:"paths"`
}

// ListTreeIn is the input for list_tree.
type ListTreeIn struct {
	Root      string `json:"root,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	Depth     int    `json:"depth,omitempty"`
}

// ListTreeOut is the output for list_tree.
type ListTreeOut struct {
	Tree any `json:"tree"`
}

// ListZonesIn is the input for list_zones.
type ListZonesIn struct {
	ProjectID string `json:"project_id" jsonschema:"required"`
}

// ListZonesOut is the output for list_zones.
type ListZonesOut struct {
	Zones []*ZoneDTO `json:"zones"`
}

// ListProjectsOut is the output for list_projects.
type ListProjectsOut struct {
	Projects []*ProjectDTO `json:"projects"`
}

// GetProjectIn is the input for get_project.
type GetProjectIn struct {
	ProjectID string `json:"project_id" jsonschema:"required"`
}

// GetProjectOut is the output for get_project.
type GetProjectOut struct {
	Project *ProjectDTO `json:"project"`
}

// CreateProjectIn is the input for create_project.
type CreateProjectIn struct {
	Name    string `json:"name,omitempty"`
	RootDir string `json:"root_dir" jsonschema:"required"`
}

// CreateProjectOut is the output for create_project.
type CreateProjectOut struct {
	Project *ProjectDTO `json:"project"`
}

// UpdateProjectIn is the input for update_project.
type UpdateProjectIn struct {
	ProjectID string `json:"project_id" jsonschema:"required"`
	Name      string `json:"name,omitempty"`
	RootDir   string `json:"root_dir,omitempty"`
}

// UpdateProjectOut is the output for update_project.
type UpdateProjectOut struct {
	Project *ProjectDTO `json:"project"`
}

// AddIgnoredPathIn is the input for add_ignored_path.
type AddIgnoredPathIn struct {
	ProjectID string `json:"project_id" jsonschema:"required"`
	Path      string `json:"path" jsonschema:"required"`
}

// AddIgnoredPathOut is the output for add_ignored_path.
type AddIgnoredPathOut struct {
	Project *ProjectDTO `json:"project"`
}

// RemoveIgnoredPathIn is the input for remove_ignored_path.
type RemoveIgnoredPathIn struct {
	ProjectID string `json:"project_id" jsonschema:"required"`
	Path      string `json:"path" jsonschema:"required"`
}

// RemoveIgnoredPathOut is the output for remove_ignored_path.
type RemoveIgnoredPathOut struct {
	Project *ProjectDTO `json:"project"`
}

// GetZoneIn is the input for get_zone.
type GetZoneIn struct {
	ZoneID string `json:"zone_id" jsonschema:"required"`
}

// GetZoneOut is the output for get_zone.
type GetZoneOut struct {
	Zone *ZoneDTO `json:"zone"`
}

// CreateZoneIn is the input for create_zone.
type CreateZoneIn struct {
	ProjectID      string     `json:"project_id" jsonschema:"required"`
	Name           string     `json:"name" jsonschema:"required"`
	Pattern        string     `json:"pattern,omitempty"`
	Purpose        string     `json:"purpose,omitempty"`
	Constraints    []string   `json:"constraints,omitempty"`
	AssignedAgents []AgentDTO `json:"assigned_agents,omitempty"`
}

// CreateZoneOut is the output for create_zone.
type CreateZoneOut struct {
	Zone *ZoneDTO `json:"zone"`
}

// UpdateZoneIn is the input for update_zone.
type UpdateZoneIn struct {
	ZoneID         string     `json:"zone_id" jsonschema:"required"`
	Name           string     `json:"name,omitempty"`
	Pattern        string     `json:"pattern,omitempty"`
	Purpose        string     `json:"purpose,omitempty"`
	Constraints    []string   `json:"constraints,omitempty"`
	AssignedAgents []AgentDTO `json:"assigned_agents,omitempty"`
}

// UpdateZoneOut is the output for update_zone.
type UpdateZoneOut struct {
	Zone *ZoneDTO `json:"zone"`
}

// AssignPathToZoneIn is the input for assign_path_to_zone.
type AssignPathToZoneIn struct {
	ZoneID string `json:"zone_id" jsonschema:"required"`
	Path   string `json:"path" jsonschema:"required"`
}

// AssignPathToZoneOut is the output for assign_path_to_zone.
type AssignPathToZoneOut struct {
	Zone *ZoneDTO `json:"zone"`
}

// RegisterTools registers all blueprint MCP tools on the server, wiring them to the application service.
func RegisterTools(server *sdkmcp.Server, svc *blueprint.Service) {
	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_projects",
		Description: "Return all projects. A project defines the directory root that everything (tree, zones, paths) is based on.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in struct{}) (*sdkmcp.CallToolResult, ListProjectsOut, error) {
		projects := svc.ListProjects()
		return nil, ListProjectsOut{Projects: ProjectsToDTO(projects)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_project",
		Description: "Return one project by id.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in GetProjectIn) (*sdkmcp.CallToolResult, GetProjectOut, error) {
		p := svc.GetProject(in.ProjectID)
		if p == nil {
			return nil, GetProjectOut{}, &domain.StructuredError{Code: "PROJECT_NOT_FOUND", Message: "project not found"}
		}
		return nil, GetProjectOut{Project: ProjectToDTO(p)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "create_project",
		Description: "Create a project with a name and root directory. The root is the base path for list_tree, list_matching_paths, and zones.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in CreateProjectIn) (*sdkmcp.CallToolResult, CreateProjectOut, error) {
		p, err := svc.CreateProject(in.Name, in.RootDir)
		if err != nil {
			return nil, CreateProjectOut{}, err
		}
		return nil, CreateProjectOut{Project: ProjectToDTO(p)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "update_project",
		Description: "Update a project's name and/or root_dir.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in UpdateProjectIn) (*sdkmcp.CallToolResult, UpdateProjectOut, error) {
		p, err := svc.UpdateProject(in.ProjectID, in.Name, in.RootDir)
		if err != nil {
			return nil, UpdateProjectOut{}, err
		}
		return nil, UpdateProjectOut{Project: ProjectToDTO(p)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "add_ignored_path",
		Description: "Add a file or directory path to the project's ignore list. Ignored paths are hidden from the tree view.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in AddIgnoredPathIn) (*sdkmcp.CallToolResult, AddIgnoredPathOut, error) {
		p, err := svc.AddIgnoredPath(in.ProjectID, in.Path)
		if err != nil {
			return nil, AddIgnoredPathOut{}, err
		}
		return nil, AddIgnoredPathOut{Project: ProjectToDTO(p)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "remove_ignored_path",
		Description: "Remove a path from the project's ignore list so it is shown again in the tree view.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in RemoveIgnoredPathIn) (*sdkmcp.CallToolResult, RemoveIgnoredPathOut, error) {
		p, err := svc.RemoveIgnoredPath(in.ProjectID, in.Path)
		if err != nil {
			return nil, RemoveIgnoredPathOut{}, err
		}
		return nil, RemoveIgnoredPathOut{Project: ProjectToDTO(p)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_matching_paths",
		Description: "Return paths under project root that match the given regex pattern. Use project_id or root to specify the base directory.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in ListMatchingPathsIn) (*sdkmcp.CallToolResult, ListMatchingPathsOut, error) {
		paths, err := svc.ListMatchingPaths(in.Root, in.ProjectID, in.Pattern)
		if err != nil {
			return nil, ListMatchingPathsOut{}, err
		}
		return nil, ListMatchingPathsOut{Paths: paths}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_tree",
		Description: "Return the project's folder structure as a hierarchical tree. Use project_id or root to specify the base directory.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in ListTreeIn) (*sdkmcp.CallToolResult, ListTreeOut, error) {
		tree, err := svc.ListTree(in.Root, in.ProjectID)
		if err != nil {
			return nil, ListTreeOut{}, err
		}
		return nil, ListTreeOut{Tree: TreeNodeToDTO(tree)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "list_zones",
		Description: "Return all zones for the given project.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in ListZonesIn) (*sdkmcp.CallToolResult, ListZonesOut, error) {
		zones := svc.ListZones(in.ProjectID)
		return nil, ListZonesOut{Zones: ZonesToDTO(zones)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "get_zone",
		Description: "Return one zone by id.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in GetZoneIn) (*sdkmcp.CallToolResult, GetZoneOut, error) {
		z := svc.GetZone(in.ZoneID)
		if z == nil {
			return nil, GetZoneOut{}, &domain.StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
		}
		return nil, GetZoneOut{Zone: ZoneToDTO(z)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "create_zone",
		Description: "Create a zone in the given project with optional metadata and pattern.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in CreateZoneIn) (*sdkmcp.CallToolResult, CreateZoneOut, error) {
		z, err := svc.CreateZone(in.ProjectID, in.Name, in.Pattern, in.Purpose, in.Constraints, DTOToAgents(in.AssignedAgents))
		if err != nil {
			return nil, CreateZoneOut{}, err
		}
		return nil, CreateZoneOut{Zone: ZoneToDTO(z)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "update_zone",
		Description: "Update zone name, pattern, purpose, constraints, assigned_agents.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in UpdateZoneIn) (*sdkmcp.CallToolResult, UpdateZoneOut, error) {
		z, err := svc.UpdateZone(in.ZoneID, in.Name, in.Pattern, in.Purpose, in.Constraints, DTOToAgents(in.AssignedAgents))
		if err != nil {
			return nil, UpdateZoneOut{}, err
		}
		return nil, UpdateZoneOut{Zone: ZoneToDTO(z)}, nil
	})

	sdkmcp.AddTool(server, &sdkmcp.Tool{
		Name:        "assign_path_to_zone",
		Description: "Add a path to a zone's explicit path set.",
	}, func(ctx context.Context, req *sdkmcp.CallToolRequest, in AssignPathToZoneIn) (*sdkmcp.CallToolResult, AssignPathToZoneOut, error) {
		z, err := svc.AssignPathToZone(in.ZoneID, in.Path)
		if err != nil {
			return nil, AssignPathToZoneOut{}, err
		}
		return nil, AssignPathToZoneOut{Zone: ZoneToDTO(z)}, nil
	})
}
