package diff

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func buildTrendResults(statuses ...Status) []Result {
	results := make([]Result, len(statuses))
	for i, s := range statuses {
		results[i] = Result{Key: fmt.Sprintf("KEY_%d", i), Status: s}
	}
	return results
}

func baseTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestTrend_AddRun_CountsProblems(t *testing.T) {
	tr := &Trend{}
	results := []Result{
		{Key: "A", Status: StatusMatch},
		{Key: "B", Status: StatusMismatch},
		{Key: "C", Status: StatusMissingInB},
	}
	tr.AddRun("run1", baseTime(), results)
	if len(tr.Points) != 1 {
		t.Fatalf("expected 1 point, got %d", len(tr.Points))
	}
	p := tr.Points[0]
	if p.Problems != 2 {
		t.Errorf("expected 2 problems, got %d", p.Problems)
	}
	if p.Total != 3 {
		t.Errorf("expected total 3, got %d", p.Total)
	}
}

func TestTrend_Sorted_ByTimestamp(t *testing.T) {
	tr := &Trend{}
	t2 := baseTime().Add(48 * time.Hour)
	t1 := baseTime()
	tr.AddRun("later", t2, []Result{{Key: "X", Status: StatusMatch}})
	tr.AddRun("earlier", t1, []Result{{Key: "Y", Status: StatusMismatch}})

	sorted := tr.Sorted()
	if sorted[0].Label != "earlier" {
		t.Errorf("expected 'earlier' first, got %s", sorted[0].Label)
	}
}

func TestTrend_Direction_Improving(t *testing.T) {
	tr := &Trend{}
	tr.AddRun("r1", baseTime(), []Result{
		{Key: "A", Status: StatusMismatch},
		{Key: "B", Status: StatusMismatch},
	})
	tr.AddRun("r2", baseTime().Add(24*time.Hour), []Result{
		{Key: "A", Status: StatusMatch},
	})
	if tr.Direction() != -1 {
		t.Errorf("expected direction -1 (improving), got %d", tr.Direction())
	}
}

func TestTrend_Direction_Worsening(t *testing.T) {
	tr := &Trend{}
	tr.AddRun("r1", baseTime(), []Result{
		{Key: "A", Status: StatusMatch},
	})
	tr.AddRun("r2", baseTime().Add(24*time.Hour), []Result{
		{Key: "A", Status: StatusMismatch},
		{Key: "B", Status: StatusMissingInB},
	})
	if tr.Direction() != 1 {
		t.Errorf("expected direction 1 (worsening), got %d", tr.Direction())
	}
}

func TestWriteTrend_EmptyPrintsMessage(t *testing.T) {
	var buf bytes.Buffer
	WriteTrend(&buf, &Trend{})
	if !strings.Contains(buf.String(), "no trend data") {
		t.Errorf("expected 'no trend data' message, got: %s", buf.String())
	}
}

func TestWriteTrend_ContainsLabelsAndArrow(t *testing.T) {
	tr := &Trend{}
	tr.AddRun("prod", baseTime(), []Result{
		{Key: "A", Status: StatusMismatch},
	})
	tr.AddRun("staging", baseTime().Add(24*time.Hour), []Result{
		{Key: "A", Status: StatusMatch},
	})
	var buf bytes.Buffer
	WriteTrend(&buf, tr)
	out := buf.String()
	if !strings.Contains(out, "prod") {
		t.Errorf("expected 'prod' in output")
	}
	if !strings.Contains(out, "improving") {
		t.Errorf("expected 'improving' trend arrow, got: %s", out)
	}
}
