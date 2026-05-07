package diff

import (
	"testing"
)

func makeResult(key string, status Status) Result {
	return Result{Key: key, ValueA: "a", ValueB: "b", Status: status}
}

func TestSortResults_ByKey(t *testing.T) {
	input := []Result{
		makeResult("ZEBRA", StatusMatch),
		makeResult("APPLE", StatusMismatch),
		makeResult("MANGO", StatusMissingInB),
	}

	sorted := SortResults(input)

	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, r := range sorted {
		if r.Key != expected[i] {
			t.Errorf("index %d: expected key %q, got %q", i, expected[i], r.Key)
		}
	}
}

func TestSortResults_SameKeyOrderedByStatus(t *testing.T) {
	input := []Result{
		{Key: "FOO", Status: StatusMissingInB},
		{Key: "FOO", Status: StatusMatch},
		{Key: "FOO", Status: StatusMismatch},
	}

	sorted := SortResults(input)

	if sorted[0].Status != StatusMatch {
		t.Errorf("expected first to be Match, got %v", sorted[0].Status)
	}
	if sorted[1].Status != StatusMismatch {
		t.Errorf("expected second to be Mismatch, got %v", sorted[1].Status)
	}
	if sorted[2].Status != StatusMissingInB {
		t.Errorf("expected third to be MissingInB, got %v", sorted[2].Status)
	}
}

func TestSortResults_DoesNotMutateOriginal(t *testing.T) {
	input := []Result{
		makeResult("Z", StatusMatch),
		makeResult("A", StatusMatch),
	}

	_ = SortResults(input)

	if input[0].Key != "Z" {
		t.Errorf("original slice was mutated")
	}
}

func TestGroupByStatus_Partitions(t *testing.T) {
	input := []Result{
		makeResult("A", StatusMatch),
		makeResult("B", StatusMismatch),
		makeResult("C", StatusMissingInA),
		makeResult("D", StatusMatch),
	}

	groups := GroupByStatus(input)

	if len(groups[StatusMatch]) != 2 {
		t.Errorf("expected 2 Match results, got %d", len(groups[StatusMatch]))
	}
	if len(groups[StatusMismatch]) != 1 {
		t.Errorf("expected 1 Mismatch result, got %d", len(groups[StatusMismatch]))
	}
	if len(groups[StatusMissingInA]) != 1 {
		t.Errorf("expected 1 MissingInA result, got %d", len(groups[StatusMissingInA]))
	}
}

func TestProblemResults_FiltersMatches(t *testing.T) {
	input := []Result{
		makeResult("A", StatusMatch),
		makeResult("B", StatusMismatch),
		makeResult("C", StatusMissingInB),
	}

	problems := ProblemResults(input)

	if len(problems) != 2 {
		t.Errorf("expected 2 problems, got %d", len(problems))
	}
	for _, r := range problems {
		if !r.IsProblem() {
			t.Errorf("expected problem result, got Match for key %q", r.Key)
		}
	}
}
