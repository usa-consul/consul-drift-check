package diff

import (
	"testing"
	"time"
)

func makeStaleResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: "modified"}
	}
	return out
}

func TestDetectStale_EmptyInput_ReturnsNil(t *testing.T) {
	got := DetectStale(nil, nil, time.Now(), StaleOptions{MaxAge: time.Hour})
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestDetectStale_NoMaxAge_NeverStale(t *testing.T) {
	now := time.Now()
	seenAt := map[string]time.Time{
		"app/db": now.Add(-48 * time.Hour),
	}
	results := makeStaleResults("app/db")
	got := DetectStale(results, seenAt, now, StaleOptions{})
	if got[0].IsStale {
		t.Fatal("expected not stale when MaxAge is zero")
	}
}

func TestDetectStale_OldKey_MarkedStale(t *testing.T) {
	now := time.Now()
	seenAt := map[string]time.Time{
		"app/db": now.Add(-25 * time.Hour),
	}
	results := makeStaleResults("app/db")
	got := DetectStale(results, seenAt, now, StaleOptions{MaxAge: 24 * time.Hour})
	if !got[0].IsStale {
		t.Fatal("expected stale")
	}
	if got[0].Age < 24*time.Hour {
		t.Fatalf("unexpected age: %v", got[0].Age)
	}
}

func TestDetectStale_FreshKey_NotStale(t *testing.T) {
	now := time.Now()
	seenAt := map[string]time.Time{
		"app/cache": now.Add(-1 * time.Hour),
	}
	results := makeStaleResults("app/cache")
	got := DetectStale(results, seenAt, now, StaleOptions{MaxAge: 24 * time.Hour})
	if got[0].IsStale {
		t.Fatal("expected not stale")
	}
}

func TestDetectStale_ExcludedPrefix_SkipsCheck(t *testing.T) {
	now := time.Now()
	seenAt := map[string]time.Time{
		"infra/network": now.Add(-100 * time.Hour),
	}
	results := makeStaleResults("infra/network")
	got := DetectStale(results, seenAt, now, StaleOptions{
		MaxAge:          24 * time.Hour,
		ExcludePrefixes: []string{"infra/"},
	})
	if got[0].IsStale {
		t.Fatal("excluded prefix should not be marked stale")
	}
}

func TestDetectStale_OnlyStatus_FiltersOtherStatuses(t *testing.T) {
	now := time.Now()
	seenAt := map[string]time.Time{
		"app/key": now.Add(-48 * time.Hour),
	}
	results := []Result{{Key: "app/key", Status: "ok"}}
	got := DetectStale(results, seenAt, now, StaleOptions{
		MaxAge:     24 * time.Hour,
		OnlyStatus: "modified",
	})
	if got[0].IsStale {
		t.Fatal("status mismatch should not be flagged stale")
	}
}

func TestDetectStale_UnseenKey_UsesReferenceTime(t *testing.T) {
	now := time.Now()
	results := makeStaleResults("new/key")
	got := DetectStale(results, map[string]time.Time{}, now, StaleOptions{MaxAge: time.Millisecond})
	// age == 0 which is < 1ms, so not stale
	if got[0].IsStale {
		t.Fatal("unseen key at reference time should not be stale")
	}
}
