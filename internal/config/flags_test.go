package config_test

import (
	"testing"

	"github.com/user/envdiff/internal/config"
)

func TestParseFlags_Defaults(t *testing.T) {
	cfg, files, err := config.ParseFlags([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Format != config.FormatPretty {
		t.Errorf("expected FormatPretty, got %q", cfg.Format)
	}
	if cfg.Verbose {
		t.Error("expected Verbose=false by default")
	}
	if !cfg.Color {
		t.Error("expected Color=true by default")
	}
	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}
}

func TestParseFlags_FormatMarkdown(t *testing.T) {
	cfg, _, err := config.ParseFlags([]string{"-format", "markdown", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != config.FormatMarkdown {
		t.Errorf("expected FormatMarkdown, got %q", cfg.Format)
	}
}

func TestParseFlags_VerboseAndExitOnDiff(t *testing.T) {
	cfg, _, err := config.ParseFlags([]string{"-verbose", "-exit-on-diff", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Verbose {
		t.Error("expected Verbose=true")
	}
	if !cfg.ExitOnDiff {
		t.Error("expected ExitOnDiff=true")
	}
}

func TestParseFlags_MultipleKeys(t *testing.T) {
	cfg, _, err := config.ParseFlags([]string{"-key", "FOO", "-key", "BAR", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.FilterKeys) != 2 {
		t.Errorf("expected 2 filter keys, got %d", len(cfg.FilterKeys))
	}
}

func TestParseFlags_TooFewFiles(t *testing.T) {
	_, _, err := config.ParseFlags([]string{"only.env"})
	if err == nil {
		t.Fatal("expected error when only one file provided")
	}
}

func TestParseFlags_PrefixFlag(t *testing.T) {
	cfg, _, err := config.ParseFlags([]string{"-prefix", "APP_", "a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.FilterPrefix != "APP_" {
		t.Errorf("expected FilterPrefix=APP_, got %q", cfg.FilterPrefix)
	}
}
