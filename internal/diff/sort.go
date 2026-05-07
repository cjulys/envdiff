package diff

import "sort"

// SortResults returns a sorted copy of the given results slice.
// Results are sorted by key name alphabetically. Results with the
// same key are ordered by status: Match < Mismatch < MissingInA < MissingInB.
func SortResults(results []Result) []Result {
	copy_ := make([]Result, len(results))
	copy(copy_, results)

	sort.Slice(copy_, func(i, j int) bool {
		if copy_[i].Key != copy_[j].Key {
			return copy_[i].Key < copy_[j].Key
		}
		return statusOrder(copy_[i].Status) < statusOrder(copy_[j].Status)
	})

	return copy_
}

// GroupByStatus partitions results into a map keyed by Status.
func GroupByStatus(results []Result) map[Status][]Result {
	groups := make(map[Status][]Result)
	for _, r := range results {
		groups[r.Status] = append(groups[r.Status], r)
	}
	return groups
}

// ProblemResults returns only results that represent a diff (non-Match).
func ProblemResults(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.IsProblem() {
			out = append(out, r)
		}
	}
	return out
}

func statusOrder(s Status) int {
	switch s {
	case StatusMatch:
		return 0
	case StatusMismatch:
		return 1
	case StatusMissingInA:
		return 2
	case StatusMissingInB:
		return 3
	default:
		return 99
	}
}
