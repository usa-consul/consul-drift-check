package diff

import (
	"testing"
)

func makeMaskResults() []Result {
	return []Result{
		{Key: "app/db/password", SourceValue: []byte("secret"), DestValue: []byte("other"), Status: "modified"},
		{Key: "app/db/host", SourceValue: []byte("localhost"), DestValue: []byte("prod-db"), Status: "modified"},
		{Key: "secrets/token", SourceValue: []byte("tok123"), DestValue: nil, Status: "only_in_source"},
		{Key: "config/timeout", SourceValue: []byte("30s"), DestValue: []byte("30s"), Status: "ok"},
	}
}

func TestMask_NoSensitivePrefixes_ReturnsUnchanged(t *testing.T) {
	results := makeMaskResults()
	out := Mask(results, MaskOptions{})
	if len(out) != len(results) {
		t.Fatalf("expected %d results, got %d", len(results), len(out))
	}
	if out[0].SourceValue != "secret" {
		t.Errorf("expected source value unchanged, got %q", out[0].SourceValue)
	}
}

func TestMask_MatchingPrefix_MasksValues(t *testing.T) {
	results := makeMaskResults()
	out := Mask(results, MaskOptions{SensitivePrefixes: []string{"app/db/password", "secrets/"}})

	if out[0].SourceValue != "***" {
		t.Errorf("expected masked source, got %q", out[0].SourceValue)
	}
	if out[0].DestValue != "***" {
		t.Errorf("expected masked dest, got %q", out[0].DestValue)
	}
	if out[2].SourceValue != "***" {
		t.Errorf("expected masked secrets token, got %q", out[2].SourceValue)
	}
}

func TestMask_NonMatchingKey_Unchanged(t *testing.T) {
	results := makeMaskResults()
	out := Mask(results, MaskOptions{SensitivePrefixes: []string{"secrets/"}})
	if out[1].SourceValue != "localhost" {
		t.Errorf("expected unchanged host, got %q", out[1].SourceValue)
	}
	if out[3].SourceValue != "30s" {
		t.Errorf("expected unchanged timeout, got %q", out[3].SourceValue)
	}
}

func TestMask_CustomMask_UsesCustomString(t *testing.T) {
	results := makeMaskResults()
	out := Mask(results, MaskOptions{SensitivePrefixes: []string{"app/db/password"}, Mask: "[REDACTED]"})
	if out[0].SourceValue != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out[0].SourceValue)
	}
}

func TestMask_EmptyValue_NotMasked(t *testing.T) {
	results := []Result{
		{Key: "secrets/empty", SourceValue: nil, DestValue: []byte("val"), Status: "only_in_dest"},
	}
	out := Mask(results, MaskOptions{SensitivePrefixes: []string{"secrets/"}})
	if out[0].SourceValue != "" {
		t.Errorf("expected empty source to stay empty, got %q", out[0].SourceValue)
	}
	if out[0].DestValue != "***" {
		t.Errorf("expected dest masked, got %q", out[0].DestValue)
	}
}

func TestMask_CaseInsensitivePrefix(t *testing.T) {
	results := []Result{
		{Key: "Secrets/APIKey", SourceValue: []byte("key"), DestValue: []byte("key2"), Status: "modified"},
	}
	out := Mask(results, MaskOptions{SensitivePrefixes: []string{"secrets/"}})
	if out[0].SourceValue != "***" {
		t.Errorf("expected case-insensitive mask, got %q", out[0].SourceValue)
	}
}
