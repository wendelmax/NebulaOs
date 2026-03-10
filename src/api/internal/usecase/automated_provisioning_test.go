package usecase_test

import (
	"context"
	"testing"

	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/infrastructure"
	"github.com/wendelmax/nebulaos/src/api/internal/usecase"
)

func TestAutomatedProvisioningUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("Success_K8sPreset", func(t *testing.T) {
		bpRepo := infrastructure.NewInMemoryBlueprintRepository()
		resRepo := infrastructure.NewInMemoryResourceRepository()
		factory := &mockProviderFactory{provider: &mockCloudProvider{}}
		
		deployUC := usecase.NewDeployBlueprintUseCase(bpRepo, resRepo, factory)
		uc := usecase.NewAutomatedProvisioningUseCase(deployUC, bpRepo)

		// Create the blueprint that the preset resolves to
		bpRepo.Create(ctx, &domain.Blueprint{
			ID: "bp-k8s-ha",
			Resources: []domain.ResourceDefinition{
				{Name: "k8s-master", Type: domain.ComputeResource, Provider: "proxmox"},
			},
		})

		req := domain.AutomaticProvisioningRequest{
			PresetID:  "k8s-standard",
			ProjectID: "proj-1",
		}

		err := uc.Execute(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Verify resource creation
		resources, _ := resRepo.GetByProject(ctx, "proj-1")
		found := false
		for _, res := range resources {
			if res.BlueprintID == "bp-k8s-ha" {
				found = true
				break
			}
		}
		if !found {
			t.Error("expected resources to be deployed from bp-k8s-ha")
		}
	})

	t.Run("UnknownPreset", func(t *testing.T) {
		bpRepo := infrastructure.NewInMemoryBlueprintRepository()
		resRepo := infrastructure.NewInMemoryResourceRepository()
		factory := &mockProviderFactory{provider: &mockCloudProvider{}}
		deployUC := usecase.NewDeployBlueprintUseCase(bpRepo, resRepo, factory)
		uc := usecase.NewAutomatedProvisioningUseCase(deployUC, bpRepo)

		req := domain.AutomaticProvisioningRequest{
			PresetID: "invalid-preset",
		}

		err := uc.Execute(ctx, req)
		if err == nil {
			t.Fatal("expected error for unknown preset, got nil")
		}
	})
}
