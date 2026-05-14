package diff

import (
	"fmt"
	"io"
)

// DriftOptions controls the behaviour of RunDrift.
type DriftOptions struct {
	// Labels maps each file path to a human-readable environment name.
	Labels map[string][]Result
	// Verbose includes value details in output (reserved for future use).
	Verbose bool
}

// DefaultDriftOptions returns a DriftOptions with sensible defaults.
func DefaultDriftOptions() DriftOptions {
	return DriftOptions{
		Labels:  make(map[string][]Result),
		Verbose: false,
	}
}

// RunDrift builds and writes a DriftReport from the provided options.
// It returns true when any environment has at least one problem.
func RunDrift(w io.Writer, opts DriftOptions) bool {
	if len(opts.Labels) == 0 {
		fmt.Fprintln(w, "No environments provided for drift analysis.")
		return false
	}

	report := &DriftReport{}
	for label, results := range opts.Labels {
		report.AddEntry(label, results)
	}

	WriteDrift(w, report)

	for _, e := range report.Entries {
		if e.Problems > 0 {
			return true
		}
	}
	return false
}
