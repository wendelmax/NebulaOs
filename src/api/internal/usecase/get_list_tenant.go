package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type GetTenantUseCase struct {
	repo domain.TenantRepository
}

func NewGetTenantUseCase(repo domain.TenantRepository) *GetTenantUseCase {
	return &GetTenantUseCase{repo: repo}
}

func (uc *GetTenantUseCase) Execute(ctx context.Context, id string) (*domain.Tenant, error) {
	return uc.repo.GetByID(ctx, id)
}

type ListTenantsUseCase struct {
	repo domain.TenantRepository
}

func NewListTenantsUseCase(repo domain.TenantRepository) *ListTenantsUseCase {
	return &ListTenantsUseCase{repo: repo}
}

func (uc *ListTenantsUseCase) Execute(ctx context.Context) ([]*domain.Tenant, error) {
	return uc.repo.List(ctx)
}
