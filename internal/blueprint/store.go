package blueprint

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

// Zone holds zone state (pattern, metadata, explicit paths).
type Zone struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Pattern        string   `json:"pattern"`
	Purpose        string   `json:"purpose"`
	Constraints    []string `json:"constraints"`
	AssignedAgent  string   `json:"assigned_agent"`
	ExplicitPaths  []string `json:"explicit_paths"`
}

// Store holds in-memory zones keyed by id.
type Store struct {
	mu    sync.RWMutex
	zones map[string]*Zone
}

// NewStore returns a new in-memory zone store.
func NewStore() *Store {
	return &Store{zones: make(map[string]*Zone)}
}

// Create creates a zone and returns it with generated id. Name must be non-empty.
func (s *Store) Create(name, pattern, purpose string, constraints []string, assignedAgent string) (*Zone, error) {
	if name == "" {
		return nil, &StructuredError{Code: "INVALID_NAME", Message: "zone name is required"}
	}
	id, err := genID()
	if err != nil {
		return nil, err
	}
	z := &Zone{
		ID:            id,
		Name:          name,
		Pattern:       pattern,
		Purpose:       purpose,
		Constraints:   append([]string(nil), constraints...),
		AssignedAgent: assignedAgent,
		ExplicitPaths: nil,
	}
	s.mu.Lock()
	s.zones[id] = z
	s.mu.Unlock()
	return z, nil
}

// Get returns the zone by id, or nil if not found.
func (s *Store) Get(id string) *Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	z, ok := s.zones[id]
	if !ok {
		return nil
	}
	return cloneZone(z)
}

// List returns all zones.
func (s *Store) List() []*Zone {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Zone, 0, len(s.zones))
	for _, z := range s.zones {
		out = append(out, cloneZone(z))
	}
	return out
}

// Update updates a zone by id. Returns StructuredError if not found or invalid.
func (s *Store) Update(id, name, pattern, purpose string, constraints []string, assignedAgent string) (*Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, ok := s.zones[id]
	if !ok {
		return nil, &StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
	}
	if name != "" {
		z.Name = name
	}
	z.Pattern = pattern
	z.Purpose = purpose
	z.Constraints = append([]string(nil), constraints...)
	z.AssignedAgent = assignedAgent
	return cloneZone(z), nil
}

// AssignPath adds path to zone's explicit_paths. Returns updated zone or error.
func (s *Store) AssignPath(zoneID, path string) (*Zone, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	z, ok := s.zones[zoneID]
	if !ok {
		return nil, &StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
	}
	for _, p := range z.ExplicitPaths {
		if p == path {
			return cloneZone(z), nil
		}
	}
	z.ExplicitPaths = append(z.ExplicitPaths, path)
	return cloneZone(z), nil
}

func cloneZone(z *Zone) *Zone {
	c := *z
	c.Constraints = append([]string(nil), z.Constraints...)
	c.ExplicitPaths = append([]string(nil), z.ExplicitPaths...)
	return &c
}

func genID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
