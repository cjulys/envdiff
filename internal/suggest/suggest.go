// Package suggest provides heuristic-based key suggestions for missing
// or mismatched entries found during an env diff comparison.
package suggest

import (
	"sort"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion pairs a result with a candidate key name that closely resembles
// the missing key, helping users spot typos or casing inconsistencies.
type Suggestion struct {
	Result    diff.Result
	Candidate string
	Score     int // higher is a closer match (0-100)
}

// For returns suggestions for all problem results in src, searching for
// similar key names within the pool of known keys.
func For(results []diff.Result, knownKeys []string) []Suggestion {
	var suggestions []Suggestion
	for _, r := range results {
		if !r.IsProblem() {
			continue
		}
		best, score := bestMatch(r.Key, knownKeys)
		if score >= 50 {
			suggestions = append(suggestions, Suggestion{
				Result:    r,
				Candidate: best,
				Score:     score,
			})
		}
	}
	sort.Slice(suggestions, func(i, j int) bool {
		if suggestions[i].Score != suggestions[j].Score {
			return suggestions[i].Score > suggestions[j].Score
		}
		return suggestions[i].Result.Key < suggestions[j].Result.Key
	})
	return suggestions
}

// bestMatch returns the key from candidates that most closely resembles
// target, along with a similarity score in the range [0, 100].
func bestMatch(target string, candidates []string) (string, int) {
	best := ""
	bestScore := 0
	norm := strings.ToUpper(target)
	for _, c := range candidates {
		if strings.EqualFold(c, target) {
			continue // exact (case-insensitive) match — not a suggestion
		}
		s := similarity(norm, strings.ToUpper(c))
		if s > bestScore {
			bestScore = s
			best = c
		}
	}
	return best, bestScore
}

// similarity returns a score in [0,100] based on the longest common
// subsequence length relative to the longer of the two strings.
func similarity(a, b string) int {
	if len(a) == 0 || len(b) == 0 {
		return 0
	}
	lcs := lcsLength(a, b)
	max := len(a)
	if len(b) > max {
		max = len(b)
	}
	return (lcs * 100) / max
}

func lcsLength(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	m, n := len(ra), len(rb)
	dp := make([]int, n+1)
	for i := 1; i <= m; i++ {
		prev := 0
		for j := 1; j <= n; j++ {
			tmp := dp[j]
			if ra[i-1] == rb[j-1] {
				dp[j] = prev + 1
			} else if dp[j] < dp[j-1] {
				dp[j] = dp[j-1]
			}
			prev = tmp
		}
	}
	return dp[n]
}
