package mcp

import (
	"github.com/google/jsonschema-go/jsonschema"
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

// emptyIn is used for ListTools schema (HTTP /api/tools).
type emptyIn struct{}

// ToolDescriptor describes an MCP tool for discovery (e.g. HTTP GET /api/tools).
type ToolDescriptor struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema,omitempty"`
}

// ListTools returns all blueprint tool descriptors with input schemas for MCP discovery.
func ListTools() []ToolDescriptor {
	schemaEmpty, _ := jsonschema.For[emptyIn](nil)
	schemaGetProject, _ := jsonschema.For[GetProjectIn](nil)
	schemaCreateProject, _ := jsonschema.For[CreateProjectIn](nil)
	schemaUpdateProject, _ := jsonschema.For[UpdateProjectIn](nil)
	schemaAddIgnoredPath, _ := jsonschema.For[AddIgnoredPathIn](nil)
	schemaRemoveIgnoredPath, _ := jsonschema.For[RemoveIgnoredPathIn](nil)
	schemaListMatchingPaths, _ := jsonschema.For[ListMatchingPathsIn](nil)
	schemaListTree, _ := jsonschema.For[ListTreeIn](nil)
	schemaListZones, _ := jsonschema.For[ListZonesIn](nil)
	schemaGetZone, _ := jsonschema.For[GetZoneIn](nil)
	schemaCreateZone, _ := jsonschema.For[CreateZoneIn](nil)
	schemaUpdateZone, _ := jsonschema.For[UpdateZoneIn](nil)
	schemaAssignPathToZone, _ := jsonschema.For[AssignPathToZoneIn](nil)

	return []ToolDescriptor{
		{"list_projects", "Return all projects. A project defines the directory root that everything (tree, zones, paths) is based on.", schemaEmpty},
		{"get_project", "Return one project by id.", schemaGetProject},
		{"create_project", "Create a project with a name and root directory. The root is the base path for list_tree, list_matching_paths, and zones.", schemaCreateProject},
		{"update_project", "Update a project's name and/or root_dir.", schemaUpdateProject},
		{"add_ignored_path", "Add a file or directory path to the project's ignore list. Ignored paths are hidden from the tree view.", schemaAddIgnoredPath},
		{"remove_ignored_path", "Remove a path from the project's ignore list so it is shown again in the tree view.", schemaRemoveIgnoredPath},
		{"list_matching_paths", "Return paths under project root that match the given regex pattern. Use project_id or root to specify the base directory.", schemaListMatchingPaths},
		{"list_tree", "Return the project's folder structure as a hierarchical tree. Use project_id or root to specify the base directory.", schemaListTree},
		{"list_zones", "Return all zones for the given project.", schemaListZones},
		{"get_zone", "Return one zone by id.", schemaGetZone},
		{"create_zone", "Create a zone in the given project with optional metadata and pattern.", schemaCreateZone},
		{"update_zone", "Update zone name, pattern, purpose, constraints, assigned_agents.", schemaUpdateZone},
		{"assign_path_to_zone", "Add a path to a zone's explicit path set.", schemaAssignPathToZone},
	}
}
