package env_test

import (
	"os"
	"testing"

	"github.com/yourorg/envdiff/internal/env"
)

func TestExpand_BracketSyntax(t *testing.T) {
	r := env.NewResolver(map[string]string{"HOST": "localhost", "PORT": "5432"}, false)
	got := r.Expand("postgres://${HOST}:${PORT}/db")
	want := "postgres://localhost:5432/db"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestExpand_BareVariable(t *testing.T) {
	r := env.NewResolver(map[string]string{"SCHEME": "https"}, false)
	got := r.Expand("$SCHEME://example.com")
	if got != "https://example.com" {
		t.Errorf("unexpected: %q", got)
	}
}

func TestExpand_UnresolvableFallsToEmpty(t *testing.T) {
	r := env.NewResolver(map[string]string{}, false)
	got := r.Expand("${MISSING}_value")
	if got != "_value" {
		t.Errorf("got %q, want %q", got, "_value")
	}
}

func TestExpand_FallsBackToSysEnv(t *testing.T) {
	t.Setenv("SYS_VAR", "from_os")
	r := env.NewResolver(map[string]string{}, true)
	got := r.Expand("${SYS_VAR}")
	if got != "from_os" {
		t.Errorf("got %q, want %q", got, "from_os")
	}
}

func TestExpand_MapTakesPrecedenceOverSys(t *testing.T) {
	os.Setenv("PRIORITY", "os_value")
	t.Cleanup(func() { os.Unsetenv("PRIORITY") })
	r := env.NewResolver(map[string]string{"PRIORITY": "map_value"}, true)
	got := r.Expand("${PRIORITY}")
	if got != "map_value" {
		t.Errorf("got %q, want %q", got, "map_value")
	}
}

func TestExpandAll_ReturnsExpandedMap(t *testing.T) {
	m := map[string]string{
		"BASE": "http://localhost",
		"URL":  "${BASE}/api",
	}
	r := env.NewResolver(m, false)
	out := r.ExpandAll()
	if out["URL"] != "http://localhost/api" {
		t.Errorf("URL not expanded: %q", out["URL"])
	}
	if out["BASE"] != "http://localhost" {
		t.Errorf("BASE changed unexpectedly: %q", out["BASE"])
	}
}
