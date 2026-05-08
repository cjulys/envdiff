package diff

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteSummaryTable_NoProblems(t *testing.T) {
	var buf bytes.Buffer
	s := Stats{Total: 3, Matched: 3}
	WriteSummaryTable(&buf, s, "")
	out := buf.String()
	if !strings.Contains(out, "No problems found.") {
		t.Errorf("expected 'No problems found.' in output, got:\n%s", out)
	}
}

func TestWriteSummaryTable_WithLabel(t *testing.T) {
	var buf bytes.Buffer
	s := Stats{Total: 1, Matched: 1}
	WriteSummaryTable(&buf, s, "prod vs staging")
	if !strings.Contains(buf.String(), "prod vs staging") {
		t.Errorf("expected label in output")
	}
}

func TestWriteSummaryTable_ShowsProblemCount(t *testing.T) {
	var buf bytes.Buffer
	s := Stats{Total: 5, Matched: 2, Missing: 2, Mismatch: 1}
	WriteSummaryTable(&buf, s, "")
	out := buf.String()
	if !strings.Contains(out, "Problems: 3") {
		t.Errorf("expected 'Problems: 3' in output, got:\n%s", out)
	}
}

func TestWriteSummaryTable_ContainsAllMetrics(t *testing.T) {
	var buf bytes.Buffer
	s := Stats{Total: 4, Matched: 1, Missing: 2, Mismatch: 1}
	WriteSummaryTable(&buf, s, "")
	out := buf.String()
	for _, want := range []string{"Total", "Matched", "Missing", "Mismatch"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in table output", want)
		}
	}
}

func TestWriteSummaryTable_TableBorders(t *testing.T) {
	var buf bytes.Buffer
	WriteSummaryTable(&buf, Stats{}, "")
	out := buf.String()
	if !strings.Contains(out, "+----------+-------+") {
		t.Errorf("expected table border in output, got:\n%s", out)
	}
}
