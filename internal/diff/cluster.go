package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ClusterEntry represents a group of keys that share a common naming pattern.
type ClusterEntry struct {
	Pattern  string
	Keys     []string
	Problems int
	Total    int
}

// ClusterReport holds all discovered clusters across a set of results.
type ClusterReport struct {
	Entries []ClusterEntry
}

// BuildCluster groups result keys by their common prefix segment (e.g. "DB",
// "AWS", "REDIS") and counts how many of each cluster have problems.
func BuildCluster(results []Result) ClusterReport {
	type bucket struct {
		keys     []string
		problems int
	}
	buckets := make(map[string]*bucket)

	for _, r := range results {
		pattern := clusterPrefix(r.Key)
		b, ok := buckets[pattern]
		if !ok {
			b = &bucket{}
			buckets[pattern] = b
		}
		b.keys = append(b.keys, r.Key)
		if r.Status != StatusMatch {
			b.problems++
		}
	}

	entries := make([]ClusterEntry, 0, len(buckets))
	for pattern, b := range buckets {
		sort.Strings(b.keys)
		entries = append(entries, ClusterEntry{
			Pattern:  pattern,
			Keys:     b.keys,
			Problems: b.problems,
			Total:    len(b.keys),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Problems != entries[j].Problems {
			return entries[i].Problems > entries[j].Problems
		}
		return entries[i].Pattern < entries[j].Pattern
	})

	return ClusterReport{Entries: entries}
}

// clusterPrefix returns the first underscore-delimited segment of a key,
// or the full key if no underscore is present.
func clusterPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

// WriteCluster writes a human-readable cluster report to w.
func WriteCluster(w io.Writer, report ClusterReport) {
	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "no keys to cluster")
		return
	}
	fmt.Fprintf(w, "%-20s  %6s  %6s  %s\n", "CLUSTER", "KEYS", "ISSUES", "KEYS")
	fmt.Fprintln(w, strings.Repeat("-", 70))
	for _, e := range report.Entries {
		preview := strings.Join(e.Keys, ", ")
		if len(preview) > 40 {
			preview = preview[:37] + "..."
		}
		fmt.Fprintf(w, "%-20s  %6d  %6d  %s\n", e.Pattern, e.Total, e.Problems, preview)
	}
}
