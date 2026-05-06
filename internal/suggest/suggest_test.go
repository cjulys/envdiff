package suggest_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/suggest"
)

func makeResult(key, status string) diff.Result {
	return diff.Result{Key: key, Status: status}
}

func TestFor_NoProblems_ReturnsEmpty(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOST", "match"),
	}
	known := []string{"DB_HOST", "DB_PORT"}
	got := suggest.For(results, known)
	if len(got) != 0 {
		t.Errorf("expected no suggestions, got %d", len(got))
	}
}

func TestFor_TypoInKey_ReturnsSuggestion(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOOST", "missing_in_b"),
	}
	known := []string{"DB_HOST", "DB_PORT", "API_KEY"}
	got := suggest.For(results, known)
	if len(got) == 0 {
		t.Fatal("expected at least one suggestion")
	}
	if got[0].Candidate != "DB_HOST" {
		t.Errorf("expected candidate DB_HOST, got %s", got[0].Candidate)
	}
	if got[0].Score < 50 {
		t.Errorf("expected score >= 50, got %d", got[0].Score)
	}
}

func TestFor_CasingDifference_ReturnsSuggestion(t *testing.T) {
	results := []diff.Result{
		makeResult("api_key", "missing_in_a"),
	}
	known := []string{"API_KEY", "DB_HOST"}
	got := suggest.For(results, known)
	if len(got) == 0 {
		t.Fatal("expected suggestion for casing difference")
	}
	if got[0].Candidate != "API_KEY" {
		t.Errorf("expected API_KEY, got %s", got[0].Candidate)
	}
}

func TestFor_NoCloseMatch_ReturnsEmpty(t *testing.T) {
	results := []diff.Result{
		makeResult("ZZZZZ", "mismatch"),
	}
	known := []string{"DB_HOST", "API_KEY", "REDIS_URL"}
	got := suggest.For(results, known)
	if len(got) != 0 {
		t.Errorf("expected no suggestions for unrelated key, got %d", len(got))
	}
}

func TestFor_SortedByScoreDescending(t *testing.T) {
	results := []diff.Result{
		makeResult("DB_HOOST", "missing_in_b"),
		makeResult("REDIS_ULR", "missing_in_b"),
	}
	known := []string{"DB_HOST", "REDIS_URL"}
	got := suggest.For(results, known)
	if len(got) < 2 {
		t.Fatalf("expected 2 suggestions, got %d", len(got))
	}
	if got[0].Score < got[1].Score {
		t.Errorf("expected descending score order, got %d then %d", got[0].Score, got[1].Score)
	}
}
