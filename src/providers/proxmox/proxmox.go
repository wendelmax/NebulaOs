package proxmox

import (
	"context"
	"fmt"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type ProxmoxProvider struct {
	BaseURL string
	Token   string
}

func NewProxmoxProvider(baseURL, token string) *ProxmoxProvider {
	return &ProxmoxProvider{
		BaseURL: baseURL,
		Token:   token,
	}
}

func (p *ProxmoxProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("Proxmox: Provisioning VM %s on NebulaOS hypervisor...\n", resource.Name)
	resource.ProviderID = fmt.Sprintf("pve-%s", resource.Name)
	resource.State = "provisioning"
	return nil
}

func (p *ProxmoxProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("Proxmox: Deleting VM %s...\n", resource.ProviderID)
	resource.State = "deleted"
	return nil
}

func (p *ProxmoxProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	return "running", nil
}
func (p *ProxmoxProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	fmt.Printf("[Proxmox] Attaching Security Group %s to VM %s\n", sgID, resourceID)
	return nil
}
