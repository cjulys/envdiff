package schema_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/schema"
)

func TestWriteReport_NoViolations(t *testing.T) {
	var buf strings.Builder
	schema.WriteReport(&buf, nil, "prod.env")
	out := buf.String()
	if !strings.Contains(out, "pass schema validation") {
		t.Errorf("expected success message, got: %s", out)
	}
}

func TestWriteReport_WithViolations(t *testing.T) {
	violations := []schema.Violation{
		{Key: "DB_URL", Message: "required key is missing"},
		{Key: "PORT", Message: `value "abc" does not match pattern ^[0-9]+$`},
	}
	var buf strings.Builder
	schema.WriteReport(&buf, violations, "staging.env")
	out := buf.String()

	if !strings.Contains(out, "2 violation(s)") {
		t.Errorf("expected violation count, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output")
	}
}

func TestWriteReport_SortedOutput(t *testing.T) {
	violations := []schema.Violation{
		{Key: "Z_KEY", Message: "required key is missing"},
		{Key: "A_KEY", Message: "required key is missing"},
	}
	var buf strings.Builder
	schema.WriteReport(&buf, violations, "test.env")
	out := buf.String()

	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA > idxZ {
		t.Error("expected A_KEY to appear before Z_KEY")
	}
}

func TestWriteReport_IncludesFilename(t *testing.T) {
	var buf strings.Builder
	schema.WriteReport(&buf, nil, "my-special.env")
	if !strings.Contains(buf.String(), "my-special.env") {
		t.Error("expected filename in report header")
	}
}
