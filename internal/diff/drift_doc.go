// Package diff provides utilities for comparing .env files and surfacing
// differences across environments.
//
// # Drift
//
// The drift sub-feature tracks how "far" each named environment has drifted
// from a clean state by recording the number of problematic keys (missing or
// mismatched) relative to the total key count.
//
// Basic usage:
//
//	report := &diff.DriftReport{}
//	report.AddEntry("production", results)
//	report.AddEntry("staging", stagingResults)
//	diff.WriteDrift(os.Stdout, report)
//
// Or use the higher-level command helper:
//
//	opts := diff.DefaultDriftOptions()
//	opts.Labels["production"] = results
//	hasProblems := diff.RunDrift(os.Stdout, opts)
package diff
