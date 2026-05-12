package diff_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/your/envdiff/internal/diff"
)

func makeRuns() [][]diff.Result {
	return [][]diff.Result{
		{
			{Key: "DB_HOST", Status: diff.StatusMismatch},
			{Key: "API_KEY", Status: diff.StatusMissingInB},
			{Key: "HOST", Status: diff.StatusMatch},
		},
		{
			{Key: "DB_HOST", Status: diff.StatusMismatch},
			{Key: "SECRET", Status: diff.StatusMissingInA},
		},
		{
			{Key: "DB_HOST", Status: diff.StatusMismatch},
			{Key: "API_KEY", Status: diff.StatusMissingInB},
		},
	}
}

func TestBuildHeatmap_CountsProblems(t *testing.T) {
	entries := diff.BuildHeatmap(makeRuns())
	found := map[string]int{}
	for _, e := range entries {
		found[e.Key] = e.Count
	}
	if found["DB_HOST"] != 3 {
		t.Errorf("expected DB_HOST count=3, got %d", found["DB_HOST"])
	}
	if found["API_KEY"] != 2 {
		t.Errorf("expected API_KEY count=2, got %d", found["API_KEY"])
	}
	if found["HOST"] != 0 {
		t.Errorf("expected HOST not in heatmap (match), got %d", found["HOST"])
	}
}

func TestBuildHeatmap_SortedByCountDesc(t *testing.T) {
	entries := diff.BuildHeatmap(makeRuns())
	if len(entries) == 0 {
		t.Fatal("expected entries")
	}
	if entries[0].Key != "DB_HOST" {
		t.Errorf("expected first entry to be DB_HOST, got %s", entries[0].Key)
	}
}

func TestBuildHeatmap_EmptyRuns(t *testing.T) {
	entries := diff.BuildHeatmap([][]diff.Result{})
	if len(entries) != 0 {
		t.Errorf("expected empty heatmap, got %d entries", len(entries))
	}
}

func TestWriteHeatmap_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	diff.WriteHeatmap(nil, &buf)
	if !strings.Contains(buf.String(), "no problem keys") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteHeatmap_ContainsKeys(t *testing.T) {
	entries := diff.BuildHeatmap(makeRuns())
	var buf bytes.Buffer
	diff.WriteHeatmap(entries, &buf)
	out := buf.String()
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected DB_HOST in output, got:\n%s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got:\n%s", out)
	}
}
