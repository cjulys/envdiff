// Package diff provides comparison utilities for .env files.
// This file exposes a high-level RunContext helper for use by cmd/envdiff.
package diff

import (
	"fmt"
	"io"
)

// ContextOptions controls how context output is rendered.
type ContextOptions struct {
	Lines   int  // number of neighbouring keys to show around each problem
	Verbose bool // whether to include values in output
}

// DefaultContextOptions returns sensible defaults for context display.
func DefaultContextOptions() ContextOptions {
	return ContextOptions{
		Lines:   2,
		Verbose: false,
	}
}

// RunContext compares envA and envB, builds context blocks for all problems,
// and writes the output to w. It returns true if any problems were found.
func RunContext(w io.Writer, envA, envB map[string]string, opts ContextOptions) (bool, error) {
	results := Compare(envA, envB)
	blocks := BuildContext(results, envA, envB, opts.Lines)

	if len(blocks) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return false, nil
	}

	fmt.Fprintf(w, "Found %d problem(s) with context (±%d keys):\n\n", len(blocks), opts.Lines)
	WriteContext(w, blocks, opts.Verbose)
	return true, nil
}
