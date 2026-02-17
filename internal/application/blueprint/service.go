package blueprint

import (
	"operators-mcp/internal/application/ports"
	"operators-mcp/internal/domain"
)

// Service implements blueprint use cases by delegating to the outbound ports.
// It is the application (use-case) layer in hexagonal architecture.
type Service struct {
	Projects     ports.ProjectRepository
	Zones        ports.ZoneRepository
	PathMatcher  ports.PathMatcher
	TreeLister   ports.TreeLister
	DefaultRoot  string
}

// NewService returns a blueprint application service with the given ports.
// defaultRoot is used when no project_id is provided for list_tree/list_matching_paths.
func NewService(projects ports.ProjectRepository, zones ports.ZoneRepository, pathMatcher ports.PathMatcher, treeLister ports.TreeLister, defaultRoot string) *Service {
	return &Service{
		Projects:    projects,
		Zones:       zones,
		PathMatcher: pathMatcher,
		TreeLister:  treeLister,
		DefaultRoot: defaultRoot,
	}
}

// resolveRoot returns the root path for tree/path operations. If root is non-empty it is used;
// else if projectID is non-empty the project's RootDir is used; otherwise DefaultRoot.
func (s *Service) resolveRoot(root, projectID string) (string, error) {
	if root != "" {
		return root, nil
	}
	if projectID != "" {
		p := s.Projects.Get(projectID)
		if p == nil {
			return "", &domain.StructuredError{Code: "PROJECT_NOT_FOUND", Message: "project not found"}
		}
		return p.RootDir, nil
	}
	return s.DefaultRoot, nil
}

// ListProjects returns all projects.
func (s *Service) ListProjects() []*domain.Project {
	return s.Projects.List()
}

// GetProject returns one project by id, or nil if not found.
func (s *Service) GetProject(projectID string) *domain.Project {
	return s.Projects.Get(projectID)
}

// CreateProject creates a project with the given name and root directory.
func (s *Service) CreateProject(name, rootDir string) (*domain.Project, error) {
	return s.Projects.Create(name, rootDir)
}

// UpdateProject updates an existing project.
func (s *Service) UpdateProject(projectID, name, rootDir string) (*domain.Project, error) {
	return s.Projects.Update(projectID, name, rootDir)
}

// AddIgnoredPath adds a path to the project's ignored list (hidden in tree view).
func (s *Service) AddIgnoredPath(projectID, path string) (*domain.Project, error) {
	return s.Projects.AddIgnoredPath(projectID, path)
}

// RemoveIgnoredPath removes a path from the project's ignored list.
func (s *Service) RemoveIgnoredPath(projectID, path string) (*domain.Project, error) {
	return s.Projects.RemoveIgnoredPath(projectID, path)
}

// ListMatchingPaths returns paths under root that match the regex pattern.
// root and projectID are optional; if both empty, DefaultRoot is used.
func (s *Service) ListMatchingPaths(root, projectID, pattern string) ([]string, error) {
	r, err := s.resolveRoot(root, projectID)
	if err != nil {
		return nil, err
	}
	return s.PathMatcher.ListMatchingPaths(r, pattern)
}

// ListTree returns the directory tree from root.
// root and projectID are optional; if both empty, DefaultRoot is used.
func (s *Service) ListTree(root, projectID string) (*domain.TreeNode, error) {
	r, err := s.resolveRoot(root, projectID)
	if err != nil {
		return nil, err
	}
	return s.TreeLister.ListTree(r)
}

// ListZones returns all zones for the given project.
func (s *Service) ListZones(projectID string) []*domain.Zone {
	return s.Zones.ListByProject(projectID)
}

// GetZone returns one zone by id, or nil if not found.
func (s *Service) GetZone(zoneID string) *domain.Zone {
	return s.Zones.Get(zoneID)
}

// CreateZone creates a zone in the given project with the given metadata.
func (s *Service) CreateZone(projectID, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	return s.Zones.Create(projectID, name, pattern, purpose, constraints, agents)
}

// UpdateZone updates an existing zone.
func (s *Service) UpdateZone(zoneID, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	return s.Zones.Update(zoneID, name, pattern, purpose, constraints, agents)
}

// AssignPathToZone adds a path to a zone's explicit paths (path is normalized).
func (s *Service) AssignPathToZone(zoneID, path string) (*domain.Zone, error) {
	return s.Zones.AssignPath(zoneID, domain.NormalizePath(path))
}
