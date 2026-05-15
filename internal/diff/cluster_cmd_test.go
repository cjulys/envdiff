package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildClusterCmdResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMatch},
		{Key: "DB_PORT", Status: StatusMismatch},
		{Key: "AWS_KEY", Status: StatusMissingInB},
		{Key: "AWS_SECRET", Status: StatusMatch},
		{Key: "PORT", Status: StatusMatch},
	}
}

func TestRunCluster_NoDifferences_ReturnsFalse(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Status: StatusMatch},
		{Key: "DB_PORT", Status: StatusMatch},
	}
	var buf bytes.Buffer
	got := RunCluster(&buf, results, DefaultClusterOptions())
	if got {
		t.Error("expected false when no problems")
	}
}

func TestRunCluster_WithProblems_ReturnsTrue(t *testing.T) {
	var buf bytes.Buffer
	got := RunCluster(&buf, buildClusterCmdResults(), DefaultClusterOptions())
	if !got {
		t.Error("expected true when problems exist")
	}
}

func TestRunCluster_OutputContainsCluster(t *testing.T) {
	var buf bytes.Buffer
	RunCluster(&buf, buildClusterCmdResults(), DefaultClusterOptions())
	out := buf.String()
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB in output, got: %s", out)
	}
}

func TestRunCluster_OnlyProblems_FiltersClean(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Status: StatusMatch},
		{Key: "DB_PORT", Status: StatusMatch},
		{Key: "AWS_KEY", Status: StatusMismatch},
	}
	opts := DefaultClusterOptions()
	opts.OnlyProblems = true
	var buf bytes.Buffer
	RunCluster(&buf, results, opts)
	out := buf.String()
	if strings.Contains(out, "DB") {
		t.Errorf("DB cluster should be filtered out: %s", out)
	}
	if !strings.Contains(out, "AWS") {
		t.Errorf("AWS cluster should appear: %s", out)
	}
}

func TestRunCluster_MinKeys_FiltersSmallClusters(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Status: StatusMismatch},
		{Key: "DB_PORT", Status: StatusMismatch},
		{Key: "PORT", Status: StatusMismatch},
	}
	opts := DefaultClusterOptions()
	opts.MinKeys = 2
	var buf bytes.Buffer
	RunCluster(&buf, results, opts)
	out := buf.String()
	if strings.Contains(out, "PORT ") {
		t.Errorf("single-key cluster PORT should be filtered: %s", out)
	}
}

func TestClusterSummary_Format(t *testing.T) {
	report := BuildCluster(buildClusterCmdResults())
	summary := ClusterSummary(report)
	if !strings.Contains(summary, "cluster") {
		t.Errorf("expected 'cluster' in summary: %s", summary)
	}
	if !strings.Contains(summary, "issue") {
		t.Errorf("expected 'issue' in summary: %s", summary)
	}
}
