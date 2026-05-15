package diff

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

// EntropyEntry holds the Shannon entropy score for a single key's value
// across multiple environments.
type EntropyEntry struct {
	Key     string
	Entropy float64
	Unique  int
	Total   int
}

// BuildEntropy computes per-key Shannon entropy across a set of labeled
// environment maps. High entropy indicates values differ widely across envs.
func BuildEntropy(envs map[string]map[string]string) []EntropyEntry {
	// Collect all keys
	keySet := map[string]struct{}{}
	for _, env := range envs {
		for k := range env {
			keySet[k] = struct{}{}
		}
	}

	entries := make([]EntropyEntry, 0, len(keySet))

	for key := range keySet {
		freq := map[string]int{}
		total := 0
		for _, env := range envs {
			if v, ok := env[key]; ok {
				freq[v]++
				total++
			}
		}
		if total == 0 {
			continue
		}
		entropy := 0.0
		for _, count := range freq {
			p := float64(count) / float64(total)
			entropy -= p * math.Log2(p)
		}
		entries = append(entries, EntropyEntry{
			Key:     key,
			Entropy: entropy,
			Unique:  len(freq),
			Total:   total,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Entropy != entries[j].Entropy {
			return entries[i].Entropy > entries[j].Entropy
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// WriteEntropy writes a human-readable entropy report to w.
func WriteEntropy(w io.Writer, entries []EntropyEntry, topN int) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no keys found across environments")
		return
	}
	limit := len(entries)
	if topN > 0 && topN < limit {
		limit = topN
	}
	fmt.Fprintf(w, "%-40s %8s %7s %6s\n", "KEY", "ENTROPY", "UNIQUE", "TOTAL")
	fmt.Fprintln(w, strings.Repeat("-", 65))
	for _, e := range entries[:limit] {
		fmt.Fprintf(w, "%-40s %8.4f %7d %6d\n", e.Key, e.Entropy, e.Unique, e.Total)
	}
}
