package usecase

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type CreateVolumeInput struct {
	ID        string
	ProjectID string
	Name      string
	SizeGB    int
}

type CreateVolumeUseCase struct {
	repo     domain.VolumeRepository
	provider domain.StorageProvider
}

func NewCreateVolumeUseCase(repo domain.VolumeRepository, provider domain.StorageProvider) *CreateVolumeUseCase {
	return &CreateVolumeUseCase{repo: repo, provider: provider}
}

func (uc *CreateVolumeUseCase) Execute(ctx context.Context, input CreateVolumeInput) error {
	vol := &domain.Volume{
		ID:        input.ID,
		ProjectID: input.ProjectID,
		Name:      input.Name,
		SizeGB:    input.SizeGB,
		State:     "creating",
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, vol); err != nil {
		return err
	}

	return uc.provider.CreateVolume(ctx, vol)
}

type CreateBucketInput struct {
	ID        string
	ProjectID string
	Name      string
	Region    string
}

type CreateBucketUseCase struct {
	repo     domain.BucketRepository
	provider domain.StorageProvider
}

func NewCreateBucketUseCase(repo domain.BucketRepository, provider domain.StorageProvider) *CreateBucketUseCase {
	return &CreateBucketUseCase{repo: repo, provider: provider}
}

func (uc *CreateBucketUseCase) Execute(ctx context.Context, input CreateBucketInput) error {
	b := &domain.Bucket{
		ID:        input.ID,
		ProjectID: input.ProjectID,
		Name:      input.Name,
		Region:    input.Region,
		State:     "creating",
		CreatedAt: time.Now(),
	}

	if err := uc.repo.Create(ctx, b); err != nil {
		return err
	}

	return uc.provider.CreateBucket(ctx, b)
}
