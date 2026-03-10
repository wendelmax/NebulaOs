package vault

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wendelmax/nebulaos/src/api/domain"
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

func (p *VaultProvider) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL+"/v1/sys/health", nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests { // Vault returns 429 if uninitialized/sealed but still alive
		return fmt.Errorf("vault health check failed with status: %d", resp.StatusCode)
	}
	return nil
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
