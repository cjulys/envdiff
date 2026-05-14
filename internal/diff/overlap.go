package diff

import (
	"fmt"
	"io"
	"sort"
)

// OverlapEntry records how many keys two environments share and how many conflict.
type OverlapEntry struct {
	LabelA    string
	LabelB    string
	Shared    int
	Conflicts int
	OnlyInA   int
	OnlyInB   int
}

// OverlapScore returns a value in [0,1] representing key-level similarity.
func (e OverlapEntry) OverlapScore() float64 {
	total := e.Shared + e.OnlyInA + e.OnlyInB
	if total == 0 {
		return 1.0
	}
	return float64(e.Shared) / float64(total)
}

// BuildOverlap computes pairwise overlap statistics from a set of named env maps.
// Each entry in envs maps a label to its key→value map.
func BuildOverlap(envs map[string]map[string]string) []OverlapEntry {
	labels := make([]string, 0, len(envs))
	for l := range envs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	var entries []OverlapEntry
	for i := 0; i < len(labels); i++ {
		for j := i + 1; j < len(labels); j++ {
			a, b := labels[i], labels[j]
			entry := computeOverlap(a, envs[a], b, envs[b])
			entries = append(entries, entry)
		}
	}
	return entries
}

func computeOverlap(labelA string, mapA map[string]string, labelB string, mapB map[string]string) OverlapEntry {
	e := OverlapEntry{LabelA: labelA, LabelB: labelB}
	for k, va := range mapA {
		if vb, ok := mapB[k]; ok {
			e.Shared++
			if va != vb {
				e.Conflicts++
			}
		} else {
			e.OnlyInA++
		}
	}
	for k := range mapB {
		if _, ok := mapA[k]; !ok {
			e.OnlyInB++
		}
	}
	return e
}

// WriteOverlap writes a human-readable overlap table to w.
func WriteOverlap(w io.Writer, entries []OverlapEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "no environment pairs to compare")
		return
	}
	fmt.Fprintf(w, "%-20s %-20s %7s %9s %7s %7s %8s\n",
		"ENV A", "ENV B", "SHARED", "CONFLICTS", "ONLY-A", "ONLY-B", "SCORE")
	fmt.Fprintln(w, "--------------------------------------------------------------------------------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-20s %-20s %7d %9d %7d %7d %7.0f%%\n",
			e.LabelA, e.LabelB, e.Shared, e.Conflicts, e.OnlyInA, e.OnlyInB,
			e.OverlapScore()*100)
	}
}
