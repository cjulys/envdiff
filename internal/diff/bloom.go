package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// BloomEntry records how frequently a key appears as a problem across multiple runs.
type BloomEntry struct {
	Key      string
	Hits     int
	Total    int
	Freq     float64 // Hits / Total
}

// BloomReport aggregates problem frequency across a set of labeled runs.
type BloomReport struct {
	entries map[string]*BloomEntry
	total   int
}

// NewBloomReport creates an empty BloomReport.
func NewBloomReport() *BloomReport {
	return &BloomReport{entries: make(map[string]*BloomEntry)}
}

// AddRun records one run's results into the bloom report.
func (b *BloomReport) AddRun(results []Result) {
	b.total++
	seen := make(map[string]bool)
	for _, r := range results {
		if !r.IsProblem() || seen[r.Key] {
			continue
		}
		seen[r.Key] = true
		e, ok := b.entries[r.Key]
		if !ok {
			e = &BloomEntry{Key: r.Key, Total: 0}
			b.entries[r.Key] = e
		}
		e.Hits++
	}
	for _, e := range b.entries {
		e.Total = b.total
		e.Freq = float64(e.Hits) / float64(e.Total)
	}
}

// Sorted returns entries ordered by frequency descending, then key ascending.
func (b *BloomReport) Sorted() []BloomEntry {
	out := make([]BloomEntry, 0, len(b.entries))
	for _, e := range b.entries {
		out = append(out, *e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Freq != out[j].Freq {
			return out[i].Freq > out[j].Freq
		}
		return out[i].Key < out[j].Key
	})
	return out
}

// WriteBloom writes a human-readable bloom frequency report to w.
func WriteBloom(w io.Writer, b *BloomReport, threshold float64) {
	entries := b.Sorted()
	if len(entries) == 0 {
		fmt.Fprintln(w, "No problem keys recorded.")
		return
	}
	fmt.Fprintf(w, "Bloom Frequency Report (%d run(s), threshold=%.0f%%)\n", b.total, threshold*100)
	fmt.Fprintln(w, strings.Repeat("-", 52))
	for _, e := range entries {
		marker := " "
		if e.Freq >= threshold {
			marker = "!"
		}
		fmt.Fprintf(w, "%s %-30s %3d/%3d  (%.0f%%)\n", marker, e.Key, e.Hits, e.Total, e.Freq*100)
	}
}
