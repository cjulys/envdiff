package diff

import (
	"strings"
	"testing"
)

func TestWritePretty_ContainsSummary(t *testing.T) {
	results := []Result{
		{Key: "APP_NAME", Status: StatusMatch, ValueA: "myapp", ValueB: "myapp"},
		{Key: "DB_PASS", Status: StatusMismatch, ValueA: "secret", ValueB: "other"},
		{Key: "API_KEY", Status: StatusMissingInB, ValueA: "abc123", ValueB: ""},
		{Key: "NEW_FLAG", Status: StatusMissingInA, ValueA: "", ValueB: "true"},
	}

	var sb strings.Builder
	WritePretty(&sb, results, "dev.env", "prod.env", FormatOptions{Color: false, Verbose: false})
	out := sb.String()

	if !strings.Contains(out, "Summary:") {
		t.Error("expected Summary line in output")
	}
	if !strings.Contains(out, "dev.env") || !strings.Contains(out, "prod.env") {
		t.Error("expected file names in output")
	}
}

func TestWritePretty_VerboseShowsValues(t *testing.T) {
	results := []Result{
		{Key: "DB_PASS", Status: StatusMismatch, ValueA: "secret", ValueB: "other"},
	}

	var sb strings.Builder
	WritePretty(&sb, results, "a.env", "b.env", FormatOptions{Color: false, Verbose: true})
	out := sb.String()

	if !strings.Contains(out, "secret") {
		t.Error("expected ValueA in verbose output")
	}
	if !strings.Contains(out, "other") {
		t.Error("expected ValueB in verbose output")
	}
}

func TestWritePretty_NonVerboseHidesValues(t *testing.T) {
	results := []Result{
		{Key: "DB_PASS", Status: StatusMismatch, ValueA: "topsecret", ValueB: "anothersecret"},
	}

	var sb strings.Builder
	WritePretty(&sb, results, "a.env", "b.env", FormatOptions{Color: false, Verbose: false})
	out := sb.String()

	if strings.Contains(out, "topsecret") {
		t.Error("did not expect ValueA in non-verbose output")
	}
}

func TestWritePretty_SortedOutput(t *testing.T) {
	results := []Result{
		{Key: "Z_KEY", Status: StatusMatch, ValueA: "z", ValueB: "z"},
		{Key: "A_KEY", Status: StatusMissingInB, ValueA: "a", ValueB: ""},
	}

	var sb strings.Builder
	WritePretty(&sb, results, "a.env", "b.env", FormatOptions{Color: false, Verbose: false})
	out := sb.String()

	idxA := strings.Index(out, "A_KEY")
	idxZ := strings.Index(out, "Z_KEY")
	if idxA == -1 || idxZ == -1 {
		t.Fatal("expected both keys in output")
	}
	if idxA > idxZ {
		t.Error("expected A_KEY to appear before Z_KEY (sorted output)")
	}
}

func TestWritePretty_ColorEnabled(t *testing.T) {
	results := []Result{
		{Key: "API_KEY", Status: StatusMissingInB, ValueA: "abc", ValueB: ""},
	}

	var sb strings.Builder
	WritePretty(&sb, results, "a.env", "b.env", FormatOptions{Color: true, Verbose: false})
	out := sb.String()

	if !strings.Contains(out, "\033[") {
		t.Error("expected ANSI escape codes when Color is true")
	}
}
