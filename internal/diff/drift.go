package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// DriftEntry records the problem count for a named environment at a point in time.
type DriftEntry struct {
	Label     string
	Timestamp time.Time
	Problems  int
	Total     int
}

// DriftReport holds drift entries across multiple environments.
type DriftReport struct {
	Entries []DriftEntry
}

// AddEntry appends a drift entry derived from a slice of Results.
func (d *DriftReport) AddEntry(label string, results []Result) {
	problems := 0
	for _, r := range results {
		if r.Status != StatusMatch {
			problems++
		}
	}
	d.Entries = append(d.Entries, DriftEntry{
		Label:     label,
		Timestamp: time.Now(),
		Problems:  problems,
		Total:     len(results),
	})
}

// SortedEntries returns entries sorted by label then timestamp.
func (d *DriftReport) SortedEntries() []DriftEntry {
	copy := make([]DriftEntry, len(d.Entries))
	for i, e := range d.Entries {
		copy[i] = e
	}
	sort.Slice(copy, func(i, j int) bool {
		if copy[i].Label != copy[j].Label {
			return copy[i].Label < copy[j].Label
		}
		return copy[i].Timestamp.Before(copy[j].Timestamp)
	})
	return copy
}

// WriteDrift writes a human-readable drift report to w.
func WriteDrift(w io.Writer, report *DriftReport) {
	entries := report.SortedEntries()
	if len(entries) == 0 {
		fmt.Fprintln(w, "No drift data recorded.")
		return
	}
	fmt.Fprintln(w, "Environment Drift Report")
	fmt.Fprintln(w, "========================")
	for _, e := range entries {
		pct := 0.0
		if e.Total > 0 {
			pct = float64(e.Problems) / float64(e.Total) * 100
		}
		fmt.Fprintf(w, "  %-20s  problems: %3d / %3d  (%.1f%%)  [%s]\n",
			e.Label, e.Problems, e.Total, pct,
			e.Timestamp.Format("2006-01-02 15:04:05"))
	}
}
