package diff

import (
	"strings"
	"testing"
)

func TestWritePatch_MissingInA(t *testing.T) {
	results := []Result{
		{Key: "NEW_KEY", Status: StatusMissingInA, ValueA: "", ValueB: "hello"},
	}
	var buf strings.Builder
	if err := WritePatch(&buf, results, PatchOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+NEW_KEY=hello") {
		t.Errorf("expected +NEW_KEY=hello in patch, got:\n%s", out)
	}
}

func TestWritePatch_MissingInB(t *testing.T) {
	results := []Result{
		{Key: "OLD_KEY", Status: StatusMissingInB, ValueA: "old", ValueB: ""},
	}
	var buf strings.Builder
	_ = WritePatch(&buf, results, PatchOptions{})
	out := buf.String()
	if !strings.Contains(out, "-OLD_KEY=old") {
		t.Errorf("expected -OLD_KEY=old in patch, got:\n%s", out)
	}
}

func TestWritePatch_Mismatch(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", Status: StatusMismatch, ValueA: "localhost", ValueB: "prod.db"},
	}
	var buf strings.Builder
	_ = WritePatch(&buf, results, PatchOptions{})
	out := buf.String()
	if !strings.Contains(out, "-DB_HOST=localhost") {
		t.Errorf("expected removal line, got:\n%s", out)
	}
	if !strings.Contains(out, "+DB_HOST=prod.db") {
		t.Errorf("expected addition line, got:\n%s", out)
	}
}

func TestWritePatch_MatchShownByDefault(t *testing.T) {
	results := []Result{
		{Key: "PORT", Status: StatusMatch, ValueA: "8080", ValueB: "8080"},
	}
	var buf strings.Builder
	_ = WritePatch(&buf, results, PatchOptions{OnlyProblems: false})
	out := buf.String()
	if !strings.Contains(out, " PORT=8080") {
		t.Errorf("expected context line for match, got:\n%s", out)
	}
}

func TestWritePatch_OnlyProblems_HidesMatch(t *testing.T) {
	results := []Result{
		{Key: "PORT", Status: StatusMatch, ValueA: "8080", ValueB: "8080"},
	}
	var buf strings.Builder
	_ = WritePatch(&buf, results, PatchOptions{OnlyProblems: true})
	out := buf.String()
	if strings.Contains(out, "PORT") {
		t.Errorf("expected match key to be hidden, got:\n%s", out)
	}
}

func TestWritePatch_CustomLabel(t *testing.T) {
	var buf strings.Builder
	_ = WritePatch(&buf, nil, PatchOptions{TargetLabel: "prod/.env"})
	out := buf.String()
	if !strings.Contains(out, "+++ prod/.env") {
		t.Errorf("expected custom target label, got:\n%s", out)
	}
}

func TestWritePatch_SortedOutput(t *testing.T) {
	results := []Result{
		{Key: "Z_KEY", Status: StatusMissingInA, ValueB: "z"},
		{Key: "A_KEY", Status: StatusMissingInA, ValueB: "a"},
	}
	var buf strings.Builder
	_ = WritePatch(&buf, results, PatchOptions{})
	out := buf.String()
	aPos := strings.Index(out, "A_KEY")
	zPos := strings.Index(out, "Z_KEY")
	if aPos > zPos {
		t.Errorf("expected A_KEY before Z_KEY in sorted output")
	}
}
