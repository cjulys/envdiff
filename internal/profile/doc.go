// Package profile provides a registry of named environment profiles.
//
// A profile associates a logical name (such as "staging" or "production")
// with one or more .env file paths. Profiles are persisted as a JSON file
// and can be loaded by envdiff commands to avoid repeating file paths on
// every invocation.
//
// Example usage:
//
//	reg, err := profile.LoadFile(".envdiff-profiles.json")
//	reg.Add(profile.Profile{Name: "staging", Files: []string{".env", ".env.staging"}})
//	reg.SaveFile(".envdiff-profiles.json")
package profile
