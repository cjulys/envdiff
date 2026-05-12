package diff

import (
	"strings"
	"testing"
)

func buildScoreResult(key string, status Status) Result {
	return Result{Key: key, Status: status}
}

func TestComputeScore_Empty(t *testing.T) {
	s := ComputeScore(nil)
	if s.Value != 100.0 {
		t.Errorf("expected 100, got %.1f", s.Value)
	}
	if s.Grade != "A" {
		t.Errorf("expected grade A, got %s", s.Grade)
	}
}

func TestComputeScore_AllMatch(t *testing.T) {
	results := []Result{
		buildScoreResult("A", StatusMatch),
		buildScoreResult("B", StatusMatch),
	}
	s := ComputeScore(results)
	if s.Value != 100.0 {
		t.Errorf("expected 100, got %.1f", s.Value)
	}
	if s.Problems != 0 {
		t.Errorf("expected 0 problems, got %d", s.Problems)
	}
}

func TestComputeScore_AllProblems(t *testing.T) {
	results := []Result{
		buildScoreResult("A", StatusMissingInB),
		buildScoreResult("B", StatusMismatch),
	}
	s := ComputeScore(results)
	if s.Value != 0.0 {
		t.Errorf("expected 0, got %.1f", s.Value)
	}
	if s.Grade != "F" {
		t.Errorf("expected grade F, got %s", s.Grade)
	}
}

func TestComputeScore_Mixed(t *testing.T) {
	results := []Result{
		buildScoreResult("A", StatusMatch),
		buildScoreResult("B", StatusMatch),
		buildScoreResult("C", StatusMatch),
		buildScoreResult("D", StatusMismatch),
	}
	s := ComputeScore(results)
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
	if s.Problems != 1 {
		t.Errorf("expected 1 problem, got %d", s.Problems)
	}
	if s.Value != 75.0 {
		t.Errorf("expected 75.0, got %.1f", s.Value)
	}
	if s.Grade != "C" {
		t.Errorf("expected grade C, got %s", s.Grade)
	}
}

func TestWriteScore_ContainsGrade(t *testing.T) {
	results := []Result{
		buildScoreResult("X", StatusMatch),
	}
	s := ComputeScore(results)
	var sb strings.Builder
	WriteScore(&sb, s, "production")
	out := sb.String()
	if !strings.Contains(out, "Grade: A") {
		t.Errorf("expected grade A in output, got:\n%s", out)
	}
	if !strings.Contains(out, "production") {
		t.Errorf("expected label 'production' in output, got:\n%s", out)
	}
}

func TestWriteScore_DefaultLabel(t *testing.T) {
	var sb strings.Builder
	WriteScore(&sb, Score{Value: 80, Grade: "B", Total: 5, Problems: 1}, "")
	out := sb.String()
	if !strings.Contains(out, "overall") {
		t.Errorf("expected 'overall' label, got:\n%s", out)
	}
}
