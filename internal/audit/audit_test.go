package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/audit"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", StatusA: "present", StatusB: "present", ValueA: "localhost", ValueB: "prod.db"},
		{Key: "API_KEY", StatusA: "present", StatusB: "missing", ValueA: "abc", ValueB: ""},
		{Key: "PORT", StatusA: "missing", StatusB: "present", ValueA: "", ValueB: "8080"},
	}
}

func TestLog_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	err := audit.Log(path, []string{".env", ".env.prod"}, sampleResults())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected log file to be created")
	}
}

func TestLog_AppendMultipleEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	results := sampleResults()
	if err := audit.Log(path, []string{".env"}, results); err != nil {
		t.Fatalf("first log: %v", err)
	}
	if err := audit.Log(path, []string{".env.staging"}, results); err != nil {
		t.Fatalf("second log: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read all: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestReadAll_NotFound_ReturnsNil(t *testing.T) {
	entries, err := audit.ReadAll("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entries != nil {
		t.Fatal("expected nil entries for missing file")
	}
}

func TestLog_RecordsCorrectProblemCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	if err := audit.Log(path, []string{".env"}, sampleResults()); err != nil {
		t.Fatalf("log: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read all: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if entries[0].Problems != 2 {
		t.Errorf("expected 2 problems, got %d", entries[0].Problems)
	}
	if entries[0].TotalKeys != 3 {
		t.Errorf("expected 3 total keys, got %d", entries[0].TotalKeys)
	}
}

func TestLog_TimestampIsRecent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	before := time.Now().UTC()
	if err := audit.Log(path, []string{".env"}, sampleResults()); err != nil {
		t.Fatalf("log: %v", err)
	}
	after := time.Now().UTC()

	entries, _ := audit.ReadAll(path)
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestLog_RecordsFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	files := []string{".env", ".env.prod"}
	if err := audit.Log(path, files, sampleResults()); err != nil {
		t.Fatalf("log: %v", err)
	}

	entries, err := audit.ReadAll(path)
	if err != nil {
		t.Fatalf("read all: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if len(entries[0].Files) != len(files) {
		t.Fatalf("expected %d files, got %d", len(files), len(entries[0].Files))
	}
	for i, f := range files {
		if entries[0].Files[i] != f {
			t.Errorf("expected file %q at index %d, got %q", f, i, entries[0].Files[i])
		}
	}
}
