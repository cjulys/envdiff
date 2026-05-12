package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// TimelineEntry records a snapshot of diff results at a point in time.
type TimelineEntry struct {
	Timestamp time.Time
	Label     string
	Stats     Stats
}

// Timeline holds an ordered sequence of diff snapshots.
type Timeline struct {
	Entries []TimelineEntry
}

// Add appends a new entry to the timeline.
func (t *Timeline) Add(label string, results []Result) {
	t.Entries = append(t.Entries, TimelineEntry{
		Timestamp: time.Now(),
		Label:     label,
		Stats:     ComputeStats(results),
	})
}

// Sorted returns entries ordered by timestamp ascending.
func (t *Timeline) Sorted() []TimelineEntry {
	copy := make([]TimelineEntry, len(t.Entries))
	for i, e := range t.Entries {
		copy[i] = e
	}
	sort.Slice(copy, func(i, j int) bool {
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// WriteTimeline writes a human-readable timeline table to w.
func WriteTimeline(w io.Writer, t *Timeline) {
	entries := t.Sorted()
	if len(entries) == 0 {
		fmt.Fprintln(w, "No timeline entries recorded.")
		return
	}
	fmt.Fprintf(w, "%-30s %-20s %8s %8s %8s %8s\n",
		"Timestamp", "Label", "Total", "Match", "Missing", "Mismatch")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 90))
	for _, e := range entries {
		fmt.Fprintf(w, "%-30s %-20s %8d %8d %8d %8d\n",
			e.Timestamp.Format(time.RFC3339),
			e.Label,
			e.Stats.Total,
			e.Stats.Match,
			e.Stats.Missing,
			e.Stats.Mismatch,
		)
	}
}

func repeatChar(ch rune, n int) string {
	out := make([]rune, n)
	for i := range out {
		out[i] = ch
	}
	return string(out)
}
