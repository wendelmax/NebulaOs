package keycloak

import (
	"context"
	"testing"
)

func TestNewKeycloakProvider(t *testing.T) {
	p := NewKeycloakProvider("http://localhost:8080", "nebula-api")
	if p == nil {
		t.Fatal("NewKeycloakProvider() returned nil")
	}
	if p.ClientID != "nebula-api" {
		t.Errorf("ClientID = %q, want %q", p.ClientID, "nebula-api")
	}
}

func TestKeycloakProvider_Ping_NoConnection(t *testing.T) {
	p := NewKeycloakProvider("http://localhost:1", "nebula-api")
	err := p.Ping(context.Background())
	if err == nil {
		t.Error("expected error when keycloak is not running")
	}
}

func TestKeycloakProvider_ValidateToken_Empty(t *testing.T) {
	p := NewKeycloakProvider("http://localhost:8080", "nebula-api")
	_, err := p.ValidateToken(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty token")
	}
}

func TestKeycloakProvider_Authenticate_NoConnection(t *testing.T) {
	p := NewKeycloakProvider("http://localhost:1", "nebula-api")
	_, err := p.Authenticate(context.Background(), "admin", "admin")
	if err == nil {
		t.Error("expected error when keycloak is not running")
	}
}
