package window_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/window"
)

func makeResults(keys ...string) []diff.Result {
	out := make([]diff.Result, len(keys))
	for i, k := range keys {
		out[i] = diff.Result{Key: k, Status: diff.Modified}
	}
	return out
}

func TestNew_ZeroSize_DefaultsToOne(t *testing.T) {
	w := window.New(0)
	w.Push(makeResults("a"))
	w.Push(makeResults("b"))
	if w.Len() != 1 {
		t.Fatalf("expected 1 slot, got %d", w.Len())
	}
}

func TestPush_RespectsWindowSize(t *testing.T) {
	w := window.New(3)
	for i := 0; i < 5; i++ {
		w.Push(makeResults("k"))
	}
	if w.Len() != 3 {
		t.Fatalf("expected 3 slots, got %d", w.Len())
	}
}

func TestSlots_ReturnsCopy(t *testing.T) {
	w := window.New(2)
	w.Push(makeResults("a", "b"))
	slots := w.Slots()
	slots[0][0].Key = "mutated"
	slots2 := w.Slots()
	if slots2[0][0].Key == "mutated" {
		t.Fatal("Slots should return independent copies")
	}
}

func TestPersistentKeys_AllSlots(t *testing.T) {
	w := window.New(3)
	w.Push(makeResults("alpha", "beta"))
	w.Push(makeResults("alpha", "gamma"))
	w.Push(makeResults("alpha"))
	keys := w.PersistentKeys()
	if len(keys) != 1 || keys[0] != "alpha" {
		t.Fatalf("expected [alpha], got %v", keys)
	}
}

func TestPersistentKeys_EmptyWindow(t *testing.T) {
	w := window.New(3)
	if keys := w.PersistentKeys(); keys != nil {
		t.Fatalf("expected nil, got %v", keys)
	}
}

func TestLen_Empty(t *testing.T) {
	w := window.New(5)
	if w.Len() != 0 {
		t.Fatalf("expected 0, got %d", w.Len())
	}
}
