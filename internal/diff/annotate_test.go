package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeAnnotateResults() []Result {
	return []Result{
		{Key: "service/api/timeout", Status: "modified"},
		{Key: "infra/network/cidr", Status: "only_in_source"},
		{Key: "unknown/key", Status: "only_in_destination"},
		{Key: "service/db/pool", Status: "modified"},
	}
}

func TestAnnotate_EmptyResults_ReturnsNil(t *testing.T) {
	out := Annotate(nil, AnnotateOptions{})
	assert.Nil(t, out)
}

func TestAnnotate_MatchingPrefix_AssignsLabel(t *testing.T) {
	results := makeAnnotateResults()
	opts := AnnotateOptions{
		PrefixLabels: map[string]string{
			"service": "services",
			"infra":   "infrastructure",
		},
		DefaultLabel: "general",
	}
	out := Annotate(results, opts)
	assert.Len(t, out, 4)
	assert.Equal(t, "services", out[0].Label)
	assert.Equal(t, "infrastructure", out[1].Label)
	assert.Equal(t, "general", out[2].Label)
	assert.Equal(t, "services", out[3].Label)
}

func TestAnnotate_NoRules_UsesDefaultLabel(t *testing.T) {
	results := makeAnnotateResults()
	opts := AnnotateOptions{DefaultLabel: "fallback"}
	out := Annotate(results, opts)
	for _, r := range out {
		assert.Equal(t, "fallback", r.Label)
	}
}

func TestAnnotate_EmptyDefaultLabel_ReturnsEmptyString(t *testing.T) {
	results := []Result{{Key: "x/y", Status: "modified"}}
	out := Annotate(results, AnnotateOptions{})
	assert.Equal(t, "", out[0].Label)
}

func TestAnnotate_PreservesOriginalResult(t *testing.T) {
	results := []Result{{Key: "service/a", Status: "modified", SourceValue: []byte("v1"), DestValue: []byte("v2")}}
	opts := AnnotateOptions{PrefixLabels: map[string]string{"service": "svc"}, DefaultLabel: "x"}
	out := Annotate(results, opts)
	assert.Equal(t, results[0], out[0].Result)
}
