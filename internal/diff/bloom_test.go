package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildBloomResults(keys ...string) []Result {
	out := make([]Result, len(keys))
	for i, k := range keys {
		out[i] = Result{Key: k, Status: StatusMissing}
	}
	return out
}

func TestNewBloomReport_IsEmpty(t *testing.T) {
	b := NewBloomReport()
	if len(b.Sorted()) != 0 {
		t.Fatal("expected empty report")
	}
}

func TestBloom_AddRun_CountsHits(t *testing.T) {
	b := NewBloomReport()
	b.AddRun(buildBloomResults("KEY_A", "KEY_B"))
	b.AddRun(buildBloomResults("KEY_A"))

	entries := b.Sorted()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "KEY_A" {
		t.Errorf("expected KEY_A first, got %s", entries[0].Key)
	}
	if entries[0].Hits != 2 {
		t.Errorf("expected 2 hits for KEY_A, got %d", entries[0].Hits)
	}
}

func TestBloom_FrequencyCalculation(t *testing.T) {
	b := NewBloomReport()
	b.AddRun(buildBloomResults("KEY_A"))
	b.AddRun(buildBloomResults("KEY_A"))
	b.AddRun(buildBloomResults())

	entries := b.Sorted()
	if len(entries) == 0 {
		t.Fatal("expected entries")
	}
	if entries[0].Total != 3 {
		t.Errorf("expected total=3, got %d", entries[0].Total)
	}
	want := 2.0 / 3.0
	if entries[0].Freq < want-0.001 || entries[0].Freq > want+0.001 {
		t.Errorf("expected freq ~%.3f, got %.3f", want, entries[0].Freq)
	}
}

func TestBloom_MatchResultsNotCounted(t *testing.T) {
	b := NewBloomReport()
	b.AddRun([]Result{{Key: "KEY_A", Status: StatusMatch}})
	if len(b.Sorted()) != 0 {
		t.Fatal("match results should not be counted")
	}
}

func TestWriteBloom_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteBloom(&buf, NewBloomReport(), 0.5)
	if !strings.Contains(buf.String(), "No problem keys") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteBloom_ThresholdMarker(t *testing.T) {
	b := NewBloomReport()
	b.AddRun(buildBloomResults("KEY_A"))
	b.AddRun(buildBloomResults("KEY_A"))

	var buf bytes.Buffer
	WriteBloom(&buf, b, 0.5)
	out := buf.String()
	if !strings.Contains(out, "!") {
		t.Errorf("expected threshold marker '!', got:\n%s", out)
	}
	if !strings.Contains(out, "KEY_A") {
		t.Errorf("expected KEY_A in output, got:\n%s", out)
	}
}
