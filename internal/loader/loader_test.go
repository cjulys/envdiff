package loader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadFile_Basic(t *testing.T) {
	path := writeTempEnv(t, "KEY=value\nFOO=bar\n")

	ef, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ef.Vars["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %q", ef.Vars["KEY"])
	}
	if ef.Vars["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", ef.Vars["FOO"])
	}
	if ef.Name != filepath.Base(path) {
		t.Errorf("expected Name=%q, got %q", filepath.Base(path), ef.Name)
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFiles_Multiple(t *testing.T) {
	p1 := writeTempEnv(t, "A=1\n")
	p2 := writeTempEnv(t, "B=2\n")

	files, err := LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Vars["A"] != "1" {
		t.Errorf("expected A=1")
	}
	if files[1].Vars["B"] != "2" {
		t.Errorf("expected B=2")
	}
}

func TestLoadFiles_PartialError(t *testing.T) {
	p1 := writeTempEnv(t, "X=10\n")

	files, err := LoadFiles([]string{p1, "/no/such/file.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if len(files) != 1 {
		t.Errorf("expected 1 successfully loaded file, got %d", len(files))
	}
}
