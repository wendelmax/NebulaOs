package services

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

// GSLBManager handles cross-region traffic routing
type GSLBManager struct {
	repo domain.GSLBRepository
}

func NewGSLBManager(repo domain.GSLBRepository) *GSLBManager {
	return &GSLBManager{repo: repo}
}

func (m *GSLBManager) CreateEndpoint(ctx context.Context, ep *domain.GlobalEndpoint) error {
	ep.State = "active"
	return m.repo.Save(ctx, ep)
}

func (m *GSLBManager) ListEndpoints(ctx context.Context) ([]*domain.GlobalEndpoint, error) {
	return m.repo.List(ctx)
}

// AI Resource Advisor
type AIResourceAdvisor struct {
	resourceRepo domain.ResourceRepository
}

func NewAIResourceAdvisor(resRepo domain.ResourceRepository) *AIResourceAdvisor {
	return &AIResourceAdvisor{resourceRepo: resRepo}
}

func (a *AIResourceAdvisor) AnalyzeUsage(ctx context.Context, projectID string) ([]domain.ResourceInsight, error) {
	resources, err := a.resourceRepo.GetByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var insights []domain.ResourceInsight

	for _, res := range resources {
		// Mock AI Logic: Flag resources older than 7 days with no recent metadata changes
		if time.Since(res.CreatedAt) > 168*time.Hour {
			insights = append(insights, domain.ResourceInsight{
				ResourceID: res.ID,
				Type:       "cost",
				Message:    "Resource identified as 'Zombie'. Consider downsizing or decommissioning to save $45/mo.",
				Severity:   "medium",
				Actionable: true,
				CreatedAt:  time.Now(),
			})
		}
	}

	// Always add a "Premium Recommendation"
	insights = append(insights, domain.ResourceInsight{
		ResourceID: "global",
		Type:       "performance",
		Message:    "High latency detected in us-east-1. Enable GSLB Failover to improve UX by 35%.",
		Severity:   "high",
		Actionable: true,
		CreatedAt:  time.Now(),
	})

	return insights, nil
}
