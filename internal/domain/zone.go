package domain

// Zone holds zone state (pattern, metadata, explicit paths).
// It is the core entity for the blueprint/pattern-management domain.
// A zone belongs to a project and paths are relative to that project's root.
type Zone struct {
	ID             string
	ProjectID      string
	Name           string
	Pattern        string
	Purpose        string
	Constraints    []string
	AssignedAgents []Agent
	ExplicitPaths  []string
}
