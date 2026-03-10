package usecase

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type AutomatedProvisioningUseCase struct {
	deployBlueprintUC *DeployBlueprintUseCase
	blueprintRepo     domain.BlueprintRepository
}

func NewAutomatedProvisioningUseCase(deployUC *DeployBlueprintUseCase, bpRepo domain.BlueprintRepository) *AutomatedProvisioningUseCase {
	return &AutomatedProvisioningUseCase{
		deployBlueprintUC: deployUC,
		blueprintRepo:     bpRepo,
	}
}

func (uc *AutomatedProvisioningUseCase) Execute(ctx context.Context, req domain.AutomaticProvisioningRequest) error {
	fmt.Printf("[AutomatedProvisioning] Resolving preset %s for project %s\n", req.PresetID, req.ProjectID)

	// 1. Resolve preset to Blueprint ID
	// For this version, we map preset IDs directly or use a lookup
	blueprintID := ""
	switch req.PresetID {
	case "k8s-standard":
		blueprintID = "bp-k8s-ha"
	case "openstack-basic":
		blueprintID = "bp-openstack-aio"
	default:
		return fmt.Errorf("unknown preset ID: %s", req.PresetID)
	}

	// 2. Delegate to DeployBlueprintUseCase
	deployInput := DeployBlueprintInput{
		BlueprintID: blueprintID,
		ProjectID:   req.ProjectID,
		Variables:   req.Variables,
	}

	fmt.Printf("[AutomatedProvisioning] Triggering deployment of blueprint %s\n", blueprintID)
	return uc.deployBlueprintUC.Execute(ctx, deployInput)
}

func (uc *AutomatedProvisioningUseCase) ListPresets(ctx context.Context) ([]domain.InfraPreset, error) {
	return []domain.InfraPreset{
		{
			ID:          "k8s-standard",
			Name:        "Kubernetes Standard Cluster",
			Description: "HA Control Plane + 3 Workers on Proxmox/BareMetal",
			BlueprintID: "bp-k8s-ha",
		},
		{
			ID:          "openstack-basic",
			Name:        "OpenStack Private Cloud",
			Description: "All-in-one OpenStack deployment for local labs",
			BlueprintID: "bp-openstack-aio",
		},
	}, nil
}
