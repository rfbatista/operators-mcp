package ports

import "operators-mcp/internal/domain"

// ZoneRepository is the outbound port for persisting and retrieving zones.
// Implemented by adapters (e.g. in-memory store, future DB).
type ZoneRepository interface {
	Get(id string) *domain.Zone
	List() []*domain.Zone
	Create(name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error)
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
