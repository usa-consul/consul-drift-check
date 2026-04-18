package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatch_NoDrift_ReturnsEmpty(t *testing.T) {
	results := []Result{
		{Key: "a/b", Status: StatusMatch},
	}
	entries := Patch(results)
	assert.Empty(t, entries)
}

func TestPatch_OnlyInSource_ReturnsSet(t *testing.T) {
	results := []Result{
		{Key: "a/b", Status: StatusOnlyInSource, SourceValue: []byte("v1")},
	}
	entries := Patch(results)
	assert.Len(t, entries, 1)
	assert.Equal(t, "set", entries[0].Action)
	assert.Equal(t, "a/b", entries[0].Key)
	assert.Equal(t, []byte("v1"), entries[0].Value)
}

func TestPatch_OnlyInDestination_ReturnsDelete(t *testing.T) {
	results := []Result{
		{Key: "a/c", Status: StatusOnlyInDestination, DestValue: []byte("old")},
	}
	entries := Patch(results)
	assert.Len(t, entries, 1)
	assert.Equal(t, "delete", entries[0].Action)
	assert.Equal(t, "a/c", entries[0].Key)
}

func TestPatch_Modified_ReturnsSet(t *testing.T) {
	results := []Result{
		{Key: "x/y", Status: StatusModified, SourceValue: []byte("new"), DestValue: []byte("old")},
	}
	entries := Patch(results)
	assert.Len(t, entries, 1)
	assert.Equal(t, "set", entries[0].Action)
	assert.Equal(t, []byte("new"), entries[0].Value)
}

func TestPatch_SortedByKey(t *testing.T) {
	results := []Result{
		{Key: "z/z", Status: StatusOnlyInSource, SourceValue: []byte("1")},
		{Key: "a/a", Status: StatusOnlyInSource, SourceValue: []byte("2")},
	}
	entries := Patch(results)
	assert.Equal(t, "a/a", entries[0].Key)
	assert.Equal(t, "z/z", entries[1].Key)
}

func TestToKVPairs_OnlySetEntries(t *testing.T) {
	entries := []PatchEntry{
		{Key: "a", Action: "set", Value: []byte("v")},
		{Key: "b", Action: "delete"},
	}
	pairs := ToKVPairs(entries)
	assert.Len(t, pairs, 1)
	assert.Equal(t, "a", pairs[0].Key)
}
