package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildChurnResults(statuses map[string]Status) []Result {
	results := make([]Result, 0, len(statuses))
	for key, status := range statuses {
		results = append(results, Result{Key: key, Status: status})
	}
	return results
}

func TestNewChurnReport_IsEmpty(t *testing.T) {
	r := NewChurnReport()
	entries := r.Build()
	if len(entries) != 0 {
		t.Fatalf("expected empty report, got %d entries", len(entries))
	}
}

func TestChurn_AddRun_CountsProblems(t *testing.T) {
	r := NewChurnReport()
	r.AddRun(buildChurnResults(map[string]Status{
		"DB_HOST": StatusMismatch,
		"API_KEY": StatusMissingInB,
		"PORT":    StatusMatch,
	}))
	entries := r.Build()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestChurn_RateCalculation(t *testing.T) {
	r := NewChurnReport()
	// Run 1: DB_HOST is a problem
	r.AddRun(buildChurnResults(map[string]Status{"DB_HOST": StatusMismatch}))
	// Run 2: DB_HOST is a problem again
	r.AddRun(buildChurnResults(map[string]Status{"DB_HOST": StatusMismatch}))
	// Run 3: DB_HOST is fine
	r.AddRun(buildChurnResults(map[string]Status{"DB_HOST": StatusMatch}))

	entries := r.Build()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	got := entries[0]
	if got.Key != "DB_HOST" {
		t.Errorf("expected DB_HOST, got %s", got.Key)
	}
	if got.Problems != 2 {
		t.Errorf("expected 2 problems, got %d", got.Problems)
	}
	if got.Runs != 3 {
		t.Errorf("expected 3 runs, got %d", got.Runs)
	}
	want := 2.0 / 3.0
	if got.Rate < want-0.001 || got.Rate > want+0.001 {
		t.Errorf("expected rate ~%.4f, got %.4f", want, got.Rate)
	}
}

func TestChurn_SortedByRateDesc(t *testing.T) {
	r := NewChurnReport()
	r.AddRun(buildChurnResults(map[string]Status{
		"RARE_KEY":   StatusMismatch,
		"COMMON_KEY": StatusMismatch,
	}))
	r.AddRun(buildChurnResults(map[string]Status{
		"COMMON_KEY": StatusMismatch,
	}))
	entries := r.Build()
	if len(entries) < 2 {
		t.Fatal("expected at least 2 entries")
	}
	if entries[0].Key != "COMMON_KEY" {
		t.Errorf("expected COMMON_KEY first, got %s", entries[0].Key)
	}
}

func TestWriteChurn_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteChurn(&buf, nil)
	if !strings.Contains(buf.String(), "No churn") {
		t.Errorf("expected no-churn message, got: %s", buf.String())
	}
}

func TestWriteChurn_ContainsKey(t *testing.T) {
	r := NewChurnReport()
	r.AddRun(buildChurnResults(map[string]Status{"SECRET_KEY": StatusMissingInA}))
	var buf bytes.Buffer
	WriteChurn(&buf, r.Build())
	if !strings.Contains(buf.String(), "SECRET_KEY") {
		t.Errorf("expected SECRET_KEY in output, got: %s", buf.String())
	}
}
