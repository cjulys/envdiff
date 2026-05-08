package baseline

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// WriteReport writes a human-readable summary of new vs suppressed results
// to w, given the full current result set and the loaded baseline.
func WriteReport(w io.Writer, current []diff.Result, b *Baseline) {
	new := NewResults(current, b)
	suppressed := len(current) - len(new)

	fmt.Fprintf(w, "Baseline recorded: %s\n", b.CreatedAt.Format(time.RFC3339))
	fmt.Fprintf(w, "Suppressed (known): %d  New: %d\n\n", suppressed, len(new))

	if len(new) == 0 {
		fmt.Fprintln(w, "No new differences detected.")
		return
	}

	sort.Slice(new, func(i, j int) bool {
		return new[i].Key < new[j].Key
	})

	fmt.Fprintln(w, "New differences:")
	for _, r := range new {
		switch r.Status {
		case diff.StatusMismatch:
			fmt.Fprintf(w, "  ~ %s (mismatch)\n", r.Key)
		case diff.StatusMissingInA:
			fmt.Fprintf(w, "  - %s (missing in first file)\n", r.Key)
		case diff.StatusMissingInB:
			fmt.Fprintf(w, "  + %s (missing in second file)\n", r.Key)
		default:
			fmt.Fprintf(w, "  ? %s\n", r.Key)
		}
	}
}
