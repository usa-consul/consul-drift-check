package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeSupersedResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified"}
	}
	return out
}

func supersedKeys(rs []Result) []string {
	keys := make([]string, len(rs))
	for i, r := range rs {
		keys[i] = r.Key
	}
	return keys
}

func TestSupersede_NoPrefixes_ReturnsAll(t *testing.T) {
	in := makeSupersedResults("a/x", "a/y", "b/z")
	out := Supersede(in, SupersedeOptions{})
	assert.Equal(t, supersedKeys(in), supersedKeys(out))
}

func TestSupersede_EmptyInput_ReturnsNil(t *testing.T) {
	out := Supersede(nil, SupersedeOptions{Prefixes: []string{"config/app"}})
	assert.Nil(t, out)
}

func TestSupersede_PriorityPrefixDropsOtherSegmentEntries(t *testing.T) {
	in := makeSupersedResults(
		"config/app/db",   // matches priority prefix — keep
		"config/legacy/x", // same top segment, no match — drop
		"other/key",       // different segment — keep
	)
	out := Supersede(in, SupersedeOptions{Prefixes: []string{"config/app"}})
	assert.Equal(t, []string{"config/app/db", "other/key"}, supersedKeys(out))
}

func TestSupersede_NoMatchingSegment_KeepsAll(t *testing.T) {
	in := makeSupersedResults("service/a", "service/b")
	out := Supersede(in, SupersedeOptions{Prefixes: []string{"config/app"}})
	assert.Equal(t, []string{"service/a", "service/b"}, supersedKeys(out))
}

func TestSupersede_CaseInsensitive_MatchesMixedCase(t *testing.T) {
	in := makeSupersedResults(
		"Config/App/key",
		"Config/Other/key",
	)
	out := Supersede(in, SupersedeOptions{
		Prefixes:      []string{"config/app"},
		CaseSensitive: false,
	})
	assert.Equal(t, []string{"Config/App/key"}, supersedKeys(out))
}

func TestSupersede_CaseSensitive_NoFold(t *testing.T) {
	in := makeSupersedResults(
		"Config/App/key",
		"Config/Other/key",
	)
	// With case-sensitive matching, "config/app" does NOT match "Config/App".
	out := Supersede(in, SupersedeOptions{
		Prefixes:      []string{"config/app"},
		CaseSensitive: true,
	})
	// No priority segment activated, so all results are kept.
	assert.Equal(t, []string{"Config/App/key", "Config/Other/key"}, supersedKeys(out))
}

func TestSupersede_ResultsSortedByKey(t *testing.T) {
	in := makeSupersedResults("z/app/c", "a/other/x", "z/app/a", "z/legacy/b")
	out := Supersede(in, SupersedeOptions{Prefixes: []string{"z/app"}})
	assert.Equal(t, []string{"a/other/x", "z/app/a", "z/app/c"}, supersedKeys(out))
}
