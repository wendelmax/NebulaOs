package vault

import (
	"context"
	"testing"
)

func TestNewVaultProvider(t *testing.T) {
	p := NewVaultProvider("http://localhost:8200", "test-token")
	if p == nil {
		t.Fatal("NewVaultProvider() returned nil")
	}
	if p.URL != "http://localhost:8200" {
		t.Errorf("URL = %q, want %q", p.URL, "http://localhost:8200")
	}
	if p.Token != "test-token" {
		t.Errorf("Token = %q, want %q", p.Token, "test-token")
	}
}

func TestVaultProvider_StoreSecret_NoConnection(t *testing.T) {
	p := NewVaultProvider("http://localhost:1", "test-token")
	err := p.StoreSecret(context.Background(), "test/key", "test-value")
	if err == nil {
		t.Error("expected error when vault is not running")
	}
}

func TestVaultProvider_GetSecret_NoConnection(t *testing.T) {
	p := NewVaultProvider("http://localhost:1", "test-token")
	_, err := p.GetSecret(context.Background(), "test/key")
	if err == nil {
		t.Error("expected error when vault is not running")
	}
}
