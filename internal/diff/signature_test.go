package diff

import (
	"strings"
	"testing"
)

func buildSigResults(statuses map[string]Status) []Result {
	out := make([]Result, 0, len(statuses))
	for k, s := range statuses {
		out = append(out, Result{Key: k, Status: s})
	}
	return out
}

func TestBuildSignatures_ReturnsOnePerLabel(t *testing.T) {
	runs := map[string][]Result{
		"prod": buildSigResults(map[string]Status{"A": StatusMatch, "B": StatusMissingInB}),
		"staging": buildSigResults(map[string]Status{"A": StatusMatch}),
	}
	sigs := BuildSignatures(runs)
	if len(sigs) != 2 {
		t.Fatalf("expected 2 signatures, got %d", len(sigs))
	}
}

func TestBuildSignatures_SortedByLabel(t *testing.T) {
	runs := map[string][]Result{
		"zzz": buildSigResults(map[string]Status{"X": StatusMatch}),
		"aaa": buildSigResults(map[string]Status{"X": StatusMatch}),
	}
	sigs := BuildSignatures(runs)
	if sigs[0].Label != "aaa" || sigs[1].Label != "zzz" {
		t.Errorf("expected labels sorted, got %s, %s", sigs[0].Label, sigs[1].Label)
	}
}

func TestBuildSignatures_ProblemCount(t *testing.T) {
	runs := map[string][]Result{
		"env": buildSigResults(map[string]Status{
			"A": StatusMatch,
			"B": StatusMismatch,
			"C": StatusMissingInA,
		}),
	}
	sigs := BuildSignatures(runs)
	if sigs[0].Problems != 2 {
		t.Errorf("expected 2 problems, got %d", sigs[0].Problems)
	}
}

func TestBuildSignatures_HashIsDeterministic(t *testing.T) {
	results := buildSigResults(map[string]Status{"FOO": StatusMatch, "BAR": StatusMismatch})
	runs := map[string][]Result{"e": results}
	a := BuildSignatures(runs)
	b := BuildSignatures(runs)
	if a[0].Hash != b[0].Hash {
		t.Errorf("hash not deterministic: %s vs %s", a[0].Hash, b[0].Hash)
	}
}

func TestBuildSignatures_DifferentResultsDifferentHash(t *testing.T) {
	runA := map[string][]Result{
		"e": buildSigResults(map[string]Status{"FOO": StatusMatch}),
	}
	runB := map[string][]Result{
		"e": buildSigResults(map[string]Status{"FOO": StatusMismatch}),
	}
	sigA := BuildSignatures(runA)
	sigB := BuildSignatures(runB)
	if sigA[0].Hash == sigB[0].Hash {
		t.Error("expected different hashes for different statuses")
	}
}

func TestBuildSignatures_EmptyRunsReturnsNil(t *testing.T) {
	sigs := BuildSignatures(nil)
	if sigs != nil {
		t.Errorf("expected nil, got %v", sigs)
	}
}

func TestWriteSignatures_ContainsLabel(t *testing.T) {
	sigs := []Signature{
		{Label: "production", Hash: "abcd1234", Problems: 3, Keys: []string{"A", "B"}},
	}
	var sb strings.Builder
	WriteSignatures(&sb, sigs)
	out := sb.String()
	if !strings.Contains(out, "production") {
		t.Error("expected label 'production' in output")
	}
	if !strings.Contains(out, "abcd1234") {
		t.Error("expected hash in output")
	}
}

func TestWriteSignatures_EmptyPrintsMessage(t *testing.T) {
	var sb strings.Builder
	WriteSignatures(&sb, nil)
	if !strings.Contains(sb.String(), "no signatures") {
		t.Error("expected 'no signatures' message for empty input")
	}
}
