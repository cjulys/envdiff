package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildClusterResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMatch},
		{Key: "DB_PORT", Status: StatusMismatch},
		{Key: "DB_NAME", Status: StatusMissingInB},
		{Key: "AWS_KEY", Status: StatusMatch},
		{Key: "AWS_SECRET", Status: StatusMismatch},
		{Key: "PORT", Status: StatusMatch},
	}
}

func TestBuildCluster_EntryCount(t *testing.T) {
	results := buildClusterResults()
	report := BuildCluster(results)
	// DB, AWS, PORT → 3 clusters
	if len(report.Entries) != 3 {
		t.Fatalf("expected 3 clusters, got %d", len(report.Entries))
	}
}

func TestBuildCluster_DBProblems(t *testing.T) {
	report := BuildCluster(buildClusterResults())
	var db *ClusterEntry
	for i := range report.Entries {
		if report.Entries[i].Pattern == "DB" {
			db = &report.Entries[i]
			break
		}
	}
	if db == nil {
		t.Fatal("DB cluster not found")
	}
	if db.Problems != 2 {
		t.Errorf("expected 2 DB problems, got %d", db.Problems)
	}
	if db.Total != 3 {
		t.Errorf("expected 3 DB keys, got %d", db.Total)
	}
}

func TestBuildCluster_SortedByProblemsDesc(t *testing.T) {
	report := BuildCluster(buildClusterResults())
	for i := 1; i < len(report.Entries); i++ {
		if report.Entries[i].Problems > report.Entries[i-1].Problems {
			t.Errorf("entries not sorted by problems desc at index %d", i)
		}
	}
}

func TestBuildCluster_KeysAreSorted(t *testing.T) {
	report := BuildCluster(buildClusterResults())
	for _, e := range report.Entries {
		for i := 1; i < len(e.Keys); i++ {
			if e.Keys[i] < e.Keys[i-1] {
				t.Errorf("keys not sorted in cluster %s", e.Pattern)
			}
		}
	}
}

func TestBuildCluster_NoUnderscore(t *testing.T) {
	results := []Result{
		{Key: "PORT", Status: StatusMatch},
		{Key: "HOST", Status: StatusMismatch},
	}
	report := BuildCluster(results)
	if len(report.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(report.Entries))
	}
}

func TestWriteCluster_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteCluster(&buf, ClusterReport{})
	if !strings.Contains(buf.String(), "no keys") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteCluster_ContainsPattern(t *testing.T) {
	report := BuildCluster(buildClusterResults())
	var buf bytes.Buffer
	WriteCluster(&buf, report)
	out := buf.String()
	if !strings.Contains(out, "DB") {
		t.Errorf("expected DB cluster in output, got: %s", out)
	}
	if !strings.Contains(out, "AWS") {
		t.Errorf("expected AWS cluster in output, got: %s", out)
	}
}
