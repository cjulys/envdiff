// Package baseline provides functionality to record and compare a known-good
// set of diff results, allowing users to suppress already-acknowledged
// differences and focus only on newly introduced changes.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// ErrNotFound is returned when no baseline file exists at the given path.
var ErrNotFound = errors.New("baseline: file not found")

// Baseline holds a recorded snapshot of diff results.
type Baseline struct {
	CreatedAt time.Time     `json:"created_at"`
	Results   []diff.Result `json:"results"`
}

// Save writes the given results to path as a JSON baseline file.
func Save(path string, results []diff.Result) error {
	b := Baseline{
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(b)
}

// Load reads a baseline file from path and returns it.
// Returns ErrNotFound if the file does not exist.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	defer f.Close()
	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// NewResults returns only the results from current that are not present
// in the baseline. A result is considered known if its Key and Status match.
func NewResults(current []diff.Result, b *Baseline) []diff.Result {
	known := make(map[string]bool, len(b.Results))
	for _, r := range b.Results {
		known[resultKey(r)] = true
	}
	var out []diff.Result
	for _, r := range current {
		if !known[resultKey(r)] {
			out = append(out, r)
		}
	}
	return out
}

func resultKey(r diff.Result) string {
	return r.Key + "|" + string(r.Status)
}
