package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/wendelmax/nebulaos/src/api/domain"
)

type KeycloakProvider struct {
	URL        string
	ClientID   string
	httpClient *http.Client
}

func NewKeycloakProvider(url, clientID string) *KeycloakProvider {
	return &KeycloakProvider{
		URL:        url,
		ClientID:   clientID,
		httpClient: http.DefaultClient,
	}
}

func (p *KeycloakProvider) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL+"/health/live", nil)
	if err != nil {
		return err
	}
	resp, err := p.httpClient.Do(req)
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
	data := url.Values{
		"client_id":  {p.ClientID},
		"grant_type": {"password"},
		"username":   {username},
		"password":   {password},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", p.URL+"/realms/master/protocol/openid-connect/token", bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		return "", fmt.Errorf("keycloak token request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("keycloak token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("keycloak authentication failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("keycloak response parse failed: %w", err)
	}
	return result.AccessToken, nil
}

func (p *KeycloakProvider) ValidateToken(ctx context.Context, token string) (*domain.User, error) {
	if token == "" {
		return nil, fmt.Errorf("invalid token")
	}
	data := url.Values{
		"client_id": {p.ClientID},
		"token":     {token},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", p.URL+"/realms/master/protocol/openid-connect/token/introspect", bytes.NewReader([]byte(data.Encode())))
	if err != nil {
		return nil, fmt.Errorf("keycloak introspect request failed: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("keycloak introspect request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak token introspection failed with status: %d", resp.StatusCode)
	}

	var result struct {
		Active   bool   `json:"active"`
		Sub      string `json:"sub"`
		Username string `json:"preferred_username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("keycloak response parse failed: %w", err)
	}
	if !result.Active {
		return nil, fmt.Errorf("token is not active")
	}

	return &domain.User{
		ID:       result.Sub,
		Username: result.Username,
		TenantID: "tenant-root",
	}, nil
}
