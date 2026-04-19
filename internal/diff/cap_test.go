package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeCapResults() []Result {
	return []Result{
		{Key: "delta/x", Status: "modified"},
		{Key: "alpha/y", Status: "only_in_source"},
		{Key: "gamma/z", Status: "only_in_destination"},
		{Key: "beta/w", Status: "modified"},
	}
}

func TestCap_NoLimit_ReturnsAll(t *testing.T) {
	results := makeCapResults()
	out := Cap(results, CapOptions{})
	require.Len(t, out, 4)
}

func TestCap_LimitsToMaxResults(t *testing.T) {
	out := Cap(makeCapResults(), CapOptions{MaxResults: 2})
	require.Len(t, out, 2)
}

func TestCap_MaxResultsExceedsInput_ReturnsAll(t *testing.T) {
	out := Cap(makeCapResults(), CapOptions{MaxResults: 100})
	require.Len(t, out, 4)
}

func TestCap_OrderByKey_SortedBeforeCap(t *testing.T) {
	out := Cap(makeCapResults(), CapOptions{MaxResults: 2, OrderByKey: true})
	require.Len(t, out, 2)
	assert.Equal(t, "alpha/y", out[0].Key)
	assert.Equal(t, "beta/w", out[1].Key)
}

func TestCap_OrderByKey_NoLimit_AllSorted(t *testing.T) {
	out := Cap(makeCapResults(), CapOptions{OrderByKey: true})
	require.Len(t, out, 4)
	assert.Equal(t, "alpha/y", out[0].Key)
	assert.Equal(t, "delta/x", out[3].Key)
}

func TestCap_EmptyInput_ReturnsNil(t *testing.T) {
	out := Cap(nil, CapOptions{MaxResults: 5})
	assert.Nil(t, out)
}

func TestCap_DoesNotMutateOriginal(t *testing.T) {
	results := makeCapResults()
	origFirst := results[0].Key
	Cap(results, CapOptions{OrderByKey: true, MaxResults: 2})
	assert.Equal(t, origFirst, results[0].Key)
}
