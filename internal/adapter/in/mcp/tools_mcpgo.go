package mcp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"operators-mcp/internal/application/blueprint"
	"operators-mcp/internal/domain"
)

// RegisterTools registers all blueprint tools on the mcp-go server.
func RegisterTools(s *server.MCPServer, svc *blueprint.Service) {
	// list_projects
	s.AddTool(mcp.NewTool("list_projects",
		mcp.WithDescription("Return all projects. A project defines the directory root that everything (tree, zones, paths) is based on."),
	), toolListProjects(svc))

	// get_project
	s.AddTool(mcp.NewTool("get_project",
		mcp.WithDescription("Return one project by id."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
	), toolGetProject(svc))

	// create_project
	s.AddTool(mcp.NewTool("create_project",
		mcp.WithDescription("Create a project with a name and root directory. The root is the base path for list_tree, list_matching_paths, and zones."),
		mcp.WithString("name", mcp.Description("Project name")),
		mcp.WithString("root_dir", mcp.Required(), mcp.Description("Root directory path")),
	), toolCreateProject(svc))

	// update_project
	s.AddTool(mcp.NewTool("update_project",
		mcp.WithDescription("Update a project's name and/or root_dir."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("name", mcp.Description("Project name")),
		mcp.WithString("root_dir", mcp.Description("Root directory path")),
	), toolUpdateProject(svc))

	// add_ignored_path
	s.AddTool(mcp.NewTool("add_ignored_path",
		mcp.WithDescription("Add a file or directory path to the project's ignore list. Ignored paths are hidden from the tree view."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to ignore")),
	), toolAddIgnoredPath(svc))

	// remove_ignored_path
	s.AddTool(mcp.NewTool("remove_ignored_path",
		mcp.WithDescription("Remove a path from the project's ignore list so it is shown again in the tree view."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to remove from ignore list")),
	), toolRemoveIgnoredPath(svc))

	// list_matching_paths
	s.AddTool(mcp.NewTool("list_matching_paths",
		mcp.WithDescription("Return paths under project root that match the given regex pattern. Use project_id or root to specify the base directory."),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Regex pattern")),
		mcp.WithString("root", mcp.Description("Root path (optional)")),
		mcp.WithString("project_id", mcp.Description("Project ID (optional)")),
	), toolListMatchingPaths(svc))

	// list_tree
	s.AddTool(mcp.NewTool("list_tree",
		mcp.WithDescription("Return the project's folder structure as a hierarchical tree. Use project_id or root to specify the base directory."),
		mcp.WithString("root", mcp.Description("Root path (optional)")),
		mcp.WithString("project_id", mcp.Description("Project ID (optional)")),
		mcp.WithNumber("depth", mcp.Description("Max depth (optional)")),
	), toolListTree(svc))

	// list_zones
	s.AddTool(mcp.NewTool("list_zones",
		mcp.WithDescription("Return all zones for the given project."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
	), toolListZones(svc))

	// get_zone
	s.AddTool(mcp.NewTool("get_zone",
		mcp.WithDescription("Return one zone by id."),
		mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
	), toolGetZone(svc))

	// create_zone
	s.AddTool(mcp.NewTool("create_zone",
		mcp.WithDescription("Create a zone in the given project with optional metadata and pattern."),
		mcp.WithString("project_id", mcp.Required(), mcp.Description("Project ID")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Zone name")),
		mcp.WithString("pattern", mcp.Description("Regex pattern")),
		mcp.WithString("purpose", mcp.Description("Purpose")),
		mcp.WithArray("constraints", mcp.Description("Constraints"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithAny("assigned_agents", mcp.Description("Assigned agents (array of {id, name})")),
	), toolCreateZone(svc))

	// update_zone
	s.AddTool(mcp.NewTool("update_zone",
		mcp.WithDescription("Update zone name, pattern, purpose, constraints, assigned_agents."),
		mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
		mcp.WithString("name", mcp.Description("Zone name")),
		mcp.WithString("pattern", mcp.Description("Regex pattern")),
		mcp.WithString("purpose", mcp.Description("Purpose")),
		mcp.WithArray("constraints", mcp.Description("Constraints"), mcp.Items(map[string]any{"type": "string"})),
		mcp.WithAny("assigned_agents", mcp.Description("Assigned agents (array of {id, name})")),
	), toolUpdateZone(svc))

	// assign_path_to_zone
	s.AddTool(mcp.NewTool("assign_path_to_zone",
		mcp.WithDescription("Add a path to a zone's explicit path set."),
		mcp.WithString("zone_id", mcp.Required(), mcp.Description("Zone ID")),
		mcp.WithString("path", mcp.Required(), mcp.Description("Path to assign")),
	), toolAssignPathToZone(svc))
}

func toolListProjects(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects := svc.ListProjects()
		out := ListProjectsOut{Projects: ProjectsToDTO(projects)}
		return jsonResult(out)
	}
}

func toolGetProject(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		p := svc.GetProject(projectID)
		if p == nil {
			return mcp.NewToolResultError("project not found"), nil
		}
		return jsonResult(GetProjectOut{Project: ProjectToDTO(p)})
	}
}

func toolCreateProject(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		rootDir, err := req.RequireString("root_dir")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name := req.GetString("name", "")
		p, err := svc.CreateProject(name, rootDir)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(CreateProjectOut{Project: ProjectToDTO(p)})
	}
}

func toolUpdateProject(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name := req.GetString("name", "")
		rootDir := req.GetString("root_dir", "")
		p, err := svc.UpdateProject(projectID, name, rootDir)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(UpdateProjectOut{Project: ProjectToDTO(p)})
	}
}

func toolAddIgnoredPath(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		p, err := svc.AddIgnoredPath(projectID, path)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(AddIgnoredPathOut{Project: ProjectToDTO(p)})
	}
}

func toolRemoveIgnoredPath(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		p, err := svc.RemoveIgnoredPath(projectID, path)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(RemoveIgnoredPathOut{Project: ProjectToDTO(p)})
	}
}

func toolListMatchingPaths(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pattern, err := req.RequireString("pattern")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		root := req.GetString("root", "")
		projectID := req.GetString("project_id", "")
		paths, err := svc.ListMatchingPaths(root, projectID, pattern)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(ListMatchingPathsOut{Paths: paths})
	}
}

func toolListTree(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		root := req.GetString("root", "")
		projectID := req.GetString("project_id", "")
		tree, err := svc.ListTree(root, projectID)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(ListTreeOut{Tree: TreeNodeToDTO(tree)})
	}
}

func toolListZones(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		zones := svc.ListZones(projectID)
		return jsonResult(ListZonesOut{Zones: ZonesToDTO(zones)})
	}
}

func toolGetZone(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, err := req.RequireString("zone_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		z := svc.GetZone(zoneID)
		if z == nil {
			return mcp.NewToolResultError("zone not found"), nil
		}
		return jsonResult(GetZoneOut{Zone: ZoneToDTO(z)})
	}
}

func toolCreateZone(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID, err := req.RequireString("project_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		pattern := req.GetString("pattern", "")
		purpose := req.GetString("purpose", "")
		constraints := req.GetStringSlice("constraints", []string{})
		args := req.GetArguments()
		agentsRaw := args["assigned_agents"]
		var agents []AgentDTO
		if agentsRaw != nil {
			if slice, ok := agentsRaw.([]any); ok {
				for _, v := range slice {
					if m, ok := v.(map[string]any); ok {
						id, _ := m["id"].(string)
						n, _ := m["name"].(string)
						agents = append(agents, AgentDTO{ID: id, Name: n})
					}
				}
			}
		}
		z, err := svc.CreateZone(projectID, name, pattern, purpose, constraints, DTOToAgents(agents))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(CreateZoneOut{Zone: ZoneToDTO(z)})
	}
}

func toolUpdateZone(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, err := req.RequireString("zone_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		name := req.GetString("name", "")
		pattern := req.GetString("pattern", "")
		purpose := req.GetString("purpose", "")
		constraints := req.GetStringSlice("constraints", []string{})
		args := req.GetArguments()
		agentsRaw := args["assigned_agents"]
		var agents []AgentDTO
		if agentsRaw != nil {
			if slice, ok := agentsRaw.([]any); ok {
				for _, v := range slice {
					if m, ok := v.(map[string]any); ok {
						id, _ := m["id"].(string)
						n, _ := m["name"].(string)
						agents = append(agents, AgentDTO{ID: id, Name: n})
					}
				}
			}
		}
		z, err := svc.UpdateZone(zoneID, name, pattern, purpose, constraints, DTOToAgents(agents))
		if err != nil {
			return toolError(err)
		}
		return jsonResult(UpdateZoneOut{Zone: ZoneToDTO(z)})
	}
}

func toolAssignPathToZone(svc *blueprint.Service) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		zoneID, err := req.RequireString("zone_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		path, err := req.RequireString("path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		z, err := svc.AssignPathToZone(zoneID, path)
		if err != nil {
			return toolError(err)
		}
		return jsonResult(AssignPathToZoneOut{Zone: ZoneToDTO(z)})
	}
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func toolError(err error) (*mcp.CallToolResult, error) {
	var se *domain.StructuredError
	if errors.As(err, &se) {
		return mcp.NewToolResultError(se.Message), nil
	}
	return mcp.NewToolResultError(err.Error()), nil
}
