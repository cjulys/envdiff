package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildDriftResults(match, missing, mismatch int) []Result {
	var results []Result
	for i := 0; i < match; i++ {
		results = append(results, Result{Key: "KEY", Status: StatusMatch})
	}
	for i := 0; i < missing; i++ {
		results = append(results, Result{Key: "MISS", Status: StatusMissingInB})
	}
	for i := 0; i < mismatch; i++ {
		results = append(results, Result{Key: "BAD", Status: StatusMismatch})
	}
	return results
}

func TestDrift_AddEntry_CountsProblems(t *testing.T) {
	dr := &DriftReport{}
	dr.AddEntry("staging", buildDriftResults(8, 2, 1))
	if len(dr.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(dr.Entries))
	}
	if dr.Entries[0].Problems != 3 {
		t.Errorf("expected 3 problems, got %d", dr.Entries[0].Problems)
	}
	if dr.Entries[0].Total != 11 {
		t.Errorf("expected total 11, got %d", dr.Entries[0].Total)
	}
}

func TestDrift_AddEntry_ZeroProblems(t *testing.T) {
	dr := &DriftReport{}
	dr.AddEntry("prod", buildDriftResults(5, 0, 0))
	if dr.Entries[0].Problems != 0 {
		t.Errorf("expected 0 problems, got %d", dr.Entries[0].Problems)
	}
}

func TestDrift_SortedEntries_ByLabel(t *testing.T) {
	dr := &DriftReport{}
	dr.AddEntry("prod", buildDriftResults(3, 0, 0))
	dr.AddEntry("dev", buildDriftResults(3, 1, 0))
	dr.AddEntry("staging", buildDriftResults(3, 0, 1))
	sorted := dr.SortedEntries()
	if sorted[0].Label != "dev" {
		t.Errorf("expected first label 'dev', got '%s'", sorted[0].Label)
	}
	if sorted[2].Label != "staging" {
		t.Errorf("expected last label 'staging', got '%s'", sorted[2].Label)
	}
}

func TestWriteDrift_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteDrift(&buf, &DriftReport{})
	if !strings.Contains(buf.String(), "No drift data") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteDrift_ContainsLabel(t *testing.T) {
	dr := &DriftReport{}
	dr.AddEntry("production", buildDriftResults(10, 2, 1))
	var buf bytes.Buffer
	WriteDrift(&buf, dr)
	out := buf.String()
	if !strings.Contains(out, "production") {
		t.Errorf("expected label 'production' in output, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected problem count in output, got: %s", out)
	}
}

func TestWriteDrift_PercentageShown(t *testing.T) {
	dr := &DriftReport{}
	dr.AddEntry("env", buildDriftResults(5, 5, 0))
	var buf bytes.Buffer
	WriteDrift(&buf, dr)
	out := buf.String()
	if !strings.Contains(out, "50.0%") {
		t.Errorf("expected 50.0%% in output, got: %s", out)
	}
}
