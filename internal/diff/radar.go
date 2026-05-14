package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// RadarEntry holds the problem rate for a single environment label.
type RadarEntry struct {
	Label       string
	Total       int
	Problems    int
	Rate        float64 // 0.0–1.0
}

// RadarReport is a collection of entries across environments.
type RadarReport struct {
	Entries []RadarEntry
}

// BuildRadar computes per-environment problem rates from named result sets.
func BuildRadar(envs map[string][]Result) RadarReport {
	var entries []RadarEntry
	for label, results := range envs {
		total := len(results)
		problems := 0
		for _, r := range results {
			if r.IsProblem() {
				problems++
			}
		}
		rate := 0.0
		if total > 0 {
			rate = float64(problems) / float64(total)
		}
		entries = append(entries, RadarEntry{
			Label:    label,
			Total:    total,
			Problems: problems,
			Rate:     rate,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Rate != entries[j].Rate {
			return entries[i].Rate > entries[j].Rate
		}
		return entries[i].Label < entries[j].Label
	})
	return RadarReport{Entries: entries}
}

const radarBarWidth = 20

// WriteRadar writes a text-based radar (bar chart) of problem rates per environment.
func WriteRadar(w io.Writer, report RadarReport) {
	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "no environments to display")
		return
	}
	fmt.Fprintln(w, "Environment Problem Rate Radar")
	fmt.Fprintln(w, strings.Repeat("─", 50))
	for _, e := range report.Entries {
		filled := int(e.Rate * radarBarWidth)
		bar := strings.Repeat("█", filled) + strings.Repeat("░", radarBarWidth-filled)
		fmt.Fprintf(w, "  %-20s [%s] %5.1f%%  (%d/%d)\n",
			e.Label, bar, e.Rate*100, e.Problems, e.Total)
	}
	fmt.Fprintln(w, strings.Repeat("─", 50))
}
