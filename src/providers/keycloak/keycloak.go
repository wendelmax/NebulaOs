package keycloak

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type KeycloakProvider struct {
	URL      string
	ClientID string
}

func NewKeycloakProvider(url, clientID string) *KeycloakProvider {
	return &KeycloakProvider{
		URL:      url,
		ClientID: clientID,
	}
}

func (p *KeycloakProvider) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL+"/health/live", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("keycloak health check failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (p *KeycloakProvider) Authenticate(ctx context.Context, username, password string) (string, error) {
	log.Printf("[KeycloakProvider] Authenticating user: %s", username)
	// Placeholder for OIDC getToken call
	return "mock-jwt-token", nil
}

func (p *KeycloakProvider) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	log.Printf("[KeycloakProvider] Validating token")
	if token == "" {
		return nil, fmt.Errorf("invalid token")
	}
	// Placeholder for token introspection
	return &domain.User{
		ID:       "user-123",
		Username: "admin",
		TenantID: "tenant-root",
	}, nil
}
