// Package audit provides functionality to record and retrieve audit logs
// of envdiff comparison runs, enabling history tracking over time.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Entry represents a single audit log record for one comparison run.
type Entry struct {
	Timestamp  time.Time     `json:"timestamp"`
	Files      []string      `json:"files"`
	TotalKeys  int           `json:"total_keys"`
	Problems   int           `json:"problems"`
	Results    []diff.Result `json:"results"`
}

// Log appends a new audit entry to the given log file path.
// The file is created if it does not exist.
func Log(path string, files []string, results []diff.Result) error {
	problems := 0
	for _, r := range results {
		if r.IsProblem() {
			problems++
		}
	}

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Files:     files,
		TotalKeys: len(results),
		Problems:  problems,
		Results:   results,
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("audit: create directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("audit: encode entry: %w", err)
	}
	return nil
}

// ReadAll reads all audit entries from the given log file.
func ReadAll(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	var entries []Entry
	dec := json.NewDecoder(f)
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
