package diff

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeRebaseTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "rebase.yaml")
	require.NoError(t, os.WriteFile(p, []byte(content), 0o644))
	return p
}

func TestLoadRebaseOptions_EmptyPath_ReturnsZero(t *testing.T) {
	opts, err := LoadRebaseOptions("")
	require.NoError(t, err)
	assert.Equal(t, RebaseOptions{}, opts)
}

func TestLoadRebaseOptions_ValidFile(t *testing.T) {
	p := writeRebaseTemp(t, "strip_prefix: prod/svc\nnew_prefix: staging/svc\nskip_unmatched: true\n")
	opts, err := LoadRebaseOptions(p)
	require.NoError(t, err)
	assert.Equal(t, "prod/svc", opts.StripPrefix)
	assert.Equal(t, "staging/svc", opts.NewPrefix)
	assert.True(t, opts.SkipUnmatched)
}

func TestLoadRebaseOptions_MissingFile_ReturnsError(t *testing.T) {
	_, err := LoadRebaseOptions("/nonexistent/rebase.yaml")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rebase: read config")
}

func TestLoadRebaseOptions_InvalidYAML_ReturnsError(t *testing.T) {
	p := writeRebaseTemp(t, ": : invalid")
	_, err := LoadRebaseOptions(p)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "rebase: parse config")
}

func TestLoadRebaseOptions_PartialFile_DefaultsRemainder(t *testing.T) {
	p := writeRebaseTemp(t, "strip_prefix: prod\n")
	opts, err := LoadRebaseOptions(p)
	require.NoError(t, err)
	assert.Equal(t, "prod", opts.StripPrefix)
	assert.Equal(t, "", opts.NewPrefix)
	assert.False(t, opts.SkipUnmatched)
}
