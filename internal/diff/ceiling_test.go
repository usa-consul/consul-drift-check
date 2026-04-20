package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCeilingResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified"}
	}
	return out
}

func TestCeiling_ZeroMax_ReturnsAll(t *testing.T) {
	results := makeCeilingResults("app/a", "app/b", "app/c")
	got := Ceiling(results, CeilingOptions{MaxPerPrefix: 0})
	assert.Len(t, got, 3)
}

func TestCeiling_EmptyInput_ReturnsNil(t *testing.T) {
	got := Ceiling(nil, CeilingOptions{MaxPerPrefix: 2})
	assert.Nil(t, got)
}

func TestCeiling_LimitsPerPrefix(t *testing.T) {
	results := makeCeilingResults(
		"app/a", "app/b", "app/c",
		"db/x", "db/y",
	)
	got := Ceiling(results, CeilingOptions{MaxPerPrefix: 2})
	assert.Len(t, got, 4)
	prefixCount := map[string]int{}
	for _, r := range got {
		seg := topCeilingSegment(r.Key)
		prefixCount[seg]++
	}
	assert.Equal(t, 2, prefixCount["app"])
	assert.Equal(t, 2, prefixCount["db"])
}

func TestCeiling_OrderByKey_SortsBeforeCap(t *testing.T) {
	results := makeCeilingResults("app/z", "app/a", "app/m")
	got := Ceiling(results, CeilingOptions{MaxPerPrefix: 2, OrderByKey: true})
	assert.Len(t, got, 2)
	assert.Equal(t, "app/a", got[0].Key)
	assert.Equal(t, "app/m", got[1].Key)
}

func TestCeiling_LeadingSlash_Stripped(t *testing.T) {
	results := makeCeilingResults("/svc/a", "/svc/b", "/svc/c")
	got := Ceiling(results, CeilingOptions{MaxPerPrefix: 2})
	assert.Len(t, got, 2)
}

func TestCeiling_MultiplePrefixes_IndependentLimits(t *testing.T) {
	results := makeCeilingResults(
		"alpha/1", "alpha/2", "alpha/3",
		"beta/1",
	)
	got := Ceiling(results, CeilingOptions{MaxPerPrefix: 2})
	assert.Len(t, got, 3)
}
