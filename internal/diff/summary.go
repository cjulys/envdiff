package diff

import "fmt"

// Stats holds aggregate counts derived from a slice of Results.
type Stats struct {
	Total    int
	Matched  int
	Missing  int
	Mismatch int
}

// ComputeStats returns a Stats summary for the provided results.
func ComputeStats(results []Result) Stats {
	s := Stats{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			s.Matched++
		case StatusMissingInA, StatusMissingInB:
			s.Missing++
		case StatusMismatch:
			s.Mismatch++
		}
	}
	return s
}

// String returns a human-readable one-line summary.
func (s Stats) String() string {
	return fmt.Sprintf(
		"total=%d matched=%d missing=%d mismatch=%d",
		s.Total, s.Matched, s.Missing, s.Mismatch,
	)
}

// HasProblems returns true when there are any missing or mismatched keys.
func (s Stats) HasProblems() bool {
	return s.Missing > 0 || s.Mismatch > 0
}

// ProblemCount returns the total number of non-matching results.
func (s Stats) ProblemCount() int {
	return s.Missing + s.Mismatch
}
