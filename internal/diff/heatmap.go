package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// HeatmapEntry records how many times a key has appeared as a problem
// across multiple comparison runs.
type HeatmapEntry struct {
	Key   string
	Count int
}

// BuildHeatmap counts problem occurrences per key across the given result
// slices. Each slice typically represents one comparison run.
func BuildHeatmap(runs [][]Result) []HeatmapEntry {
	counts := make(map[string]int)
	for _, results := range runs {
		for _, r := range results {
			if r.Status != StatusMatch {
				counts[r.Key]++
			}
		}
	}

	entries := make([]HeatmapEntry, 0, len(counts))
	for k, v := range counts {
		entries = append(entries, HeatmapEntry{Key: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// WriteHeatmap writes a text heatmap to w showing keys ranked by problem
// frequency. Each row includes a simple bar proportional to its count.
func WriteHeatmap(entries []HeatmapEntry, w io.Writer) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no problem keys recorded")
		return
	}

	max := entries[0].Count
	barWidth := 20

	fmt.Fprintln(w, "Key Heatmap (most frequent problems)")
	fmt.Fprintln(w, strings.Repeat("-", 50))
	for _, e := range entries {
		filled := 0
		if max > 0 {
			filled = (e.Count * barWidth) / max
		}
		bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
		fmt.Fprintf(w, "  %-30s %s %d\n", e.Key, bar, e.Count)
	}
}
