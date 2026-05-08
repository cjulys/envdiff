package diff

import (
	"testing"
)

func buildResult(key, source string, status Status) Result {
	return Result{Key: key, Source: source, Status: status}
}

func TestGroupByFile_NoLabels_ReturnsSingleDefaultGroup(t *testing.T) {
	results := []Result{
		buildResult("A", "prod", StatusMatch),
		buildResult("B", "staging", StatusMissingInB),
	}
	groups := GroupByFile(results, nil)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if groups[0].Label != "default" {
		t.Errorf("expected label 'default', got %q", groups[0].Label)
	}
	if len(groups[0].Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(groups[0].Results))
	}
}

func TestGroupByFile_MatchesLabelsBySource(t *testing.T) {
	results := []Result{
		buildResult("KEY1", "prod", StatusMatch),
		buildResult("KEY2", "staging", StatusMismatch),
		buildResult("KEY3", "prod", StatusMissingInA),
	}
	groups := GroupByFile(results, []string{"prod", "staging"})
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Label != "prod" || len(groups[0].Results) != 2 {
		t.Errorf("prod group: got label=%q count=%d", groups[0].Label, len(groups[0].Results))
	}
	if groups[1].Label != "staging" || len(groups[1].Results) != 1 {
		t.Errorf("staging group: got label=%q count=%d", groups[1].Label, len(groups[1].Results))
	}
}

func TestGroupByFile_UnknownSourceGoesToDefault(t *testing.T) {
	results := []Result{
		buildResult("KEY1", "prod", StatusMatch),
		buildResult("KEY2", "unknown", StatusMismatch),
	}
	groups := GroupByFile(results, []string{"prod"})
	// Should have 'prod' group + 'default' fallback
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	last := groups[len(groups)-1]
	if last.Label != "default" {
		t.Errorf("expected fallback label 'default', got %q", last.Label)
	}
	if len(last.Results) != 1 || last.Results[0].Key != "KEY2" {
		t.Errorf("unexpected fallback results: %+v", last.Results)
	}
}

func TestProblemGroups_FiltersMatchResults(t *testing.T) {
	groups := []GroupedResults{
		{
			Label: "prod",
			Results: []Result{
				buildResult("A", "prod", StatusMatch),
				buildResult("B", "prod", StatusMismatch),
			},
		},
		{
			Label: "staging",
			Results: []Result{
				buildResult("C", "staging", StatusMatch),
			},
		},
	}
	problems := ProblemGroups(groups)
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem group, got %d", len(problems))
	}
	if problems[0].Label != "prod" {
		t.Errorf("expected label 'prod', got %q", problems[0].Label)
	}
	if len(problems[0].Results) != 1 || problems[0].Results[0].Key != "B" {
		t.Errorf("unexpected problem results: %+v", problems[0].Results)
	}
}

func TestProblemGroups_EmptyWhenAllMatch(t *testing.T) {
	groups := []GroupedResults{
		{Label: "prod", Results: []Result{buildResult("X", "prod", StatusMatch)}},
	}
	if got := ProblemGroups(groups); len(got) != 0 {
		t.Errorf("expected no problem groups, got %d", len(got))
	}
}
