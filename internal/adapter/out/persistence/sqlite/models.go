package sqlite

import "operators-mcp/internal/domain"

// AgentModel is the GORM model for domain.Agent.
type AgentModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	Prompt      string
}

// TableName overrides the table name.
func (AgentModel) TableName() string { return "agents" }

// ToDomain converts the model to a domain.Agent.
func (m *AgentModel) ToDomain() *domain.Agent {
	if m == nil {
		return nil
	}
	return &domain.Agent{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Prompt:      m.Prompt,
	}
}

// ProjectModel is the GORM model for domain.Project.
type ProjectModel struct {
	ID           string `gorm:"primaryKey"`
	Name         string
	RootDir      string      `gorm:"column:root_dir"`
	IgnoredPaths stringSlice `gorm:"column:ignored_paths"`
}

// TableName overrides the table name.
func (ProjectModel) TableName() string { return "projects" }

// ToDomain converts the model to a domain.Project.
func (m *ProjectModel) ToDomain() *domain.Project {
	if m == nil {
		return nil
	}
	paths := []string(m.IgnoredPaths)
	if paths == nil {
		paths = []string{}
	}
	return &domain.Project{
		ID:           m.ID,
		Name:         m.Name,
		RootDir:      m.RootDir,
		IgnoredPaths: paths,
	}
}

// ZoneModel is the GORM model for domain.Zone.
type ZoneModel struct {
	ID             string `gorm:"primaryKey"`
	ProjectID      string `gorm:"column:project_id;index"`
	Name           string
	Pattern        string
	Purpose        string
	Constraints    stringSlice `gorm:"column:constraints"`
	AssignedAgents agentSlice  `gorm:"column:assigned_agents"`
	ExplicitPaths  stringSlice `gorm:"column:explicit_paths"`
}

// TableName overrides the table name.
func (ZoneModel) TableName() string { return "zones" }

// ToDomain converts the model to a domain.Zone.
func (m *ZoneModel) ToDomain() *domain.Zone {
	if m == nil {
		return nil
	}
	return &domain.Zone{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		Name:           m.Name,
		Pattern:        m.Pattern,
		Purpose:        m.Purpose,
		Constraints:    sliceOrNil([]string(m.Constraints)),
		AssignedAgents: sliceAgentsOrNil([]domain.Agent(m.AssignedAgents)),
		ExplicitPaths:  sliceOrNil([]string(m.ExplicitPaths)),
	}
}

func sliceOrNil(s []string) []string {
	if s == nil {
		return nil
	}
	return append([]string(nil), s...)
}

func sliceAgentsOrNil(a []domain.Agent) []domain.Agent {
	if a == nil {
		return nil
	}
	return append([]domain.Agent(nil), a...)
}
