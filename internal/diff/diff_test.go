package diff

import (
	"testing"
)

func makeMap(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func findResult(results []Result, key string) (Result, bool) {
	for _, r := range results {
		if r.Key == key {
			return r, true
		}
	}
	return Result{}, false
}

func TestCompare_Match(t *testing.T) {
	a := makeMap("FOO", "bar")
	b := makeMap("FOO", "bar")
	results := Compare(a, b)
	r, ok := findResult(results, "FOO")
	if !ok {
		t.Fatal("expected result for FOO")
	}
	if r.Status != StatusMatch {
		t.Errorf("expected StatusMatch, got %v", r.Status)
	}
}

func TestCompare_Mismatch(t *testing.T) {
	a := makeMap("FOO", "bar")
	b := makeMap("FOO", "baz")
	results := Compare(a, b)
	r, ok := findResult(results, "FOO")
	if !ok {
		t.Fatal("expected result for FOO")
	}
	if r.Status != StatusMismatch {
		t.Errorf("expected StatusMismatch, got %v", r.Status)
	}
	if r.ValueA != "bar" || r.ValueB != "baz" {
		t.Errorf("unexpected values: %q %q", r.ValueA, r.ValueB)
	}
}

func TestCompare_MissingInB(t *testing.T) {
	a := makeMap("ONLY_A", "1")
	b := makeMap[string, string]()
	results := Compare(a, b)
	r, ok := findResult(results, "ONLY_A")
	if !ok {
		t.Fatal("expected result for ONLY_A")
	}
	if r.Status != StatusMissingB {
		t.Errorf("expected StatusMissingB, got %v", r.Status)
	}
}

func TestCompare_MissingInA(t *testing.T) {
	a := makeMap[string, string]()
	b := makeMap("ONLY_B", "2")
	results := Compare(a, b)
	r, ok := findResult(results, "ONLY_B")
	if !ok {
		t.Fatal("expected result for ONLY_B")
	}
	if r.Status != StatusMissingA {
		t.Errorf("expected StatusMissingA, got %v", r.Status)
	}
}

func TestSummarize(t *testing.T) {
	a := makeMap("A", "1", "B", "2", "C", "3")
	b := makeMap("A", "1", "B", "9", "D", "4")
	results := Compare(a, b)
	s := Summarize(results)
	if s.Matches != 1 {
		t.Errorf("Matches: want 1, got %d", s.Matches)
	}
	if s.Mismatches != 1 {
		t.Errorf("Mismatches: want 1, got %d", s.Mismatches)
	}
	if s.MissingInB != 1 {
		t.Errorf("MissingInB: want 1, got %d", s.MissingInB)
	}
	if s.MissingInA != 1 {
		t.Errorf("MissingInA: want 1, got %d", s.MissingInA)
	}
}
