package diff

import (
	"strings"
	"testing"
)

func buildFingerprintCmdResults(problems bool) []Result {
	if problems {
		return []Result{
			{Key: "HOST", Status: StatusMismatch},
			{Key: "PORT", Status: StatusMatch},
		}
	}
	return []Result{
		{Key: "HOST", Status: StatusMatch},
		{Key: "PORT", Status: StatusMatch},
	}
}

func TestRunFingerprint_NoDifferences_ReturnsFalse(t *testing.T) {
	results := buildFingerprintCmdResults(false)
	var sb strings.Builder
	got := RunFingerprint(&sb, results, nil, DefaultFingerprintOptions())
	if got {
		t.Error("expected false when no issues")
	}
}

func TestRunFingerprint_WithProblems_ReturnsTrue(t *testing.T) {
	results := buildFingerprintCmdResults(true)
	var sb strings.Builder
	got := RunFingerprint(&sb, results, nil, DefaultFingerprintOptions())
	if !got {
		t.Error("expected true when issues exist")
	}
}

func TestRunFingerprint_OutputContainsLabel(t *testing.T) {
	results := buildFingerprintCmdResults(false)
	opts := DefaultFingerprintOptions()
	opts.LabelA = "production"
	var sb strings.Builder
	RunFingerprint(&sb, results, nil, opts)
	if !strings.Contains(sb.String(), "production") {
		t.Error("expected label in output")
	}
}

func TestRunFingerprint_CompareMatchingEnvs(t *testing.T) {
	a := buildFingerprintCmdResults(false)
	b := buildFingerprintCmdResults(false)
	opts := DefaultFingerprintOptions()
	opts.Compare = true
	var sb strings.Builder
	RunFingerprint(&sb, a, b, opts)
	if !strings.Contains(sb.String(), "match") {
		t.Error("expected 'match' in compare output")
	}
}

func TestRunFingerprint_CompareDivergingEnvs(t *testing.T) {
	a := buildFingerprintCmdResults(false)
	b := buildFingerprintCmdResults(true)
	opts := DefaultFingerprintOptions()
	opts.Compare = true
	var sb strings.Builder
	RunFingerprint(&sb, a, b, opts)
	if !strings.Contains(sb.String(), "differ") {
		t.Error("expected 'differ' in compare output")
	}
}

func TestRunFingerprint_DefaultOptions(t *testing.T) {
	opts := DefaultFingerprintOptions()
	if opts.LabelA == "" || opts.LabelB == "" {
		t.Error("expected non-empty default labels")
	}
	if opts.Compare {
		t.Error("expected Compare to default to false")
	}
}
