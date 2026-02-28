package mock

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type MockProvider struct {
	resources map[string]string // ID -> State
}

func NewMockProvider() *MockProvider {
	return &MockProvider{
		resources: make(map[string]string),
	}
}

func (p *MockProvider) Provision(ctx context.Context, res *domain.Resource) error {
	log.Printf("[MockProvider] Provisioning resource %s (%s) on mock infrastructure...", res.ID, res.Type)
	p.resources[res.ID] = "provisioning"

	// Simulate async provisioning
	go func() {
		time.Sleep(2 * time.Second)
		p.resources[res.ID] = "active"
		log.Printf("[MockProvider] Resource %s is now active", res.ID)
	}()

	return nil
}

func (p *MockProvider) Decommission(ctx context.Context, res *domain.Resource) error {
	log.Printf("[MockProvider] Decommissioning resource %s...", res.ID)
	delete(p.resources, res.ID)
	return nil
}

func (p *MockProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	status, ok := p.resources[resourceID]
	if !ok {
		return "", fmt.Errorf("resource not found in mock provider")
	}
	return status, nil
}
