package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type CreateResourceInput struct {
	ID        string
	ProjectID string
	Type      domain.ResourceType
	Provider  string
	Metadata  map[string]interface{}
}

type CreateResourceUseCase struct {
	resourceRepo  domain.ResourceRepository
	projectRepo   domain.ProjectRepository
	quotaRepo     domain.QuotaRepository
	policyService domain.PolicyService
	provider      domain.CloudProvider
}

func NewCreateResourceUseCase(
	resRepo domain.ResourceRepository,
	projRepo domain.ProjectRepository,
	qRepo domain.QuotaRepository,
	pService domain.PolicyService,
	prov domain.CloudProvider,
) *CreateResourceUseCase {
	return &CreateResourceUseCase{
		resourceRepo:  resRepo,
		projectRepo:   projRepo,
		quotaRepo:     qRepo,
		policyService: pService,
		provider:      prov,
	}
}

func (uc *CreateResourceUseCase) Execute(ctx context.Context, input CreateResourceInput) error {
	// 1. Get Project -> Tenant
	project, err := uc.projectRepo.GetByID(ctx, input.ProjectID)
	if err != nil {
		return fmt.Errorf("project not found: %w", err)
	}

	// 2. Get Tenant Quota
	quota, err := uc.quotaRepo.GetByTenant(ctx, project.TenantID)
	if err != nil {
		return fmt.Errorf("failed to fetch quota: %w", err)
	}

	// 3. Calculate current usage (Simplified: count resources as units)
	existingResources, err := uc.resourceRepo.List(ctx)
	if err != nil {
		return err
	}

	// 4. Enforce quota
	totalCPUs := 0
	totalRAM := 0
	totalDisk := 0

	for _, res := range existingResources {
		if res.ProjectID == input.ProjectID {
			// In a real system, we'd extract CPU/RAM from resource metadata
			// For this demo, let's assume each resource uses 1 CPU, 2GB RAM, 20GB Disk
			totalCPUs += 1
			totalRAM += 2048
			totalDisk += 20
		}
	}

	if totalCPUs+1 > quota.MaxCPUs || totalRAM+2048 > quota.MaxRAM || totalDisk+20 > quota.MaxDisk {
		return fmt.Errorf("resource quota exceeded for project %s (Usage: %d CPUs, %dMB RAM, %dGB Disk)",
			input.ProjectID, totalCPUs, totalRAM, totalDisk)
	}

	// 5. Enforce Sovereignty Policy (NEW)
	region, ok := input.Metadata["region"].(string)
	if ok && region != "" {
		if err := uc.policyService.ValidateRegion(ctx, project.TenantID, region); err != nil {
			return err
		}
	}

	// 6. Provision
	resource := &domain.Resource{
		ID:        input.ID,
		ProjectID: input.ProjectID,
		Type:      input.Type,
		Provider:  input.Provider,
		State:     "provisioning",
		Metadata:  input.Metadata,
		CreatedAt: time.Now(),
	}

	if err := uc.resourceRepo.Create(ctx, resource); err != nil {
		return err
	}

	return uc.provider.Provision(ctx, resource)
}
