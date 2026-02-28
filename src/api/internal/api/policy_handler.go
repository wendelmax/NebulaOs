package api

import (
	"encoding/json"
	"net/http"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type PolicyHandler struct {
	policyService domain.PolicyService
}

func NewPolicyHandler(service domain.PolicyService) *PolicyHandler {
	return &PolicyHandler{policyService: service}
}

func (h *PolicyHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	var policy domain.SovereigntyPolicy
	if err := json.NewDecoder(r.Body).Decode(&policy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.policyService.UpdatePolicy(r.Context(), &policy); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "policy updated", "sovereignty": "enforced"})
}

func (h *PolicyHandler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	policy, err := h.policyService.GetPolicy(r.Context(), tenantID)
	if err != nil {
		http.Error(w, "policy not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(policy)
}
