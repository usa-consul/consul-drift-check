package transform

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config describes a serialisable pipeline definition.
type Config struct {
	Stages []StageConfig `yaml:"stages"`
}

// StageConfig holds the name and optional parameter for a built-in stage.
type StageConfig struct {
	Type  string `yaml:"type"`
	Param string `yaml:"param,omitempty"`
}

// LoadConfig reads a YAML file and returns a Config.
// An empty path returns an empty Config without error.
func LoadConfig(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("transform: read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("transform: parse config: %w", err)
	}
	return cfg, nil
}

// Build constructs a Pipeline from a Config, returning an error for unknown
// stage types.
func Build(cfg Config) (*Pipeline, error) {
	p := New()
	for _, s := range cfg.Stages {
		switch s.Type {
		case "strip_prefix":
			p.Add(s.Type, StripPrefix(s.Param))
		case "lower_keys":
			p.Add(s.Type, LowerKeys())
		default:
			return nil, fmt.Errorf("transform: unknown stage type %q", s.Type)
		}
	}
	return p, nil
}
