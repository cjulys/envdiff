package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/export"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMismatch, ValueA: "localhost", ValueB: "prod.db"},
		{Key: "API_KEY", Status: diff.StatusMissingInB, ValueA: "secret", ValueB: ""},
		{Key: "PORT", Status: diff.StatusMatch, ValueA: "8080", ValueB: "8080"},
	}
}

func TestWriteJSON_Structure(t *testing.T) {
	var buf bytes.Buffer
	if err := export.WriteJSON(&buf, sampleResults(), true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.JSONRecord
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 3 {
		t.Errorf("expected 3 records, got %d", len(records))
	}
}

func TestWriteJSON_VerboseIncludesValues(t *testing.T) {
	var buf bytes.Buffer
	_ = export.WriteJSON(&buf, sampleResults(), true)
	if !strings.Contains(buf.String(), "localhost") {
		t.Error("expected verbose output to contain value_a")
	}
}

func TestWriteJSON_NonVerboseOmitsValues(t *testing.T) {
	var buf bytes.Buffer
	_ = export.WriteJSON(&buf, sampleResults(), false)
	if strings.Contains(buf.String(), "localhost") {
		t.Error("expected non-verbose output to omit values")
	}
}

func TestWriteCSV_Header(t *testing.T) {
	var buf bytes.Buffer
	if err := export.WriteCSV(&buf, sampleResults(), false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "key,status,value_a,value_b" {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestWriteCSV_RowCount(t *testing.T) {
	var buf bytes.Buffer
	_ = export.WriteCSV(&buf, sampleResults(), false)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// 1 header + 3 data rows
	if len(lines) != 4 {
		t.Errorf("expected 4 lines, got %d", len(lines))
	}
}

func TestWriteCSV_SortedOutput(t *testing.T) {
	var buf bytes.Buffer
	_ = export.WriteCSV(&buf, sampleResults(), false)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if !strings.HasPrefix(lines[1], "API_KEY") {
		t.Errorf("expected first data row to be API_KEY, got: %s", lines[1])
	}
}
