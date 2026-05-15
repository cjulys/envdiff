package diff

import (
	"fmt"
	"io"
)

// ClusterOptions controls the behaviour of RunCluster.
type ClusterOptions struct {
	// MinKeys filters out clusters with fewer than MinKeys keys.
	MinKeys int
	// OnlyProblems restricts output to clusters that have at least one problem.
	OnlyProblems bool
}

// DefaultClusterOptions returns sensible defaults for RunCluster.
func DefaultClusterOptions() ClusterOptions {
	return ClusterOptions{
		MinKeys:      1,
		OnlyProblems: false,
	}
}

// RunCluster builds a cluster report from results, applies the given options,
// writes the report to w and returns true if any cluster contains problems.
func RunCluster(w io.Writer, results []Result, opts ClusterOptions) bool {
	report := BuildCluster(results)

	filtered := report.Entries[:0]
	for _, e := range report.Entries {
		if e.Total < opts.MinKeys {
			continue
		}
		if opts.OnlyProblems && e.Problems == 0 {
			continue
		}
		filtered = append(filtered, e)
	}
	report.Entries = filtered

	WriteCluster(w, report)

	for _, e := range report.Entries {
		if e.Problems > 0 {
			return true
		}
	}
	return false
}

// ClusterSummary returns a one-line summary string for the report.
func ClusterSummary(report ClusterReport) string {
	totalProblems := 0
	for _, e := range report.Entries {
		totalProblems += e.Problems
	}
	return fmt.Sprintf("%d cluster(s), %d issue(s) across %d key(s)",
		len(report.Entries), totalProblems, totalKeys(report))
}

func totalKeys(report ClusterReport) int {
	n := 0
	for _, e := range report.Entries {
		n += e.Total
	}
	return n
}
