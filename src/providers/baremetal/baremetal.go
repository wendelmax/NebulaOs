package baremetal

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type BareMetalProvider struct {
	IPMIUser     string
	IPMIPassword string
}

func NewBareMetalProvider(user, password string) *BareMetalProvider {
	return &BareMetalProvider{
		IPMIUser:     user,
		IPMIPassword: password,
	}
}

func (p *BareMetalProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("BareMetal: Sending IPMI Power On command to node %s...\n", resource.Name)
	resource.ProviderID = fmt.Sprintf("bm-%s", resource.Name)
	resource.State = "provisioning"
	return nil
}

func (p *BareMetalProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("BareMetal: Sending IPMI Power Off command to node %s...\n", resource.ProviderID)
	resource.State = "decommissioned"
	return nil
}

func (p *BareMetalProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	return "power_on", nil
}
func (p *BareMetalProvider) Ping(ctx context.Context) error {
	return nil
}

func (p *BareMetalProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return nil, nil
}

func (p *BareMetalProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	return fmt.Errorf("baremetal does not support direct container deployment")
}

func (p *BareMetalProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	fmt.Printf("[BareMetal] Configuring switch VLAN for %s\n", n.Name)
	return nil
}

func (p *BareMetalProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	return nil
}

func (p *BareMetalProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	fmt.Printf("[BareMetal] Applying ACLs/VLAN tags for Security Group %s to switch port of node %s\n", sgID, resourceID)
	return nil
}
