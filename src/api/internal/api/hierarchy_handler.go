package api

import (
	"encoding/json"
	"net/http"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type HierarchyHandler struct {
	orgRepo  domain.OrganizationRepository
	deptRepo domain.DepartmentRepository
}

func NewHierarchyHandler(orgRepo domain.OrganizationRepository, deptRepo domain.DepartmentRepository) *HierarchyHandler {
	return &HierarchyHandler{orgRepo: orgRepo, deptRepo: deptRepo}
}

func (h *HierarchyHandler) ListOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs, err := h.orgRepo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orgs)
}

func (h *HierarchyHandler) ListDepartments(w http.ResponseWriter, r *http.Request) {
	orgID := r.URL.Query().Get("org_id")
	if orgID != "" {
		depts, err := h.deptRepo.GetByOrganization(r.Context(), orgID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(depts)
		return
	}
	// Fallback to all (if needed, or error)
	http.Error(w, "missing org_id", http.StatusBadRequest)
}

func (h *HierarchyHandler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	var org domain.Organization
	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	org.ID = domain.NewID()
	if err := h.orgRepo.Create(r.Context(), &org); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(org)
}
