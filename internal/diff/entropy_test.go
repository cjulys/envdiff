package diff

import (
	"strings"
	"testing"
)

func buildEntropyEnvs() map[string]map[string]string {
	return map[string]map[string]string{
		"prod": {"DB_HOST": "prod.db", "APP_ENV": "production", "SECRET": "abc"},
		"staging": {"DB_HOST": "staging.db", "APP_ENV": "staging", "SECRET": "abc"},
		"dev": {"DB_HOST": "localhost", "APP_ENV": "development", "SECRET": "abc"},
	}
}

func TestBuildEntropy_EntryCount(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestBuildEntropy_ZeroEntropyForConstantValue(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	var secretEntry *EntropyEntry
	for i := range entries {
		if entries[i].Key == "SECRET" {
			secretEntry = &entries[i]
			break
		}
	}
	if secretEntry == nil {
		t.Fatal("expected SECRET entry")
	}
	if secretEntry.Entropy != 0.0 {
		t.Errorf("expected entropy 0 for constant value, got %f", secretEntry.Entropy)
	}
	if secretEntry.Unique != 1 {
		t.Errorf("expected 1 unique value, got %d", secretEntry.Unique)
	}
}

func TestBuildEntropy_HighEntropyForAllDifferent(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	var dbEntry *EntropyEntry
	for i := range entries {
		if entries[i].Key == "DB_HOST" {
			dbEntry = &entries[i]
			break
		}
	}
	if dbEntry == nil {
		t.Fatal("expected DB_HOST entry")
	}
	if dbEntry.Entropy <= 1.0 {
		t.Errorf("expected high entropy for all-different values, got %f", dbEntry.Entropy)
	}
	if dbEntry.Unique != 3 {
		t.Errorf("expected 3 unique values, got %d", dbEntry.Unique)
	}
}

func TestBuildEntropy_SortedByEntropyDesc(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	for i := 1; i < len(entries); i++ {
		if entries[i].Entropy > entries[i-1].Entropy {
			t.Errorf("entries not sorted by entropy desc at index %d", i)
		}
	}
}

func TestWriteEntropy_NoEntries(t *testing.T) {
	var sb strings.Builder
	WriteEntropy(&sb, []EntropyEntry{}, 0)
	if !strings.Contains(sb.String(), "no keys") {
		t.Errorf("expected 'no keys' message, got: %s", sb.String())
	}
}

func TestWriteEntropy_TopNLimitsOutput(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	var sb strings.Builder
	WriteEntropy(&sb, entries, 1)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	// header + separator + 1 data row
	if len(lines) != 3 {
		t.Errorf("expected 3 lines with topN=1, got %d: %v", len(lines), lines)
	}
}

func TestWriteEntropy_ContainsKeyName(t *testing.T) {
	envs := buildEntropyEnvs()
	entries := BuildEntropy(envs)
	var sb strings.Builder
	WriteEntropy(&sb, entries, 0)
	if !strings.Contains(sb.String(), "DB_HOST") {
		t.Errorf("expected DB_HOST in output")
	}
}
