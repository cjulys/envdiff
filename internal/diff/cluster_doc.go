// Package diff provides utilities for comparing .env files and surfacing
// differences between environments.
//
// # Cluster
//
// The cluster feature groups environment variable keys by their common
// prefix segment — the part before the first underscore — and reports
// how many keys in each group have problems (mismatches or missing values).
//
// This is useful for quickly identifying which subsystem (e.g. "DB",
// "AWS", "REDIS") is most affected by configuration drift.
//
// Usage:
//
//	report := diff.BuildCluster(results)
//	diff.WriteCluster(os.Stdout, report)
//
// Or via the command helper:
//
//	hasProblems := diff.RunCluster(os.Stdout, results, diff.DefaultClusterOptions())
//
// ClusterOptions allows filtering by minimum cluster size or restricting
// output to only clusters that contain at least one problem.
package diff
