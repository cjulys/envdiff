// Package env provides utilities for resolving and expanding environment
// variable references within .env file values.
package env

import (
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR} and $VAR style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Resolver expands variable references in values using a combined lookup
// of the provided map and the process environment.
type Resolver struct {
	env     map[string]string
	useSys  bool
}

// NewResolver creates a Resolver backed by the given map. If useSys is true,
// variables not found in the map are looked up in the process environment.
func NewResolver(env map[string]string, useSys bool) *Resolver {
	return &Resolver{env: env, useSys: useSys}
}

// Expand replaces all variable references in value with their resolved values.
// Unresolvable references are replaced with an empty string.
func (r *Resolver) Expand(value string) string {
	return varPattern.ReplaceAllStringFunc(value, func(match string) string {
		name := extractName(match)
		if v, ok := r.env[name]; ok {
			return v
		}
		if r.useSys {
			return os.Getenv(name)
		}
		return ""
	})
}

// ExpandAll returns a new map with all values expanded.
func (r *Resolver) ExpandAll() map[string]string {
	out := make(map[string]string, len(r.env))
	for k, v := range r.env {
		out[k] = r.Expand(v)
	}
	return out
}

// Lookup resolves a single variable name using the resolver's lookup order:
// first the env map, then the process environment (if useSys is true).
// The second return value reports whether the variable was found.
func (r *Resolver) Lookup(name string) (string, bool) {
	if v, ok := r.env[name]; ok {
		return v, true
	}
	if r.useSys {
		if v, ok := os.LookupEnv(name); ok {
			return v, true
		}
	}
	return "", false
}

// extractName pulls the variable name out of a $VAR or ${VAR} token.
func extractName(token string) string {
	token = strings.TrimPrefix(token, "$")
	token = strings.TrimPrefix(token, "{")
	token = strings.TrimSuffix(token, "}")
	return token
}
