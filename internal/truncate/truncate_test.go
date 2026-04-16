package truncate_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"consul-drift-check/internal/truncate"
)

func makePairs(kvs map[string]string) []*api.KVPair {
	out := make([]*api.KVPair, 0, len(kvs))
	for k, v := range kvs {
		out = append(out, &api.KVPair{Key: k, Value: []byte(v)})
	}
	return out
}

func TestApply_NoLimit_ReturnsUnchanged(t *testing.T) {
	pairs := makePairs(map[string]string{"a": "hello", "b": "world"})
	result := truncate.Apply(pairs, truncate.Options{MaxBytes: 0})
	assert.Equal(t, pairs, result)
}

func TestApply_ShortValues_NotTruncated(t *testing.T) {
	pairs := makePairs(map[string]string{"key": "hi"})
	result := truncate.Apply(pairs, truncate.Options{MaxBytes: 10})
	require.Len(t, result, 1)
	assert.Equal(t, "hi", string(result[0].Value))
}

func TestApply_LongValue_TruncatedWithDefaultSuffix(t *testing.T) {
	pairs := makePairs(map[string]string{"key": "abcdefghij"})
	result := truncate.Apply(pairs, truncate.Options{MaxBytes: 5})
	require.Len(t, result, 1)
	assert.Equal(t, "abcde...", string(result[0].Value))
}

func TestApply_LongValue_TruncatedWithCustomSuffix(t *testing.T) {
	pairs := makePairs(map[string]string{"key": "abcdefghij"})
	result := truncate.Apply(pairs, truncate.Options{MaxBytes: 4, Suffix: "[cut]"})
	require.Len(t, result, 1)
	assert.Equal(t, "abcd[cut]", string(result[0].Value))
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := []byte("abcdefghij")
	pairs := []*api.KVPair{{Key: "k", Value: original}}
	truncate.Apply(pairs, truncate.Options{MaxBytes: 3})
	assert.Equal(t, "abcdefghij", string(pairs[0].Value))
}

func TestApply_MultiplePairs_OnlySomeTruncated(t *testing.T) {
	pairs := []*api.KVPair{
		{Key: "short", Value: []byte("hi")},
		{Key: "long", Value: []byte("verylongvalue")},
	}
	result := truncate.Apply(pairs, truncate.Options{MaxBytes: 5})
	require.Len(t, result, 2)
	assert.Equal(t, "hi", string(result[0].Value))
	assert.Equal(t, "veryl...", string(result[1].Value))
}
