// Package audit records comparison run history for envdiff.
//
// Each time envdiff compares two or more .env files, an audit entry can be
// appended to a newline-delimited JSON log file. This allows users to track
// how their environment configuration drift changes over time.
//
// Usage:
//
//	err := audit.Log(".envdiff-audit.log", []string{".env", ".env.prod"}, results)
//
//	entries, err := audit.ReadAll(".envdiff-audit.log")
//
// Each Entry records the timestamp, compared file paths, total key count,
// number of problems detected, and the full set of diff results.
package audit
