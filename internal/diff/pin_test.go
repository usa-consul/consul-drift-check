package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makePinResults(pairs [][3]string) []Result {
	out := make([]Result, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, Result{Key: p[0], Status: p[1], SourceValue: []byte(p[2])})
	}
	return out
}

func TestPin_EmptyResults_ReturnsNil(t *testing.T) {
	violations := Pin(nil, PinOptions{PinnedKeys: []string{"config/key"}})
	assert.Nil(t, violations)
}

func TestPin_NoViolations_AllMatch(t *testing.T) {
	results := makePinResults([][3]string{
		{"config/tls/cert", "match", "v1"},
		{"config/bootstrap/token", "match", "abc"},
	})
	violations := Pin(results, PinOptions{
		PinnedKeys:     []string{"config/tls/cert"},
		PinnedPrefixes: []string{"config/bootstrap/"},
	})
	assert.Nil(t, violations)
}

func TestPin_ExactKeyDrifted_ReturnsViolation(t *testing.T) {
	results := makePinResults([][3]string{
		{"config/tls/cert", "modified", "v1"},
		{"config/other", "modified", "x"},
	})
	violations := Pin(results, PinOptions{
		PinnedKeys: []string{"config/tls/cert"},
	})
	require.Len(t, violations, 1)
	assert.Equal(t, "config/tls/cert", violations[0].Key)
	assert.Equal(t, "modified", violations[0].Status)
	assert.Equal(t, "exact key pin", violations[0].Reason)
}

func TestPin_PrefixDrifted_ReturnsViolation(t *testing.T) {
	results := makePinResults([][3]string{
		{"config/bootstrap/token", "only_in_source", ""},
	})
	violations := Pin(results, PinOptions{
		PinnedPrefixes: []string{"config/bootstrap/"},
	})
	require.Len(t, violations, 1)
	assert.Contains(t, violations[0].Reason, "prefix pin")
}

func TestPin_SortedByKey(t *testing.T) {
	results := makePinResults([][3]string{
		{"z/key", "modified", ""},
		{"a/key", "modified", ""},
		{"m/key", "modified", ""},
	})
	violations := Pin(results, PinOptions{
		PinnedPrefixes: []string{"a/", "m/", "z/"},
	})
	require.Len(t, violations, 3)
	assert.Equal(t, "a/key", violations[0].Key)
	assert.Equal(t, "m/key", violations[1].Key)
	assert.Equal(t, "z/key", violations[2].Key)
}

func TestPin_CaseInsensitiveKeyMatch(t *testing.T) {
	results := makePinResults([][3]string{
		{"Config/TLS/Cert", "modified", ""},
	})
	violations := Pin(results, PinOptions{
		PinnedKeys: []string{"config/tls/cert"},
	})
	require.Len(t, violations, 1)
	assert.Equal(t, "Config/TLS/Cert", violations[0].Key)
}
