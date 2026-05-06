package export

import (
	"fmt"
	"io"

	"github.com/user/envdiff/internal/diff"
)

// Options controls export behaviour.
type Options struct {
	Format  Format
	Verbose bool
}

// Write dispatches to the correct exporter based on opts.Format.
func Write(w io.Writer, results []diff.Result, opts Options) error {
	switch opts.Format {
	case FormatJSON:
		return WriteJSON(w, results, opts.Verbose)
	case FormatCSV:
		return WriteCSV(w, results, opts.Verbose)
	default:
		return fmt.Errorf("export: unsupported format %q", opts.Format)
	}
}

// IsSupported returns true when f is a recognised export format.
func IsSupported(f Format) bool {
	switch f {
	case FormatJSON, FormatCSV:
		return true
	}
	return false
}
