package diff

import (
	"strings"
	"testing"
)

func buildDensityResults() []Result {
	return []Result{
		{Key: "DB_HOST", Status: StatusMatch},
		{Key: "DB_PORT", Status: StatusMismatch},
		{Key: "DB_PASS", Status: StatusMissingInB},
		{Key: "APP_ENV", Status: StatusMatch},
		{Key: "APP_DEBUG", Status: StatusMismatch},
		{Key: "TOKEN", Status: StatusMissingInA},
	}
}

func findDensity(entries []DensityEntry, prefix string) (DensityEntry, bool) {
	for _, e := range entries {
		if e.Prefix == prefix {
			return e, true
		}
	}
	return DensityEntry{}, false
}

func TestBuildDensity_EntryCount(t *testing.T) {
	report := BuildDensity(buildDensityResults())
	if len(report.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(report.Entries))
	}
}

func TestBuildDensity_DBPrefix(t *testing.T) {
	report := BuildDensity(buildDensityResults())
	e, ok := findDensity(report.Entries, "DB")
	if !ok {
		t.Fatal("expected DB prefix entry")
	}
	if e.Total != 3 {
		t.Errorf("expected total 3, got %d", e.Total)
	}
	if e.Problems != 2 {
		t.Errorf("expected 2 problems, got %d", e.Problems)
	}
}

func TestBuildDensity_RootPrefix(t *testing.T) {
	report := BuildDensity(buildDensityResults())
	e, ok := findDensity(report.Entries, "(root)")
	if !ok {
		t.Fatal("expected (root) prefix entry")
	}
	if e.Total != 1 {
		t.Errorf("expected total 1, got %d", e.Total)
	}
	if e.Problems != 1 {
		t.Errorf("expected 1 problem, got %d", e.Problems)
	}
}

func TestBuildDensity_SortedByDensityDesc(t *testing.T) {
	report := BuildDensity(buildDensityResults())
	for i := 1; i < len(report.Entries); i++ {
		if report.Entries[i].Density > report.Entries[i-1].Density {
			t.Errorf("entries not sorted by density desc at index %d", i)
		}
	}
}

func TestBuildDensity_EmptyResults(t *testing.T) {
	report := BuildDensity(nil)
	if len(report.Entries) != 0 {
		t.Errorf("expected empty entries, got %d", len(report.Entries))
	}
}

func TestWriteDensity_NoEntries(t *testing.T) {
	var sb strings.Builder
	WriteDensity(&sb, DensityReport{})
	if !strings.Contains(sb.String(), "no data") {
		t.Error("expected 'no data' message")
	}
}

func TestWriteDensity_ContainsHeader(t *testing.T) {
	var sb strings.Builder
	WriteDensity(&sb, BuildDensity(buildDensityResults()))
	out := sb.String()
	for _, col := range []string{"PREFIX", "TOTAL", "PROBLEMS", "DENSITY"} {
		if !strings.Contains(out, col) {
			t.Errorf("expected column %q in output", col)
		}
	}
}

func TestWriteDensity_ContainsPrefixes(t *testing.T) {
	var sb strings.Builder
	WriteDensity(&sb, BuildDensity(buildDensityResults()))
	out := sb.String()
	for _, p := range []string{"DB", "APP", "(root)"} {
		if !strings.Contains(out, p) {
			t.Errorf("expected prefix %q in output", p)
		}
	}
}
