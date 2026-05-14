package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildDriftCmdResults(problems int) []Result {
	var r []Result
	for i := 0; i < problems; i++ {
		r = append(r, Result{Key: "K", Status: StatusMismatch})
	}
	r = append(r, Result{Key: "OK", Status: StatusMatch})
	return r
}

func TestRunDrift_NoEnvironments(t *testing.T) {
	var buf bytes.Buffer
	got := RunDrift(&buf, DefaultDriftOptions())
	if got {
		t.Error("expected false when no environments provided")
	}
	if !strings.Contains(buf.String(), "No environments") {
		t.Errorf("expected no-env message, got: %s", buf.String())
	}
}

func TestRunDrift_NoDifferences_ReturnsFalse(t *testing.T) {
	opts := DefaultDriftOptions()
	opts.Labels["prod"] = []Result{{Key: "A", Status: StatusMatch}}
	var buf bytes.Buffer
	got := RunDrift(&buf, opts)
	if got {
		t.Error("expected false when no problems")
	}
}

func TestRunDrift_WithProblems_ReturnsTrue(t *testing.T) {
	opts := DefaultDriftOptions()
	opts.Labels["staging"] = buildDriftCmdResults(3)
	var buf bytes.Buffer
	got := RunDrift(&buf, opts)
	if !got {
		t.Error("expected true when problems exist")
	}
}

func TestRunDrift_OutputContainsLabel(t *testing.T) {
	opts := DefaultDriftOptions()
	opts.Labels["myenv"] = buildDriftCmdResults(1)
	var buf bytes.Buffer
	RunDrift(&buf, opts)
	if !strings.Contains(buf.String(), "myenv") {
		t.Errorf("expected label 'myenv' in output, got: %s", buf.String())
	}
}

func TestRunDrift_MultipleEnvironments(t *testing.T) {
	opts := DefaultDriftOptions()
	opts.Labels["dev"] = buildDriftCmdResults(0)
	opts.Labels["prod"] = buildDriftCmdResults(2)
	var buf bytes.Buffer
	got := RunDrift(&buf, opts)
	if !got {
		t.Error("expected true because prod has problems")
	}
	out := buf.String()
	if !strings.Contains(out, "dev") || !strings.Contains(out, "prod") {
		t.Errorf("expected both labels in output, got: %s", out)
	}
}
