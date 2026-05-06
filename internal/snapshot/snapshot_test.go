package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.Match, ValueA: "localhost", ValueB: "localhost"},
		{Key: "DB_PASS", Status: diff.Mismatch, ValueA: "secret", ValueB: "other"},
		{Key: "API_KEY", Status: diff.MissingInB, ValueA: "abc123", ValueB: ""},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	err := snapshot.Save(path, ".env.dev", ".env.prod", sampleResults())
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("expected snapshot file to exist")
	}
}

func TestLoad_ReturnsSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := sampleResults()
	if err := snapshot.Save(path, ".env.dev", ".env.prod", results); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if s.FileA != ".env.dev" {
		t.Errorf("FileA = %q, want .env.dev", s.FileA)
	}
	if s.FileB != ".env.prod" {
		t.Errorf("FileB = %q, want .env.prod", s.FileB)
	}
	if len(s.Results) != len(results) {
		t.Errorf("Results len = %d, want %d", len(s.Results), len(results))
	}
	if s.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if s.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt should not be in the future")
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestFilterProblems_ReturnsOnlyDiffs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	if err := snapshot.Save(path, "a", "b", sampleResults()); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	s, _ := snapshot.Load(path)
	problems := s.FilterProblems()

	if len(problems) != 2 {
		t.Errorf("FilterProblems() len = %d, want 2", len(problems))
	}
	for _, r := range problems {
		if !r.IsProblem() {
			t.Errorf("expected problem result, got status %v", r.Status)
		}
	}
}
