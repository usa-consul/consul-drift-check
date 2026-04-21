package diff

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type rebaseConfigFile struct {
	StripPrefix   string `yaml:"strip_prefix"`
	NewPrefix     string `yaml:"new_prefix"`
	SkipUnmatched bool   `yaml:"skip_unmatched"`
}

// LoadRebaseOptions reads RebaseOptions from a YAML file at path.
// If path is empty, a zero-value RebaseOptions is returned with no error.
func LoadRebaseOptions(path string) (RebaseOptions, error) {
	if path == "" {
		return RebaseOptions{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return RebaseOptions{}, fmt.Errorf("rebase: read config: %w", err)
	}

	var cfg rebaseConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return RebaseOptions{}, fmt.Errorf("rebase: parse config: %w", err)
	}

	return RebaseOptions{
		StripPrefix:   cfg.StripPrefix,
		NewPrefix:     cfg.NewPrefix,
		SkipUnmatched: cfg.SkipUnmatched,
	}, nil
}
