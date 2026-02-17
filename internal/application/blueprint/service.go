package blueprint

import (
	"operators-mcp/internal/application/ports"
	"operators-mcp/internal/domain"
)

// Service implements blueprint use cases by delegating to the outbound ports.
// It is the application (use-case) layer in hexagonal architecture.
type Service struct {
	Zones       ports.ZoneRepository
	PathMatcher ports.PathMatcher
	TreeLister  ports.TreeLister
}

// NewService returns a blueprint application service with the given ports.
func NewService(zones ports.ZoneRepository, pathMatcher ports.PathMatcher, treeLister ports.TreeLister) *Service {
	return &Service{
		Zones:       zones,
		PathMatcher: pathMatcher,
		TreeLister:  treeLister,
	}
}

// ListMatchingPaths returns paths under root that match the regex pattern.
func (s *Service) ListMatchingPaths(root, pattern string) ([]string, error) {
	return s.PathMatcher.ListMatchingPaths(root, pattern)
}

// ListTree returns the directory tree from root.
func (s *Service) ListTree(root string) (*domain.TreeNode, error) {
	return s.TreeLister.ListTree(root)
}

// ListZones returns all zones.
func (s *Service) ListZones() []*domain.Zone {
	return s.Zones.List()
}

// GetZone returns one zone by id, or nil if not found.
func (s *Service) GetZone(zoneID string) *domain.Zone {
	return s.Zones.Get(zoneID)
}

// CreateZone creates a zone with the given metadata.
func (s *Service) CreateZone(name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	return s.Zones.Create(name, pattern, purpose, constraints, agents)
}

// UpdateZone updates an existing zone.
func (s *Service) UpdateZone(zoneID, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	return s.Zones.Update(zoneID, name, pattern, purpose, constraints, agents)
}

// AssignPathToZone adds a path to a zone's explicit paths (path is normalized).
func (s *Service) AssignPathToZone(zoneID, path string) (*domain.Zone, error) {
	return s.Zones.AssignPath(zoneID, domain.NormalizePath(path))
}
