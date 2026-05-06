package merge_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/merge"
)

func TestMerge_UnionOfKeys(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAR": "3", "BAZ": "4"}

	r := merge.Merge([]string{"a", "b"}, []map[string]string{a, b})

	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
	expected := []string{"BAR", "BAZ", "FOO"}
	for i, k := range expected {
		if r.Keys[i] != k {
			t.Errorf("expected key[%d]=%q, got %q", i, k, r.Keys[i])
		}
	}
}

func TestMerge_FirstValueWins(t *testing.T) {
	a := map[string]string{"KEY": "from_a"}
	b := map[string]string{"KEY": "from_b"}

	r := merge.Merge([]string{"a", "b"}, []map[string]string{a, b})

	if r.Values["KEY"] != "from_a" {
		t.Errorf("expected value from_a, got %q", r.Values["KEY"])
	}
}

func TestMerge_SourceTracking(t *testing.T) {
	a := map[string]string{"SHARED": "1", "ONLY_A": "x"}
	b := map[string]string{"SHARED": "2", "ONLY_B": "y"}

	r := merge.Merge([]string{"a", "b"}, []map[string]string{a, b})

	if len(r.Sources["SHARED"]) != 2 {
		t.Errorf("expected SHARED in 2 sources, got %d", len(r.Sources["SHARED"]))
	}
	if len(r.Sources["ONLY_A"]) != 1 || r.Sources["ONLY_A"][0] != "a" {
		t.Errorf("expected ONLY_A source to be [a], got %v", r.Sources["ONLY_A"])
	}
}

func TestMerge_EmptyInput(t *testing.T) {
	r := merge.Merge(nil, nil)
	if len(r.Keys) != 0 {
		t.Errorf("expected no keys, got %d", len(r.Keys))
	}
}

func TestUniqueKeys_ReturnsOnlyExclusiveKeys(t *testing.T) {
	a := map[string]string{"SHARED": "1", "ONLY_A": "x"}
	b := map[string]string{"SHARED": "2", "ONLY_B": "y"}

	unique := merge.UniqueKeys([]string{"a", "b"}, []map[string]string{a, b})

	if _, ok := unique["SHARED"]; ok {
		t.Error("SHARED should not be in unique keys")
	}
	if unique["ONLY_A"] != "a" {
		t.Errorf("expected ONLY_A -> a, got %q", unique["ONLY_A"])
	}
	if unique["ONLY_B"] != "b" {
		t.Errorf("expected ONLY_B -> b, got %q", unique["ONLY_B"])
	}
}
