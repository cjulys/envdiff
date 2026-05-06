package template_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/template"
)

func TestGenerate_SortedKeys(t *testing.T) {
	envs := []map[string]string{
		{"ZEBRA": "z", "APPLE": "a"},
	}
	var buf strings.Builder
	if err := template.Generate(&buf, envs, template.DefaultOptions()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	applePos := strings.Index(out, "APPLE=")
	zebraPos := strings.Index(out, "ZEBRA=")
	if applePos == -1 || zebraPos == -1 {
		t.Fatal("expected both keys in output")
	}
	if applePos > zebraPos {
		t.Error("expected APPLE to appear before ZEBRA")
	}
}

func TestGenerate_EmptyValues(t *testing.T) {
	envs := []map[string]string{
		{"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	var buf strings.Builder
	opts := template.Options{Placeholder: "", IncludeComments: false}
	if err := template.Generate(&buf, envs, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DB_HOST=") {
		t.Error("expected DB_HOST= in output")
	}
	if strings.Contains(out, "localhost") {
		t.Error("original value should not appear in template")
	}
}

func TestGenerate_CustomPlaceholder(t *testing.T) {
	envs := []map[string]string{{"API_KEY": "secret"}}
	var buf strings.Builder
	opts := template.Options{Placeholder: "CHANGEME", IncludeComments: false}
	if err := template.Generate(&buf, envs, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "API_KEY=CHANGEME") {
		t.Errorf("expected API_KEY=CHANGEME, got: %s", out)
	}
}

func TestGenerate_UnionOfKeys(t *testing.T) {
	envs := []map[string]string{
		{"FOO": "1"},
		{"BAR": "2"},
		{"FOO": "3"},
	}
	var buf strings.Builder
	opts := template.Options{Placeholder: "", IncludeComments: false}
	if err := template.Generate(&buf, envs, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Count(out, "FOO=") != 1 {
		t.Error("expected FOO to appear exactly once")
	}
	if !strings.Contains(out, "BAR=") {
		t.Error("expected BAR in output")
	}
}

func TestGenerate_IncludesHeaderComment(t *testing.T) {
	envs := []map[string]string{{"X": "1"}}
	var buf strings.Builder
	opts := template.DefaultOptions()
	if err := template.Generate(&buf, envs, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(buf.String(), "#") {
		t.Error("expected output to start with a comment")
	}
}

func TestGenerate_NoHeaderComment(t *testing.T) {
	envs := []map[string]string{{"X": "1"}}
	var buf strings.Builder
	opts := template.Options{IncludeComments: false}
	if err := template.Generate(&buf, envs, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.HasPrefix(buf.String(), "#") {
		t.Error("expected output to not start with a comment")
	}
}
