package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeDelta_EmptyMaps(t *testing.T) {
	result := ComputeDelta(nil, nil)
	assert.Empty(t, result)
}

func TestComputeDelta_NewPrefixes(t *testing.T) {
	before := map[string]int{}
	after := map[string]int{"app": 3, "db": 1}

	entries := ComputeDelta(before, after)
	assert.Len(t, entries, 2)
	assert.Equal(t, "app", entries[0].Prefix)
	assert.Equal(t, 3, entries[0].Delta)
	assert.Equal(t, 0, entries[0].Before)
	assert.Equal(t, 3, entries[0].After)
}

func TestComputeDelta_RemovedPrefixes(t *testing.T) {
	before := map[string]int{"app": 4}
	after := map[string]int{}

	entries := ComputeDelta(before, after)
	assert.Len(t, entries, 1)
	assert.Equal(t, -4, entries[0].Delta)
}

func TestComputeDelta_MixedChanges(t *testing.T) {
	before := map[string]int{"app": 2, "db": 5}
	after := map[string]int{"app": 4, "db": 3}

	entries := ComputeDelta(before, after)
	assert.Len(t, entries, 2)
	// sorted descending by delta: app (+2) before db (-2)
	assert.Equal(t, "app", entries[0].Prefix)
	assert.Equal(t, 2, entries[0].Delta)
	assert.Equal(t, "db", entries[1].Prefix)
	assert.Equal(t, -2, entries[1].Delta)
}

func TestTotalDelta_Positive(t *testing.T) {
	before := map[string]int{"app": 1, "db": 2}
	after := map[string]int{"app": 3, "db": 4}
	assert.Equal(t, 4, TotalDelta(before, after))
}

func TestTotalDelta_Zero(t *testing.T) {
	before := map[string]int{"app": 5}
	after := map[string]int{"app": 5}
	assert.Equal(t, 0, TotalDelta(before, after))
}

func TestTotalDelta_EmptyBefore(t *testing.T) {
	assert.Equal(t, 7, TotalDelta(nil, map[string]int{"x": 7}))
}
