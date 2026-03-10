package providers

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
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

func (f *ProviderFactory) ListImages(ctx context.Context) ([]domain.Image, error) {
	var allImages []domain.Image
	for _, p := range f.providers {
		imgs, _ := p.ListImages(ctx)
		allImages = append(allImages, imgs...)
	}
	return allImages, nil
}

func (f *ProviderFactory) DeployContainer(ctx context.Context, container *domain.Container) error {
	// Simple proxy based on some logic or default to proxmox for containers
	if p, exists := f.providers["proxmox"]; exists {
		return p.DeployContainer(ctx, container)
	}
	return fmt.Errorf("no container provider found")
}

func (f *ProviderFactory) ConfigureNetwork(ctx context.Context, network *domain.Network) error {
	if p, ok := f.providers[network.Provider]; ok {
		return p.ConfigureNetwork(ctx, network)
	}
	return fmt.Errorf("provider %s not found for network configuration", network.Provider)
}

func (f *ProviderFactory) AddRoute(ctx context.Context, route *domain.Route) error {
	// Proxy to provider factory logic or specific provider
	return fmt.Errorf("use NetworkManager for cross-provider routing")
}

func (f *ProviderFactory) Ping(ctx context.Context) error {
	return nil
}
