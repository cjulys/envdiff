package diff

import (
	"fmt"
	"io"
)

// FingerprintOptions controls the behaviour of RunFingerprint.
type FingerprintOptions struct {
	// LabelA is the display name for the first environment.
	LabelA string
	// LabelB is the display name for the second environment.
	LabelB string
	// Compare, when true, prints both fingerprints and notes whether they match.
	Compare bool
}

// DefaultFingerprintOptions returns sensible defaults.
func DefaultFingerprintOptions() FingerprintOptions {
	return FingerprintOptions{
		LabelA:  "env-a",
		LabelB:  "env-b",
		Compare: false,
	}
}

// RunFingerprint computes and writes a fingerprint report for results.
// When opts.Compare is true a second set of results (resultsB) is also
// fingerprinted and the two hashes are compared.
// Returns true when issues were detected.
func RunFingerprint(w io.Writer, resultsA, resultsB []Result, opts FingerprintOptions) bool {
	fpA := ComputeFingerprint(resultsA)
	WriteFingerprint(w, fpA, opts.LabelA)

	if !opts.Compare || resultsB == nil {
		return fpA.Issues > 0
	}

	fmt.Fprintln(w)
	fpB := ComputeFingerprint(resultsB)
	WriteFingerprint(w, fpB, opts.LabelB)

	fmt.Fprintln(w)
	if fpA.Equal(fpB) {
		fmt.Fprintln(w, "Result: fingerprints match — environments are identical")
	} else {
		fmt.Fprintln(w, "Result: fingerprints differ — environments have diverged")
	}

	return fpA.Issues > 0 || fpB.Issues > 0
}
