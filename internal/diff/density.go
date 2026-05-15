package diff

import (
	"fmt"
	"io"
	"sort"
)

// DensityEntry holds the problem density for a single key prefix.
type DensityEntry struct {
	Prefix  string
	Total   int
	Problems int
	Density float64 // Problems / Total
}

// DensityReport aggregates density entries across prefixes.
type DensityReport struct {
	Entries []DensityEntry
}

// BuildDensity computes problem density grouped by key prefix (up to the first
// underscore). Keys without an underscore are grouped under "(root)".
func BuildDensity(results []Result) DensityReport {
	totals := map[string]int{}
	problems := map[string]int{}

	for _, r := range results {
		prefix := keyPrefix(r.Key)
		totals[prefix]++
		if r.Status != StatusMatch {
			problems[prefix]++
		}
	}

	entries := make([]DensityEntry, 0, len(totals))
	for prefix, total := range totals {
		p := problems[prefix]
		density := 0.0
		if total > 0 {
			density = float64(p) / float64(total)
		}
		entries = append(entries, DensityEntry{
			Prefix:   prefix,
			Total:    total,
			Problems: p,
			Density:  density,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Density != entries[j].Density {
			return entries[i].Density > entries[j].Density
		}
		return entries[i].Prefix < entries[j].Prefix
	})

	return DensityReport{Entries: entries}
}

// WriteDensity writes a density report table to w.
func WriteDensity(w io.Writer, report DensityReport) {
	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "no data")
		return
	}

	fmt.Fprintf(w, "%-20s  %6s  %8s  %7s\n", "PREFIX", "TOTAL", "PROBLEMS", "DENSITY")
	fmt.Fprintf(w, "%-20s  %6s  %8s  %7s\n", "------", "-----", "--------", "-------")
	for _, e := range report.Entries {
		fmt.Fprintf(w, "%-20s  %6d  %8d  %6.1f%%\n",
			e.Prefix, e.Total, e.Problems, e.Density*100)
	}
}
