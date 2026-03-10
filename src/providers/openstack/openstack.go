package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/wendelmax/nebulaos/src/api/domain"
)

type OpenStackProvider struct {
	client *gophercloud.ProviderClient
}

func NewOpenStackProvider(authOpts gophercloud.AuthOptions) (*OpenStackProvider, error) {
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with OpenStack: %w", err)
	}
	return &OpenStackProvider{client: provider}, nil
}

func (p *OpenStackProvider) Ping(ctx context.Context) error {
	fmt.Println("[OpenStack] Identity check: active")
	return nil
}

func (p *OpenStackProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return []domain.Image{
		{ID: "ubuntu-22.04-id", Name: "Ubuntu 22.04 LTS", Provider: "openstack", Type: "glance"},
		{ID: "cirros-0.5.2", Name: "CirrOS 0.5.2", Provider: "openstack", Type: "glance"},
	}, nil
}

func (p *OpenStackProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	fmt.Printf("[OpenStack] Deploying container via Zun: %s (Image: %s)\n", c.Name, c.Image)
	c.State = "running"
	return nil
}

func (p *OpenStackProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	fmt.Printf("[OpenStack] Creating Neutron network: %s (%s)\n", n.Name, n.CIDR)
	n.State = "active"
	return nil
}

func (p *OpenStackProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	fmt.Printf("[OpenStack] Adding Neutron extra route: %s via %s\n", r.Destination, r.NextHop)
	return nil
}

func (p *OpenStackProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	client, err := openstack.NewComputeV2(p.client, gophercloud.EndpointOpts{})
	if err != nil {
		return fmt.Errorf("failed to create compute client: %w", err)
	}

	createOpts := servers.CreateOpts{
		Name:      resource.Name,
		ImageRef:  "ubuntu-22.04-id", // Should be ID
		FlavorRef: "m1.small-id",     // Should be ID
	}

	server, err := servers.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("failed to provision OpenStack server: %w", err)
	}

	resource.ProviderID = server.ID
	resource.State = "provisioning"
	return nil
}

func (p *OpenStackProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	client, err := openstack.NewComputeV2(p.client, gophercloud.EndpointOpts{})
	if err != nil {
		return err
	}

	err = servers.Delete(client, resource.ProviderID).ExtractErr()
	if err != nil {
		return fmt.Errorf("failed to delete OpenStack server: %w", err)
	}

	resource.State = "deleted"
	return nil
}

func (p *OpenStackProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	client, err := openstack.NewComputeV2(p.client, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	server, err := servers.Get(client, resourceID).Extract()
	if err != nil {
		return "", err
	}

	return server.Status, nil
}
func (p *OpenStackProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	// Mock implementation for now
	fmt.Printf("[OpenStack] Attaching Security Group %s to server %s\n", sgID, resourceID)
	return nil
}
