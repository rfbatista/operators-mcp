package memory

import (
	"crypto/rand"
	"encoding/hex"
	"sync"

	"operators-mcp/internal/application/ports"
	"operators-mcp/internal/domain"
)

// Ensure Store implements ports.ZoneRepository at compile time.
var _ ports.ZoneRepository = (*Store)(nil)

// Store holds in-memory zones keyed by id.
type Store struct {
	mu    sync.RWMutex
	zones map[string]*domain.Zone
}

// NewStore returns a new in-memory zone store.
func NewStore() *Store {
	return &Store{zones: make(map[string]*domain.Zone)}
}

// Get returns the zone by id, or nil if not found.
func (s *Store) Get(id string) *domain.Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	z, ok := s.zones[id]
	if !ok {
		return nil
	}
	return cloneZone(z)
}

// List returns all zones.
func (s *Store) List() []*domain.Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*domain.Zone, 0, len(s.zones))
	for _, z := range s.zones {
		out = append(out, cloneZone(z))
	}
	return out
}

// Create creates a zone and returns it with generated id. Name must be non-empty.
func (s *Store) Create(name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	if name == "" {
		return nil, &domain.StructuredError{Code: "INVALID_NAME", Message: "zone name is required"}
	}
	id, err := genID()
	if err != nil {
		return nil, err
	}
	z := &domain.Zone{
		ID:             id,
		Name:           name,
		Pattern:        pattern,
		Purpose:        purpose,
		Constraints:    append([]string(nil), constraints...),
		AssignedAgents: cloneAgents(agents),
		ExplicitPaths:  nil,
	}
	s.mu.Lock()
	s.zones[id] = z
	s.mu.Unlock()
	return cloneZone(z), nil
}

// Update updates a zone by id. Returns StructuredError if not found or invalid.
func (s *Store) Update(id, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, ok := s.zones[id]
	if !ok {
		return nil, &domain.StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
	}
	if name != "" {
		z.Name = name
	}
	z.Pattern = pattern
	z.Purpose = purpose
	z.Constraints = append([]string(nil), constraints...)
	z.AssignedAgents = cloneAgents(agents)
	return cloneZone(z), nil
}

// AssignPath adds path to zone's explicit paths. Returns updated zone or error.
func (s *Store) AssignPath(zoneID, path string) (*domain.Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, ok := s.zones[zoneID]
	if !ok {
		return nil, &domain.StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
	}
	for _, p := range z.ExplicitPaths {
		if p == path {
			return cloneZone(z), nil
		}
	}
	z.ExplicitPaths = append(z.ExplicitPaths, path)
	return cloneZone(z), nil
}

func cloneZone(z *domain.Zone) *domain.Zone {
	c := *z
	c.Constraints = append([]string(nil), z.Constraints...)
	c.ExplicitPaths = append([]string(nil), z.ExplicitPaths...)
	c.AssignedAgents = cloneAgents(z.AssignedAgents)
	return &c
}

func cloneAgents(a []domain.Agent) []domain.Agent {
	if a == nil {
		return nil
	}
	return append([]domain.Agent(nil), a...)
}

func genID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
