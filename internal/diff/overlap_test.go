package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"prod": {"HOST": "prod.example.com", "PORT": "443", "SECRET": "abc"},
		"staging": {"HOST": "staging.example.com", "PORT": "443", "DEBUG": "true"},
		"dev": {"HOST": "localhost", "PORT": "8080", "DEBUG": "true", "SECRET": "dev"},
	}
}

func findOverlap(entries []OverlapEntry, a, b string) *OverlapEntry {
	for i := range entries {
		if (entries[i].LabelA == a && entries[i].LabelB == b) ||
			(entries[i].LabelA == b && entries[i].LabelB == a) {
			return &entries[i]
		}
	}
	return nil
}

func TestBuildOverlap_PairCount(t *testing.T) {
	envs := buildEnvs()
	entries := BuildOverlap(envs)
	if len(entries) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(entries))
	}
}

func TestBuildOverlap_SharedKeys(t *testing.T) {
	entries := BuildOverlap(buildEnvs())
	e := findOverlap(entries, "prod", "staging")
	if e == nil {
		t.Fatal("prod/staging pair not found")
	}
	// HOST and PORT are shared
	if e.Shared != 2 {
		t.Errorf("expected 2 shared keys, got %d", e.Shared)
	}
}

func TestBuildOverlap_Conflicts(t *testing.T) {
	entries := BuildOverlap(buildEnvs())
	e := findOverlap(entries, "prod", "staging")
	if e == nil {
		t.Fatal("prod/staging pair not found")
	}
	// HOST differs, PORT matches → 1 conflict
	if e.Conflicts != 1 {
		t.Errorf("expected 1 conflict, got %d", e.Conflicts)
	}
}

func TestBuildOverlap_OnlyKeys(t *testing.T) {
	entries := BuildOverlap(buildEnvs())
	e := findOverlap(entries, "prod", "staging")
	if e == nil {
		t.Fatal("prod/staging pair not found")
	}
	if e.OnlyInA+e.OnlyInB == 0 {
		t.Error("expected some exclusive keys")
	}
}

func TestOverlapScore_AllShared(t *testing.T) {
	e := OverlapEntry{Shared: 5, Conflicts: 2, OnlyInA: 0, OnlyInB: 0}
	if e.OverlapScore() != 1.0 {
		t.Errorf("expected 1.0, got %f", e.OverlapScore())
	}
}

func TestOverlapScore_Empty(t *testing.T) {
	e := OverlapEntry{}
	if e.OverlapScore() != 1.0 {
		t.Errorf("expected 1.0 for empty entry, got %f", e.OverlapScore())
	}
}

func TestWriteOverlap_NoEntries(t *testing.T) {
	var buf bytes.Buffer
	WriteOverlap(&buf, nil)
	if !strings.Contains(buf.String(), "no environment") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestWriteOverlap_ContainsLabels(t *testing.T) {
	entries := BuildOverlap(buildEnvs())
	var buf bytes.Buffer
	WriteOverlap(&buf, entries)
	out := buf.String()
	for _, label := range []string{"prod", "staging", "dev"} {
		if !strings.Contains(out, label) {
			t.Errorf("expected label %q in output", label)
		}
	}
}

func TestWriteOverlap_ContainsScore(t *testing.T) {
	entries := BuildOverlap(buildEnvs())
	var buf bytes.Buffer
	WriteOverlap(&buf, entries)
	if !strings.Contains(buf.String(), "%") {
		t.Error("expected percentage score in output")
	}
}
