package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeSampleResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified"}
	}
	return out
}

func TestSample_EmptyInput_ReturnsNil(t *testing.T) {
	result := Sample(nil, SampleOptions{MaxResults: 5})
	assert.Nil(t, result)
}

func TestSample_NoLimit_ReturnsAll(t *testing.T) {
	input := makeSampleResults("a", "b", "c")
	out := Sample(input, SampleOptions{})
	require.Len(t, out, 3)
}

func TestSample_LimitExceedsInput_ReturnsAll(t *testing.T) {
	input := makeSampleResults("a", "b", "c")
	out := Sample(input, SampleOptions{MaxResults: 100})
	require.Len(t, out, 3)
}

func TestSample_LimitsToMaxResults(t *testing.T) {
	input := makeSampleResults("a", "b", "c", "d", "e")
	out := Sample(input, SampleOptions{MaxResults: 3, Seed: 42})
	require.Len(t, out, 3)
}

func TestSample_ResultIsSortedByKey(t *testing.T) {
	input := makeSampleResults("z", "m", "a", "b", "c")
	out := Sample(input, SampleOptions{MaxResults: 3, Seed: 7})
	for i := 1; i < len(out); i++ {
		assert.LessOrEqual(t, out[i-1].Key, out[i].Key)
	}
}

func TestSample_DeterministicSameSeedSameSubset(t *testing.T) {
	input := makeSampleResults("a", "b", "c", "d", "e", "f")
	opts := SampleOptions{MaxResults: 3, Seed: 99, Deterministic: true}
	first := Sample(input, opts)
	second := Sample(input, opts)
	require.Equal(t, first, second)
}

func TestSample_DoesNotMutateInput(t *testing.T) {
	input := makeSampleResults("a", "b", "c", "d")
	orig := make([]Result, len(input))
	copy(orig, input)
	Sample(input, SampleOptions{MaxResults: 2, Seed: 1})
	assert.Equal(t, orig, input)
}
