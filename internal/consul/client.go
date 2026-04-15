package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

// Client wraps the Consul API client with datacenter context.
type Client struct {
	api        *consulapi.Client
	Datacenter string
}

// Config holds connection parameters for a Consul datacenter.
type Config struct {
	Address    string
	Datacenter string
	Token      string
	Scheme     string
}

// NewClient creates a new Consul client for the given config.
func NewClient(cfg Config) (*Client, error) {
	apiCfg := consulapi.DefaultConfig()
	apiCfg.Address = cfg.Address
	apiCfg.Datacenter = cfg.Datacenter
	apiCfg.Token = cfg.Token
	if cfg.Scheme != "" {
		apiCfg.Scheme = cfg.Scheme
	}

	client, err := consulapi.NewClient(apiCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client for %s: %w", cfg.Datacenter, err)
	}

	return &Client{
		api:        client,
		Datacenter: cfg.Datacenter,
	}, nil
}

// ListKeys retrieves all KV pairs under the given prefix.
func (c *Client) ListKeys(prefix string) (map[string]string, error) {
	kv := c.api.KV()
	pairs, _, err := kv.List(prefix, &consulapi.QueryOptions{
		Datacenter: c.Datacenter,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys under %q in %s: %w", prefix, c.Datacenter, err)
	}

	result := make(map[string]string, len(pairs))
	for _, p := range pairs {
		result[p.Key] = string(p.Value)
	}
	return result, nil
}
