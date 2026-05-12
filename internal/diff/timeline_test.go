package diff_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
)

func buildTimelineResults(statuses ...string) []diff.Result {
	results := make([]diff.Result, len(statuses))
	for i, s := range statuses {
		results[i] = diff.Result{
			Key:    fmt.Sprintf("KEY_%d", i),
			Status: s,
		}
	}
	return results
}

func TestTimeline_AddAndSorted(t *testing.T) {
	tl := &diff.Timeline{}
	tl.Add("prod", []diff.Result{{Key: "A", Status: diff.StatusMatch}})
	tl.Add("staging", []diff.Result{{Key: "B", Status: diff.StatusMismatch}})

	entries := tl.Sorted()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Timestamp.Before(entries[1].Timestamp) && entries[0].Timestamp != entries[1].Timestamp {
		t.Error("entries not in ascending order")
	}
}

func TestTimeline_StatsAreRecorded(t *testing.T) {
	tl := &diff.Timeline{}
	results := []diff.Result{
		{Key: "A", Status: diff.StatusMatch},
		{Key: "B", Status: diff.StatusMismatch},
		{Key: "C", Status: diff.StatusMissingInB},
	}
	tl.Add("env", results)

	if len(tl.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	s := tl.Entries[0].Stats
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
	if s.Match != 1 {
		t.Errorf("expected Match=1, got %d", s.Match)
	}
	if s.Mismatch != 1 {
		t.Errorf("expected Mismatch=1, got %d", s.Mismatch)
	}
}

func TestWriteTimeline_EmptyPrintsMessage(t *testing.T) {
	tl := &diff.Timeline{}
	var buf bytes.Buffer
	diff.WriteTimeline(&buf, tl)
	if !strings.Contains(buf.String(), "No timeline") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteTimeline_ContainsLabel(t *testing.T) {
	tl := &diff.Timeline{}
	tl.Add("production", []diff.Result{{Key: "X", Status: diff.StatusMatch}})
	var buf bytes.Buffer
	diff.WriteTimeline(&buf, tl)
	if !strings.Contains(buf.String(), "production") {
		t.Errorf("expected label 'production' in output, got:\n%s", buf.String())
	}
}

func TestWriteTimeline_ContainsTimestamp(t *testing.T) {
	tl := &diff.Timeline{}
	before := time.Now()
	tl.Add("dev", nil)
	var buf bytes.Buffer
	diff.WriteTimeline(&buf, tl)
	year := before.Format("2006")
	if !strings.Contains(buf.String(), year) {
		t.Errorf("expected year %s in output", year)
	}
}
