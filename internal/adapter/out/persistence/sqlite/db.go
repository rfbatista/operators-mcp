package sqlite

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Open opens a SQLite database at the given path (e.g. "file:data.db" or ":memory:").
// It runs AutoMigrate for ProjectModel and ZoneModel.
func Open(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.AutoMigrate(&ProjectModel{}, &ZoneModel{}, &AgentModel{}); err != nil {
		return nil, fmt.Errorf("sqlite migrate: %w", err)
	}
	return db, nil
}
