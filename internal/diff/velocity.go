package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// VelocityRun represents a single snapshot of diff results at a point in time.
type VelocityRun struct {
	Label     string
	Timestamp time.Time
	Results   []Result
}

// VelocityEntry holds computed rate-of-change metrics for a key.
type VelocityEntry struct {
	Key            string
	Appearances    int
	FirstSeen      time.Time
	LastSeen       time.Time
	ChangeRate     float64 // problems per day
	TotalProblems  int
}

// BuildVelocity computes per-key change velocity across a series of runs.
func BuildVelocity(runs []VelocityRun) []VelocityEntry {
	type keyStats struct {
		appearances   int
		totalProblems int
		firstSeen     time.Time
		lastSeen      time.Time
	}

	index := make(map[string]*keyStats)

	for _, run := range runs {
		for _, r := range run.Results {
			if !r.IsProblem() {
				continue
			}
			s, ok := index[r.Key]
			if !ok {
				s = &keyStats{firstSeen: run.Timestamp, lastSeen: run.Timestamp}
				index[r.Key] = s
			}
			s.appearances++
			s.totalProblems++
			if run.Timestamp.Before(s.firstSeen) {
				s.firstSeen = run.Timestamp
			}
			if run.Timestamp.After(s.lastSeen) {
				s.lastSeen = run.Timestamp
			}
		}
	}

	entries := make([]VelocityEntry, 0, len(index))
	for key, s := range index {
		span := s.lastSeen.Sub(s.firstSeen).Hours() / 24.0
		rate := 0.0
		if span > 0 {
			rate = float64(s.totalProblems) / span
		}
		entries = append(entries, VelocityEntry{
			Key:           key,
			Appearances:   s.appearances,
			FirstSeen:     s.firstSeen,
			LastSeen:      s.lastSeen,
			ChangeRate:    rate,
			TotalProblems: s.totalProblems,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].ChangeRate != entries[j].ChangeRate {
			return entries[i].ChangeRate > entries[j].ChangeRate
		}
		return entries[i].Key < entries[j].Key
	})

	return entries
}

// WriteVelocity writes a velocity report to w.
func WriteVelocity(w io.Writer, entries []VelocityEntry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No problem keys detected across runs.")
		return
	}
	fmt.Fprintf(w, "%-30s  %12s  %11s  %s\n", "KEY", "PROBLEMS", "RATE/DAY", "FIRST SEEN")
	fmt.Fprintf(w, "%-30s  %12s  %11s  %s\n", "---", "--------", "--------", "----------")
	for _, e := range entries {
		fmt.Fprintf(w, "%-30s  %12d  %11.3f  %s\n",
			e.Key, e.TotalProblems, e.ChangeRate,
			e.FirstSeen.Format("2006-01-02"),
		)
	}
}
