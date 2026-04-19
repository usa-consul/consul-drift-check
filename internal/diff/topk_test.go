package diff

import (
	"testing"
)

func makeTopKResults() []WeightedResult {
	return []WeightedResult{
		{Key: "app/db/host", Status: "modified", Weight: 9.0, Source: []byte("old"), Dest: []byte("new")},
		{Key: "app/db/port", Status: "modified", Weight: 4.5, Source: []byte("5432"), Dest: []byte("5433")},
		{Key: "app/cache/ttl", Status: "only_in_source", Weight: 7.2, Source: []byte("300"), Dest: nil},
		{Key: "infra/region", Status: "only_in_dest", Weight: 2.1, Source: nil, Dest: []byte("us-east-1")},
		{Key: "app/secret", Status: "modified", Weight: 11.0, Source: []byte("x"), Dest: []byte("y")},
	}
}

func TestTopK_DefaultK_ReturnsAll(t *testing.T) {
	input := makeTopKResults()[:3]
	got := TopK(input, TopKOptions{})
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
}

func TestTopK_LimitsToK(t *testing.T) {
	got := TopK(makeTopKResults(), TopKOptions{K: 2})
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestTopK_OrderByWeight_HighestFirst(t *testing.T) {
	got := TopK(makeTopKResults(), TopKOptions{K: 3})
	if got[0].Key != "app/secret" {
		t.Errorf("expected app/secret first, got %s", got[0].Key)
	}
	if got[1].Key != "app/db/host" {
		t.Errorf("expected app/db/host second, got %s", got[1].Key)
	}
}

func TestTopK_OrderByKey_Alphabetical(t *testing.T) {
	got := TopK(makeTopKResults(), TopKOptions{K: 5, OrderBy: "key"})
	if got[0].Key != "app/cache/ttl" {
		t.Errorf("expected app/cache/ttl first, got %s", got[0].Key)
	}
}

func TestTopK_RankIsOneBased(t *testing.T) {
	got := TopK(makeTopKResults(), TopKOptions{K: 3})
	for i, r := range got {
		if r.Rank != i+1 {
			t.Errorf("rank mismatch at index %d: got %d", i, r.Rank)
		}
	}
}

func TestTopK_EmptyInput_ReturnsEmpty(t *testing.T) {
	got := TopK(nil, TopKOptions{K: 5})
	if len(got) != 0 {
		t.Errorf("expected empty, got %d", len(got))
	}
}

func TestTopK_KLargerThanInput_ReturnsAll(t *testing.T) {
	got := TopK(makeTopKResults(), TopKOptions{K: 100})
	if len(got) != len(makeTopKResults()) {
		t.Errorf("expected %d results, got %d", len(makeTopKResults()), len(got))
	}
}
