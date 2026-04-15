package diff

import (
	"testing"
)

func TestCompare_NoDrift(t *testing.T) {
	src := map[string][]byte{
		"config/app/port": []byte("8080"),
		"config/app/host": []byte("localhost"),
	}
	dst := map[string][]byte{
		"config/app/port": []byte("8080"),
		"config/app/host": []byte("localhost"),
	}

	result := Compare(src, dst)

	if result.HasDrift() {
		t.Errorf("expected no drift, got: %+v", result)
	}
}

func TestCompare_OnlyInSource(t *testing.T) {
	src := map[string][]byte{
		"config/app/port": []byte("8080"),
		"config/app/debug": []byte("true"),
	}
	dst := map[string][]byte{
		"config/app/port": []byte("8080"),
	}

	result := Compare(src, dst)

	if len(result.OnlyInSource) != 1 {
		t.Fatalf("expected 1 key only in source, got %d", len(result.OnlyInSource))
	}
	if result.OnlyInSource[0].Key != "config/app/debug" {
		t.Errorf("unexpected key: %s", result.OnlyInSource[0].Key)
	}
}

func TestCompare_OnlyInDestination(t *testing.T) {
	src := map[string][]byte{
		"config/app/port": []byte("8080"),
	}
	dst := map[string][]byte{
		"config/app/port": []byte("8080"),
		"config/app/timeout": []byte("30s"),
	}

	result := Compare(src, dst)

	if len(result.OnlyInDestination) != 1 {
		t.Fatalf("expected 1 key only in destination, got %d", len(result.OnlyInDestination))
	}
	if result.OnlyInDestination[0].Key != "config/app/timeout" {
		t.Errorf("unexpected key: %s", result.OnlyInDestination[0].Key)
	}
}

func TestCompare_ModifiedValues(t *testing.T) {
	src := map[string][]byte{
		"config/app/port": []byte("8080"),
	}
	dst := map[string][]byte{
		"config/app/port": []byte("9090"),
	}

	result := Compare(src, dst)

	if len(result.Modified) != 1 {
		t.Fatalf("expected 1 modified key, got %d", len(result.Modified))
	}
	if result.Modified[0].Key != "config/app/port" {
		t.Errorf("unexpected key: %s", result.Modified[0].Key)
	}
}

func TestCompare_EmptyMaps(t *testing.T) {
	result := Compare(map[string][]byte{}, map[string][]byte{})
	if result.HasDrift() {
		t.Error("expected no drift for empty maps")
	}
}
