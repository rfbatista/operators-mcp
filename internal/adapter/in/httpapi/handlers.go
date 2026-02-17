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
	svc         *blueprint.Service
	defaultRoot string
}

// NewHandler returns an HTTP handler that serves /api/list_tree, /api/list_zones, etc.
func NewHandler(svc *blueprint.Service, defaultRoot string) *Handler {
	return &Handler{svc: svc, defaultRoot: defaultRoot}
}

// Mount registers all API routes on mux under the given prefix (e.g. "/api").
func (h *Handler) Mount(mux *http.ServeMux, prefix string) {
	mux.HandleFunc(prefix+"/list_tree", h.handleListTree)
	mux.HandleFunc(prefix+"/list_zones", h.handleListZones)
	mux.HandleFunc(prefix+"/list_matching_paths", h.handleListMatchingPaths)
	mux.HandleFunc(prefix+"/get_zone", h.handleGetZone)
	mux.HandleFunc(prefix+"/create_zone", h.handleCreateZone)
	mux.HandleFunc(prefix+"/update_zone", h.handleUpdateZone)
	mux.HandleFunc(prefix+"/assign_path_to_zone", h.handleAssignPathToZone)
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
	}
	root := h.defaultRoot
	if in.Root != "" {
		root = in.Root
	}
	tree, err := h.svc.ListTree(root)
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
	zones := h.svc.ListZones()
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
	}
	root := h.defaultRoot
	if in.Root != "" {
		root = in.Root
	}
	paths, err := h.svc.ListMatchingPaths(root, in.Pattern)
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
	z, err := h.svc.CreateZone(in.Name, in.Pattern, in.Purpose, in.Constraints, mcp.DTOToAgents(in.AssignedAgents))
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
		case "ZONE_NOT_FOUND":
			writeJSONError(w, se.Message, http.StatusNotFound)
			return
		case "INVALID_PATTERN":
			writeJSONError(w, se.Message, http.StatusBadRequest)
			return
		}
	}
	writeJSONError(w, err.Error(), http.StatusInternalServerError)
}
