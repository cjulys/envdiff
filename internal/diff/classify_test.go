package diff

import "testing"

func TestClassify_Match(t *testing.T) {
	r := Result{Key: "APP_NAME", Status: StatusMatch}
	if got := Classify(r); got != SeverityNone {
		t.Errorf("expected none, got %s", got)
	}
}

func TestClassify_Mismatch_NonSensitive(t *testing.T) {
	r := Result{Key: "APP_NAME", Status: StatusMismatch}
	if got := Classify(r); got != SeverityLow {
		t.Errorf("expected low, got %s", got)
	}
}

func TestClassify_Mismatch_SensitiveKey(t *testing.T) {
	r := Result{Key: "DB_PASSWORD", Status: StatusMismatch}
	if got := Classify(r); got != SeverityHigh {
		t.Errorf("expected high, got %s", got)
	}
}

func TestClassify_MissingInB_NonSensitive(t *testing.T) {
	r := Result{Key: "LOG_LEVEL", Status: StatusMissingInB}
	if got := Classify(r); got != SeverityMedium {
		t.Errorf("expected medium, got %s", got)
	}
}

func TestClassify_MissingInA_SensitiveKey(t *testing.T) {
	r := Result{Key: "API_TOKEN", Status: StatusMissingInA}
	if got := Classify(r); got != SeverityHigh {
		t.Errorf("expected high, got %s", got)
	}
}

func TestClassify_SeverityString(t *testing.T) {
	cases := map[Severity]string{
		SeverityNone:   "none",
		SeverityLow:    "low",
		SeverityMedium: "medium",
		SeverityHigh:   "high",
	}
	for s, want := range cases {
		if got := s.String(); got != want {
			t.Errorf("Severity(%d).String() = %q, want %q", s, got, want)
		}
	}
}

func TestClassifyAll_ReturnsMapForAllKeys(t *testing.T) {
	results := []Result{
		{Key: "APP_NAME", Status: StatusMatch},
		{Key: "DB_SECRET", Status: StatusMismatch},
		{Key: "PORT", Status: StatusMissingInB},
	}
	m := ClassifyAll(results)
	if len(m) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(m))
	}
	if m["APP_NAME"] != SeverityNone {
		t.Errorf("APP_NAME: expected none")
	}
	if m["DB_SECRET"] != SeverityHigh {
		t.Errorf("DB_SECRET: expected high")
	}
	if m["PORT"] != SeverityMedium {
		t.Errorf("PORT: expected medium")
	}
}

func TestIsSensitiveKey_Variants(t *testing.T) {
	sensitive := []string{"DB_PASSWORD", "api_secret", "AWS_ACCESS_KEY", "AUTH_TOKEN", "CREDENTIALS"}
	for _, k := range sensitive {
		if !isSensitiveKey(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
	safe := []string{"APP_NAME", "LOG_LEVEL", "PORT", "HOST"}
	for _, k := range safe {
		if isSensitiveKey(k) {
			t.Errorf("expected %q to be safe", k)
		}
	}
}
