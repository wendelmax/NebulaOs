package usecase

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type CreateTenantInput struct {
	ID   string
	Name string
}

type CreateTenantUseCase struct {
	repo domain.TenantRepository
}

func NewCreateTenantUseCase(repo domain.TenantRepository) *CreateTenantUseCase {
	return &CreateTenantUseCase{repo: repo}
}

func (uc *CreateTenantUseCase) Execute(ctx context.Context, input CreateTenantInput) error {
	tenant := &domain.Tenant{
		ID:        input.ID,
		Name:      input.Name,
		CreatedAt: time.Now(),
	}
	return uc.repo.Create(ctx, tenant)
}
