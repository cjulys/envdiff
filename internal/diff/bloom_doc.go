// Package diff provides utilities for comparing .env files and surfacing
// differences between environments.
//
// # Bloom Frequency Analysis
//
// The bloom sub-feature tracks how often individual keys appear as problems
// (missing or mismatched) across multiple comparison runs. This is useful for
// identifying keys that are chronically inconsistent across deployments.
//
// Basic usage:
//
//	report := diff.NewBloomReport()
//	for _, run := range historicalRuns {
//	    report.AddRun(run)
//	}
//	diff.WriteBloom(os.Stdout, report, 0.5)
//
// RunBloom provides a higher-level entry point that accepts options and returns
// a boolean indicating whether any key exceeded the configured threshold:
//
//	hasProblems := diff.RunBloom(os.Stdout, runs, diff.DefaultBloomOptions())
package diff
