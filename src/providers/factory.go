package providers

import (
	"context"
	"fmt"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type ProviderFactory struct {
	providers map[string]domain.CloudProvider
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]domain.CloudProvider),
	}
}

func (f *ProviderFactory) Register(name string, provider domain.CloudProvider) {
	f.providers[name] = provider
}

func (f *ProviderFactory) Provision(ctx context.Context, resource *domain.Resource) error {
	provider, ok := f.providers[string(resource.Type)]
	if !ok {
		if mock, exists := f.providers["mock"]; exists {
			return mock.Provision(ctx, resource)
		}
		return fmt.Errorf("no provider found for resource type: %s", resource.Type)
	}
	return provider.Provision(ctx, resource)
}

func (f *ProviderFactory) Decommission(ctx context.Context, resource *domain.Resource) error {
	provider, ok := f.providers[string(resource.Type)]
	if !ok {
		if mock, exists := f.providers["mock"]; exists {
			return mock.Decommission(ctx, resource)
		}
		return fmt.Errorf("no provider found for resource type: %s", resource.Type)
	}
	return provider.Decommission(ctx, resource)
}

func (f *ProviderFactory) GetStatus(ctx context.Context, resourceID string) (string, error) {
	return "active", nil
}
