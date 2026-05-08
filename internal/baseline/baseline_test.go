package baseline_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/baseline"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "APP_ENV", Status: diff.StatusMatch},
		{Key: "DB_HOST", Status: diff.StatusMismatch, ValueA: "localhost", ValueB: "prod-db"},
		{Key: "SECRET", Status: diff.StatusMissingInB},
	}
}

func TestSave_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	if err := baseline.Save(path, sampleResults()); err != nil {
		t.Fatalf("Save() error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestLoad_ReturnsBaseline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := sampleResults()
	_ = baseline.Save(path, results)

	b, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if len(b.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(b.Results))
	}
	if b.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if b.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := baseline.Load("/nonexistent/baseline.json")
	if err != baseline.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestNewResults_FiltersKnown(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = baseline.Save(path, sampleResults())
	b, _ := baseline.Load(path)

	current := append(sampleResults(), diff.Result{
		Key: "NEW_KEY", Status: diff.StatusMissingInA,
	})
	new := baseline.NewResults(current, b)
	if len(new) != 1 {
		t.Errorf("expected 1 new result, got %d", len(new))
	}
	if new[0].Key != "NEW_KEY" {
		t.Errorf("unexpected key: %s", new[0].Key)
	}
}

func TestNewResults_AllKnown_ReturnsNil(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = baseline.Save(path, sampleResults())
	b, _ := baseline.Load(path)

	new := baseline.NewResults(sampleResults(), b)
	if len(new) != 0 {
		t.Errorf("expected 0 new results, got %d", len(new))
	}
}
