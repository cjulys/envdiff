package diff

import (
	"fmt"
	"io"
)

// TimelineOptions controls timeline rendering behaviour.
type TimelineOptions struct {
	// ShowEmpty includes entries that have no problems.
	ShowEmpty bool
	// MaxEntries limits the number of entries shown (0 = unlimited).
	MaxEntries int
}

// DefaultTimelineOptions returns sensible defaults.
func DefaultTimelineOptions() TimelineOptions {
	return TimelineOptions{
		ShowEmpty:  true,
		MaxEntries: 0,
	}
}

// RunTimeline renders the timeline to w, applying options, and returns true
// if any entry in the timeline contains at least one problem.
func RunTimeline(w io.Writer, tl *Timeline, opts TimelineOptions) bool {
	entries := tl.Sorted()

	if !opts.ShowEmpty {
		filtered := entries[:0]
		for _, e := range entries {
			if e.Stats.HasProblems {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	if opts.MaxEntries > 0 && len(entries) > opts.MaxEntries {
		entries = entries[len(entries)-opts.MaxEntries:]
	}

	trimmed := &Timeline{Entries: entries}
	WriteTimeline(w, trimmed)

	hasProblems := false
	for _, e := range entries {
		if e.Stats.HasProblems {
			hasProblems = true
			break
		}
	}

	if hasProblems {
		fmt.Fprintln(w, "\n⚠  One or more timeline entries contain differences.")
	}
	return hasProblems
}
