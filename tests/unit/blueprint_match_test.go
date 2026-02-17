package unit

import (
	"os"
	"path/filepath"
	"testing"

	"operators-mcp/internal/blueprint"
)

func TestListMatchingPaths_ValidPattern_ReturnsPaths(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "cmd", "server"), 0755)
	_ = os.MkdirAll(filepath.Join(root, "internal", "mcp"), 0755)

	paths, err := blueprint.ListMatchingPaths(root, "cmd")
	if err != nil {
		t.Fatalf("ListMatchingPaths: %v", err)
	}
	if len(paths) == 0 {
		t.Fatal("expected at least one path")
	}
	found := false
	for _, p := range paths {
		if p == "cmd" || p == "cmd/server" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected cmd or cmd/server in %v", paths)
	}
}

func TestListMatchingPaths_InvalidPattern_StructuredError(t *testing.T) {
	root := t.TempDir()
	_, err := blueprint.ListMatchingPaths(root, "[")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
	se, ok := err.(*blueprint.StructuredError)
	if !ok {
		t.Fatalf("expected StructuredError, got %T", err)
	}
	if se.Code != "INVALID_PATTERN" {
		t.Errorf("code: got %q", se.Code)
	}
}

func TestListMatchingPaths_NonexistentRoot_StructuredError(t *testing.T) {
	_, err := blueprint.ListMatchingPaths("/nonexistent/path/12345", ".")
	if err == nil {
		t.Fatal("expected error")
	}
	se, ok := err.(*blueprint.StructuredError)
	if !ok || se.Code != "ROOT_UNREADABLE" {
		t.Errorf("expected ROOT_UNREADABLE, got %v", err)
	}
}
