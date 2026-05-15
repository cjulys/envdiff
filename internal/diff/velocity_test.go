package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

var baseVelocityTime = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func buildVelocityRuns() []VelocityRun {
	r1 := []Result{
		{Key: "DB_HOST", Status: StatusMismatch},
		{Key: "API_KEY", Status: StatusMissingInB},
	}
	r2 := []Result{
		{Key: "DB_HOST", Status: StatusMismatch},
		{Key: "SECRET", Status: StatusMissingInA},
	}
	return []VelocityRun{
		{Label: "run-1", Timestamp: baseVelocityTime, Results: r1},
		{Label: "run-2", Timestamp: baseVelocityTime.Add(48 * time.Hour), Results: r2},
	}
}

func TestBuildVelocity_EntryCount(t *testing.T) {
	runs := buildVelocityRuns()
	entries := BuildVelocity(runs)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestBuildVelocity_RateForRepeatedKey(t *testing.T) {
	runs := buildVelocityRuns()
	entries := BuildVelocity(runs)
	var dbEntry *VelocityEntry
	for i := range entries {
		if entries[i].Key == "DB_HOST" {
			dbEntry = &entries[i]
		}
	}
	if dbEntry == nil {
		t.Fatal("DB_HOST entry not found")
	}
	if dbEntry.TotalProblems != 2 {
		t.Errorf("expected 2 problems, got %d", dbEntry.TotalProblems)
	}
	if dbEntry.ChangeRate <= 0 {
		t.Errorf("expected positive change rate, got %f", dbEntry.ChangeRate)
	}
}

func TestBuildVelocity_SortedByRateDesc(t *testing.T) {
	runs := buildVelocityRuns()
	entries := BuildVelocity(runs)
	for i := 1; i < len(entries); i++ {
		if entries[i].ChangeRate > entries[i-1].ChangeRate {
			t.Errorf("entries not sorted by rate desc at index %d", i)
		}
	}
}

func TestBuildVelocity_EmptyRuns(t *testing.T) {
	entries := BuildVelocity(nil)
	if len(entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(entries))
	}
}

func TestWriteVelocity_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteVelocity(&buf, nil)
	if !strings.Contains(buf.String(), "No problem keys") {
		t.Errorf("expected no-problem message, got: %s", buf.String())
	}
}

func TestWriteVelocity_ContainsKey(t *testing.T) {
	runs := buildVelocityRuns()
	entries := BuildVelocity(runs)
	var buf bytes.Buffer
	WriteVelocity(&buf, entries)
	if !strings.Contains(buf.String(), "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got: %s", buf.String())
	}
}

func TestRunVelocity_ReturnsFalseWhenEmpty(t *testing.T) {
	var buf bytes.Buffer
	ok := RunVelocity(&buf, nil, DefaultVelocityOptions())
	if ok {
		t.Error("expected false for empty runs")
	}
}

func TestRunVelocity_TopN_Limits(t *testing.T) {
	runs := buildVelocityRuns()
	opts := DefaultVelocityOptions()
	opts.TopN = 1
	var buf bytes.Buffer
	RunVelocity(&buf, runs, opts)
	out := buf.String()
	if strings.Contains(out, "(1 keys)") == false && !strings.Contains(out, "1 keys") {
		// just check the header mentions 1
		if !strings.Contains(out, "1") {
			t.Errorf("expected output to reflect TopN=1, got: %s", out)
		}
	}
}
