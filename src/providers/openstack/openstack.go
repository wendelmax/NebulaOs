package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/jacksonwendel/nebulaos/src/api/domain"
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
