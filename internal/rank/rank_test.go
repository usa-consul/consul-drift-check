package rank_test

import (
	"testing"

	"github.com/example/consul-drift-check/internal/diff"
	"github.com/example/consul-drift-check/internal/rank"
)

func makeResults() []diff.Result {
	return []diff.Result{
		{Key: "b/key", Status: diff.StatusMatch},
		{Key: "a/key", Status: diff.StatusOnlyInSrc},
		{Key: "c/key", Status: diff.StatusModified},
		{Key: "d/key", Status: diff.StatusOnlyInDst},
	}
}

func TestApply_OrderBySeverity_ModifiedFirst(t *testing.T) {
	out := rank.Apply(makeResults(), rank.OrderSeverity)
	if out[0].Status != diff.StatusModified {
		t.Fatalf("expected modified first, got %s", out[0].Status)
	}
}

func TestApply_OrderBySeverity_MatchLast(t *testing.T) {
	out := rank.Apply(makeResults(), rank.OrderSeverity)
	if out[len(out)-1].Status != diff.StatusMatch {
		t.Fatalf("expected match last, got %s", out[len(out)-1].Status)
	}
}

func TestApply_OrderByKey_Alphabetical(t *testing.T) {
	out := rank.Apply(makeResults(), rank.OrderKey)
	for i := 1; i < len(out); i++ {
		if out[i].Key < out[i-1].Key {
			t.Fatalf("keys not sorted: %s before %s", out[i-1].Key, out[i].Key)
		}
	}
}

func TestApply_OrderByStatus_Alphabetical(t *testing.T) {
	out := rank.Apply(makeResults(), rank.OrderStatus)
	for i := 1; i < len(out); i++ {
		if out[i].Status < out[i-1].Status {
			t.Fatalf("statuses not sorted at index %d", i)
		}
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	input := makeResults()
	first := input[0].Key
	rank.Apply(input, rank.OrderKey)
	if input[0].Key != first {
		t.Fatal("original slice was mutated")
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := rank.Apply([]diff.Result{}, rank.OrderSeverity)
	if len(out) != 0 {
		t.Fatalf("expected empty, got %d results", len(out))
	}
}
