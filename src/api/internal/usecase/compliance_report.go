package usecase

import (
	"context"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type ComplianceReport struct {
	TenantID string `json:"tenant_id"`
	Usage    struct {
		CPUs int `json:"cpus"`
		RAM  int `json:"ram"`
		Disk int `json:"disk"`
	} `json:"usage"`
	Quota domain.Quota `json:"quota"`
}

type GetComplianceReportUseCase struct {
	resourceRepo domain.ResourceRepository
	projectRepo  domain.ProjectRepository
	quotaRepo    domain.QuotaRepository
}

func NewGetComplianceReportUseCase(resRepo domain.ResourceRepository, projRepo domain.ProjectRepository, qRepo domain.QuotaRepository) *GetComplianceReportUseCase {
	return &GetComplianceReportUseCase{
		resourceRepo: resRepo,
		projectRepo:  projRepo,
		quotaRepo:    qRepo,
	}
}

func (uc *GetComplianceReportUseCase) Execute(ctx context.Context, tenantID string) (ComplianceReport, error) {
	quota, _ := uc.quotaRepo.GetByTenant(ctx, tenantID)
	projects, _ := uc.projectRepo.GetByTenant(ctx, tenantID)
	allResources, _ := uc.resourceRepo.List(ctx)

	report := ComplianceReport{
		TenantID: tenantID,
		Quota:    *quota,
	}

	for _, p := range projects {
		for _, res := range allResources {
			if res.ProjectID == p.ID {
				report.Usage.CPUs += 1
				report.Usage.RAM += 2048
				report.Usage.Disk += 20
			}
		}
	}

	return report, nil
}
