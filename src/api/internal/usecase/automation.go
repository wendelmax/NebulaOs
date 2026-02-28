package usecase

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

// Terraform State Use Cases
type SaveTerraformStateInput struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Version   int    `json:"version"`
	State     []byte `json:"state"`
}

type SaveTerraformStateUseCase struct {
	repo domain.TerraformStateRepository
}

func NewSaveTerraformStateUseCase(repo domain.TerraformStateRepository) *SaveTerraformStateUseCase {
	return &SaveTerraformStateUseCase{repo: repo}
}

func (uc *SaveTerraformStateUseCase) Execute(ctx context.Context, input SaveTerraformStateInput) error {
	state := &domain.TerraformState{
		ID:        input.ID,
		ProjectID: input.ProjectID,
		Version:   input.Version,
		State:     input.State,
		UpdatedAt: time.Now(),
	}
	return uc.repo.Save(ctx, state)
}

type GetTerraformStateUseCase struct {
	repo domain.TerraformStateRepository
}

func NewGetTerraformStateUseCase(repo domain.TerraformStateRepository) *GetTerraformStateUseCase {
	return &GetTerraformStateUseCase{repo: repo}
}

func (uc *GetTerraformStateUseCase) Execute(ctx context.Context, projectID string) (*domain.TerraformState, error) {
	return uc.repo.GetByProjectID(ctx, projectID)
}

// Blueprint Use Cases
type ListBlueprintsUseCase struct {
	repo domain.BlueprintRepository
}

func NewListBlueprintsUseCase(repo domain.BlueprintRepository) *ListBlueprintsUseCase {
	return &ListBlueprintsUseCase{repo: repo}
}

func (uc *ListBlueprintsUseCase) Execute(ctx context.Context) ([]*domain.Blueprint, error) {
	return uc.repo.List(ctx)
}

type CreateBlueprintUseCase struct {
	repo domain.BlueprintRepository
}

func NewCreateBlueprintUseCase(repo domain.BlueprintRepository) *CreateBlueprintUseCase {
	return &CreateBlueprintUseCase{repo: repo}
}

func (uc *CreateBlueprintUseCase) Execute(ctx context.Context, b *domain.Blueprint) error {
	b.CreatedAt = time.Now()
	return uc.repo.Create(ctx, b)
}

type DeployBlueprintInput struct {
	BlueprintID string                 `json:"blueprint_id"`
	ProjectID   string                 `json:"project_id"`
	Variables   map[string]interface{} `json:"variables"`
}

type DeployBlueprintUseCase struct {
	blueprintRepo domain.BlueprintRepository
	resourceRepo  domain.ResourceRepository
	factory       domain.ProviderFactory
}

func NewDeployBlueprintUseCase(bpRepo domain.BlueprintRepository, resRepo domain.ResourceRepository, factory domain.ProviderFactory) *DeployBlueprintUseCase {
	return &DeployBlueprintUseCase{
		blueprintRepo: bpRepo,
		resourceRepo:  resRepo,
		factory:       factory,
	}
}

func (uc *DeployBlueprintUseCase) Execute(ctx context.Context, input DeployBlueprintInput) error {
	bp, err := uc.blueprintRepo.GetByID(ctx, input.BlueprintID)
	if err != nil {
		return err
	}

	// Orchestrate resource creation
	for _, resDef := range bp.Resources {
		provider, err := uc.factory.GetProvider(resDef.Provider)
		if err != nil {
			return err
		}

		res := &domain.Resource{
			ID:          domain.NewID(),
			ProjectID:   input.ProjectID,
			Name:        resDef.Name,
			Type:        resDef.Type,
			Provider:    resDef.Provider,
			State:       "provisioning",
			Metadata:    resDef.Metadata,
			BlueprintID: bp.ID,
			CreatedAt:   time.Now(),
		}

		if err := provider.Provision(ctx, res); err != nil {
			return err
		}

		if err := uc.resourceRepo.Create(ctx, res); err != nil {
			return err
		}
	}

	return nil
}
