package diff

import (
	"strings"
	"testing"
)

func makeClampResults(pairs [][2]string) []Result {
	out := make([]Result, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, Result{
			Key:         p[0],
			SourceValue: []byte(p[1]),
			DestValue:   []byte(p[1] + "_dest"),
			Status:      "modified",
		})
	}
	return out
}

func TestClamp_ZeroMax_ReturnsAll(t *testing.T) {
	results := makeClampResults([][2]string{{"a/key", "value"}})
	out := Clamp(results, ClampOptions{})
	if len(out) != 1 {
		t.Fatalf("expected 1 result, got %d", len(out))
	}
	if string(out[0].SourceValue) != "value" {
		t.Errorf("expected value unchanged, got %q", out[0].SourceValue)
	}
}

func TestClamp_ShortValues_NotTruncated(t *testing.T) {
	results := makeClampResults([][2]string{{"a/key", "hi"}})
	out := Clamp(results, ClampOptions{MaxValueLen: 10})
	if string(out[0].SourceValue) != "hi" {
		t.Errorf("expected 'hi', got %q", out[0].SourceValue)
	}
}

func TestClamp_LongValue_TruncatedWithDefaultSuffix(t *testing.T) {
	long := strings.Repeat("x", 20)
	results := makeClampResults([][2]string{{"a/key", long}})
	out := Clamp(results, ClampOptions{MaxValueLen: 10})
	got := string(out[0].SourceValue)
	if len([]byte(got)) > 10 {
		t.Errorf("expected at most 10 bytes, got %d: %q", len(got), got)
	}
	if !strings.HasSuffix(got, "\u2026") {
		t.Errorf("expected default suffix '…', got %q", got)
	}
}

func TestClamp_LongValue_TruncatedWithCustomSuffix(t *testing.T) {
	long := strings.Repeat("y", 30)
	results := makeClampResults([][2]string{{"b/key", long}})
	out := Clamp(results, ClampOptions{MaxValueLen: 8, Suffix: "..."})
	got := string(out[0].SourceValue)
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected custom suffix '...', got %q", got)
	}
	if len(got) > 8 {
		t.Errorf("expected at most 8 bytes, got %d", len(got))
	}
}

func TestClamp_EmptyInput_ReturnsNil(t *testing.T) {
	out := Clamp(nil, ClampOptions{MaxValueLen: 5})
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d", len(out))
	}
}

func TestClamp_ResultsSortedByKey(t *testing.T) {
	results := makeClampResults([][2]string{
		{"z/last", "val"},
		{"a/first", "val"},
		{"m/mid", "val"},
	})
	out := Clamp(results, ClampOptions{MaxValueLen: 100})
	if out[0].Key != "a/first" || out[1].Key != "m/mid" || out[2].Key != "z/last" {
		t.Errorf("results not sorted: %v", []string{out[0].Key, out[1].Key, out[2].Key})
	}
}
