package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntersect_EmptyMaps_ReturnsNil(t *testing.T) {
	assert.Nil(t, Intersect(nil, nil))
	assert.Nil(t, Intersect(map[string][]byte{"a": {}}, nil))
	assert.Nil(t, Intersect(nil, map[string][]byte{"a": {}}))
}

func TestIntersect_NoCommonKeys_ReturnsNil(t *testing.T) {
	src := map[string][]byte{"a": []byte("1")}
	dst := map[string][]byte{"b": []byte("2")}
	assert.Nil(t, Intersect(src, dst))
}

func TestIntersect_CommonKeys_Returned(t *testing.T) {
	src := map[string][]byte{"a": []byte("1"), "b": []byte("2"), "c": []byte("3")}
	dst := map[string][]byte{"b": []byte("2"), "c": []byte("changed")}

	results := Intersect(src, dst)
	assert.Len(t, results, 2)
	assert.Equal(t, "b", results[0].Key)
	assert.True(t, results[0].Equal)
	assert.Equal(t, "c", results[1].Key)
	assert.False(t, results[1].Equal)
}

func TestIntersect_SortedByKey(t *testing.T) {
	src := map[string][]byte{"z": []byte("1"), "a": []byte("2"), "m": []byte("3")}
	dst := map[string][]byte{"z": []byte("1"), "a": []byte("2"), "m": []byte("3")}

	results := Intersect(src, dst)
	assert.Equal(t, "a", results[0].Key)
	assert.Equal(t, "m", results[1].Key)
	assert.Equal(t, "z", results[2].Key)
}

func TestIntersectKeys_ReturnsOnlyKeys(t *testing.T) {
	src := map[string][]byte{"x": []byte("v1"), "y": []byte("v2")}
	dst := map[string][]byte{"y": []byte("v2"), "z": []byte("v3")}

	keys := IntersectKeys(src, dst)
	assert.Equal(t, []string{"y"}, keys)
}
