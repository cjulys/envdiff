package diff

import (
	"fmt"
	"io"
)

// ScoreOptions controls the behaviour of RunScore.
type ScoreOptions struct {
	// Label is displayed in the score card header.
	Label string
	// FailBelow causes RunScore to return true (indicating failure)
	// when the computed score is strictly below this threshold.
	// A value of 0 disables the threshold check.
	FailBelow float64
}

// DefaultScoreOptions returns a ScoreOptions with sensible defaults.
func DefaultScoreOptions() ScoreOptions {
	return ScoreOptions{
		Label:     "",
		FailBelow: 0,
	}
}

// RunScore computes and prints the health score for results.
// It returns true if the score falls below opts.FailBelow, signalling
// that the caller should exit with a non-zero status code.
func RunScore(w io.Writer, results []Result, opts ScoreOptions) bool {
	s := ComputeScore(results)
	WriteScore(w, s, opts.Label)

	if opts.FailBelow > 0 && s.Value < opts.FailBelow {
		fmt.Fprintf(w, "FAIL: score %.1f is below threshold %.1f\n", s.Value, opts.FailBelow)
		return true
	}
	return false
}
