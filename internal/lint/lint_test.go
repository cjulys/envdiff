package lint_test

import (
	"testing"

	"github.com/user/envdiff/internal/lint"
)

func findFinding(findings []lint.Finding, key string) *lint.Finding {
	for _, f := range findings {
		if f.Key == key {
			return &f
		}
	}
	return nil
}

func TestCheck_EmptyValue_ReturnsWarning(t *testing.T) {
	env := map[string]string{"API_KEY": ""}
	findings := lint.Check(env)
	f := findFinding(findings, "API_KEY")
	if f == nil {
		t.Fatal("expected a finding for empty value, got none")
	}
	if f.Severity != lint.Warning {
		t.Errorf("expected Warning, got %s", f.Severity)
	}
}

func TestCheck_LeadingWhitespace_ReturnsWarning(t *testing.T) {
	env := map[string]string{"DB_HOST": "  localhost"}
	findings := lint.Check(env)
	f := findFinding(findings, "DB_HOST")
	if f == nil {
		t.Fatal("expected a finding for leading whitespace")
	}
	if f.Severity != lint.Warning {
		t.Errorf("expected Warning, got %s", f.Severity)
	}
}

func TestCheck_TrailingWhitespace_ReturnsWarning(t *testing.T) {
	env := map[string]string{"DB_PORT": "5432  "}
	findings := lint.Check(env)
	f := findFinding(findings, "DB_PORT")
	if f == nil {
		t.Fatal("expected a finding for trailing whitespace")
	}
}

func TestCheck_WellKnownKey_ReturnsError(t *testing.T) {
	env := map[string]string{"PATH": "/custom/bin"}
	findings := lint.Check(env)
	f := findFinding(findings, "PATH")
	if f == nil {
		t.Fatal("expected a finding for well-known key PATH")
	}
	if f.Severity != lint.Error {
		t.Errorf("expected Error severity, got %s", f.Severity)
	}
}

func TestCheck_TabInValue_ReturnsWarning(t *testing.T) {
	env := map[string]string{"NOTES": "hello\tworld"}
	findings := lint.Check(env)
	f := findFinding(findings, "NOTES")
	if f == nil {
		t.Fatal("expected a finding for tab in value")
	}
	if f.Severity != lint.Warning {
		t.Errorf("expected Warning, got %s", f.Severity)
	}
}

func TestCheck_CleanEnv_ReturnsNoFindings(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "envdiff",
		"LOG_LEVEL": "info",
		"PORT":      "8080",
	}
	findings := lint.Check(env)
	if len(findings) != 0 {
		t.Errorf("expected no findings, got %d: %v", len(findings), findings)
	}
}

func TestFinding_String_IncludesFields(t *testing.T) {
	f := lint.Finding{Key: "FOO", Message: "some issue", Severity: lint.Warning}
	s := f.String()
	for _, want := range []string{"FOO", "some issue", "warning"} {
		if !contains(s, want) {
			t.Errorf("String() missing %q, got: %s", want, s)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
