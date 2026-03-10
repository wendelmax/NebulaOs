package services

import (
	"context"
	"fmt"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type NetworkManager struct {
	factory domain.ProviderFactory
}

func NewNetworkManager(factory domain.ProviderFactory) *NetworkManager {
	return &NetworkManager{factory: factory}
}

func (m *NetworkManager) CreateInterProviderRoute(ctx context.Context, route *domain.Route, sourceProvider, targetProvider string) error {
	// 1. Get source provider
	srcProv, err := m.factory.GetProvider(sourceProvider)
	if err != nil {
		return fmt.Errorf("source provider %s not found: %w", sourceProvider, err)
	}

	// 2. Configure source provider to route traffic to target provider
	route.TargetProvider = targetProvider
	if err := srcProv.AddRoute(ctx, route); err != nil {
		return fmt.Errorf("failed to add route to source provider %s: %w", sourceProvider, err)
	}

	fmt.Printf("[NetworkManager] Successfully established route from %s to %s for destination %s\n", 
		sourceProvider, targetProvider, route.Destination)
	
	return nil
}

func (m *NetworkManager) ProximoxToOpenStackBridge(ctx context.Context, proxmoxNet, openstackNet string) error {
	fmt.Printf("[NetworkManager] Establishing SDN Bridge between Proxmox (%s) and OpenStack (%s)\n", proxmoxNet, openstackNet)
	return nil
}
