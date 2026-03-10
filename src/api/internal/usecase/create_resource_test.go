package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/infrastructure"
	"github.com/wendelmax/nebulaos/src/api/internal/usecase"
)

type mockPolicyService struct{}

func (s *mockPolicyService) ValidateRegion(ctx context.Context, tenantID string, region string) error {
	return nil
}
func (s *mockPolicyService) GetPolicy(ctx context.Context, tenantID string) (*domain.SovereigntyPolicy, error) {
	return nil, nil
}
func (s *mockPolicyService) UpdatePolicy(ctx context.Context, policy *domain.SovereigntyPolicy) error {
	return nil
}

type mockCloudProvider struct{}

func (p *mockCloudProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	return nil
}
func (p *mockCloudProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	return nil
}
func (p *mockCloudProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	return "running", nil
}
func (p *mockCloudProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	return nil
}

func TestCreateResourceUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	tenantID := "tenant-1"
	projectID := "proj-1"

	t.Run("Success", func(t *testing.T) {
		resRepo := infrastructure.NewInMemoryResourceRepository()
		projRepo := infrastructure.NewInMemoryProjectRepository()
		quotaRepo := infrastructure.NewInMemoryQuotaRepository()
		policyService := &mockPolicyService{}
		provider := &mockCloudProvider{}

		uc := usecase.NewCreateResourceUseCase(resRepo, projRepo, quotaRepo, policyService, provider)
		projRepo.Create(ctx, &domain.Project{ID: projectID, TenantID: tenantID, Name: "Project 1", CreatedAt: time.Now()})

		input := usecase.CreateResourceInput{
			ID:        "res-1",
			ProjectID: projectID,
			Type:      domain.ComputeResource,
			Provider:  "aws",
			Metadata:  map[string]interface{}{"region": "us-east-1"},
		}

		err := uc.Execute(ctx, input)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		res, err := resRepo.GetByID(ctx, "res-1")
		if err != nil {
			t.Fatalf("expected resource to be created, got %v", err)
		}
		if res.State != "provisioning" {
			t.Errorf("expected state to be provisioning, got %s", res.State)
		}
	})

	t.Run("Quota Exceeded", func(t *testing.T) {
		resRepo := infrastructure.NewInMemoryResourceRepository()
		projRepo := infrastructure.NewInMemoryProjectRepository()
		quotaRepo := infrastructure.NewInMemoryQuotaRepository()
		policyService := &mockPolicyService{}
		provider := &mockCloudProvider{}

		uc := usecase.NewCreateResourceUseCase(resRepo, projRepo, quotaRepo, policyService, provider)
		projRepo.Create(ctx, &domain.Project{ID: projectID, TenantID: tenantID, Name: "Project 1", CreatedAt: time.Now()})

		// Set a low quota
		quotaRepo.Update(ctx, &domain.Quota{
			TenantID: tenantID,
			MaxCPUs:  1,
			MaxRAM:   2048,
			MaxDisk:  20,
		})

		// First resource should pass
		input1 := usecase.CreateResourceInput{
			ID:        "res-2",
			ProjectID: projectID,
			Type:      domain.ComputeResource,
			Provider:  "aws",
			Metadata:  map[string]interface{}{"region": "us-east-1"},
		}
		if err := uc.Execute(ctx, input1); err != nil {
			t.Fatalf("expected no error for first resource, got %v", err)
		}

		// Second resource should fail quota (each resource uses 1 CPU, 2048MB RAM, 20GB Disk in simplified calculation)
		input2 := usecase.CreateResourceInput{
			ID:        "res-3",
			ProjectID: projectID,
			Type:      domain.ComputeResource,
			Provider:  "aws",
			Metadata:  map[string]interface{}{"region": "us-east-1"},
		}
		err := uc.Execute(ctx, input2)
		if err == nil {
			t.Fatal("expected quota exceeded error, got nil")
		}
	})
}
