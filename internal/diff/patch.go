package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// PatchOptions controls how a patch is generated.
type PatchOptions struct {
	// OnlyProblems omits matching keys from the patch output.
	OnlyProblems bool
	// TargetLabel is the name printed as the target file header.
	TargetLabel string
}

// WritePatch writes a unified-diff-style patch to w that, when applied to the
// "A" environment, would bring it in line with the "B" environment.
//
// Lines prefixed with '-' are keys present in A that differ from B.
// Lines prefixed with '+' are keys that should be added or updated in A.
// Lines prefixed with ' ' are matching keys (only when OnlyProblems is false).
func WritePatch(w io.Writer, results []Result, opts PatchOptions) error {
	label := opts.TargetLabel
	if label == "" {
		label = "b/.env"
	}

	fmt.Fprintf(w, "--- a/.env\n")
	fmt.Fprintf(w, "+++ %s\n\n", label)

	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		switch {
		case r.Status == StatusMissingInA:
			// Key exists in B but not A — add it.
			fmt.Fprintf(w, "+%s=%s\n", r.Key, r.ValueB)
		case r.Status == StatusMissingInB:
			// Key exists in A but not B — remove it.
			fmt.Fprintf(w, "-%s=%s\n", r.Key, r.ValueA)
		case r.Status == StatusMismatch:
			// Key exists in both but values differ.
			fmt.Fprintf(w, "-%s=%s\n", r.Key, r.ValueA)
			fmt.Fprintf(w, "+%s=%s\n", r.Key, r.ValueB)
		default:
			if !opts.OnlyProblems {
				fmt.Fprintf(w, " %s=%s\n", r.Key, r.ValueA)
			}
		}
	}

	_ = strings.TrimSpace // imported for potential future use
	return nil
}
