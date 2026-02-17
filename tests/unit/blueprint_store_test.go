package unit

import (
	"testing"

	"operators-mcp/internal/blueprint"
)

func TestStore_CreateListGetUpdateAssignPath(t *testing.T) {
	s := blueprint.NewStore()

	z, err := s.Create("backend", "cmd/.*", "Server code", []string{"no UI"}, "agent-1")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if z.ID == "" {
		t.Error("expected id")
	}
	if z.Name != "backend" {
		t.Errorf("name: got %q", z.Name)
	}

	list := s.List()
	if len(list) != 1 {
		t.Fatalf("List: got %d zones", len(list))
	}
	if list[0].Name != "backend" {
		t.Errorf("List[0].Name: got %q", list[0].Name)
	}

	got := s.Get(z.ID)
	if got == nil {
		t.Fatal("Get: nil")
	}
	if got.Name != "backend" {
		t.Errorf("Get.Name: got %q", got.Name)
	}

	updated, err := s.Update(z.ID, "backend-updated", "internal/.*", "Internal pkgs", nil, "agent-2")
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "backend-updated" {
		t.Errorf("Update name: got %q", updated.Name)
	}

	assigned, err := s.AssignPath(z.ID, "internal/blueprint")
	if err != nil {
		t.Fatalf("AssignPath: %v", err)
	}
	if len(assigned.ExplicitPaths) != 1 || assigned.ExplicitPaths[0] != "internal/blueprint" {
		t.Errorf("ExplicitPaths: got %v", assigned.ExplicitPaths)
	}
}

func TestStore_CreateEmptyName_Error(t *testing.T) {
	s := blueprint.NewStore()
	_, err := s.Create("", "x", "", nil, "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if se, ok := err.(*blueprint.StructuredError); !ok || se.Code != "INVALID_NAME" {
		t.Errorf("expected INVALID_NAME, got %v", err)
	}
}

func TestStore_GetNotFound_Nil(t *testing.T) {
	s := blueprint.NewStore()
	if s.Get("nonexistent") != nil {
		t.Error("Get nonexistent should return nil")
	}
}

func TestStore_UpdateNotFound_Error(t *testing.T) {
	s := blueprint.NewStore()
	_, err := s.Update("nonexistent", "x", "", "", nil, "")
	if err == nil {
		t.Fatal("expected error")
	}
	if se, ok := err.(*blueprint.StructuredError); !ok || se.Code != "ZONE_NOT_FOUND" {
		t.Errorf("expected ZONE_NOT_FOUND, got %v", err)
	}
}
