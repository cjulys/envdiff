package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildPivotResult(key, valA, valB string, status Status) Result {
	return Result{Key: key, ValueA: valA, ValueB: valB, Status: status}
}

func TestBuildPivot_AllMatch(t *testing.T) {
	results := []Result{
		buildPivotResult("DB_HOST", "localhost", "localhost", StatusMatch),
	}
	pt := BuildPivot(results, []string{"dev", "prod"}, false)

	if len(pt.Rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(pt.Rows))
	}
	row := pt.Rows[0]
	if row.Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", row.Key)
	}
	if row.Status != string(StatusMatch) {
		t.Errorf("expected status match, got %s", row.Status)
	}
}

func TestBuildPivot_VerboseShowsValues(t *testing.T) {
	results := []Result{
		buildPivotResult("API_KEY", "abc123", "xyz789", StatusMismatch),
	}
	pt := BuildPivot(results, []string{"dev", "prod"}, true)

	row := pt.Rows[0]
	if row.Cells["dev"] != "abc123" {
		t.Errorf("expected abc123, got %s", row.Cells["dev"])
	}
	if row.Cells["prod"] != "xyz789" {
		t.Errorf("expected xyz789, got %s", row.Cells["prod"])
	}
}

func TestBuildPivot_NonVerboseHidesValues(t *testing.T) {
	results := []Result{
		buildPivotResult("SECRET", "abc", "def", StatusMismatch),
	}
	pt := BuildPivot(results, []string{"dev", "prod"}, false)

	row := pt.Rows[0]
	if row.Cells["dev"] != "(set)" {
		t.Errorf("expected (set), got %s", row.Cells["dev"])
	}
}

func TestBuildPivot_MissingInB(t *testing.T) {
	results := []Result{
		buildPivotResult("ONLY_DEV", "value", "", StatusMissingInB),
	}
	pt := BuildPivot(results, []string{"dev", "prod"}, true)

	row := pt.Rows[0]
	if row.Cells["prod"] != "(missing)" {
		t.Errorf("expected (missing), got %s", row.Cells["prod"])
	}
}

func TestBuildPivot_SortedKeys(t *testing.T) {
	results := []Result{
		buildPivotResult("Z_KEY", "1", "1", StatusMatch),
		buildPivotResult("A_KEY", "2", "2", StatusMatch),
		buildPivotResult("M_KEY", "3", "3", StatusMatch),
	}
	pt := BuildPivot(results, []string{"dev", "prod"}, false)

	if pt.Rows[0].Key != "A_KEY" || pt.Rows[1].Key != "M_KEY" || pt.Rows[2].Key != "Z_KEY" {
		t.Errorf("rows not sorted: %v", []string{pt.Rows[0].Key, pt.Rows[1].Key, pt.Rows[2].Key})
	}
}

func TestPivotTable_Write_ContainsHeaders(t *testing.T) {
	results := []Result{
		buildPivotResult("PORT", "8080", "9090", StatusMismatch),
	}
	pt := BuildPivot(results, []string{"staging", "prod"}, true)

	var buf bytes.Buffer
	pt.Write(&buf)
	out := buf.String()

	for _, want := range []string{"KEY", "staging", "prod", "STATUS", "PORT", "mismatch"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}
