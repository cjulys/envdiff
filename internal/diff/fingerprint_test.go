package diff

import (
	"strings"
	"testing"
)

func buildFingerprintResults() []Result {
	return []Result{
		{Key: "APP_ENV", Status: StatusMatch},
		{Key: "DB_HOST", Status: StatusMismatch},
		{Key: "SECRET", Status: StatusMissingInB},
	}
}

func TestComputeFingerprint_Deterministic(t *testing.T) {
	results := buildFingerprintResults()
	f1 := ComputeFingerprint(results)
	f2 := ComputeFingerprint(results)
	if f1.Hash != f2.Hash {
		t.Errorf("expected same hash, got %s vs %s", f1.Hash, f2.Hash)
	}
}

func TestComputeFingerprint_OrderIndependent(t *testing.T) {
	a := []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMismatch},
	}
	b := []Result{
		{Key: "B", Status: StatusMismatch},
		{Key: "A", Status: StatusMatch},
	}
	if ComputeFingerprint(a).Hash != ComputeFingerprint(b).Hash {
		t.Error("expected order-independent fingerprint")
	}
}

func TestComputeFingerprint_DifferentResults(t *testing.T) {
	a := []Result{{Key: "X", Status: StatusMatch}}
	b := []Result{{Key: "X", Status: StatusMismatch}}
	if ComputeFingerprint(a).Equal(ComputeFingerprint(b)) {
		t.Error("expected different fingerprints for different statuses")
	}
}

func TestComputeFingerprint_IssueCount(t *testing.T) {
	results := buildFingerprintResults()
	fp := ComputeFingerprint(results)
	if fp.Issues != 2 {
		t.Errorf("expected 2 issues, got %d", fp.Issues)
	}
	if fp.Keys != 3 {
		t.Errorf("expected 3 keys, got %d", fp.Keys)
	}
}

func TestComputeFingerprint_Empty(t *testing.T) {
	fp := ComputeFingerprint(nil)
	if fp.Keys != 0 || fp.Issues != 0 {
		t.Error("expected zero keys and issues for empty input")
	}
	if fp.Hash == "" {
		t.Error("expected non-empty hash even for empty input")
	}
}

func TestWriteFingerprint_ContainsHash(t *testing.T) {
	fp := ComputeFingerprint(buildFingerprintResults())
	var sb strings.Builder
	WriteFingerprint(&sb, fp, "staging")
	out := sb.String()
	if !strings.Contains(out, "staging") {
		t.Error("expected label in output")
	}
	if !strings.Contains(out, "issues") {
		t.Error("expected issues line in output")
	}
	if !strings.Contains(out, "hash") {
		t.Error("expected hash line in output")
	}
}

func TestWriteFingerprint_CleanStatus(t *testing.T) {
	results := []Result{{Key: "A", Status: StatusMatch}}
	fp := ComputeFingerprint(results)
	var sb strings.Builder
	WriteFingerprint(&sb, fp, "")
	if !strings.Contains(sb.String(), "clean") {
		t.Error("expected 'clean' status for zero issues")
	}
}
