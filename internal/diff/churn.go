package diff

import (
	"fmt"
	"io"
	"sort"
)

// ChurnEntry records how frequently a key has appeared as a problem
// across multiple comparison runs.
type ChurnEntry struct {
	Key    string
	Runs   int
	Problems int
	Rate   float64
}

// ChurnReport accumulates problem counts per key across runs.
type ChurnReport struct {
	runs   int
	counts map[string]int
}

// NewChurnReport returns an empty ChurnReport.
func NewChurnReport() *ChurnReport {
	return &ChurnReport{counts: make(map[string]int)}
}

// AddRun records one comparison run's results into the report.
func (c *ChurnReport) AddRun(results []Result) {
	c.runs++
	seen := make(map[string]bool)
	for _, r := range results {
		if r.Status != StatusMatch && !seen[r.Key] {
			c.counts[r.Key]++
			seen[r.Key] = true
		}
	}
}

// Build returns a sorted slice of ChurnEntry values, highest rate first.
func (c *ChurnReport) Build() []ChurnEntry {
	entries := make([]ChurnEntry, 0, len(c.counts))
	for key, problems := range c.counts {
		rate := 0.0
		if c.runs > 0 {
			rate = float64(problems) / float64(c.runs)
		}
		entries = append(entries, ChurnEntry{
			Key:      key,
			Runs:     c.runs,
			Problems: problems,
			Rate:     rate,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rate != entries[j].Rate {
			return entries[i].Rate > entries[j].Rate
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// WriteChurn writes a churn report table to w.
func WriteChurn(w io.Writer, entries []ChurnEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No churn data recorded.")
		return
	}
	fmt.Fprintf(w, "%-40s %6s %8s %8s\n", "KEY", "RUNS", "PROBLEMS", "RATE")
	fmt.Fprintf(w, "%-40s %6s %8s %8s\n", "---", "----", "--------", "----")
	for _, e := range entries {
		fmt.Fprintf(w, "%-40s %6d %8d %7.0f%%\n", e.Key, e.Runs, e.Problems, e.Rate*100)
	}
}
