package stats_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/stats"
)

func makeResult(key string, status diff.Status) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestCompute_EmptyResults(t *testing.T) {
	r := stats.Compute(nil)
	if r.Total != 0 || r.Problems != 0 {
		t.Errorf("expected zero report, got %+v", r)
	}
}

func TestCompute_AllMatch(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", diff.StatusMatch),
		makeResult("DB_PORT", diff.StatusMatch),
	}
	r := stats.Compute(results)
	if r.Total != 2 || r.Matched != 2 || r.Problems != 0 {
		t.Errorf("unexpected report: %+v", r)
	}
}

func TestCompute_MissingAndMismatch(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", diff.StatusMatch),
		makeResult("DB_PASS", diff.StatusMismatch),
		makeResult("API_KEY", diff.StatusMissingInB),
		makeResult("LOG_LEVEL", diff.StatusMissingInA),
	}
	r := stats.Compute(results)
	if r.Total != 4 {
		t.Errorf("Total: want 4, got %d", r.Total)
	}
	if r.Matched != 1 {
		t.Errorf("Matched: want 1, got %d", r.Matched)
	}
	if r.Mismatched != 1 {
		t.Errorf("Mismatched: want 1, got %d", r.Mismatched)
	}
	if r.Missing != 2 {
		t.Errorf("Missing: want 2, got %d", r.Missing)
	}
	if r.Problems != 3 {
		t.Errorf("Problems: want 3, got %d", r.Problems)
	}
}

func TestCompute_ByPrefix(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", diff.StatusMatch),
		makeResult("DB_PORT", diff.StatusMismatch),
		makeResult("API_KEY", diff.StatusMissingInB),
		makeResult("NOUNDERSCORE", diff.StatusMatch),
	}
	r := stats.Compute(results)
	if r.ByPrefix["DB"] != 2 {
		t.Errorf("DB prefix: want 2, got %d", r.ByPrefix["DB"])
	}
	if r.ByPrefix["API"] != 1 {
		t.Errorf("API prefix: want 1, got %d", r.ByPrefix["API"])
	}
	if _, ok := r.ByPrefix["NOUNDERSCORE"]; ok {
		t.Error("expected no entry for key without underscore")
	}
}

func TestWrite_ContainsExpectedLines(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", diff.StatusMatch),
		makeResult("DB_PASS", diff.StatusMismatch),
	}
	r := stats.Compute(results)
	var buf bytes.Buffer
	stats.Write(&buf, r)
	out := buf.String()

	for _, want := range []string{"Total", "Matched", "Mismatched", "Problems", "By prefix", "DB"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}
