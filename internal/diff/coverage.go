package diff

import (
	"fmt"
	"io"
	"sort"
)

// CoverageEntry holds the key-presence statistics for a single environment.
type CoverageEntry struct {
	Label   string
	Total   int
	Present int
	Rate    float64 // Present / Total
}

// BuildCoverage computes, for each labelled environment, how many of the
// union of all keys are actually present (i.e. not missing).
//
// runs is a map of label → key→value pairs for that environment.
// The union of all keys across every run forms the denominator.
func BuildCoverage(runs map[string]map[string]string) []CoverageEntry {
	if len(runs) == 0 {
		return nil
	}

	// Build union key set.
	union := map[string]struct{}{}
	for _, env := range runs {
		for k := range env {
			union[k] = struct{}{}
		}
	}
	total := len(union)

	entries := make([]CoverageEntry, 0, len(runs))
	for label, env := range runs {
		present := 0
		for k := range union {
			if _, ok := env[k]; ok {
				present++
			}
		}
		rate := 0.0
		if total > 0 {
			rate = float64(present) / float64(total)
		}
		entries = append(entries, CoverageEntry{
			Label:   label,
			Total:   total,
			Present: present,
			Rate:    rate,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rate != entries[j].Rate {
			return entries[i].Rate > entries[j].Rate
		}
		return entries[i].Label < entries[j].Label
	})
	return entries
}

// WriteCoverage writes a human-readable coverage table to w.
func WriteCoverage(w io.Writer, entries []CoverageEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no environments to report")
		return
	}
	fmt.Fprintf(w, "%-24s  %7s  %7s  %7s\n", "Environment", "Present", "Total", "Coverage")
	fmt.Fprintf(w, "%-24s  %7s  %7s  %7s\n",
		"------------------------", "-------", "-------", "--------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-24s  %7d  %7d  %6.1f%%\n",
			e.Label, e.Present, e.Total, e.Rate*100)
	}
}
