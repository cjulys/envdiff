package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildRadarEnvs() map[string][]Result {
	return map[string][]Result{
		"production": {
			{Key: "A", Status: StatusMatch},
			{Key: "B", Status: StatusMismatch},
			{Key: "C", Status: StatusMissingInB},
			{Key: "D", Status: StatusMatch},
		},
		"staging": {
			{Key: "A", Status: StatusMatch},
			{Key: "B", Status: StatusMatch},
		},
		"dev": {
			{Key: "A", Status: StatusMissingInA},
			{Key: "B", Status: StatusMissingInB},
			{Key: "C", Status: StatusMismatch},
		},
	}
}

func TestBuildRadar_EntryCount(t *testing.T) {
	envs := buildRadarEnvs()
	report := BuildRadar(envs)
	if len(report.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(report.Entries))
	}
}

func TestBuildRadar_RatesAreCorrect(t *testing.T) {
	envs := buildRadarEnvs()
	report := BuildRadar(envs)
	for _, e := range report.Entries {
		switch e.Label {
		case "production":
			if e.Problems != 2 || e.Total != 4 {
				t.Errorf("production: want 2/4, got %d/%d", e.Problems, e.Total)
			}
		case "staging":
			if e.Problems != 0 || e.Total != 2 {
				t.Errorf("staging: want 0/2, got %d/%d", e.Problems, e.Total)
			}
		case "dev":
			if e.Problems != 3 || e.Total != 3 {
				t.Errorf("dev: want 3/3, got %d/%d", e.Problems, e.Total)
			}
		}
	}
}

func TestBuildRadar_SortedByRateDesc(t *testing.T) {
	envs := buildRadarEnvs()
	report := BuildRadar(envs)
	for i := 1; i < len(report.Entries); i++ {
		if report.Entries[i].Rate > report.Entries[i-1].Rate {
			t.Errorf("entries not sorted descending at index %d", i)
		}
	}
}

func TestBuildRadar_EmptyEnvs(t *testing.T) {
	report := BuildRadar(map[string][]Result{})
	if len(report.Entries) != 0 {
		t.Fatalf("expected empty report")
	}
}

func TestWriteRadar_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteRadar(&buf, RadarReport{})
	if !strings.Contains(buf.String(), "no environments") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteRadar_ContainsLabels(t *testing.T) {
	envs := buildRadarEnvs()
	report := BuildRadar(envs)
	var buf bytes.Buffer
	WriteRadar(&buf, report)
	out := buf.String()
	for _, label := range []string{"production", "staging", "dev"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestWriteRadar_ContainsPercentage(t *testing.T) {
	envs := map[string][]Result{
		"alpha": {
			{Key: "X", Status: StatusMismatch},
			{Key: "Y", Status: StatusMatch},
		},
	}
	report := BuildRadar(envs)
	var buf bytes.Buffer
	WriteRadar(&buf, report)
	if !strings.Contains(buf.String(), "50.0%") {
		t.Errorf("expected 50.0%% in output, got: %s", buf.String())
	}
}
