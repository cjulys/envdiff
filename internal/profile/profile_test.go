package profile_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envdiff/internal/profile"
)

func writeTempRegistry(t *testing.T, reg *profile.Registry) string {
	t.Helper()
	data, err := json.Marshal(reg)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "profiles.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestNew_IsEmpty(t *testing.T) {
	r := profile.New()
	if len(r.List()) != 0 {
		t.Errorf("expected empty registry, got %d profiles", len(r.List()))
	}
}

func TestAdd_And_Get(t *testing.T) {
	r := profile.New()
	p := profile.Profile{Name: "staging", Files: []string{".env.staging"}}
	r.Add(p)
	got, ok := r.Get("staging")
	if !ok {
		t.Fatal("expected profile to exist")
	}
	if got.Name != "staging" || len(got.Files) != 1 {
		t.Errorf("unexpected profile: %+v", got)
	}
}

func TestRemove_ExistingProfile(t *testing.T) {
	r := profile.New()
	r.Add(profile.Profile{Name: "prod", Files: []string{".env.prod"}})
	if !r.Remove("prod") {
		t.Error("expected Remove to return true")
	}
	if _, ok := r.Get("prod"); ok {
		t.Error("expected profile to be removed")
	}
}

func TestRemove_NonExistent(t *testing.T) {
	r := profile.New()
	if r.Remove("ghost") {
		t.Error("expected Remove to return false for missing profile")
	}
}

func TestList_SortedByName(t *testing.T) {
	r := profile.New()
	r.Add(profile.Profile{Name: "z-env", Files: nil})
	r.Add(profile.Profile{Name: "a-env", Files: nil})
	r.Add(profile.Profile{Name: "m-env", Files: nil})
	list := r.List()
	if list[0].Name != "a-env" || list[1].Name != "m-env" || list[2].Name != "z-env" {
		t.Errorf("unexpected order: %v", list)
	}
}

func TestLoadFile_NotFound_ReturnsEmpty(t *testing.T) {
	r, err := profile.LoadFile("/nonexistent/profiles.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.List()) != 0 {
		t.Error("expected empty registry")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	r := profile.New()
	r.Add(profile.Profile{Name: "dev", Files: []string{".env", ".env.local"}})
	path := filepath.Join(t.TempDir(), "profiles.json")
	if err := r.SaveFile(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := profile.LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := loaded.Get("dev")
	if !ok || len(got.Files) != 2 {
		t.Errorf("round-trip failed: %+v", got)
	}
}
