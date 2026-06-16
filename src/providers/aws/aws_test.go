package aws

import (
	"context"
	"testing"
)

func TestNewAWSProvider(t *testing.T) {
	provider, err := NewAWSProvider(context.Background(), "us-east-1", "")
	if err != nil {
		t.Fatalf("NewAWSProvider() error = %v", err)
	}
	if provider == nil {
		t.Fatal("NewAWSProvider() returned nil")
	}
}

func TestNewAWSProvider_CustomEndpoint(t *testing.T) {
	provider, err := NewAWSProvider(context.Background(), "us-east-1", "http://moto:4566")
	if err != nil {
		t.Fatalf("NewAWSProvider() error = %v", err)
	}
	if provider == nil {
		t.Fatal("NewAWSProvider() returned nil")
	}
}

func TestAWSProvider_Ping(t *testing.T) {
	provider, err := NewAWSProvider(context.Background(), "us-east-1", "")
	if err != nil {
		t.Fatalf("NewAWSProvider() error = %v", err)
	}
	if err := provider.Ping(context.Background()); err != nil {
		t.Errorf("Ping() error = %v", err)
	}
}
