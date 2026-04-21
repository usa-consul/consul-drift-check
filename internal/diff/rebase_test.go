package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeRebaseResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified", SourceValue: []byte("a"), DestValue: []byte("b")}
	}
	return out
}

func rebaseKeys(results []Result) []string {
	keys := make([]string, len(results))
	for i, r := range results {
		keys[i] = r.Key
	}
	return keys
}

func TestRebase_EmptyInput_ReturnsNil(t *testing.T) {
	result := Rebase(nil, RebaseOptions{})
	assert.Nil(t, result)
}

func TestRebase_NoOptions_ReturnsSortedCopy(t *testing.T) {
	in := makeRebaseResults("z/key", "a/key", "m/key")
	out := Rebase(in, RebaseOptions{})
	require.Len(t, out, 3)
	assert.Equal(t, []string{"a/key", "m/key", "z/key"}, rebaseKeys(out))
}

func TestRebase_StripPrefix_RemovesLeadingSegment(t *testing.T) {
	in := makeRebaseResults("prod/svc/config", "prod/svc/secret")
	out := Rebase(in, RebaseOptions{StripPrefix: "prod/svc"})
	require.Len(t, out, 2)
	assert.Equal(t, []string{"config", "secret"}, rebaseKeys(out))
}

func TestRebase_NewPrefix_PrependedToKeys(t *testing.T) {
	in := makeRebaseResults("config", "secret")
	out := Rebase(in, RebaseOptions{NewPrefix: "staging/svc"})
	require.Len(t, out, 2)
	assert.Equal(t, []string{"staging/svc/config", "staging/svc/secret"}, rebaseKeys(out))
}

func TestRebase_StripAndReplace_CombinesCorrectly(t *testing.T) {
	in := makeRebaseResults("prod/svc/config", "prod/svc/db")
	out := Rebase(in, RebaseOptions{StripPrefix: "prod/svc", NewPrefix: "staging/svc"})
	require.Len(t, out, 2)
	assert.Equal(t, []string{"staging/svc/config", "staging/svc/db"}, rebaseKeys(out))
}

func TestRebase_SkipUnmatched_DropsNonMatchingKeys(t *testing.T) {
	in := makeRebaseResults("prod/svc/config", "other/key")
	out := Rebase(in, RebaseOptions{StripPrefix: "prod/svc", SkipUnmatched: true})
	require.Len(t, out, 1)
	assert.Equal(t, "config", out[0].Key)
}

func TestRebase_SkipUnmatched_False_KeepsNonMatchingKeys(t *testing.T) {
	in := makeRebaseResults("prod/svc/config", "other/key")
	out := Rebase(in, RebaseOptions{StripPrefix: "prod/svc", SkipUnmatched: false})
	require.Len(t, out, 2)
}

func TestRebase_PreservesStatusAndValues(t *testing.T) {
	in := []Result{{Key: "prod/k", Status: "only_in_source", SourceValue: []byte("v")}}
	out := Rebase(in, RebaseOptions{StripPrefix: "prod"})
	require.Len(t, out, 1)
	assert.Equal(t, "only_in_source", out[0].Status)
	assert.Equal(t, []byte("v"), out[0].SourceValue)
}
