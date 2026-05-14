package diff

import (
	"bytes"
	"strings"
	"testing"
)

// buildMatrixRuns returns a simple set of labelled results for matrix tests.
func buildMatrixRuns() map[string][]Result {
	// Simulate three environments each returning a flat list of results
	// that represent their values for a shared key set.
	return map[string][]Result{
		"prod": {
			{Key: "DB_HOST", Status: StatusMatch, ValueA: "db.prod", ValueB: "db.prod"},
			{Key: "API_KEY", Status: StatusMatch, ValueA: "abc", ValueB: "abc"},
			{Key: "TIMEOUT", Status: StatusMatch, ValueA: "30", ValueB: "30"},
		},
		"staging": {
			{Key: "DB_HOST", Status: StatusMatch, ValueA: "db.staging", ValueB: "db.staging"},
			{Key: "API_KEY", Status: StatusMismatch, ValueA: "xyz", ValueB: "abc"},
			{Key: "TIMEOUT", Status: StatusMatch, ValueA: "30", ValueB: "30"},
		},
		"dev": {
			{Key: "DB_HOST", Status: StatusMatch, ValueA: "localhost", ValueB: "localhost"},
			{Key: "API_KEY", Status: StatusMissingInA, ValueA: "", ValueB: "abc"},
			{Key: "TIMEOUT", Status: StatusMismatch, ValueA: "5", ValueB: "30"},
		},
	}
}

func TestBuildMatrix_LabelsAreSorted(t *testing.T) {
	runs := buildMatrixRuns()
	m := BuildMatrix(runs)

	if len(m.Labels) != 3 {
		t.Fatalf("expected 3 labels, got %d", len(m.Labels))
	}
	expected := []string{"dev", "prod", "staging"}
	for i, lbl := range expected {
		if m.Labels[i] != lbl {
			t.Errorf("label[%d]: want %q, got %q", i, lbl, m.Labels[i])
		}
	}
}

func TestBuildMatrix_NoDiagonalCells(t *testing.T) {
	runs := buildMatrixRuns()
	m := BuildMatrix(runs)

	for _, lbl := range m.Labels {
		if _, ok := m.Cells[cellKey(lbl, lbl)]; ok {
			t.Errorf("diagonal cell %q should not exist", lbl)
		}
	}
}

func TestBuildMatrix_CellCount(t *testing.T) {
	runs := buildMatrixRuns()
	m := BuildMatrix(runs)

	// N*(N-1) cells for N=3 → 6
	if len(m.Cells) != 6 {
		t.Errorf("expected 6 cells, got %d", len(m.Cells))
	}
}

func TestMatrixCell_Score_AllMatch(t *testing.T) {
	c := MatrixCell{Matched: 3, Total: 3}
	if c.Score() != 100 {
		t.Errorf("expected 100, got %d", c.Score())
	}
}

func TestMatrixCell_Score_ZeroTotal(t *testing.T) {
	c := MatrixCell{Total: 0}
	if c.Score() != 100 {
		t.Errorf("expected 100 for empty cell, got %d", c.Score())
	}
}

func TestMatrixCell_Score_Partial(t *testing.T) {
	c := MatrixCell{Matched: 1, Mismatch: 1, Missing: 1, Total: 3}
	if c.Score() != 33 {
		t.Errorf("expected 33, got %d", c.Score())
	}
}

func TestWriteMatrix_ContainsLabels(t *testing.T) {
	runs := buildMatrixRuns()
	m := BuildMatrix(runs)

	var buf bytes.Buffer
	WriteMatrix(&buf, m)
	out := buf.String()

	for _, lbl := range []string{"dev", "prod", "staging"} {
		if !strings.Contains(out, lbl) {
			t.Errorf("output missing label %q", lbl)
		}
	}
}

func TestWriteMatrix_EmptyPrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	WriteMatrix(&buf, Matrix{})
	if !strings.Contains(buf.String(), "no environments") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteMatrix_ContainsPercentSign(t *testing.T) {
	runs := buildMatrixRuns()
	m := BuildMatrix(runs)

	var buf bytes.Buffer
	WriteMatrix(&buf, m)
	if !strings.Contains(buf.String(), "%") {
		t.Error("expected percentage scores in output")
	}
}
