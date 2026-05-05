package diff

import (
	"fmt"
	"io"
	"sort"
)

const (
	symMatch    = "="
	symMismatch = "~"
	symMissingB = "-"
	symMissingA = "+"
)

// WriteReport writes a human-readable diff report to w.
// fileA and fileB are labels used in the header (typically file paths).
func WriteReport(w io.Writer, results []Result, fileA, fileB string) {
	// Sort results alphabetically by key for deterministic output.
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "envdiff: %s  vs  %s\n", fileA, fileB)
	fmt.Fprintln(w, "Legend: = match  ~ mismatch  - missing in B  + missing in A")
	fmt.Fprintln(w, "---")

	for _, r := range sorted {
		switch r.Status {
		case StatusMatch:
			fmt.Fprintf(w, "  %s  %s\n", symMatch, r.Key)
		case StatusMismatch:
			fmt.Fprintf(w, "  %s  %s\n", symMismatch, r.Key)
			fmt.Fprintf(w, "       A: %q\n", r.ValueA)
			fmt.Fprintf(w, "       B: %q\n", r.ValueB)
		case StatusMissingB:
			fmt.Fprintf(w, "  %s  %s  (only in A: %q)\n", symMissingB, r.Key, r.ValueA)
		case StatusMissingA:
			fmt.Fprintf(w, "  %s  %s  (only in B: %q)\n", symMissingA, r.Key, r.ValueB)
		}
	}

	s := Summarize(results)
	fmt.Fprintln(w, "---")
	fmt.Fprintf(w, "Summary: %d match, %d mismatch, %d missing in B, %d missing in A\n",
		s.Matches, s.Mismatches, s.MissingInB, s.MissingInA)
}
