package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunContext_NoDifferences(t *testing.T) {
	envA := map[string]string{"FOO": "bar", "BAZ": "qux"}
	envB := map[string]string{"FOO": "bar", "BAZ": "qux"}

	var buf bytes.Buffer
	hasProblems, err := RunContext(&buf, envA, envB, DefaultContextOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hasProblems {
		t.Error("expected no problems")
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message, got: %s", buf.String())
	}
}

func TestRunContext_ReturnsTrue_WhenProblemsExist(t *testing.T) {
	envA := map[string]string{"FOO": "bar"}
	envB := map[string]string{"FOO": "different"}

	var buf bytes.Buffer
	hasProblems, err := RunContext(&buf, envA, envB, DefaultContextOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !hasProblems {
		t.Error("expected problems to be detected")
	}
}

func TestRunContext_OutputContainsProblemCount(t *testing.T) {
	envA := map[string]string{"A": "1", "B": "2"}
	envB := map[string]string{"A": "9", "B": "2"}

	var buf bytes.Buffer
	_, _ = RunContext(&buf, envA, envB, DefaultContextOptions())
	out := buf.String()
	if !strings.Contains(out, "1 problem") {
		t.Errorf("expected problem count in output, got: %s", out)
	}
}

func TestRunContext_DefaultOptions(t *testing.T) {
	opts := DefaultContextOptions()
	if opts.Lines != 2 {
		t.Errorf("expected default Lines=2, got %d", opts.Lines)
	}
	if opts.Verbose {
		t.Error("expected Verbose=false by default")
	}
}

func TestRunContext_VerboseOption_ShowsValues(t *testing.T) {
	envA := map[string]string{"TOKEN": "secret123", "HOST": "localhost"}
	envB := map[string]string{"TOKEN": "other456", "HOST": "localhost"}

	opts := DefaultContextOptions()
	opts.Verbose = true

	var buf bytes.Buffer
	_, _ = RunContext(&buf, envA, envB, opts)
	out := buf.String()
	if !strings.Contains(out, "TOKEN=secret123") {
		t.Errorf("expected value in verbose output, got: %s", out)
	}
}
