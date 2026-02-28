package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type GetResourceUseCase struct {
	repo domain.ResourceRepository
}

func NewGetResourceUseCase(repo domain.ResourceRepository) *GetResourceUseCase {
	return &GetResourceUseCase{repo: repo}
}

func (uc *GetResourceUseCase) Execute(ctx context.Context, id string) (*domain.Resource, error) {
	return uc.repo.GetByID(ctx, id)
}

type ListResourcesByProjectUseCase struct {
	repo domain.ResourceRepository
}

func NewListResourcesByProjectUseCase(repo domain.ResourceRepository) *ListResourcesByProjectUseCase {
	return &ListResourcesByProjectUseCase{repo: repo}
}

func (uc *ListResourcesByProjectUseCase) Execute(ctx context.Context, projectID string) ([]*domain.Resource, error) {
	return uc.repo.GetByProject(ctx, projectID)
}

type ListAllResourcesUseCase struct {
	repo domain.ResourceRepository
}

func NewListAllResourcesUseCase(repo domain.ResourceRepository) *ListAllResourcesUseCase {
	return &ListAllResourcesUseCase{repo: repo}
}

func (uc *ListAllResourcesUseCase) Execute(ctx context.Context) ([]*domain.Resource, error) {
	return uc.repo.List(ctx)
}
