package consul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/consul/api"
)

func newTestKVClient(t *testing.T, handler http.Handler) *KVClient {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	cfg := api.DefaultConfig()
	cfg.Address = server.URL

	c, err := api.NewClient(cfg)
	if err != nil {
		t.Fatalf("creating test consul client: %v", err)
	}
	return NewKVClient(c)
}

func TestListPrefix_ReturnsEntries(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"Key":"config/app/host","Value":"bG9jYWxob3N0"},
			{"Key":"config/app/port","Value":"ODA4MA=="}
		]`))
	})

	kv := newTestKVClient(t, handler)
	entries, err := kv.ListPrefix("config/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "config/app/host" {
		t.Errorf("expected key 'config/app/host', got %q", entries[0].Key)
	}
}

func TestListPrefix_EmptyResult(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	kv := newTestKVClient(t, handler)
	entries, err := kv.ListPrefix("nonexistent/prefix")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestListPrefix_StripLeadingSlash(t *testing.T) {
	var capturedPath string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusNotFound)
	})

	kv := newTestKVClient(t, handler)
	_, _ = kv.ListPrefix("/config/app")

	if capturedPath != "/v1/kv/config/app" {
		t.Errorf("expected path '/v1/kv/config/app', got %q", capturedPath)
	}
}
