package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildCoverageRuns() map[string]map[string]string {
	return map[string]map[string]string{
		"production": {"A": "1", "B": "2", "C": "3"},
		"staging":    {"A": "1", "B": "2"},
		"dev":        {"A": "1"},
	}
}

func findCoverage(entries []CoverageEntry, label string) (CoverageEntry, bool) {
	for _, e := range entries {
		if e.Label == label {
			return e, true
		}
	}
	return CoverageEntry{}, false
}

func TestBuildCoverage_EntryCount(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestBuildCoverage_TotalIsUnion(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	for _, e := range entries {
		if e.Total != 3 {
			t.Errorf("%s: expected total 3, got %d", e.Label, e.Total)
		}
	}
}

func TestBuildCoverage_PresentCounts(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)

	cases := map[string]int{"production": 3, "staging": 2, "dev": 1}
	for label, want := range cases {
		e, ok := findCoverage(entries, label)
		if !ok {
			t.Fatalf("entry %q not found", label)
		}
		if e.Present != want {
			t.Errorf("%s: expected present %d, got %d", label, want, e.Present)
		}
	}
}

func TestBuildCoverage_SortedByRateDesc(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	for i := 1; i < len(entries); i++ {
		if entries[i].Rate > entries[i-1].Rate {
			t.Errorf("entries not sorted by rate desc at index %d", i)
		}
	}
}

func TestBuildCoverage_EmptyRuns(t *testing.T) {
	entries := BuildCoverage(nil)
	if entries != nil {
		t.Errorf("expected nil for empty input, got %v", entries)
	}
}

func TestBuildCoverage_RateCalculation(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	e, _ := findCoverage(entries, "staging")
	want := 2.0 / 3.0
	if e.Rate < want-0.001 || e.Rate > want+0.001 {
		t.Errorf("staging rate: expected %.4f, got %.4f", want, e.Rate)
	}
}

func TestWriteCoverage_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteCoverage(&buf, nil)
	if !strings.Contains(buf.String(), "no environments") {
		t.Errorf("expected 'no environments' message, got: %s", buf.String())
	}
}

func TestWriteCoverage_ContainsLabels(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	var buf bytes.Buffer
	WriteCoverage(&buf, entries)
	out := buf.String()
	for _, label := range []string{"production", "staging", "dev"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestWriteCoverage_ContainsPercentSign(t *testing.T) {
	runs := buildCoverageRuns()
	entries := BuildCoverage(runs)
	var buf bytes.Buffer
	WriteCoverage(&buf, entries)
	if !strings.Contains(buf.String(), "%") {
		t.Errorf("expected percentage sign in coverage output")
	}
}
