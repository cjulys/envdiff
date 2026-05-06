package schema_test

import (
	"os"
	"testing"

	"github.com/yourorg/envdiff/internal/schema"
)

func writeTempSchema(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.schema")
	if err != nil {
		t.Fatalf("create temp schema: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write schema: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadFile_BasicRules(t *testing.T) {
	path := writeTempSchema(t, "!DB_HOST\nPORT\n")
	s, err := schema.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(s.Rules))
	}
	if !s.Rules[0].Required {
		t.Error("DB_HOST should be required")
	}
	if s.Rules[1].Required {
		t.Error("PORT should not be required")
	}
}

func TestLoadFile_SkipsCommentsAndBlanks(t *testing.T) {
	path := writeTempSchema(t, "# comment\n\n!KEY\n")
	s, err := schema.LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(s.Rules))
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := schema.LoadFile("/nonexistent/schema.txt")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestCheck_MissingRequiredKey(t *testing.T) {
	path := writeTempSchema(t, "!DB_URL\n")
	s, _ := schema.LoadFile(path)
	violations := s.Check(map[string]string{})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DB_URL" {
		t.Errorf("unexpected key: %s", violations[0].Key)
	}
}

func TestCheck_PatternMismatch(t *testing.T) {
	path := writeTempSchema(t, "PORT=^[0-9]+$\n")
	s, _ := schema.LoadFile(path)
	violations := s.Check(map[string]string{"PORT": "not-a-number"})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_PatternMatch(t *testing.T) {
	path := writeTempSchema(t, "PORT=^[0-9]+$\n")
	s, _ := schema.LoadFile(path)
	violations := s.Check(map[string]string{"PORT": "8080"})
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_OptionalMissingKeyNoViolation(t *testing.T) {
	path := writeTempSchema(t, "OPTIONAL_KEY\n")
	s, _ := schema.LoadFile(path)
	violations := s.Check(map[string]string{})
	if len(violations) != 0 {
		t.Errorf("expected no violations for optional missing key, got %d", len(violations))
	}
}
