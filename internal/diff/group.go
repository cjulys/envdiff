package diff

// GroupedResults holds diff results organized by environment file label.
type GroupedResults struct {
	Label   string
	Results []Result
}

// GroupByFile partitions a flat list of Results into groups, one per unique
// source file pair label. The label is derived from the Key prefix convention
// used internally, but here we accept an explicit labels slice that matches
// the order of file pairs passed on the CLI.
//
// Each group retains the original ordering of results within it. Results that
// do not belong to any label fall into a catch-all group named "default".
func GroupByFile(results []Result, labels []string) []GroupedResults {
	if len(labels) == 0 {
		return []GroupedResults{{Label: "default", Results: results}}
	}

	// Build an index so we can assign results to groups in O(1).
	index := make(map[string]int, len(labels))
	for i, l := range labels {
		index[l] = i
	}

	groups := make([]GroupedResults, len(labels))
	for i, l := range labels {
		groups[i] = GroupedResults{Label: l}
	}

	var fallback GroupedResults
	fallback.Label = "default"

	for _, r := range results {
		matched := false
		for _, l := range labels {
			if r.Source == l {
				groups[index[l]].Results = append(groups[index[l]].Results, r)
				matched = true
				break
			}
		}
		if !matched {
			fallback.Results = append(fallback.Results, r)
		}
	}

	// Append the fallback group only when it has entries.
	if len(fallback.Results) > 0 {
		groups = append(groups, fallback)
	}

	return groups
}

// ProblemGroups returns only the groups that contain at least one problem result.
func ProblemGroups(groups []GroupedResults) []GroupedResults {
	out := make([]GroupedResults, 0, len(groups))
	for _, g := range groups {
		var problems []Result
		for _, r := range g.Results {
			if r.IsProblem() {
				problems = append(problems, r)
			}
		}
		if len(problems) > 0 {
			out = append(out, GroupedResults{Label: g.Label, Results: problems})
		}
	}
	return out
}
