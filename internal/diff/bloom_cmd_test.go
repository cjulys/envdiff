package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildBloomCmdResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: StatusMissing}
	}
	return out
}

func TestRunBloom_NoDifferences_ReturnsFalse(t *testing.T) {
	var buf bytes.Buffer
	got := RunBloom(&buf, [][]Result{}, DefaultBloomOptions())
	if got {
		t.Error("expected false for empty runs")
	}
}

func TestRunBloom_BelowThreshold_ReturnsFalse(t *testing.T) {
	var buf bytes.Buffer
	runs := [][]Result{
		buildBloomCmdResults("KEY_A"),
		{},
		{},
		{},
	}
	opts := DefaultBloomOptions() // threshold 0.5
	got := RunBloom(&buf, runs, opts)
	if got {
		t.Error("expected false: freq 0.25 < threshold 0.5")
	}
}

func TestRunBloom_AboveThreshold_ReturnsTrue(t *testing.T) {
	var buf bytes.Buffer
	runs := [][]Result{
		buildBloomCmdResults("KEY_A"),
		buildBloomCmdResults("KEY_A"),
	}
	got := RunBloom(&buf, runs, DefaultBloomOptions())
	if !got {
		t.Error("expected true: freq 1.0 >= threshold 0.5")
	}
}

func TestRunBloom_OutputContainsKey(t *testing.T) {
	var buf bytes.Buffer
	runs := [][]Result{
		buildBloomCmdResults("MY_KEY"),
		buildBloomCmdResults("MY_KEY"),
	}
	RunBloom(&buf, runs, DefaultBloomOptions())
	if !strings.Contains(buf.String(), "MY_KEY") {
		t.Errorf("expected MY_KEY in output:\n%s", buf.String())
	}
}

func TestRunBloom_DefaultOptions(t *testing.T) {
	opts := DefaultBloomOptions()
	if opts.Threshold != 0.5 {
		t.Errorf("expected default threshold 0.5, got %.2f", opts.Threshold)
	}
	if opts.Verbose {
		t.Error("expected verbose=false by default")
	}
}

func TestRunBloom_VerboseOption_ShowsFlagged(t *testing.T) {
	var buf bytes.Buffer
	runs := [][]Result{
		buildBloomCmdResults("KEY_X"),
		buildBloomCmdResults("KEY_X"),
	}
	opts := DefaultBloomOptions()
	opts.Verbose = true
	RunBloom(&buf, runs, opts)
	if !strings.Contains(buf.String(), "flagged") {
		t.Errorf("expected 'flagged' in verbose output:\n%s", buf.String())
	}
}
