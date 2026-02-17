package ports

import "operators-mcp/internal/domain"

// ProjectRepository is the outbound port for persisting and retrieving projects.
// A project defines the directory root that everything is based on.
type ProjectRepository interface {
	Get(id string) *domain.Project
	List() []*domain.Project
	Create(name, rootDir string) (*domain.Project, error)
	Update(id, name, rootDir string) (*domain.Project, error)
	AddIgnoredPath(projectID, path string) (*domain.Project, error)
	RemoveIgnoredPath(projectID, path string) (*domain.Project, error)
}

// ZoneRepository is the outbound port for persisting and retrieving zones.
// Zones are scoped to a project. Implemented by adapters (e.g. in-memory store, future DB).
type ZoneRepository interface {
	Get(id string) *domain.Zone
	ListByProject(projectID string) []*domain.Zone
	Create(projectID, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error)
	Update(id, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error)
	AssignPath(zoneID, path string) (*domain.Zone, error)
}

// PathMatcher is the outbound port for listing paths under a root that match a regex pattern.
// Implemented by the filesystem adapter.
type PathMatcher interface {
	ListMatchingPaths(root, pattern string) ([]string, error)
}

// TreeLister is the outbound port for building a directory tree from a root path.
// Implemented by the filesystem adapter.
type TreeLister interface {
	ListTree(root string) (*domain.TreeNode, error)
}
