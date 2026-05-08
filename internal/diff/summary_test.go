package diff

import (
	"strings"
	"testing"
)

func buildStats(statuses ...Status) []Result {
	results := make([]Result, len(statuses))
	for i, s := range statuses {
		results[i] = Result{Key: "KEY", Status: s}
	}
	return results
}

func TestComputeStats_Empty(t *testing.T) {
	s := ComputeStats(nil)
	if s.Total != 0 || s.Matched != 0 || s.Missing != 0 || s.Mismatch != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestComputeStats_AllMatch(t *testing.T) {
	s := ComputeStats(buildStats(StatusMatch, StatusMatch))
	if s.Matched != 2 || s.HasProblems() {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestComputeStats_MixedStatuses(t *testing.T) {
	results := buildStats(StatusMatch, StatusMissingInA, StatusMissingInB, StatusMismatch)
	s := ComputeStats(results)
	if s.Total != 4 {
		t.Errorf("expected total=4, got %d", s.Total)
	}
	if s.Matched != 1 {
		t.Errorf("expected matched=1, got %d", s.Matched)
	}
	if s.Missing != 2 {
		t.Errorf("expected missing=2, got %d", s.Missing)
	}
	if s.Mismatch != 1 {
		t.Errorf("expected mismatch=1, got %d", s.Mismatch)
	}
}

func TestComputeStats_HasProblems(t *testing.T) {
	if ComputeStats(buildStats(StatusMatch)).HasProblems() {
		t.Error("expected no problems for all-match")
	}
	if !ComputeStats(buildStats(StatusMismatch)).HasProblems() {
		t.Error("expected problems for mismatch")
	}
}

func TestComputeStats_ProblemCount(t *testing.T) {
	s := ComputeStats(buildStats(StatusMissingInA, StatusMissingInB, StatusMismatch, StatusMatch))
	if s.ProblemCount() != 3 {
		t.Errorf("expected problem count 3, got %d", s.ProblemCount())
	}
}

func TestStats_String(t *testing.T) {
	s := Stats{Total: 4, Matched: 1, Missing: 2, Mismatch: 1}
	str := s.String()
	for _, want := range []string{"total=4", "matched=1", "missing=2", "mismatch=1"} {
		if !strings.Contains(str, want) {
			t.Errorf("String() missing %q in %q", want, str)
		}
	}
}
