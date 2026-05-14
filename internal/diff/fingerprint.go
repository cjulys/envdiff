package diff

import (
	"crypto/sha256"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Fingerprint represents a stable hash of a set of diff results.
type Fingerprint struct {
	Hash   string
	Keys   int
	Issues int
}

// ComputeFingerprint produces a deterministic SHA-256 fingerprint from
// a slice of Results. Only the key and status are included so that the
// fingerprint is stable even when verbose values change.
func ComputeFingerprint(results []Result) Fingerprint {
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Key != sorted[j].Key {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].Status < sorted[j].Status
	})

	h := sha256.New()
	for _, r := range sorted {
		fmt.Fprintf(h, "%s:%s\n", r.Key, r.Status)
	}

	issue := 0
	for _, r := range results {
		if r.IsProblem() {
			issue++
		}
	}

	return Fingerprint{
		Hash:   fmt.Sprintf("%x", h.Sum(nil)),
		Keys:   len(results),
		Issues: issue,
	}
}

// Equal reports whether two fingerprints represent the same diff state.
func (f Fingerprint) Equal(other Fingerprint) bool {
	return f.Hash == other.Hash
}

// WriteFingerprint writes a human-readable fingerprint report to w.
func WriteFingerprint(w io.Writer, fp Fingerprint, label string) {
	if label == "" {
		label = "run"
	}
	short := fp.Hash
	if len(short) > 12 {
		short = short[:12]
	}
	fmt.Fprintf(w, "Fingerprint [%s]\n", label)
	fmt.Fprintf(w, "  hash   : %s\n", short)
	fmt.Fprintf(w, "  keys   : %d\n", fp.Keys)
	fmt.Fprintf(w, "  issues : %d\n", fp.Issues)
	if fp.Issues == 0 {
		fmt.Fprintln(w, "  status : clean")
	} else {
		fmt.Fprintln(w, "  status : "+strings.Repeat("!", min(fp.Issues, 5)))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
