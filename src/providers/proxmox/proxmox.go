package proxmox

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
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

func (p *ProxmoxProvider) Ping(ctx context.Context) error {
	fmt.Printf("[Proxmox] Health check on %s: healthy\n", p.BaseURL)
	return nil
}

func (p *ProxmoxProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return []domain.Image{
		{ID: "lxc-ubuntu-22.04", Name: "Ubuntu 22.04 LTS (LXC)", Provider: "proxmox", Type: "template"},
		{ID: "vm-k8s-node-v1.28", Name: "K8s Node Template v1.28", Provider: "proxmox", Type: "template"},
		{ID: "debian-12-iso", Name: "Debian 12 ISO", Provider: "proxmox", Type: "iso"},
	}, nil
}

func (p *ProxmoxProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	fmt.Printf("[Proxmox] Deploying LXC container: %s (Image: %s, CPU: %v, RAM: %vMB)\n", c.Name, c.Image, c.CPU, c.MemoryMB)
	c.State = "running"
	return nil
}

func (p *ProxmoxProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	fmt.Printf("[Proxmox] Creating Linux Bridge for network: %s (%s)\n", n.Name, n.CIDR)
	n.State = "active"
	return nil
}

func (p *ProxmoxProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	fmt.Printf("[Proxmox] Adding static route: %s via %s\n", r.Destination, r.NextHop)
	return nil
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
