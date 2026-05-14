package diff

import (
	"fmt"
	"io"
)

// OverlapOptions controls the behaviour of RunOverlap.
type OverlapOptions struct {
	// MinScore filters out pairs whose overlap score is at or above this threshold.
	// Set to 0 to show all pairs.
	MinScore float64
	// Verbose prints an explanatory header.
	Verbose bool
}

// DefaultOverlapOptions returns sensible defaults for RunOverlap.
func DefaultOverlapOptions() OverlapOptions {
	return OverlapOptions{
		MinScore: 0,
		Verbose:  false,
	}
}

// RunOverlap computes pairwise overlap for the supplied environments and writes
// the result table to w. It returns true when at least one pair has conflicts
// or exclusive keys (i.e. the environments are not identical).
func RunOverlap(w io.Writer, envs map[string]map[string]string, opts OverlapOptions) bool {
	entries := BuildOverlap(envs)

	// Apply minimum-score filter.
	if opts.MinScore > 0 {
		filtered := entries[:0]
		for _, e := range entries {
			if e.OverlapScore() < opts.MinScore {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if opts.Verbose {
		fmt.Fprintln(w, "# Environment Overlap Report")
		fmt.Fprintf(w, "# %d environment(s), %d pair(s) analysed\n\n", len(envs), len(entries))
	}

	WriteOverlap(w, entries)

	hasProblems := false
	for _, e := range entries {
		if e.Conflicts > 0 || e.OnlyInA > 0 || e.OnlyInB > 0 {
			hasProblems = true
			break
		}
	}
	return hasProblems
}
