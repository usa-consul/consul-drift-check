package diff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeRollupResults() []Result {
	return []Result{
		{Key: "app/db/host", Status: StatusModified},
		{Key: "app/db/port", Status: StatusOnlyInSource},
		{Key: "app/web/url", Status: StatusOnlyInDestination},
		{Key: "infra/vpc", Status: StatusModified},
		{Key: "infra/subnet", Status: StatusModified},
	}
}

func TestRollup_GroupsByTopSegment(t *testing.T) {
	results := makeRollupResults()
	out := Rollup(results)
	assert.Len(t, out, 2)
	assert.Equal(t, "app", out[0].Prefix)
	assert.Equal(t, "infra", out[1].Prefix)
}

func TestRollup_CountsCorrectly(t *testing.T) {
	results := makeRollupResults()
	out := Rollup(results)

	app := out[0]
	assert.Equal(t, 1, app.Added)
	assert.Equal(t, 1, app.Removed)
	assert.Equal(t, 1, app.Modified)
	assert.Equal(t, 3, app.Total)

	infra := out[1]
	assert.Equal(t, 0, infra.Added)
	assert.Equal(t, 0, infra.Removed)
	assert.Equal(t, 2, infra.Modified)
	assert.Equal(t, 2, infra.Total)
}

func TestRollup_EmptyResults(t *testing.T) {
	out := Rollup(nil)
	assert.Empty(t, out)
}

func TestRollup_LeadingSlashStripped(t *testing.T) {
	results := []Result{
		{Key: "/service/config", Status: StatusModified},
	}
	out := Rollup(results)
	assert.Len(t, out, 1)
	assert.Equal(t, "service", out[0].Prefix)
}

func TestRollup_FlatKey_UsesWholeKey(t *testing.T) {
	results := []Result{
		{Key: "standalone", Status: StatusOnlyInSource},
	}
	out := Rollup(results)
	assert.Len(t, out, 1)
	assert.Equal(t, "standalone", out[0].Prefix)
	assert.Equal(t, 1, out[0].Added)
}
