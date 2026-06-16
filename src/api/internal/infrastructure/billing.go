package infrastructure

import (
	"context"
	"fmt"
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

	resources, err := m.resRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources for billing: %w", err)
	}

	volumes, err := m.volRepo.ListByProject(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes for billing: %w", err)
	}

	for _, vm := range resources {
		cost := 0.0
		switch vm.Type {
		case domain.ComputeResource:
			cost = 25.00
		case domain.StorageResource:
			cost = 10.00
		default:
			cost = 5.00
		}
		report.Items = append(report.Items, domain.BillingItem{
			ResourceID: vm.ID,
			Type:       string(vm.Type),
			Cost:       cost,
			Currency:   "USD",
		})
		report.TotalCost += cost
	}

	for _, vol := range volumes {
		cost := float64(vol.SizeGB) * 0.10
		report.Items = append(report.Items, domain.BillingItem{
			ResourceID: vol.ID,
			Type:       "volume",
			Cost:       cost,
			Currency:   "USD",
		})
		report.TotalCost += cost
	}

	return report, nil
}

func (m *SovereignBillingManager) GetGlobalStats(ctx context.Context) (*domain.GlobalStats, error) {
	resources, err := m.resRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources for stats: %w", err)
	}
	tenants, err := m.tenantRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants for stats: %w", err)
	}

	var totalCPUs float64
	var totalStorage float64
	var totalEgress float64

	for _, res := range resources {
		if res.Type == domain.ComputeResource {
			if cpu, ok := res.Metadata["cpus"].(float64); ok {
				totalCPUs += cpu
			} else {
				totalCPUs += 2
			}
		}
	}

	volumes, _ := m.volRepo.ListByProject(ctx, "")
	for range volumes {
		totalStorage += 0.1
	}

	buckets, _ := m.bucketRepo.ListByProject(ctx, "")
	totalEgress = float64(len(buckets) * 50)

	return &domain.GlobalStats{
		TotalCPUs:     totalCPUs,
		TotalStorage:  totalStorage,
		TotalEgress:   totalEgress,
		ActiveTenants: len(tenants),
		TrendCPUs:     totalCPUs * 0.15,
		TrendStorage:  totalStorage * 0.10,
	}, nil
}
