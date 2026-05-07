// Package stats provides summary statistics over a set of diff results.
package stats

import (
	"io"
	"fmt"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Report holds aggregate counts derived from a slice of diff results.
type Report struct {
	Total     int
	Matched   int
	Missing   int
	Mismatched int
	Problems  int
	ByPrefix  map[string]int
}

// Compute builds a Report from the provided results.
func Compute(results []diff.Result) Report {
	r := Report{
		ByPrefix: make(map[string]int),
	}

	for _, res := range results {
		r.Total++

		switch res.Status {
		case diff.StatusMatch:
			r.Matched++
		case diff.StatusMissingInA, diff.StatusMissingInB:
			r.Missing++
			r.Problems++
		case diff.StatusMismatch:
			r.Mismatched++
			r.Problems++
		}

		if prefix := keyPrefix(res.Key); prefix != "" {
			r.ByPrefix[prefix]++
		}
	}

	return r
}

// keyPrefix returns the portion of a key before the first underscore,
// or an empty string if the key contains no underscore.
func keyPrefix(key string) string {
	for i, ch := range key {
		if ch == '_' && i > 0 {
			return key[:i]
		}
	}
	return ""
}

// Write prints a human-readable statistics summary to w.
func Write(w io.Writer, r Report) {
	fmt.Fprintf(w, "Stats:\n")
	fmt.Fprintf(w, "  Total keys : %d\n", r.Total)
	fmt.Fprintf(w, "  Matched    : %d\n", r.Matched)
	fmt.Fprintf(w, "  Missing    : %d\n", r.Missing)
	fmt.Fprintf(w, "  Mismatched : %d\n", r.Mismatched)
	fmt.Fprintf(w, "  Problems   : %d\n", r.Problems)

	if len(r.ByPrefix) == 0 {
		return
	}

	prefixes := make([]string, 0, len(r.ByPrefix))
	for p := range r.ByPrefix {
		prefixes = append(prefixes, p)
	}
	sort.Strings(prefixes)

	fmt.Fprintf(w, "  By prefix:\n")
	for _, p := range prefixes {
		fmt.Fprintf(w, "    %-20s %d\n", p, r.ByPrefix[p])
	}
}
