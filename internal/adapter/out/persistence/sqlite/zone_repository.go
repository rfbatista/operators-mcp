package sqlite

import (
	"operators-mcp/internal/application/ports"
	"operators-mcp/internal/domain"

	"gorm.io/gorm"
)

// Ensure ZoneRepository implements ports.ZoneRepository at compile time.
var _ ports.ZoneRepository = (*ZoneRepository)(nil)

// ZoneRepository persists zones in SQLite via GORM.
type ZoneRepository struct {
	db *gorm.DB
}

// NewZoneRepository returns a new zone repository.
func NewZoneRepository(db *gorm.DB) *ZoneRepository {
	return &ZoneRepository{db: db}
}

// Get returns the zone by id, or nil if not found.
func (r *ZoneRepository) Get(id string) *domain.Zone {
	var m ZoneModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return nil
	}
	return m.ToDomain()
}

// ListByProject returns all zones for the given project.
func (r *ZoneRepository) ListByProject(projectID string) []*domain.Zone {
	var models []ZoneModel
	if err := r.db.Where("project_id = ?", projectID).Find(&models).Error; err != nil {
		return nil
	}
	out := make([]*domain.Zone, 0, len(models))
	for i := range models {
		out = append(out, models[i].ToDomain())
	}
	return out
}

// Create creates a zone in the given project and returns it with generated id.
func (r *ZoneRepository) Create(projectID, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	if name == "" {
		return nil, &domain.StructuredError{Code: "INVALID_NAME", Message: "zone name is required"}
	}
	id, err := genID()
	if err != nil {
		return nil, err
	}
	agentsCopy := cloneAgents(agents)
	if agentsCopy == nil {
		agentsCopy = []domain.Agent{}
	}
	m := &ZoneModel{
		ID:             id,
		ProjectID:      projectID,
		Name:           name,
		Pattern:        pattern,
		Purpose:        purpose,
		Constraints:    stringSlice(append([]string(nil), constraints...)),
		AssignedAgents: agentSlice(agentsCopy),
		ExplicitPaths:  nil,
	}
	if err := r.db.Create(m).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

// Update updates a zone by id.
func (r *ZoneRepository) Update(id, name, pattern, purpose string, constraints []string, agents []domain.Agent) (*domain.Zone, error) {
	var m ZoneModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
		}
		return nil, err
	}
	updates := map[string]interface{}{
		"pattern":         pattern,
		"purpose":         purpose,
		"constraints":     stringSlice(append([]string(nil), constraints...)),
		"assigned_agents": agentSlice(cloneAgents(agents)),
	}
	if name != "" {
		updates["name"] = name
	}
	if err := r.db.Model(&m).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

// AssignPath adds path to zone's explicit paths.
func (r *ZoneRepository) AssignPath(zoneID, path string) (*domain.Zone, error) {
	var m ZoneModel
	if err := r.db.First(&m, "id = ?", zoneID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &domain.StructuredError{Code: "ZONE_NOT_FOUND", Message: "zone not found"}
		}
		return nil, err
	}
	for _, p := range m.ExplicitPaths {
		if p == path {
			return m.ToDomain(), nil
		}
	}
	m.ExplicitPaths = append(m.ExplicitPaths, path)
	if err := r.db.Model(&m).Update("explicit_paths", m.ExplicitPaths).Error; err != nil {
		return nil, err
	}
	return m.ToDomain(), nil
}

func cloneAgents(a []domain.Agent) []domain.Agent {
	if a == nil {
		return nil
	}
	return append([]domain.Agent(nil), a...)
}
