package domain

import (
	"context"
	"time"
)

// Terraform State API models
type TerraformState struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Version   int       `json:"version"`
	State     []byte    `json:"state"` // Scalable blob
	UpdatedAt time.Time `json:"updated_at"`
}

type TerraformStateRepository interface {
	Save(ctx context.Context, state *TerraformState) error
	GetByProjectID(ctx context.Context, projectID string) (*TerraformState, error)
}

// Marketplace Blueprint models
type Blueprint struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	Resources   []ResourceDefinition   `json:"resources"`
	Variables   map[string]interface{} `json:"variables"`
	CreatedAt   time.Time              `json:"created_at"`
}

type ResourceDefinition struct {
	Name     string                 `json:"name"`
	Type     ResourceType           `json:"type"`
	Provider string                 `json:"provider"`
	Metadata map[string]interface{} `json:"metadata"`
}

type BlueprintRepository interface {
	List(ctx context.Context) ([]*Blueprint, error)
	GetByID(ctx context.Context, id string) (*Blueprint, error)
	Create(ctx context.Context, b *Blueprint) error
}
