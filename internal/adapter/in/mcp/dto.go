package mcp

import "operators-mcp/internal/domain"

// AgentDTO is the MCP/JSON representation of an agent.
type AgentDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ProjectDTO is the MCP/JSON representation of a project (snake_case for API contract).
type ProjectDTO struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	RootDir      string   `json:"root_dir"`
	IgnoredPaths []string `json:"ignored_paths,omitempty"`
}

// ZoneDTO is the MCP/JSON representation of a zone (snake_case for API contract).
type ZoneDTO struct {
	ID             string     `json:"id"`
	ProjectID      string     `json:"project_id"`
	Name           string     `json:"name"`
	Pattern        string     `json:"pattern"`
	Purpose        string     `json:"purpose"`
	Constraints    []string   `json:"constraints"`
	AssignedAgents []AgentDTO `json:"assigned_agents"`
	ExplicitPaths  []string   `json:"explicit_paths"`
}

// TreeNodeDTO is the MCP/JSON representation of a tree node.
type TreeNodeDTO struct {
	Path     string         `json:"path"`
	Name     string         `json:"name"`
	IsDir    bool           `json:"is_dir"`
	Children []*TreeNodeDTO `json:"children"`
}

// ProjectToDTO converts a domain Project to API DTO (exported for HTTP adapter).
func ProjectToDTO(p *domain.Project) *ProjectDTO {
	if p == nil {
		return nil
	}
	ignored := p.IgnoredPaths
	if len(ignored) == 0 {
		ignored = nil
	}
	return &ProjectDTO{
		ID:           p.ID,
		Name:         p.Name,
		RootDir:      p.RootDir,
		IgnoredPaths: ignored,
	}
}

// ProjectsToDTO converts domain projects to DTOs.
func ProjectsToDTO(projects []*domain.Project) []*ProjectDTO {
	out := make([]*ProjectDTO, len(projects))
	for i, p := range projects {
		out[i] = ProjectToDTO(p)
	}
	return out
}

// ZoneToDTO converts a domain Zone to API DTO (exported for HTTP adapter).
func ZoneToDTO(z *domain.Zone) *ZoneDTO {
	if z == nil {
		return nil
	}
	return &ZoneDTO{
		ID:             z.ID,
		ProjectID:      z.ProjectID,
		Name:           z.Name,
		Pattern:        z.Pattern,
		Purpose:        z.Purpose,
		Constraints:    append([]string(nil), z.Constraints...),
		AssignedAgents: AgentsToDTO(z.AssignedAgents),
		ExplicitPaths:  append([]string(nil), z.ExplicitPaths...),
	}
}

// AgentsToDTO converts domain agents to DTOs (exported for HTTP adapter).
func AgentsToDTO(a []domain.Agent) []AgentDTO {
	if len(a) == 0 {
		return nil
	}
	out := make([]AgentDTO, len(a))
	for i := range a {
		out[i] = AgentDTO{ID: a[i].ID, Name: a[i].Name}
	}
	return out
}

// DTOToAgents converts AgentDTO slice to domain agents (exported for HTTP adapter).
func DTOToAgents(a []AgentDTO) []domain.Agent {
	if len(a) == 0 {
		return nil
	}
	out := make([]domain.Agent, len(a))
	for i := range a {
		out[i] = domain.Agent{ID: a[i].ID, Name: a[i].Name}
	}
	return out
}

// ZonesToDTO converts domain zones to DTOs (exported for HTTP adapter).
func ZonesToDTO(zones []*domain.Zone) []*ZoneDTO {
	out := make([]*ZoneDTO, len(zones))
	for i, z := range zones {
		out[i] = ZoneToDTO(z)
	}
	return out
}

// TreeNodeToDTO converts a domain TreeNode to API DTO (exported for HTTP adapter).
func TreeNodeToDTO(n *domain.TreeNode) *TreeNodeDTO {
	if n == nil {
		return nil
	}
	children := make([]*TreeNodeDTO, len(n.Children))
	for i, c := range n.Children {
		children[i] = TreeNodeToDTO(c)
	}
	return &TreeNodeDTO{
		Path:     n.Path,
		Name:     n.Name,
		IsDir:    n.IsDir,
		Children: children,
	}
}
