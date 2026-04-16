package redact_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/redact"
)

func makePairs(kv ...string) []redact.KVPair {
	pairs := make([]redact.KVPair, 0, len(kv)/2)
	for i := 0; i+1 < len(kv); i += 2 {
		pairs = append(pairs, redact.KVPair{Key: kv[i], Value: []byte(kv[i+1])})
	}
	return pairs
}

func TestApply_NoSensitiveKeys_ReturnsUnchanged(t *testing.T) {
	pairs := makePairs("config/host", "localhost", "config/port", "8500")
	out := redact.Apply(pairs, redact.Options{})
	for i, p := range out {
		if string(p.Value) != string(pairs[i].Value) {
			t.Fatalf("expected unchanged value for %s", p.Key)
		}
	}
}

func TestApply_ExactKeyMatch_MasksValue(t *testing.T) {
	pairs := makePairs("config/password", "secret", "config/host", "localhost")
	out := redact.Apply(pairs, redact.Options{SensitiveKeys: []string{"config/password"}})
	if string(out[0].Value) != "***" {
		t.Fatalf("expected masked value, got %s", out[0].Value)
	}
	if string(out[1].Value) != "localhost" {
		t.Fatalf("expected unchanged value for host")
	}
}

func TestApply_SuffixMatch_MasksValue(t *testing.T) {
	pairs := makePairs("db/token", "abc123", "app/token", "xyz")
	out := redact.Apply(pairs, redact.Options{SensitiveKeys: []string{"token"}})
	for _, p := range out {
		if string(p.Value) != "***" {
			t.Fatalf("expected masked value for %s", p.Key)
		}
	}
}

func TestApply_CustomMask(t *testing.T) {
	pairs := makePairs("secret", "value")
	out := redact.Apply(pairs, redact.Options{SensitiveKeys: []string{"secret"}, Mask: "[REDACTED]"})
	if string(out[0].Value) != "[REDACTED]" {
		t.Fatalf("expected custom mask, got %s", out[0].Value)
	}
}

func TestApply_CaseInsensitiveKeyMatch(t *testing.T) {
	pairs := makePairs("Config/Password", "topsecret")
	out := redact.Apply(pairs, redact.Options{SensitiveKeys: []string{"config/password"}})
	if string(out[0].Value) != "***" {
		t.Fatalf("expected masked value")
	}
}

func TestApply_OriginalSliceUnmodified(t *testing.T) {
	pairs := makePairs("key", "value")
	redact.Apply(pairs, redact.Options{SensitiveKeys: []string{"key"}})
	if string(pairs[0].Value) != "value" {
		t.Fatal("original slice must not be modified")
	}
}
