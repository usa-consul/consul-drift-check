package consul

import (
	"fmt"
	"strings"

	"github.com/hashicorp/consul/api"
)

// KVEntry represents a single key-value pair from Consul.
type KVEntry struct {
	Key   string
	Value []byte
}

// KVLister defines the interface for listing KV pairs under a prefix.
type KVLister interface {
	ListPrefix(prefix string) ([]KVEntry, error)
}

// KVClient wraps a Consul API client to provide KV operations.
type KVClient struct {
	client *api.Client
}

// NewKVClient creates a new KVClient from an existing api.Client.
func NewKVClient(c *api.Client) *KVClient {
	return &KVClient{client: c}
}

// ListPrefix retrieves all KV entries under the given prefix.
// The prefix is normalized to ensure it does not have a leading slash.
func (k *KVClient) ListPrefix(prefix string) ([]KVEntry, error) {
	prefix = strings.TrimPrefix(prefix, "/")

	pairs, _, err := k.client.KV().List(prefix, nil)
	if err != nil {
		return nil, fmt.Errorf("listing prefix %q: %w", prefix, err)
	}

	entries := make([]KVEntry, 0, len(pairs))
	for _, p := range pairs {
		entries = append(entries, KVEntry{
			Key:   p.Key,
			Value: p.Value,
		})
	}

	return entries, nil
}
