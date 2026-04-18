package validate_test

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/nicholasgasior/consul-drift-check/internal/validate"
)

func makePairs(kvs map[string]string) []*api.KVPair {
	pairs := make([]*api.KVPair, 0, len(kvs))
	for k, v := range kvs {
		pairs = append(pairs, &api.KVPair{Key: k, Value: []byte(v)})
	}
	return pairs
}

func TestApply_NoRules_ReturnsNoViolations(t *testing.T) {
	pairs := makePairs(map[string]string{"app/key": "value"})
	v := validate.Apply(pairs, nil)
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestApply_RequiredRule_EmptyValue(t *testing.T) {
	pairs := makePairs(map[string]string{"app/key": ""})
	rules := []validate.Rule{{Required: true, Severity: validate.SeverityError}}
	v := validate.Apply(pairs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Severity != validate.SeverityError {
		t.Errorf("expected error severity, got %s", v[0].Severity)
	}
}

func TestApply_RequiredRule_NonEmptyValue_NoViolation(t *testing.T) {
	pairs := makePairs(map[string]string{"app/key": "present"})
	rules := []validate.Rule{{Required: true}}
	v := validate.Apply(pairs, rules)
	if len(v) != 0 {
		t.Fatalf("expected 0 violations, got %d", len(v))
	}
}

func TestApply_MaxLength_Exceeded(t *testing.T) {
	pairs := makePairs(map[string]string{"cfg/token": "toolongvalue"})
	rules := []validate.Rule{{MaxLength: 5, Severity: validate.SeverityWarning}}
	v := validate.Apply(pairs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Message != "value exceeds maximum length" {
		t.Errorf("unexpected message: %s", v[0].Message)
	}
}

func TestApply_PrefixFilter_SkipsNonMatchingKeys(t *testing.T) {
	pairs := makePairs(map[string]string{"other/key": ""})
	rules := []validate.Rule{{Prefix: "app/", Required: true, Severity: validate.SeverityError}}
	v := validate.Apply(pairs, rules)
	if len(v) != 0 {
		t.Fatalf("expected 0 violations for non-matching prefix, got %d", len(v))
	}
}

func TestApply_DefaultSeverity_IsWarning(t *testing.T) {
	pairs := makePairs(map[string]string{"k": ""})
	rules := []validate.Rule{{Required: true}}
	v := validate.Apply(pairs, rules)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation")
	}
	if v[0].Severity != validate.SeverityWarning {
		t.Errorf("expected warning, got %s", v[0].Severity)
	}
}

func TestApply_NilPair_Skipped(t *testing.T) {
	pairs := []*api.KVPair{nil}
	rules := []validate.Rule{{Required: true}}
	v := validate.Apply(pairs, rules)
	if len(v) != 0 {
		t.Fatalf("nil pair should be skipped")
	}
}
