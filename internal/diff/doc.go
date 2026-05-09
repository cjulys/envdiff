// Package diff provides utilities for comparing .env file key-value maps,
// producing structured Results that describe whether keys match, mismatch,
// or are absent in one of the two environments being compared.
//
// # Core types
//
// Result represents the outcome for a single key. Status values are:
//   - StatusMatch      — key present and equal in both files
//   - StatusMismatch   — key present in both but with different values
//   - StatusMissingInA — key found in B but not in A
//   - StatusMissingInB — key found in A but not in B
//
// # Utilities
//
// Compare builds a []Result from two string maps.
// Summarize returns a short human-readable summary string.
// SortResults and GroupByStatus provide ordering and partitioning helpers.
// ComputeStats aggregates counts across a result set.
// WritePretty and WriteMarkdown format results for terminal or markdown output.
// WriteReport writes a concise plain-text report.
// WritePatch emits a unified-diff-style patch view.
// WriteSummaryTable renders a tabular summary.
// Annotate produces plain-English notes for each Result.
// WriteAnnotations writes annotations to any io.Writer.
package diff
