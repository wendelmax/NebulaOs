package kubernetes

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type KubernetesProvider struct {
	KubeConfig string
}

func NewKubernetesProvider(kubeConfig string) *KubernetesProvider {
	return &KubernetesProvider{
		KubeConfig: kubeConfig,
	}
}

func (p *KubernetesProvider) Provision(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("Kubernetes: Scaling deployment %s in NebulaOS cluster...\n", resource.Name)
	resource.ProviderID = fmt.Sprintf("k8s-%s", resource.Name)
	resource.State = "active"
	return nil
}

func (p *KubernetesProvider) Decommission(ctx context.Context, resource *domain.Resource) error {
	fmt.Printf("Kubernetes: Removing deployment %s...\n", resource.ProviderID)
	resource.State = "deleted"
	return nil
}

func (p *KubernetesProvider) GetStatus(ctx context.Context, resourceID string) (string, error) {
	return "active", nil
}
func (p *KubernetesProvider) Ping(ctx context.Context) error {
	return nil
}

func (p *KubernetesProvider) ListImages(ctx context.Context) ([]domain.Image, error) {
	return nil, nil
}

func (p *KubernetesProvider) DeployContainer(ctx context.Context, c *domain.Container) error {
	fmt.Printf("[Kubernetes] Creating Pod/Deployment for %s\n", c.Name)
	c.State = "running"
	return nil
}

func (p *KubernetesProvider) ConfigureNetwork(ctx context.Context, n *domain.Network) error {
	return nil
}

func (p *KubernetesProvider) AddRoute(ctx context.Context, r *domain.Route) error {
	return nil
}

func (p *KubernetesProvider) AttachSecurityGroup(ctx context.Context, resourceID string, sgID string) error {
	fmt.Printf("[Kubernetes] Applying NetworkPolicy %s to deployment %s\n", sgID, resourceID)
	return nil
}
