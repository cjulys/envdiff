// Package diff provides utilities for comparing .env files across environments.
//
// # Fingerprinting
//
// The fingerprint sub-feature produces a short, deterministic SHA-256 hash
// that summarises the current diff state of one or more environments.
//
// A fingerprint is computed from the sorted (key, status) pairs of a result
// set, making it:
//
//   - Order-independent: shuffling results does not change the hash.
//   - Value-independent: actual env values are excluded so secrets are never
//     embedded in the fingerprint.
//   - Stable across runs: given the same logical diff the hash is identical.
//
// Typical use-cases:
//
//   - CI pipelines that want to detect whether the diff has changed since the
//     last run without storing full snapshots.
//   - Alerting when two previously matching environments start to diverge.
//   - Quick sanity-checks in pre-deploy hooks.
//
// Usage:
//
//	fp := diff.ComputeFingerprint(results)
//	diff.WriteFingerprint(os.Stdout, fp, "staging")
package diff
