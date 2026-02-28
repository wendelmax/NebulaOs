package api

import (
	"encoding/json"
	"net/http"

	"github.com/jacksonwendel/nebulaos/src/api/internal/usecase"
)

type TenantHandler struct {
	createUseCase *usecase.CreateTenantUseCase
	getUseCase    *usecase.GetTenantUseCase
	listUseCase   *usecase.ListTenantsUseCase
}

func NewTenantHandler(createUC *usecase.CreateTenantUseCase, getUC *usecase.GetTenantUseCase, listUC *usecase.ListTenantsUseCase) *TenantHandler {
	return &TenantHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
		listUseCase:   listUC,
	}
}

func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input usecase.CreateTenantInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.createUseCase.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Tenant created successfully"})
}

func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenants, err := h.listUseCase.Execute(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tenants)
}

func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	tenant, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tenant)
}

type ProjectHandler struct {
	createUseCase *usecase.CreateProjectUseCase
	getUseCase    *usecase.GetProjectUseCase
	listUseCase   *usecase.ListProjectsByTenantUseCase
}

func NewProjectHandler(createUC *usecase.CreateProjectUseCase, getUC *usecase.GetProjectUseCase, listUC *usecase.ListProjectsByTenantUseCase) *ProjectHandler {
	return &ProjectHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
		listUseCase:   listUC,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input usecase.CreateProjectInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.createUseCase.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Project created successfully"})
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tenantID := r.URL.Query().Get("tenant_id")
	if tenantID == "" {
		http.Error(w, "Missing tenant_id parameter", http.StatusBadRequest)
		return
	}

	projects, err := h.listUseCase.Execute(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

type ResourceHandler struct {
	createUseCase *usecase.CreateResourceUseCase
	getUseCase    *usecase.GetResourceUseCase
	listUseCase   *usecase.ListResourcesByProjectUseCase
}

func NewResourceHandler(createUC *usecase.CreateResourceUseCase, getUC *usecase.GetResourceUseCase, listUC *usecase.ListResourcesByProjectUseCase) *ResourceHandler {
	return &ResourceHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
		listUseCase:   listUC,
	}
}

func (h *ResourceHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input usecase.CreateResourceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.createUseCase.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Resource created successfully"})
}

func (h *ResourceHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		http.Error(w, "Missing project_id parameter", http.StatusBadRequest)
		return
	}

	resources, err := h.listUseCase.Execute(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resources)
}

type StorageHandler struct {
	createVolumeUC *usecase.CreateVolumeUseCase
	createBucketUC *usecase.CreateBucketUseCase
}

func NewStorageHandler(cvUC *usecase.CreateVolumeUseCase, cbUC *usecase.CreateBucketUseCase) *StorageHandler {
	return &StorageHandler{
		createVolumeUC: cvUC,
		createBucketUC: cbUC,
	}
}

func (h *StorageHandler) CreateVolume(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateVolumeInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.createVolumeUC.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Volume created successfully"})
}

func (h *StorageHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateBucketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.createBucketUC.Execute(r.Context(), input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Bucket created successfully"})
}

type ComplianceHandler struct {
	complianceUC *usecase.GetComplianceReportUseCase
}

func NewComplianceHandler(uc *usecase.GetComplianceReportUseCase) *ComplianceHandler {
	return &ComplianceHandler{complianceUC: uc}
}

func (h *ComplianceHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	tenantID := r.URL.Query().Get("tenant_id")
	report, err := h.complianceUC.Execute(r.Context(), tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}
