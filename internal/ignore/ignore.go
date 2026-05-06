// Package ignore provides functionality to load and apply .envdiffignore files,
// allowing users to suppress specific keys from comparison results.
package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Set holds a collection of keys that should be ignored during diffing.
type Set struct {
	keys map[string]struct{}
}

// New returns an empty ignore Set.
func New() *Set {
	return &Set{keys: make(map[string]struct{})}
}

// LoadFile reads an ignore file from the given path and returns a populated Set.
// Each non-blank, non-comment line is treated as a key to ignore.
// Returns an empty Set and no error if the file does not exist.
func LoadFile(path string) (*Set, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, err
	}
	defer f.Close()

	s := New()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		s.keys[line] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

// Add inserts a key into the ignore Set.
func (s *Set) Add(key string) {
	s.keys[key] = struct{}{}
}

// Contains reports whether the given key is in the ignore Set.
func (s *Set) Contains(key string) bool {
	_, ok := s.keys[key]
	return ok
}

// Len returns the number of keys in the Set.
func (s *Set) Len() int {
	return len(s.keys)
}
