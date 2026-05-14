package diff

import (
	"fmt"
	"io"
	"sort"
)

// ClassifyOptions controls rendering of the severity report.
type ClassifyOptions struct {
	Verbose bool
}

// WriteClassifyReport writes a severity-annotated report of results to w.
func WriteClassifyReport(w io.Writer, results []Result, opts ClassifyOptions) {
	if len(results) == 0 {
		fmt.Fprintln(w, "No results to classify.")
		return
	}

	severities := ClassifyAll(results)

	// Group by severity for the summary header.
	counts := map[Severity]int{}
	for _, s := range severities {
		counts[s]++
	}

	fmt.Fprintf(w, "Severity Summary: high=%d  medium=%d  low=%d  none=%d\n\n",
		counts[SeverityHigh], counts[SeverityMedium], counts[SeverityLow], counts[SeverityNone])

	// Sort results for deterministic output.
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		si := severities[sorted[i].Key]
		sj := severities[sorted[j].Key]
		if si != sj {
			return si > sj // higher severity first
		}
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		sev := severities[r.Key]
		if opts.Verbose {
			fmt.Fprintf(w, "  [%-6s] %-30s  status=%-12s  a=%q  b=%q\n",
				sev, r.Key, r.Status, r.ValueA, r.ValueB)
		} else {
			fmt.Fprintf(w, "  [%-6s] %-30s  status=%s\n",
				sev, r.Key, r.Status)
		}
	}
}
