package diff

import (
	"fmt"
	"io"
)

// VelocityOptions controls the behaviour of RunVelocity.
type VelocityOptions struct {
	// MinRate filters out entries whose change rate is below this threshold.
	MinRate float64
	// TopN limits output to the N highest-velocity keys. 0 means no limit.
	TopN int
}

// DefaultVelocityOptions returns sensible defaults for RunVelocity.
func DefaultVelocityOptions() VelocityOptions {
	return VelocityOptions{
		MinRate: 0.0,
		TopN:    0,
	}
}

// RunVelocity builds a velocity report from runs and writes it to w.
// It returns true if any entries exceed the MinRate threshold.
func RunVelocity(w io.Writer, runs []VelocityRun, opts VelocityOptions) bool {
	entries := BuildVelocity(runs)

	// Apply MinRate filter.
	filtered := entries[:0]
	for _, e := range entries {
		if e.ChangeRate >= opts.MinRate {
			filtered = append(filtered, e)
		}
	}

	// Apply TopN limit.
	if opts.TopN > 0 && len(filtered) > opts.TopN {
		filtered = filtered[:opts.TopN]
	}

	if len(filtered) == 0 {
		fmt.Fprintln(w, "No high-velocity keys found.")
		return false
	}

	fmt.Fprintf(w, "Velocity report (%d keys)\n\n", len(filtered))
	WriteVelocity(w, filtered)
	return true
}
