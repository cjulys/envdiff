package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildClassifyResults() []Result {
	return []Result{
		{Key: "APP_NAME", Status: StatusMatch, ValueA: "myapp", ValueB: "myapp"},
		{Key: "DB_PASSWORD", Status: StatusMismatch, ValueA: "old", ValueB: "new"},
		{Key: "PORT", Status: StatusMissingInB, ValueA: "8080", ValueB: ""},
		{Key: "API_TOKEN", Status: StatusMissingInA, ValueA: "", ValueB: "tok"},
	}
}

func TestWriteClassifyReport_ContainsSummary(t *testing.T) {
	var buf bytes.Buffer
	WriteClassifyReport(&buf, buildClassifyResults(), ClassifyOptions{})
	out := buf.String()
	if !strings.Contains(out, "Severity Summary") {
		t.Errorf("expected summary header, got:\n%s", out)
	}
}

func TestWriteClassifyReport_EmptyResults(t *testing.T) {
	var buf bytes.Buffer
	WriteClassifyReport(&buf, nil, ClassifyOptions{})
	if !strings.Contains(buf.String(), "No results") {
		t.Errorf("expected empty message")
	}
}

func TestWriteClassifyReport_HighSeverityFirst(t *testing.T) {
	var buf bytes.Buffer
	WriteClassifyReport(&buf, buildClassifyResults(), ClassifyOptions{})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// Find first data line (after summary + blank)
	var dataLines []string
	for _, l := range lines {
		if strings.Contains(l, "[") {
			dataLines = append(dataLines, l)
		}
	}
	if len(dataLines) == 0 {
		t.Fatal("no data lines found")
	}
	if !strings.Contains(dataLines[0], "high") {
		t.Errorf("expected first data line to be high severity, got: %s", dataLines[0])
	}
}

func TestWriteClassifyReport_VerboseShowsValues(t *testing.T) {
	var buf bytes.Buffer
	WriteClassifyReport(&buf, buildClassifyResults(), ClassifyOptions{Verbose: true})
	out := buf.String()
	if !strings.Contains(out, "a=") || !strings.Contains(out, "b=") {
		t.Errorf("verbose mode should include values, got:\n%s", out)
	}
}

func TestWriteClassifyReport_NonVerboseHidesValues(t *testing.T) {
	var buf bytes.Buffer
	WriteClassifyReport(&buf, buildClassifyResults(), ClassifyOptions{Verbose: false})
	out := buf.String()
	if strings.Contains(out, "a=") {
		t.Errorf("non-verbose mode should not include values, got:\n%s", out)
	}
}
