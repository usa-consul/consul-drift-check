package pivot_test

import (
	"testing"

	"github.com/your-org/consul-drift-check/internal/diff"
	"github.com/your-org/consul-drift-check/internal/pivot"
)

func makeResults(pairs ...string) []diff.Result {
	var out []diff.Result
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, diff.Result{
			Key:              pairs[i],
			Status:           diff.Modified,
			DestinationValue: pairs[i+1],
		})
	}
	return out
}

func TestBuild_EmptyInput_ReturnsEmptyTable(t *testing.T) {
	table := pivot.Build(nil)
	if len(table) != 0 {
		t.Fatalf("expected empty table, got %d rows", len(table))
	}
}

func TestBuild_SingleDC_RowsMatchKeys(t *testing.T) {
	results := makeResults("app/port", "8080", "app/host", "localhost")
	table := pivot.Build(map[string][]diff.Result{"dc1": results})

	if len(table) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(table))
	}
	if table[0].Key != "app/host" {
		t.Errorf("expected first key app/host, got %s", table[0].Key)
	}
	if table[0].Values["dc1"] != "localhost" {
		t.Errorf("unexpected value: %s", table[0].Values["dc1"])
	}
}

func TestBuild_MultipleDCs_ValuesGroupedByKey(t *testing.T) {
	dc1 := makeResults("cfg/timeout", "30s")
	dc2 := makeResults("cfg/timeout", "60s")

	table := pivot.Build(map[string][]diff.Result{
		"dc1": dc1,
		"dc2": dc2,
	})

	if len(table) != 1 {
		t.Fatalf("expected 1 row, got %d", len(table))
	}
	row := table[0]
	if row.Values["dc1"] != "30s" {
		t.Errorf("dc1 value mismatch: %s", row.Values["dc1"])
	}
	if row.Values["dc2"] != "60s" {
		t.Errorf("dc2 value mismatch: %s", row.Values["dc2"])
	}
}

func TestDCs_ReturnsSortedLabels(t *testing.T) {
	dc1 := makeResults("k", "v1")
	dc2 := makeResults("k", "v2")
	dc3 := makeResults("k", "v3")

	table := pivot.Build(map[string][]diff.Result{
		"us-east": dc1,
		"eu-west": dc2,
		"ap-south": dc3,
	})

	dcs := table.DCs()
	if len(dcs) != 3 {
		t.Fatalf("expected 3 DCs, got %d", len(dcs))
	}
	if dcs[0] != "ap-south" || dcs[1] != "eu-west" || dcs[2] != "us-east" {
		t.Errorf("unexpected order: %v", dcs)
	}
}

func TestBuild_RowsAreSortedByKey(t *testing.T) {
	results := makeResults("z/key", "1", "a/key", "2", "m/key", "3")
	table := pivot.Build(map[string][]diff.Result{"dc1": results})

	expected := []string{"a/key", "m/key", "z/key"}
	for i, row := range table {
		if row.Key != expected[i] {
			t.Errorf("row %d: expected %s, got %s", i, expected[i], row.Key)
		}
	}
}
