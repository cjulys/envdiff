package watch_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/envdiff/internal/watch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestWatcher_DetectsChange(t *testing.T) {
	path := writeTempEnv(t, "KEY=original\n")

	w := watch.New([]string{path}, 20*time.Millisecond)
	ch := w.Start()
	defer w.Stop()

	// Overwrite the file to trigger a change.
	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(path, []byte("KEY=changed\n"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	select {
	case evt := <-ch:
		if evt.Path != path {
			t.Errorf("path = %q, want %q", evt.Path, path)
		}
		if evt.OldHash == evt.NewHash {
			t.Error("expected hashes to differ after file change")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	path := writeTempEnv(t, "KEY=stable\n")

	w := watch.New([]string{path}, 20*time.Millisecond)
	ch := w.Start()
	defer w.Stop()

	select {
	case evt := <-ch:
		t.Errorf("unexpected change event for unchanged file: %+v", evt)
	case <-time.After(120 * time.Millisecond):
		// expected — no change occurred
	}
}

func TestWatcher_StopClosesChannel(t *testing.T) {
	path := writeTempEnv(t, "A=1\n")

	w := watch.New([]string{path}, 20*time.Millisecond)
	ch := w.Start()
	w.Stop()

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed after Stop")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("channel was not closed after Stop")
	}
}

func TestWatcher_MissingFileDoesNotPanic(t *testing.T) {
	w := watch.New([]string{"/nonexistent/.env"}, 20*time.Millisecond)
	ch := w.Start()
	defer w.Stop()

	select {
	case <-ch:
		t.Error("unexpected event for nonexistent file")
	case <-time.After(100 * time.Millisecond):
		// expected — missing file is silently skipped
	}
}
