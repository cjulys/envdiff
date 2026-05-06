package export_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/export"
)

func TestWrite_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleResults(), export.Options{Format: export.FormatJSON, Verbose: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(buf.String()), "[") {
		t.Error("expected JSON array output")
	}
}

func TestWrite_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleResults(), export.Options{Format: export.FormatCSV, Verbose: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "key,status") {
		t.Error("expected CSV header in output")
	}
}

func TestWrite_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := export.Write(&buf, sampleResults(), export.Options{Format: "xml"})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "xml") {
		t.Errorf("error should mention format, got: %v", err)
	}
}

func TestIsSupported(t *testing.T) {
	if !export.IsSupported(export.FormatJSON) {
		t.Error("JSON should be supported")
	}
	if !export.IsSupported(export.FormatCSV) {
		t.Error("CSV should be supported")
	}
	if export.IsSupported("yaml") {
		t.Error("yaml should not be supported")
	}
}
