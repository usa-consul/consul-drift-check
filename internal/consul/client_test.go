package consul

import (
	"testing"
)

func TestNewClient_DefaultScheme(t *testing.T) {
	cfg := Config{
		Address:    "127.0.0.1:8500",
		Datacenter: "dc1",
		Token:      "",
		Scheme:     "",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
	if client.Datacenter != "dc1" {
		t.Errorf("expected datacenter dc1, got %s", client.Datacenter)
	}
}

func TestNewClient_WithToken(t *testing.T) {
	cfg := Config{
		Address:    "consul.example.com:8500",
		Datacenter: "dc2",
		Token:      "secret-token",
		Scheme:     "https",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client.Datacenter != "dc2" {
		t.Errorf("expected datacenter dc2, got %s", client.Datacenter)
	}
}

func TestNewClient_EmptyAddress(t *testing.T) {
	// Consul SDK accepts empty address (falls back to default),
	// so we just verify no panic and a valid client is returned.
	cfg := Config{
		Address:    "",
		Datacenter: "dc1",
	}

	client, err := NewClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}
