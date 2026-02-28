package baremetal

import (
	"context"
	"fmt"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
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
