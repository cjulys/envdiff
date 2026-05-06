// Package snapshot provides functionality to save and load env diff results
// as JSON snapshots for later comparison or auditing.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot represents a saved diff result with metadata.
type Snapshot struct {
	CreatedAt time.Time     `json:"created_at"`
	FileA     string        `json:"file_a"`
	FileB     string        `json:"file_b"`
	Results   []diff.Result `json:"results"`
}

// Save writes a snapshot of the given diff results to the specified file path.
func Save(path, fileA, fileB string, results []diff.Result) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		FileA:     fileA,
		FileB:     fileB,
		Results:   results,
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from the specified file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &s, nil
}

// FilterProblems returns only results that represent a difference (missing or mismatch).
func (s *Snapshot) FilterProblems() []diff.Result {
	var out []diff.Result
	for _, r := range s.Results {
		if r.IsProblem() {
			out = append(out, r)
		}
	}
	return out
}
