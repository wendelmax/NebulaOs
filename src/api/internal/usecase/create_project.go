package usecase

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type CreateProjectInput struct {
	ID       string
	TenantID string
	Name     string
}

type CreateProjectUseCase struct {
	repo domain.ProjectRepository
}

func NewCreateProjectUseCase(repo domain.ProjectRepository) *CreateProjectUseCase {
	return &CreateProjectUseCase{repo: repo}
}

func (uc *CreateProjectUseCase) Execute(ctx context.Context, input CreateProjectInput) error {
	project := &domain.Project{
		ID:        input.ID,
		TenantID:  input.TenantID,
		Name:      input.Name,
		CreatedAt: time.Now(),
	}
	return uc.repo.Create(ctx, project)
}
