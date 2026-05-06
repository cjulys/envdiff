package config_test

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func TestValidate_RequiresTwoFiles(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"only-one.env"}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when fewer than two files provided")
	}
}

func TestValidate_ValidFormats(t *testing.T) {
	formats := []config.OutputFormat{
		config.FormatPretty,
		config.FormatMarkdown,
		config.FormatJSON,
	}

	for _, f := range formats {
		cfg := config.Default()
		cfg.Files = []string{"a.env", "b.env"}
		cfg.Format = f

		if err := cfg.Validate(); err != nil {
			t.Errorf("format %q should be valid, got: %v", f, err)
		}
	}
}

func TestValidate_InvalidFormat(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"a.env", "b.env"}
	cfg.Format = "xml"

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown format 'xml'")
	}
}

func TestValidate_EmptyFormatDefaultsToPretty(t *testing.T) {
	cfg := &config.Config{
		Files:  []string{"a.env", "b.env"},
		Format: "",
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Format != config.FormatPretty {
		t.Errorf("expected default format %q, got %q", config.FormatPretty, cfg.Format)
	}
}

func TestValidate_InvalidFilterStatus(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"a.env", "b.env"}
	cfg.FilterStatus = []string{"unknown-status"}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for unknown filter status")
	}
}

func TestValidate_ValidFilterStatuses(t *testing.T) {
	cfg := config.Default()
	cfg.Files = []string{"a.env", "b.env"}
	cfg.FilterStatus = []string{"match", "mismatch", "missing"}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefault_Values(t *testing.T) {
	cfg := config.Default()

	if cfg.Format != config.FormatPretty {
		t.Errorf("expected FormatPretty, got %q", cfg.Format)
	}
	if !cfg.Color {
		t.Error("expected Color to be true by default")
	}
	if cfg.ExitOnDiff {
		t.Error("expected ExitOnDiff to be false by default")
	}
}
