package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeSliceResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified"}
	}
	return out
}

func TestSlice_EmptyInput_ReturnsNil(t *testing.T) {
	got := Slice(nil, SliceOptions{})
	assert.Nil(t, got)
}

func TestSlice_NoOptions_ReturnsAll(t *testing.T) {
	results := makeSliceResults("a", "b", "c")
	got := Slice(results, SliceOptions{})
	assert.Len(t, got, 3)
}

func TestSlice_LimitApplied(t *testing.T) {
	results := makeSliceResults("a", "b", "c", "d")
	got := Slice(results, SliceOptions{Limit: 2})
	assert.Len(t, got, 2)
}

func TestSlice_OffsetApplied(t *testing.T) {
	results := makeSliceResults("a", "b", "c")
	got := Slice(results, SliceOptions{Offset: 1})
	assert.Len(t, got, 2)
	assert.Equal(t, "b", got[0].Key)
}

func TestSlice_OffsetAndLimit(t *testing.T) {
	results := makeSliceResults("a", "b", "c", "d", "e")
	got := Slice(results, SliceOptions{Offset: 1, Limit: 2})
	assert.Len(t, got, 2)
	assert.Equal(t, "b", got[0].Key)
	assert.Equal(t, "c", got[1].Key)
}

func TestSlice_OffsetBeyondLength_ReturnsNil(t *testing.T) {
	results := makeSliceResults("a", "b")
	got := Slice(results, SliceOptions{Offset: 10})
	assert.Nil(t, got)
}

func TestSlice_OrderByKey_SortsBeforeSlice(t *testing.T) {
	results := makeSliceResults("c", "a", "b")
	got := Slice(results, SliceOptions{OrderByKey: true, Limit: 2})
	assert.Len(t, got, 2)
	assert.Equal(t, "a", got[0].Key)
	assert.Equal(t, "b", got[1].Key)
}

func TestSlice_DoesNotMutateInput(t *testing.T) {
	results := makeSliceResults("c", "a", "b")
	_ = Slice(results, SliceOptions{OrderByKey: true})
	assert.Equal(t, "c", results[0].Key)
}
