package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level CLI configuration.
type Config struct {
	Source      Datacenter `yaml:"source"`
	Destination Datacenter `yaml:"destination"`
	Prefix      string     `yaml:"prefix"`
	Format      string     `yaml:"format"`
}

// Datacenter holds connection settings for a single Consul datacenter.
type Datacenter struct {
	Address string `yaml:"address"`
	Token   string `yaml:"token"`
	DC      string `yaml:"dc"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("config path must not be empty")
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.Source.Address == "" {
		return errors.New("source.address is required")
	}
	if c.Destination.Address == "" {
		return errors.New("destination.address is required")
	}
	if c.Prefix == "" {
		return errors.New("prefix is required")
	}
	if c.Format != "" && c.Format != "text" && c.Format != "json" {
		return fmt.Errorf("format must be \"text\" or \"json\", got %q", c.Format)
	}
	return nil
}
