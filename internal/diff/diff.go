package diff

// KeyStatus represents the comparison result for a single key.
type KeyStatus int

const (
	StatusMatch    KeyStatus = iota // key exists in both files with the same value
	StatusMismatch                  // key exists in both files but values differ
	StatusMissingB                  // key exists in A but not in B
	StatusMissingA                  // key exists in B but not in A
)

// Result holds the diff result for a single key.
type Result struct {
	Key      string
	Status   KeyStatus
	ValueA   string
	ValueB   string
}

// Compare compares two parsed env maps (key -> value) and returns a slice
// of Result entries describing every difference and match.
func Compare(a, b map[string]string) []Result {
	seen := make(map[string]bool)
	var results []Result

	for k, va := range a {
		seen[k] = true
		vb, ok := b[k]
		switch {
		case !ok:
			results = append(results, Result{Key: k, Status: StatusMissingB, ValueA: va})
		case va == vb:
			results = append(results, Result{Key: k, Status: StatusMatch, ValueA: va, ValueB: vb})
		default:
			results = append(results, Result{Key: k, Status: StatusMismatch, ValueA: va, ValueB: vb})
		}
	}

	for k, vb := range b {
		if !seen[k] {
			results = append(results, Result{Key: k, Status: StatusMissingA, ValueB: vb})
		}
	}

	return results
}

// Summary counts results by status.
type Summary struct {
	Matches   int
	Mismatches int
	MissingInB int
	MissingInA int
}

// Summarize aggregates a slice of Results into a Summary.
func Summarize(results []Result) Summary {
	var s Summary
	for _, r := range results {
		switch r.Status {
		case StatusMatch:
			s.Matches++
		case StatusMismatch:
			s.Mismatches++
		case StatusMissingB:
			s.MissingInB++
		case StatusMissingA:
			s.MissingInA++
		}
	}
	return s
}
