package redact_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/redact"
)

func TestIsSensitive_DefaultPatterns(t *testing.T) {
	r := redact.New(nil)

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"AUTH_TOKEN",
		"AWS_SECRET",
		"PRIVATE_KEY",
		"APP_CREDENTIALS",
	}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := redact.New(nil)

	safe := []string{
		"APP_ENV",
		"PORT",
		"LOG_LEVEL",
		"DATABASE_HOST",
	}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected key %q to NOT be sensitive", key)
		}
	}
}

func TestIsSensitive_CaseInsensitive(t *testing.T) {
	r := redact.New(nil)
	if !r.IsSensitive("db_password") {
		t.Error("expected lowercase 'db_password' to be sensitive")
	}
	if !r.IsSensitive("Api_Key") {
		t.Error("expected mixed-case 'Api_Key' to be sensitive")
	}
}

func TestMaskValue_SensitiveKey(t *testing.T) {
	r := redact.New(nil)
	got := r.MaskValue("DB_PASSWORD", "supersecret")
	if got != redact.MaskedValue() {
		t.Errorf("expected %q, got %q", redact.MaskedValue(), got)
	}
}

func TestMaskValue_SafeKey(t *testing.T) {
	r := redact.New(nil)
	got := r.MaskValue("APP_ENV", "production")
	if got != "production" {
		t.Errorf("expected 'production', got %q", got)
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r := redact.New([]string{"INTERNAL", "CORP"})
	if !r.IsSensitive("INTERNAL_URL") {
		t.Error("expected INTERNAL_URL to be sensitive with custom pattern")
	}
	if !r.IsSensitive("CORP_SECRET") {
		t.Error("expected CORP_SECRET to be sensitive with custom pattern")
	}
	// default patterns should NOT apply
	if r.IsSensitive("API_KEY") {
		t.Error("expected API_KEY to NOT be sensitive with custom patterns only")
	}
}
