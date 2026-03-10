package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/infrastructure"
	"github.com/wendelmax/nebulaos/src/api/internal/services"
	"github.com/wendelmax/nebulaos/src/api/internal/usecase"
)

type mockCloudProvider struct{}

func (p *mockCloudProvider) HealthCheck(ctx context.Context) error { return nil }
func (p *mockCloudProvider) Ping(ctx context.Context) error        { return nil }
func (p *mockCloudProvider) Provision(ctx context.Context, r *domain.Resource) error {
	return nil
}
func (p *mockCloudProvider) Decommission(ctx context.Context, r *domain.Resource) error {
	return nil
}
func (p *mockCloudProvider) GetStatus(ctx context.Context, id string) (string, error) {
	return "ok", nil
}
func (p *mockCloudProvider) AttachSecurityGroup(ctx context.Context, rid, sgid string) error {
	return nil
}
func (p *mockCloudProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return nil, nil
}
func (p *mockCloudProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	return nil
}
func (p *mockCloudProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	return nil
}
func (p *mockCloudProvider) AddRoute(ctx context.Context, route *domain.Route) error {
	return nil
}

type mockProviderFactory struct {
	provider domain.CloudProvider
}

func (f *mockProviderFactory) GetProvider(name string) (domain.CloudProvider, error) {
	return f.provider, nil
}

// Helper to setup a test router similar to main.go
func setupTestRouter() http.Handler {
	mux := http.NewServeMux()
	
	// Repos
	tenantRepo := infrastructure.NewInMemoryTenantRepository()
	projectRepo := infrastructure.NewInMemoryProjectRepository()
	resourceRepo := infrastructure.NewInMemoryResourceRepository()
	blueprintRepo := infrastructure.NewInMemoryBlueprintRepository()
	regionRepo := infrastructure.NewInMemoryRegionRepository()
	
	// Providers
	factory := &mockProviderFactory{provider: &mockCloudProvider{}}
	
	// Services
	networkMgr := services.NewNetworkManager(factory)
	deployUC := usecase.NewDeployBlueprintUseCase(blueprintRepo, resourceRepo, factory)
	autoProvisionUC := usecase.NewAutomatedProvisioningUseCase(deployUC, blueprintRepo)

	// Sample data
	blueprintRepo.Create(nil, &domain.Blueprint{ID: "bp-k8s-ha"})

	// Routes
	mux.Handle("/infra/automated/provision", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var req domain.AutomaticProvisioningRequest
			json.NewDecoder(r.Body).Decode(&req)
			autoProvisionUC.Execute(r.Context(), req)
			w.WriteHeader(http.StatusAccepted)
		} else {
			presets, _ := autoProvisionUC.ListPresets(r.Context())
			json.NewEncoder(w).Encode(presets)
		}
	}))

	mux.Handle("/cloud/regions", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		regs, _ := regionRepo.List(r.Context())
		json.NewEncoder(w).Encode(regs)
	}))

	_ = tenantRepo
	_ = projectRepo
	_ = networkMgr // keep to avoid unused variable

	return mux
}

func TestAPI_E2E(t *testing.T) {
	router := setupTestRouter()

	t.Run("GetPresets", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/infra/automated/provision", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}

		var presets []domain.InfraPreset
		json.Unmarshal(rr.Body.Bytes(), &presets)
		if len(presets) == 0 {
			t.Error("expected at least one preset")
		}
	})

	t.Run("TriggerProvision", func(t *testing.T) {
		body, _ := json.Marshal(domain.AutomaticProvisioningRequest{
			PresetID:  "k8s-standard",
			ProjectID: "proj-1",
		})
		req := httptest.NewRequest("POST", "/infra/automated/provision", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusAccepted {
			t.Errorf("expected status 202, got %d", rr.Code)
		}
	})
}
