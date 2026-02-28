package usecase

import (
	"context"
	"time"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type CreateSecurityGroupInput struct {
	ID        string `json:"id"`
	ProjectID string `json:"project_id"`
	Name      string `json:"name"`
}

type CreateSecurityGroupUseCase struct {
	repo domain.SecurityGroupRepository
}

func NewCreateSecurityGroupUseCase(repo domain.SecurityGroupRepository) *CreateSecurityGroupUseCase {
	return &CreateSecurityGroupUseCase{repo: repo}
}

func (uc *CreateSecurityGroupUseCase) Execute(ctx context.Context, input CreateSecurityGroupInput) error {
	sg := &domain.SecurityGroup{
		ID:        input.ID,
		ProjectID: input.ProjectID,
		Name:      input.Name,
		CreatedAt: time.Now(),
	}
	return uc.repo.Create(ctx, sg)
}

type AddFirewallRuleInput struct {
	ID              string                  `json:"id"`
	SecurityGroupID string                  `json:"security_group_id"`
	Protocol        domain.FirewallProtocol `json:"protocol"`
	FromPort        int                     `json:"from_port"`
	ToPort          int                     `json:"to_port"`
	SourceIP        string                  `json:"source_ip"`
	Action          string                  `json:"action"`
}

type AddFirewallRuleUseCase struct {
	repo domain.SecurityGroupRepository
}

func NewAddFirewallRuleUseCase(repo domain.SecurityGroupRepository) *AddFirewallRuleUseCase {
	return &AddFirewallRuleUseCase{repo: repo}
}

func (uc *AddFirewallRuleUseCase) Execute(ctx context.Context, input AddFirewallRuleInput) error {
	rule := domain.FirewallRule{
		ID:        input.ID,
		Protocol:  input.Protocol,
		FromPort:  input.FromPort,
		ToPort:    input.ToPort,
		SourceIP:  input.SourceIP,
		Action:    input.Action,
		CreatedAt: time.Now(),
	}
	return uc.repo.AddRule(ctx, input.SecurityGroupID, rule)
}

type ListSecurityGroupsUseCase struct {
	repo domain.SecurityGroupRepository
}

func NewListSecurityGroupsUseCase(repo domain.SecurityGroupRepository) *ListSecurityGroupsUseCase {
	return &ListSecurityGroupsUseCase{repo: repo}
}

func (uc *ListSecurityGroupsUseCase) Execute(ctx context.Context, projectID string) ([]*domain.SecurityGroup, error) {
	return uc.repo.ListByProject(ctx, projectID)
}
