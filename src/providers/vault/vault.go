package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
		return fmt.Errorf("vault health check failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (p *VaultProvider) StoreSecret(ctx context.Context, key string, value string) error {
	body, _ := json.Marshal(map[string]interface{}{"data": map[string]string{"value": value}})
	req, err := http.NewRequestWithContext(ctx, "PUT", p.URL+"/v1/secret/data/"+key, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("vault store request failed: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("vault store request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		io.ReadAll(resp.Body)
		return fmt.Errorf("vault store secret failed with status: %d", resp.StatusCode)
	}
	return nil
}

func (p *VaultProvider) GetSecret(ctx context.Context, key string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.URL+"/v1/secret/data/"+key, nil)
	if err != nil {
		return "", fmt.Errorf("vault get request failed: %w", err)
	}
	req.Header.Set("X-Vault-Token", p.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("vault get request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("vault get secret failed with status: %d", resp.StatusCode)
	}
	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Data struct {
			Data map[string]string `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("vault response parse failed: %w", err)
	}
	if val, ok := result.Data.Data["value"]; ok {
		return val, nil
	}
	return "", fmt.Errorf("secret key %q not found in vault response", key)
}

var _ domain.SecretManager = (*VaultProvider)(nil)
