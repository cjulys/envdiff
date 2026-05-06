package validate_test

import (
	"testing"

	"github.com/user/envdiff/internal/validate"
)

func TestCheckKeys_ValidKeys(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost",
		"PORT":         "8080",
		"APP_SECRET_KEY": "abc",
	}
	issues := validate.CheckKeys(env)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestCheckKeys_InvalidLowercase(t *testing.T) {
	env := map[string]string{
		"database_url": "postgres://localhost",
	}
	issues := validate.CheckKeys(env)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "database_url" {
		t.Errorf("unexpected key in issue: %s", issues[0].Key)
	}
}

func TestCheckKeys_MixedCase(t *testing.T) {
	env := map[string]string{
		"MyVar": "value",
	}
	issues := validate.CheckKeys(env)
	if len(issues) == 0 {
		t.Error("expected issue for mixed-case key, got none")
	}
}

func TestCheckKeys_KeyWithWhitespace(t *testing.T) {
	env := map[string]string{
		"MY VAR": "value",
	}
	issues := validate.CheckKeys(env)
	if len(issues) == 0 {
		t.Error("expected issue for key with whitespace, got none")
	}
}

func TestCheckKeys_EmptyKey(t *testing.T) {
	env := map[string]string{
		"": "value",
	}
	issues := validate.CheckKeys(env)
	if len(issues) == 0 {
		t.Error("expected issue for empty key, got none")
	}
}

func TestCheckDuplicates_NoDuplicates(t *testing.T) {
	lines := []string{
		"PORT=8080",
		"HOST=localhost",
	}
	issues := validate.CheckDuplicates(lines)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}

func TestCheckDuplicates_WithDuplicate(t *testing.T) {
	lines := []string{
		"PORT=8080",
		"HOST=localhost",
		"PORT=9090",
	}
	issues := validate.CheckDuplicates(lines)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Key != "PORT" {
		t.Errorf("expected issue for PORT, got %s", issues[0].Key)
	}
}

func TestCheckDuplicates_SkipsCommentsAndBlanks(t *testing.T) {
	lines := []string{
		"# comment",
		"",
		"PORT=8080",
	}
	issues := validate.CheckDuplicates(lines)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %v", issues)
	}
}
