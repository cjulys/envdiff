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
// It returns an error if the format is not supported or if writing fails.
func Write(w io.Writer, results []diff.Result, opts Options) error {
	if !IsSupported(opts.Format) {
		return fmt.Errorf("export: unsupported format %q", opts.Format)
	}
	switch opts.Format {
	case FormatJSON:
		return WriteJSON(w, results, opts.Verbose)
	case FormatCSV:
		return WriteCSV(w, results, opts.Verbose)
	}
	// unreachable, but keeps the compiler happy
	return nil
}

// IsSupported returns true when f is a recognised export format.
func IsSupported(f Format) bool {
	switch f {
	case FormatJSON, FormatCSV:
		return true
	}
	return false
}

// SupportedFormats returns a slice of all recognised export formats.
func SupportedFormats() []Format {
	return []Format{FormatJSON, FormatCSV}
}
