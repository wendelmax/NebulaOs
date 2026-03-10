package services_test

import (
	"context"
	"testing"

	"github.com/wendelmax/nebulaos/src/api/domain"
	"github.com/wendelmax/nebulaos/src/api/internal/services"
)

type mockCloudProvider struct {
	lastRoute *domain.Route
}

func (p *mockCloudProvider) HealthCheck(ctx context.Context) error { return nil }
func (p *mockCloudProvider) Ping(ctx context.Context) error        { return nil }
func (p *mockCloudProvider) Provision(ctx context.Context, r *domain.Resource) error { return nil }
func (p *mockCloudProvider) Decommission(ctx context.Context, r *domain.Resource) error { return nil }
func (p *mockCloudProvider) GetStatus(ctx context.Context, id string) (string, error) { return "ok", nil }
func (p *mockCloudProvider) AttachSecurityGroup(ctx context.Context, rid, sgid string) error { return nil }
func (p *mockCloudProvider) ListImages(ctx context.Context) ([]domain.Image, error) { return nil, nil }
func (p *mockCloudProvider) DeployContainer(ctx context.Context, c *domain.Container) error { return nil }
func (p *mockCloudProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error { return nil }
func (p *mockCloudProvider) AddRoute(ctx context.Context, route *domain.Route) error {
	p.lastRoute = route
	return nil
}

type mockProviderFactory struct {
	provider *mockCloudProvider
}

func (f *mockProviderFactory) GetProvider(name string) (domain.CloudProvider, error) {
	return f.provider, nil
}

func TestNetworkManager_CreateInterProviderRoute(t *testing.T) {
	provider := &mockCloudProvider{}
	factory := &mockProviderFactory{provider: provider}
	mgr := services.NewNetworkManager(factory)

	ctx := context.Background()
	route := &domain.Route{
		Destination: "10.0.0.0/24",
		NextHop:     "192.168.1.1",
	}

	err := mgr.CreateInterProviderRoute(ctx, route, "proxmox", "openstack")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.lastRoute == nil {
		t.Fatal("expected AddRoute to be called on provider")
	}

	if provider.lastRoute.TargetProvider != "openstack" {
		t.Errorf("expected TargetProvider to be openstack, got %s", provider.lastRoute.TargetProvider)
	}
}
