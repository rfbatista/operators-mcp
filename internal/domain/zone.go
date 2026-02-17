package domain

// Zone holds zone state (pattern, metadata, explicit paths).
// It is the core entity for the blueprint/pattern-management domain.
// A zone can have multiple assigned agents.
type Zone struct {
	ID             string
	Name           string
	Pattern        string
	Purpose        string
	Constraints    []string
	AssignedAgents []Agent
	ExplicitPaths  []string
}
