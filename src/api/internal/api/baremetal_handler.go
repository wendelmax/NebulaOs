package api

import (
	"encoding/json"
	"net/http"

	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/services"
)

type BareMetalHandler struct {
	manager *services.BareMetalManager
}

func NewBareMetalHandler(manager *services.BareMetalManager) *BareMetalHandler {
	return &BareMetalHandler{manager: manager}
}

func (h *BareMetalHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.manager.ListNodes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(nodes)
}

func (h *BareMetalHandler) RegisterNode(w http.ResponseWriter, r *http.Request) {
	var node domain.BareMetalNode
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.manager.RegisterNode(r.Context(), &node); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

func (h *BareMetalHandler) ProvisionNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.URL.Query().Get("id")
	if nodeID == "" {
		http.Error(w, "missing node id", http.StatusBadRequest)
		return
	}
	if err := h.manager.ProvisionNode(r.Context(), nodeID, "default"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Provisioning started"})
}

func (h *BareMetalHandler) GetNodeLogs(w http.ResponseWriter, r *http.Request) {
	nodeID := r.URL.Query().Get("id")
	_, logs, err := h.manager.GetNodeStatus(r.Context(), nodeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(logs)
}
