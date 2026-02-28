package vault

import (
	"context"
	"log"

	"github.com/jacksonwendel/nebulaos/src/api/domain"
)

type VaultProvider struct {
	URL   string
	Token string
}

func NewVaultProvider(url, token string) *VaultProvider {
	return &VaultProvider{
		URL:   url,
		Token: token,
	}
}

func (p *VaultProvider) StoreSecret(ctx context.Context, key string, value string) error {
	log.Printf("[VaultProvider] Storing secret at: %s", key)
	return nil
}

func (p *VaultProvider) GetSecret(ctx context.Context, key string) (string, error) {
	log.Printf("[VaultProvider] Retrieving secret from: %s", key)
	return "mock-secret-value", nil
}

// Ensure VaultProvider implements domain.SecretManager
var _ domain.SecretManager = (*VaultProvider)(nil)
