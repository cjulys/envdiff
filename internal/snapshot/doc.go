// Package snapshot provides save and load functionality for envdiff results.
//
// A snapshot captures the output of a diff comparison between two .env files,
// including metadata such as the file paths and the time the snapshot was taken.
//
// Snapshots are stored as JSON files and can be loaded later for auditing,
// CI comparisons, or tracking environment drift over time.
//
// Example usage:
//
//	results := diff.Compare(mapA, mapB)
//	err := snapshot.Save("snap.json", ".env.dev", ".env.prod", results)
//
//	s, err := snapshot.Load("snap.json")
//	problems := s.FilterProblems()
package snapshot
