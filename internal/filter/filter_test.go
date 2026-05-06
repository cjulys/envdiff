package filter_test

import (
	"testing"

	"github.com/envdiff/internal/diff"
	"github.com/envdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "DB_HOST", Status: diff.Match, ValueA: "localhost", ValueB: "localhost"},
	{Key: "DB_PASS", Status: diff.Mismatch, ValueA: "secret", ValueB: "other"},
	{Key: "API_KEY", Status: diff.Missing, ValueA: "abc123", ValueB: ""},
	{Key: "APP_ENV", Status: diff.Match, ValueA: "prod", ValueB: "prod"},
	{Key: "DB_PORT", Status: diff.Missing, ValueA: "", ValueB: "5432"},
}

func TestApply_NoFilter(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{})
	if len(results) != len(sampleResults) {
		t.Errorf("expected %d results, got %d", len(sampleResults), len(results))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{OnlyMissing: true})
	for _, r := range results {
		if r.Status != diff.Missing {
			t.Errorf("expected only Missing, got %s for key %s", r.Status, r.Key)
		}
	}
	if len(results) != 2 {
		t.Errorf("expected 2 missing results, got %d", len(results))
	}
}

func TestApply_OnlyMismatch(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{OnlyMismatch: true})
	if len(results) != 1 || results[0].Key != "DB_PASS" {
		t.Errorf("expected 1 mismatch result for DB_PASS, got %+v", results)
	}
}

func TestApply_Prefix(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{Prefix: "db_"})
	if len(results) != 3 {
		t.Errorf("expected 3 DB_ results, got %d", len(results))
	}
}

func TestApply_SpecificKeys(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{Keys: []string{"api_key", "app_env"}})
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func TestApply_PrefixAndMissing(t *testing.T) {
	results := filter.Apply(sampleResults, filter.Options{Prefix: "db_", OnlyMissing: true})
	if len(results) != 1 || results[0].Key != "DB_PORT" {
		t.Errorf("expected 1 result DB_PORT, got %+v", results)
	}
}
