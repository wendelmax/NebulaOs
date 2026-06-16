package services

import (
	"context"
	"fmt"
	"log"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type MetaOrchestrator struct {
	factory    domain.ProviderFactory
	networkMgr *NetworkManager
}

func NewMetaOrchestrator(factory domain.ProviderFactory, nm *NetworkManager) *MetaOrchestrator {
	return &MetaOrchestrator{
		factory:    factory,
		networkMgr: nm,
	}
}

func (o *MetaOrchestrator) ProvisionMultiZoneStack(ctx context.Context, projectID string, regions []string) error {
	if len(regions) == 0 {
		return fmt.Errorf("at least one region must be specified")
	}
	log.Printf("[MetaOrchestrator] Provisioning global stack for project %s across %d regions", projectID, len(regions))
	for i, region := range regions {
		log.Printf("[MetaOrchestrator] [%d/%d] Setting up regional endpoint in %s", i+1, len(regions), region)
		provider, err := o.factory.GetProvider("aws")
		if err != nil {
			log.Printf("[MetaOrchestrator] AWS provider not available for region %s: %v", region, err)
			continue
		}
		resource := &domain.Resource{
			ID:       domain.NewID(),
			ProjectID: projectID,
			Name:     fmt.Sprintf("stack-%s", region),
			Type:     domain.ComputeResource,
			Provider: "aws",
			State:    "provisioning",
		}
		if err := provider.Provision(ctx, resource); err != nil {
			log.Printf("[MetaOrchestrator] Failed to provision in region %s: %v", region, err)
		} else {
			log.Printf("[MetaOrchestrator] Successfully initiated provisioning in %s (resource: %s)", region, resource.ID)
		}
	}
	return nil
}

func (o *MetaOrchestrator) MigrateResource(ctx context.Context, resourceID string, targetRegion string) error {
	log.Printf("[MetaOrchestrator] Migrating resource %s to region %s", resourceID, targetRegion)
	provider, err := o.factory.GetProvider("aws")
	if err != nil {
		return fmt.Errorf("no provider available for migration: %w", err)
	}
	resource := &domain.Resource{
		ID:       resourceID,
		Provider: "aws",
		State:    "migrating",
	}
	if err := provider.Provision(ctx, resource); err != nil {
		return fmt.Errorf("migration provisioning failed: %w", err)
	}
	log.Printf("[MetaOrchestrator] Resource %s migration to %s initiated", resourceID, targetRegion)
	return nil
}

func (o *MetaOrchestrator) SetupAvailabilityZone(ctx context.Context, regionID string, zoneName string, provider string) error {
	log.Printf("[MetaOrchestrator] Initializing Availability Zone %s in Region %s using %s provider", zoneName, regionID, provider)
	p, err := o.factory.GetProvider(provider)
	if err != nil {
		return fmt.Errorf("provider %q not found for zone setup: %w", provider, err)
	}
	network := &domain.Network{
		Name:   fmt.Sprintf("az-%s-%s", regionID, zoneName),
		RegionID: regionID,
		State:  "pending",
	}
	if err := p.ConfigureNetwork(ctx, network); err != nil {
		return fmt.Errorf("failed to configure network for zone %s: %w", zoneName, err)
	}
	log.Printf("[MetaOrchestrator] Availability Zone %s initialized in region %s", zoneName, regionID)
	return nil
}
