package diff

import (
	"fmt"
	"io"
)

// WriteSummaryTable writes a formatted ASCII table of Stats to w.
// If label is non-empty it is printed as a heading above the table.
func WriteSummaryTable(w io.Writer, s Stats, label string) {
	if label != "" {
		fmt.Fprintf(w, "Summary: %s\n", label)
	}
	fmt.Fprintln(w, "+----------+-------+")
	fmt.Fprintln(w, "| Metric   | Count |")
	fmt.Fprintln(w, "+----------+-------+")
	rows := []struct {
		name  string
		value int
	}{
		{"Total", s.Total},
		{"Matched", s.Matched},
		{"Missing", s.Missing},
		{"Mismatch", s.Mismatch},
	}
	for _, row := range rows {
		fmt.Fprintf(w, "| %-8s | %5d |\n", row.name, row.value)
	}
	fmt.Fprintln(w, "+----------+-------+")
	if s.HasProblems() {
		fmt.Fprintf(w, "Problems: %d\n", s.ProblemCount())
	} else {
		fmt.Fprintln(w, "No problems found.")
	}
}

// WriteSummaryLine writes a compact single-line summary of Stats to w.
// It is useful for logging or brief status output.
// Example output: "total=10 matched=8 missing=1 mismatch=1 problems=2"
func WriteSummaryLine(w io.Writer, s Stats) {
	fmt.Fprintf(w, "total=%d matched=%d missing=%d mismatch=%d problems=%d\n",
		s.Total, s.Matched, s.Missing, s.Mismatch, s.ProblemCount())
}
