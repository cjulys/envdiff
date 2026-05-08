package baseline_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/baseline"
	"github.com/user/envdiff/internal/diff"
)

func TestWriteReport_NoNewResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	results := sampleResults()
	_ = baseline.Save(path, results)
	b, _ := baseline.Load(path)

	var buf bytes.Buffer
	baseline.WriteReport(&buf, results, b)
	out := buf.String()

	if !strings.Contains(out, "No new differences") {
		t.Errorf("expected no-new message, got:\n%s", out)
	}
	if !strings.Contains(out, "Suppressed (known): 3") {
		t.Errorf("expected suppressed count, got:\n%s", out)
	}
}

func TestWriteReport_WithNewResults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = baseline.Save(path, sampleResults())
	b, _ := baseline.Load(path)

	current := append(sampleResults(),
		diff.Result{Key: "NEW_VAR", Status: diff.StatusMissingInB},
		diff.Result{Key: "ANOTHER", Status: diff.StatusMismatch},
	)

	var buf bytes.Buffer
	baseline.WriteReport(&buf, current, b)
	out := buf.String()

	if !strings.Contains(out, "New: 2") {
		t.Errorf("expected 2 new results, got:\n%s", out)
	}
	if !strings.Contains(out, "ANOTHER") {
		t.Errorf("expected ANOTHER in output, got:\n%s", out)
	}
	if !strings.Contains(out, "NEW_VAR") {
		t.Errorf("expected NEW_VAR in output, got:\n%s", out)
	}
}

func TestWriteReport_IncludesTimestamp(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")
	_ = baseline.Save(path, sampleResults())
	b, _ := baseline.Load(path)

	var buf bytes.Buffer
	baseline.WriteReport(&buf, sampleResults(), b)
	out := buf.String()

	if !strings.Contains(out, "Baseline recorded:") {
		t.Errorf("expected timestamp header, got:\n%s", out)
	}
}
