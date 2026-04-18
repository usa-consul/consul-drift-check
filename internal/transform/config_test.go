package transform_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"consul-drift-check/internal/transform"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	require.NoError(t, err)
	_, err = f.WriteString(content)
	require.NoError(t, err)
	require.NoError(t, f.Close())
	return filepath.Clean(f.Name())
}

func TestLoadConfig_EmptyPath_ReturnsEmpty(t *testing.T) {
	cfg, err := transform.LoadConfig("")
	require.NoError(t, err)
	assert.Empty(t, cfg.Stages)
}

func TestLoadConfig_ValidFile(t *testing.T) {
	path := writeTemp(t, "stages:\n  - type: strip_prefix\n    param: prod/\n  - type: lower_keys\n")
	cfg, err := transform.LoadConfig(path)
	require.NoError(t, err)
	require.Len(t, cfg.Stages, 2)
	assert.Equal(t, "strip_prefix", cfg.Stages[0].Type)
	assert.Equal(t, "prod/", cfg.Stages[0].Param)
	assert.Equal(t, "lower_keys", cfg.Stages[1].Type)
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := transform.LoadConfig("/nonexistent/path.yaml")
	require.Error(t, err)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	path := writeTemp(t, ": : invalid")
	_, err := transform.LoadConfig(path)
	require.Error(t, err)
}

func TestBuild_ValidConfig(t *testing.T) {
	cfg := transform.Config{
		Stages: []transform.StageConfig{
			{Type: "strip_prefix", Param: "ns"},
			{Type: "lower_keys"},
		},
	}
	p, err := transform.Build(cfg)
	require.NoError(t, err)
	assert.Equal(t, []string{"strip_prefix", "lower_keys"}, p.Stages())
}

func TestBuild_UnknownStage_ReturnsError(t *testing.T) {
	cfg := transform.Config{
		Stages: []transform.StageConfig{{Type: "unknown_op"}},
	}
	_, err := transform.Build(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown stage type")
}
