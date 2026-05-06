package schema

import (
	"fmt"
	"io"
	"sort"
)

// WriteReport writes a human-readable summary of schema violations to w.
// If there are no violations it prints a success message.
func WriteReport(w io.Writer, violations []Violation, file string) {
	fmt.Fprintf(w, "Schema check: %s\n", file)
	if len(violations) == 0 {
		fmt.Fprintln(w, "  ✓ all keys pass schema validation")
		return
	}

	// Sort for deterministic output.
	sorted := make([]Violation, len(violations))
	copy(sorted, violations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	fmt.Fprintf(w, "  %d violation(s) found:\n", len(sorted))
	for _, v := range sorted {
		fmt.Fprintf(w, "  ✗ %-30s %s\n", v.Key, v.Message)
	}
}
