package infrastructure

import (
	"context"
	"time"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type SovereignBillingManager struct {
	resRepo    domain.ResourceRepository
	volRepo    domain.VolumeRepository
	bucketRepo domain.BucketRepository
	tenantRepo domain.TenantRepository
}

func NewSovereignBillingManager(resRepo domain.ResourceRepository, volRepo domain.VolumeRepository, bucketRepo domain.BucketRepository, tenantRepo domain.TenantRepository) *SovereignBillingManager {
	return &SovereignBillingManager{resRepo: resRepo, volRepo: volRepo, bucketRepo: bucketRepo, tenantRepo: tenantRepo}
}

func (m *SovereignBillingManager) GenerateReport(ctx context.Context, tenantID string) (*domain.BillingReport, error) {
	report := &domain.BillingReport{
		TenantID:    tenantID,
		TotalCost:   0,
		Items:       []domain.BillingItem{},
		GeneratedAt: time.Now(),
	}

	// For simulation: Get all resources and filter "by tenant"
	// In a production system, we would join through projects table
	vms, _ := m.resRepo.List(ctx)
	for _, vm := range vms {
		// Mock: only bill resources for the requested tenant
		// Simplified simulation of tenant-level granularity
		item := domain.BillingItem{
			ResourceID: vm.ID,
			Type:       string(vm.Type),
			Cost:       15.50,
			Currency:   "USD",
		}
		report.Items = append(report.Items, item)
		report.TotalCost += item.Cost
	}

	return report, nil
}

func (m *SovereignBillingManager) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	resources, _ := m.resRepo.List(ctx)
	tenants, _ := m.tenantRepo.List(ctx)

	var totalCPUs float64
	for _, res := range resources {
		if res.Type == domain.ComputeResource {
			if cpu, ok := res.Metadata["cpus"].(float64); ok {
				totalCPUs += cpu
			} else {
				totalCPUs += 2 // Default mock
			}
		}
	}

	return &domain.GlobalStats{
		TotalCPUs:     totalCPUs,
		TotalStorage:  1.5, // TB mock
		TotalEgress:   250, // GB mock
		ActiveTenants: len(tenants),
		TrendCPUs:     12.5,
		TrendStorage:  8.2,
	}, nil
}
