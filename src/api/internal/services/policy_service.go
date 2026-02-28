package services

import (
	"context"
	"fmt"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type SovereignPolicyService struct {
	policyRepo domain.SovereigntyPolicyRepository
}

func NewSovereignPolicyService(repo domain.SovereigntyPolicyRepository) *SovereignPolicyService {
	return &SovereignPolicyService{policyRepo: repo}
}

func (s *SovereignPolicyService) ValidateRegion(ctx context.Context, tenantID string, region string) error {
	policy, err := s.policyRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		// If no policy exists, we allow all for now (sovereign by default or permissive)
		return nil
	}

	for _, r := range policy.AllowedRegions {
		if r == region {
			return nil
		}
	}

	return fmt.Errorf("policy violation: region %s is NOT allowed for tenant %s. Sovereign boundary enforcement active.", region, tenantID)
}

func (s *SovereignPolicyService) GetPolicy(ctx context.Context, tenantID string) (*domain.SovereigntyPolicy, error) {
	return s.policyRepo.GetByTenantID(ctx, tenantID)
}

func (s *SovereignPolicyService) UpdatePolicy(ctx context.Context, policy *domain.SovereigntyPolicy) error {
	return s.policyRepo.Save(ctx, policy)
}
