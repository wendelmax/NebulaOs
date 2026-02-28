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

func (f *ProviderFactory) GetProvider(name string) (domain.CloudProvider, error) {
	provider, ok := f.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not found", name)
	}
	return provider, nil
}

func (f *ProviderFactory) Provision(ctx context.Context, resource *domain.Resource) error {
	provider, ok := f.providers[resource.Provider]
	if !ok {
		if mock, exists := f.providers["mock"]; exists {
			return mock.Provision(ctx, resource)
		}
		return fmt.Errorf("no provider found for: %s", resource.Provider)
	}
	return provider.Provision(ctx, resource)
}

func (f *ProviderFactory) Decommission(ctx context.Context, resource *domain.Resource) error {
	provider, ok := f.providers[resource.Provider]
	if !ok {
		if mock, exists := f.providers["mock"]; exists {
			return mock.Decommission(ctx, resource)
		}
		return fmt.Errorf("no provider found for: %s", resource.Provider)
	}
	return provider.Decommission(ctx, resource)
}

func (f *ProviderFactory) GetStatus(ctx context.Context, resourceID string) (string, error) {
	// Simple proxy or fallback
	if mock, exists := f.providers["mock"]; exists {
		return mock.GetStatus(ctx, resourceID)
	}
	return "unknown", nil
}

func (f *ProviderFactory) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	// For simplicity, proxy to mock if provider not found or if we want global mock support
	if mock, exists := f.providers["mock"]; exists {
		return mock.AttachSecurityGroup(ctx, resourceID, sgID)
	}
	return fmt.Errorf("no provider available for SG attachment")
}
