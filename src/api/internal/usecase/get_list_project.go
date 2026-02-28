package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type GetProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewGetProjectUseCase(repo domain.ProjectRepository) *GetProjectUseCase {
	return &GetProjectUseCase{repo: repo}
}

func (uc *GetProjectUseCase) Execute(ctx context.Context, id string) (*domain.Project, error) {
	return uc.repo.GetByID(ctx, id)
}

type ListProjectsByTenantUseCase struct {
	repo domain.ProjectRepository
}

func NewListProjectsByTenantUseCase(repo domain.ProjectRepository) *ListProjectsByTenantUseCase {
	return &ListProjectsByTenantUseCase{repo: repo}
}

func (uc *ListProjectsByTenantUseCase) Execute(ctx context.Context, tenantID string) ([]*domain.Project, error) {
	return uc.repo.GetByTenant(ctx, tenantID)
}
