package sqlite

import (
	"database/sql/driver"
	"encoding/json"

	"operators-mcp/internal/domain"
)

// stringSlice is a []string that scans from JSON or a single plain string.
// This avoids "invalid character '.' looking for beginning of value" when
// the column contains a path like ".git" or empty/invalid JSON.
type stringSlice []string

// Scan implements sql.Scanner. Accepts nil, empty string, valid JSON array,
// or a single plain string (e.g. ".git") and normalizes to []string.
func (s *stringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		*s = nil
		return nil
	}
	if len(b) == 0 {
		*s = []string{}
		return nil
	}
	// Trim space; empty after trim -> empty slice
	if len(b) > 0 && (b[0] == '[' || b[0] == '"' || b[0] == '{') {
		// Looks like JSON: unmarshal
		var out []string
		if err := json.Unmarshal(b, &out); err != nil {
			// Invalid JSON: treat whole value as single path
			*s = []string{string(b)}
			return nil
		}
		*s = out
		return nil
	}
	// Plain string (e.g. ".git" or "node_modules")
	*s = []string{string(b)}
	return nil
}

// Value implements driver.Valuer. Always returns valid JSON.
func (s stringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	return json.Marshal(s)
}

// agentSlice is a []domain.Agent that scans from JSON or invalid/empty content.
// Uses the same defensive logic as stringSlice so malformed DB values don't break reads.
type agentSlice []domain.Agent

// Scan implements sql.Scanner.
func (a *agentSlice) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		*a = nil
		return nil
	}
	if len(b) == 0 {
		*a = []domain.Agent{}
		return nil
	}
	if b[0] != '[' {
		*a = []domain.Agent{}
		return nil
	}
	var out []domain.Agent
	if err := json.Unmarshal(b, &out); err != nil {
		*a = []domain.Agent{}
		return nil
	}
	*a = out
	return nil
}

// Value implements driver.Valuer.
func (a agentSlice) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}
