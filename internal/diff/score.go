package diff

// Score represents a numeric health score derived from diff results.
type Score struct {
	// Total is the total number of keys evaluated.
	Total int
	// Problems is the number of keys with a non-match status.
	Problems int
	// Value is the computed score in the range [0, 100].
	Value float64
	// Grade is a letter grade derived from Value.
	Grade string
}

// ComputeScore calculates a health score for a set of diff results.
// A score of 100 means all keys match; 0 means every key is a problem.
func ComputeScore(results []Result) Score {
	total := len(results)
	if total == 0 {
		return Score{Total: 0, Problems: 0, Value: 100.0, Grade: "A"}
	}

	problems := 0
	for _, r := range results {
		if r.Status != StatusMatch {
			problems++
		}
	}

	value := 100.0 * float64(total-problems) / float64(total)
	return Score{
		Total:    total,
		Problems: problems,
		Value:    value,
		Grade:    letterGrade(value),
	}
}

func letterGrade(v float64) string {
	switch {
	case v >= 95:
		return "A"
	case v >= 85:
		return "B"
	case v >= 70:
		return "C"
	case v >= 50:
		return "D"
	default:
		return "F"
	}
}
