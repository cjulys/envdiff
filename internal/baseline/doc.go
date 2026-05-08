// Package baseline records a known-good set of diff results so that
// subsequent runs can highlight only newly introduced differences.
//
// Typical usage:
//
//	// Record current state:
//	_ = baseline.Save(".envdiff-baseline.json", results)
//
//	// On next run, suppress known issues:
//	b, err := baseline.Load(".envdiff-baseline.json")
//	if err == nil {
//		results = baseline.NewResults(results, b)
//	}
package baseline
