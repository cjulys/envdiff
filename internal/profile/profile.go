// Package profile manages named environment profiles, allowing users to
// define and reference sets of .env files by a logical name (e.g. "staging",
// "production").
package profile

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
)

// Profile represents a named collection of .env file paths.
type Profile struct {
	Name  string   `json:"name"`
	Files []string `json:"files"`
}

// Registry holds a set of named profiles persisted to disk.
type Registry struct {
	Profiles map[string]Profile `json:"profiles"`
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{Profiles: make(map[string]Profile)}
}

// LoadFile reads a Registry from the given JSON file.
// If the file does not exist, an empty Registry is returned.
func LoadFile(path string) (*Registry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return New(), nil
	}
	if err != nil {
		return nil, err
	}
	var r Registry
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if r.Profiles == nil {
		r.Profiles = make(map[string]Profile)
	}
	return &r, nil
}

// SaveFile writes the Registry to the given JSON file.
func (r *Registry) SaveFile(path string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Add inserts or replaces a profile in the registry.
func (r *Registry) Add(p Profile) {
	r.Profiles[p.Name] = p
}

// Remove deletes a profile by name. Returns false if it did not exist.
func (r *Registry) Remove(name string) bool {
	_, ok := r.Profiles[name]
	delete(r.Profiles, name)
	return ok
}

// Get retrieves a profile by name.
func (r *Registry) Get(name string) (Profile, bool) {
	p, ok := r.Profiles[name]
	return p, ok
}

// List returns all profiles sorted by name.
func (r *Registry) List() []Profile {
	out := make([]Profile, 0, len(r.Profiles))
	for _, p := range r.Profiles {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
