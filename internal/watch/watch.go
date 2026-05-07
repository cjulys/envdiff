// Package watch provides file system monitoring for .env files,
// triggering a callback when any watched file changes.
package watch

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// ChangeEvent describes a file that has changed on disk.
type ChangeEvent struct {
	Path    string
	OldHash string
	NewHash string
}

// Watcher polls a set of files for content changes.
type Watcher struct {
	files    []string
	hashes   map[string]string
	interval time.Duration
	stop     chan struct{}
}

// New creates a Watcher for the given files with the given poll interval.
func New(files []string, interval time.Duration) *Watcher {
	return &Watcher{
		files:    files,
		hashes:   make(map[string]string),
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins polling and sends ChangeEvents to the returned channel.
// Call Stop to terminate the watcher.
func (w *Watcher) Start() <-chan ChangeEvent {
	ch := make(chan ChangeEvent, len(w.files))

	// Seed initial hashes so the first poll does not fire false positives.
	for _, f := range w.files {
		if h, err := hashFile(f); err == nil {
			w.hashes[f] = h
		}
	}

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		defer close(ch)
		for {
			select {
			case <-w.stop:
				return
			case <-ticker.C:
				for _, f := range w.files {
					newHash, err := hashFile(f)
					if err != nil {
						continue
					}
					oldHash := w.hashes[f]
					if newHash != oldHash {
						w.hashes[f] = newHash
						ch <- ChangeEvent{Path: f, OldHash: oldHash, NewHash: newHash}
					}
				}
			}
		}
	}()

	return ch
}

// Stop terminates the polling goroutine.
func (w *Watcher) Stop() {
	close(w.stop)
}

// hashFile returns a hex SHA-256 digest of the named file's contents.
func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
