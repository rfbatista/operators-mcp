package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"operators-mcp/internal/adapter/in/mcp"
	"operators-mcp/internal/application/blueprint"
	"operators-mcp/internal/domain"
)

// Handler exposes blueprint use cases as HTTP endpoints (same contract as MCP tools).
type Handler struct {
	svc *blueprint.Service
}

// NewHandler returns an HTTP handler that serves /api/list_tree, /api/list_zones, /api/list_projects, etc.
func NewHandler(svc *blueprint.Service) *Handler {
	return &Handler{svc: svc}
}

// Mount registers all API routes on mux under the given prefix (e.g. "/api").
func (h *Handler) Mount(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix+"/list_projects", h.handleListProjects)
	mux.HandleFunc(prefix+"/get_project", h.handleGetProject)
	mux.HandleFunc(prefix+"/create_project", h.handleCreateProject)
	mux.HandleFunc(prefix+"/update_project", h.handleUpdateProject)
	mux.HandleFunc(prefix+"/add_ignored_path", h.handleAddIgnoredPath)
	mux.HandleFunc(prefix+"/remove_ignored_path", h.handleRemoveIgnoredPath)
	mux.HandleFunc(prefix+"/list_tree", h.handleListTree)
	mux.HandleFunc(prefix+"/list_zones", h.handleListZones)
	mux.HandleFunc(prefix+"/list_matching_paths", h.handleListMatchingPaths)
	mux.HandleFunc(prefix+"/get_zone", h.handleGetZone)
	mux.HandleFunc(prefix+"/create_zone", h.handleCreateZone)
	mux.HandleFunc(prefix+"/update_zone", h.handleUpdateZone)
	mux.HandleFunc(prefix+"/assign_path_to_zone", h.handleAssignPathToZone)
}

func (h *Handler) handleListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	projects := h.svc.ListProjects()
	writeJSON(w, mcp.ListProjectsOut{Projects: mcp.ProjectsToDTO(projects)})
}

func (h *Handler) handleGetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.GetProjectIn
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, "invalid body", http.StatusBadRequest)
			return
		}
	} else {
		in.ProjectID = r.URL.Query().Get("project_id")
	}
	p := h.svc.GetProject(in.ProjectID)
	if p == nil {
		writeJSONError(w, "project not found", http.StatusNotFound)
		return
	}
	writeJSON(w, mcp.GetProjectOut{Project: mcp.ProjectToDTO(p)})
}

func (h *Handler) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.CreateProjectIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	p, err := h.svc.CreateProject(in.Name, in.RootDir)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.CreateProjectOut{Project: mcp.ProjectToDTO(p)})
}

func (h *Handler) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.UpdateProjectIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	p, err := h.svc.UpdateProject(in.ProjectID, in.Name, in.RootDir)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.UpdateProjectOut{Project: mcp.ProjectToDTO(p)})
}

func (h *Handler) handleAddIgnoredPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.AddIgnoredPathIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	p, err := h.svc.AddIgnoredPath(in.ProjectID, in.Path)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.AddIgnoredPathOut{Project: mcp.ProjectToDTO(p)})
}

func (h *Handler) handleRemoveIgnoredPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.RemoveIgnoredPathIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	p, err := h.svc.RemoveIgnoredPath(in.ProjectID, in.Path)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.RemoveIgnoredPathOut{Project: mcp.ProjectToDTO(p)})
}

func (h *Handler) handleListTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.ListTreeIn
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, "invalid body", http.StatusBadRequest)
			return
		}
	} else {
		in.Root = r.URL.Query().Get("root")
		in.ProjectID = r.URL.Query().Get("project_id")
	}
	tree, err := h.svc.ListTree(in.Root, in.ProjectID)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.ListTreeOut{Tree: mcp.TreeNodeToDTO(tree)})
}

func (h *Handler) handleListZones(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.ListZonesIn
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, "invalid body", http.StatusBadRequest)
			return
		}
	} else {
		in.ProjectID = r.URL.Query().Get("project_id")
	}
	if in.ProjectID == "" {
		writeJSONError(w, "project_id is required", http.StatusBadRequest)
		return
	}
	zones := h.svc.ListZones(in.ProjectID)
	writeJSON(w, mcp.ListZonesOut{Zones: mcp.ZonesToDTO(zones)})
}

func (h *Handler) handleListMatchingPaths(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.ListMatchingPathsIn
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, "invalid body", http.StatusBadRequest)
			return
		}
	} else {
		in.Pattern = r.URL.Query().Get("pattern")
		in.Root = r.URL.Query().Get("root")
		in.ProjectID = r.URL.Query().Get("project_id")
	}
	paths, err := h.svc.ListMatchingPaths(in.Root, in.ProjectID, in.Pattern)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.ListMatchingPathsOut{Paths: paths})
}

func (h *Handler) handleGetZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.GetZoneIn
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSONError(w, "invalid body", http.StatusBadRequest)
			return
		}
	} else {
		in.ZoneID = r.URL.Query().Get("zone_id")
	}
	z := h.svc.GetZone(in.ZoneID)
	if z == nil {
		writeJSONError(w, "zone not found", http.StatusNotFound)
		return
	}
	writeJSON(w, mcp.GetZoneOut{Zone: mcp.ZoneToDTO(z)})
}

func (h *Handler) handleCreateZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.CreateZoneIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	z, err := h.svc.CreateZone(in.ProjectID, in.Name, in.Pattern, in.Purpose, in.Constraints, mcp.DTOToAgents(in.AssignedAgents))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.CreateZoneOut{Zone: mcp.ZoneToDTO(z)})
}

func (h *Handler) handleUpdateZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.UpdateZoneIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	z, err := h.svc.UpdateZone(in.ZoneID, in.Name, in.Pattern, in.Purpose, in.Constraints, mcp.DTOToAgents(in.AssignedAgents))
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.UpdateZoneOut{Zone: mcp.ZoneToDTO(z)})
}

func (h *Handler) handleAssignPathToZone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in mcp.AssignPathToZoneIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSONError(w, "invalid body", http.StatusBadRequest)
		return
	}
	z, err := h.svc.AssignPathToZone(in.ZoneID, in.Path)
	if err != nil {
		writeDomainError(w, err)
		return
	}
	writeJSON(w, mcp.AssignPathToZoneOut{Zone: mcp.ZoneToDTO(z)})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func writeDomainError(w http.ResponseWriter, err error) {
	var se *domain.StructuredError
	if errors.As(err, &se) {
		switch se.Code {
		case "ZONE_NOT_FOUND", "PROJECT_NOT_FOUND":
			writeJSONError(w, se.Message, http.StatusNotFound)
			return
		case "INVALID_PATTERN", "INVALID_NAME", "INVALID_ROOT", "INVALID_PATH":
			writeJSONError(w, se.Message, http.StatusBadRequest)
			return
		}
	}
	writeJSONError(w, err.Error(), http.StatusInternalServerError)
}
