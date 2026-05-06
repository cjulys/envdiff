package ignore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/ignore"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".envdiffignore")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp ignore file: %v", err)
	}
	return path
}

func TestLoadFile_BasicKeys(t *testing.T) {
	path := writeTempIgnore(t, "SECRET_KEY\nDATABASE_URL\n")
	s, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Contains("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if !s.Contains("DATABASE_URL") {
		t.Error("expected DATABASE_URL to be ignored")
	}
	if s.Contains("OTHER_KEY") {
		t.Error("OTHER_KEY should not be in ignore set")
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempIgnore(t, "# this is a comment\n\nAPI_KEY\n")
	s, err := ignore.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 1 {
		t.Errorf("expected 1 key, got %d", s.Len())
	}
	if !s.Contains("API_KEY") {
		t.Error("expected API_KEY to be in ignore set")
	}
}

func TestLoadFile_NotFound_ReturnsEmptySet(t *testing.T) {
	s, err := ignore.LoadFile("/nonexistent/path/.envdiffignore")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("expected empty set, got %d keys", s.Len())
	}
}

func TestNew_IsEmpty(t *testing.T) {
	s := ignore.New()
	if s.Len() != 0 {
		t.Errorf("expected empty set, got %d", s.Len())
	}
}

func TestAdd_And_Contains(t *testing.T) {
	s := ignore.New()
	s.Add("MY_KEY")
	if !s.Contains("MY_KEY") {
		t.Error("expected MY_KEY to be contained after Add")
	}
	if s.Contains("OTHER") {
		t.Error("OTHER should not be contained")
	}
}
