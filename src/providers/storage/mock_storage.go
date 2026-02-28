package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type MockStorageProvider struct{}

func NewMockStorageProvider() *MockStorageProvider {
	return &MockStorageProvider{}
}

func (p *MockStorageProvider) CreateVolume(ctx context.Context, vol *domain.Volume) error {
	log.Printf("[MockStorage] Creating volume %s (Size: %d GB)", vol.Name, vol.SizeGB)
	vol.ProviderID = fmt.Sprintf("vol-mock-%s", vol.ID)
	vol.State = "available"
	return nil
}

func (p *MockStorageProvider) DeleteVolume(ctx context.Context, vol *domain.Volume) error {
	log.Printf("[MockStorage] Deleting volume %s", vol.ProviderID)
	vol.State = "deleted"
	return nil
}

func (p *MockStorageProvider) CreateBucket(ctx context.Context, b *domain.Bucket) error {
	log.Printf("[MockStorage] Creating bucket %s in region %s", b.Name, b.Region)
	b.State = "ready"
	return nil
}

func (p *MockStorageProvider) DeleteBucket(ctx context.Context, b *domain.Bucket) error {
	log.Printf("[MockStorage] Deleting bucket %s", b.Name)
	b.State = "deleted"
	return nil
}

var _ domain.StorageProvider = (*MockStorageProvider)(nil)
