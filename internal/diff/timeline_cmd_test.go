package diff_test

import (
	"bytes"
	"testing"

	"github.com/your/envdiff/internal/diff"
)

func buildTimelineCmdResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, ValueA: "localhost", ValueB: "localhost"},
		{Key: "API_KEY", Status: diff.StatusMismatch, ValueA: "abc", ValueB: "xyz"},
		{Key: "SECRET", Status: diff.StatusMissingInB, ValueA: "hidden", ValueB: ""},
	}
}

func TestRunTimeline_NoDifferences(t *testing.T) {
	results := []diff.Result{
		{Key: "HOST", Status: diff.StatusMatch, ValueA: "a", ValueB: "a"},
	}
	var buf bytes.Buffer
	opts := diff.DefaultTimelineOptions()
	hasDiff := diff.RunTimeline(results, []string{"v1", "v2"}, opts, &buf)
	if hasDiff {
		t.Error("expected no diff")
	}
}

func TestRunTimeline_ReturnsTrue_WhenProblemsExist(t *testing.T) {
	results := buildTimelineCmdResults()
	var buf bytes.Buffer
	opts := diff.DefaultTimelineOptions()
	hasDiff := diff.RunTimeline(results, []string{"v1", "v2"}, opts, &buf)
	if !hasDiff {
		t.Error("expected diff to be detected")
	}
}

func TestRunTimeline_OutputContainsLabels(t *testing.T) {
	results := buildTimelineCmdResults()
	var buf bytes.Buffer
	opts := diff.DefaultTimelineOptions()
	diff.RunTimeline(results, []string{"staging", "prod"}, opts, &buf)
	out := buf.String()
	if !contains(out, "staging") {
		t.Errorf("expected output to contain label 'staging', got:\n%s", out)
	}
	if !contains(out, "prod") {
		t.Errorf("expected output to contain label 'prod', got:\n%s", out)
	}
}

func TestRunTimeline_DefaultOptions(t *testing.T) {
	opts := diff.DefaultTimelineOptions()
	if opts.Verbose {
		t.Error("expected Verbose to default to false")
	}
}

func TestRunTimeline_VerboseOption_ShowsValues(t *testing.T) {
	results := buildTimelineCmdResults()
	var buf bytes.Buffer
	opts := diff.DefaultTimelineOptions()
	opts.Verbose = true
	diff.RunTimeline(results, []string{"v1", "v2"}, opts, &buf)
	out := buf.String()
	if !contains(out, "abc") && !contains(out, "xyz") {
		t.Errorf("expected verbose output to contain values, got:\n%s", out)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
