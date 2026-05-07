// Package watch implements lightweight polling-based file watching for envdiff.
//
// It periodically computes a SHA-256 hash of each registered file and emits a
// ChangeEvent whenever the hash differs from the previously recorded value.
// This approach avoids OS-specific inotify / kqueue dependencies while remaining
// accurate enough for interactive development workflows.
//
// Typical usage:
//
//	w := watch.New([]string{".env", ".env.production"}, time.Second)
//	changes := w.Start()
//	for evt := range changes {
//		fmt.Printf("%s changed\n", evt.Path)
//	}
//	w.Stop()
package watch
