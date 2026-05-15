package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Signature represents a compact structural fingerprint of an environment's
// key shape — which keys are present, missing, or mismatched — without
// exposing values. Two environments with the same Signature are structurally
// equivalent.
type Signature struct {
	Label   string
	Keys    []string // sorted list of keys contributing to the signature
	Hash    string   // short hex digest of the key shape
	Problems int
}

// BuildSignatures computes a Signature for each labelled set of results.
// Each entry in runs maps an environment label to its diff results.
func BuildSignatures(runs map[string][]Result) []Signature {
	if len(runs) == 0 {
		return nil
	}

	labels := make([]string, 0, len(runs))
	for l := range runs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	out := make([]Signature, 0, len(labels))
	for _, label := range labels {
		results := runs[label]
		problems := 0
		keySet := map[string]struct{}{}
		for _, r := range results {
			keySet[r.Key] = struct{}{}
			if r.IsProblem() {
				problems++
			}
		}
		keys := make([]string, 0, len(keySet))
		for k := range keySet {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		hash := signatureHash(keys, results)
		out = append(out, Signature{
			Label:    label,
			Keys:     keys,
			Hash:     hash,
			Problems: problems,
		})
	}
	return out
}

// signatureHash builds a short deterministic hash from the key+status pairs.
func signatureHash(keys []string, results []Result) string {
	statusMap := make(map[string]string, len(results))
	for _, r := range results {
		statusMap[r.Key] = string(r.Status)
	}
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(statusMap[k])
		sb.WriteByte(';')
	}
	// Simple djb2-style hash for a short hex string.
	h := uint32(5381)
	for _, c := range sb.String() {
		h = h*33 ^ uint32(c)
	}
	return fmt.Sprintf("%08x", h)
}

// WriteSignatures writes a human-readable signature table to w.
func WriteSignatures(w io.Writer, sigs []Signature) {
	if len(sigs) == 0 {
		fmt.Fprintln(w, "no signatures computed")
		return
	}
	fmt.Fprintf(w, "%-24s  %-10s  %8s  %s\n", "ENVIRONMENT", "HASH", "PROBLEMS", "KEYS")
	fmt.Fprintln(w, strings.Repeat("-", 60))
	for _, s := range sigs {
		fmt.Fprintf(w, "%-24s  %-10s  %8d  %d\n", s.Label, s.Hash, s.Problems, len(s.Keys))
	}
}
