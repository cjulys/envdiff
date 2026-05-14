package diff

import (
	"fmt"
	"io"
)

// BloomOptions controls the behaviour of RunBloom.
type BloomOptions struct {
	Threshold float64 // frequency threshold [0,1] above which a key is flagged
	Verbose   bool
}

// DefaultBloomOptions returns sensible defaults for BloomOptions.
func DefaultBloomOptions() BloomOptions {
	return BloomOptions{Threshold: 0.5}
}

// RunBloom builds a bloom frequency report from multiple labeled result sets
// and writes the output to w. It returns true if any key meets or exceeds the
// threshold, suitable for use as a non-zero exit code signal.
func RunBloom(w io.Writer, runs [][]Result, opts BloomOptions) bool {
	report := NewBloomReport()
	for _, r := range runs {
		report.AddRun(r)
	}
	WriteBloom(w, report, opts.Threshold)
	if opts.Verbose {
		fmt.Fprintln(w)
		for _, e := range report.Sorted() {
			if e.Freq >= opts.Threshold {
				fmt.Fprintf(w, "  flagged: %s (%.0f%%)\n", e.Key, e.Freq*100)
			}
		}
	}
	for _, e := range report.Sorted() {
		if e.Freq >= opts.Threshold {
			return true
		}
	}
	return false
}
