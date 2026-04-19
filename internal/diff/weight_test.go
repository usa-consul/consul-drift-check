package diff

import (
	"testing"
)

func makeWeightResults() []Result {
	return []Result{
		{Key: "service/web/port", Status: StatusModified},
		{Key: "config/debug", Status: StatusOnlyInSource},
		{Key: "service/db/pass", Status: StatusOnlyInDestination},
		{Key: "service/db/host", Status: StatusModified},
	}
}

func TestApply_DefaultOptions_ScoresAssigned(t *testing.T) {
	results := makeWeightResults()
	opts := DefaultWeightOptions()
	weighted := Apply(results, opts)

	if len(weighted) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(weighted))
	}
	for _, w := range weighted {
		if w.Weight <= 0 {
			t.Errorf("key %q has non-positive weight %f", w.Key, w.Weight)
		}
	}
}

func TestApply_SortedDescending(t *testing.T) {
	results := makeWeightResults()
	opts := DefaultWeightOptions()
	weighted := Apply(results, opts)

	for i := 1; i < len(weighted); i++ {
		if weighted[i].Weight > weighted[i-1].Weight {
			t.Errorf("results not sorted: index %d weight %f > index %d weight %f",
				i, weighted[i].Weight, i-1, weighted[i-1].Weight)
		}
	}
}

func TestApply_PrefixMultiplier_IncreasesWeight(t *testing.T) {
	results := []Result{
		{Key: "service/web/port", Status: StatusModified},
		{Key: "config/debug", Status: StatusModified},
	}
	opts := DefaultWeightOptions()
	opts.PrefixWeights = map[string]float64{"service/": 3.0}

	weighted := Apply(results, opts)

	var svcWeight, cfgWeight float64
	for _, w := range weighted {
		if w.Key == "service/web/port" {
			svcWeight = w.Weight
		} else {
			cfgWeight = w.Weight
		}
	}
	if svcWeight <= cfgWeight {
		t.Errorf("expected service key weight (%f) > config key weight (%f)", svcWeight, cfgWeight)
	}
}

func TestApply_EmptyResults_ReturnsEmpty(t *testing.T) {
	weighted := Apply(nil, DefaultWeightOptions())
	if len(weighted) != 0 {
		t.Errorf("expected empty, got %d", len(weighted))
	}
}

func TestApply_NoMatchingPrefix_MultiplierIsOne(t *testing.T) {
	results := []Result{
		{Key: "other/key", Status: StatusModified},
	}
	opts := DefaultWeightOptions()
	opts.PrefixWeights = map[string]float64{"service/": 5.0}
	weighted := Apply(results, opts)

	if weighted[0].Weight != opts.ModifiedScore {
		t.Errorf("expected weight %f, got %f", opts.ModifiedScore, weighted[0].Weight)
	}
}
