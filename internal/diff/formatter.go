package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ColorCode represents an ANSI color escape code.
type ColorCode string

const (
	ColorReset  ColorCode = "\033[0m"
	ColorRed    ColorCode = "\033[31m"
	ColorGreen  ColorCode = "\033[32m"
	ColorYellow ColorCode = "\033[33m"
	ColorCyan   ColorCode = "\033[36m"
)

// FormatOptions controls how output is rendered.
type FormatOptions struct {
	Color  bool
	Verbose bool
}

// WritePretty writes a human-friendly, optionally colorized diff report to w.
func WritePretty(w io.Writer, results []Result, fileA, fileB string, opts FormatOptions) {
	colorize := func(c ColorCode, s string) string {
		if opts.Color {
			return string(c) + s + string(ColorReset)
		}
		return s
	}

	fmt.Fprintf(w, "Comparing %s <-> %s\n", colorize(ColorCyan, fileA), colorize(ColorCyan, fileB))
	fmt.Fprintln(w, strings.Repeat("-", 48))

	// Sort keys for deterministic output.
	sorted := make([]Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Key < sorted[j].Key
	})

	for _, r := range sorted {
		switch r.Status {
		case StatusMatch:
			if opts.Verbose {
				fmt.Fprintf(w, "  %s %s\n", colorize(ColorGreen, "="), r.Key)
			}
		case StatusMismatch:
			fmt.Fprintf(w, "  %s %s\n", colorize(ColorYellow, "~"), r.Key)
			if opts.Verbose {
				fmt.Fprintf(w, "      %s: %q\n", fileA, r.ValueA)
				fmt.Fprintf(w, "      %s: %q\n", fileB, r.ValueB)
			}
		case StatusMissingInB:
			fmt.Fprintf(w, "  %s %s\n", colorize(ColorRed, "-"), r.Key)
			if opts.Verbose {
				fmt.Fprintf(w, "      only in %s: %q\n", fileA, r.ValueA)
			}
		case StatusMissingInA:
			fmt.Fprintf(w, "  %s %s\n", colorize(ColorGreen, "+"), r.Key)
			if opts.Verbose {
				fmt.Fprintf(w, "      only in %s: %q\n", fileB, r.ValueB)
			}
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 48))
	summary := Summarize(results)
	fmt.Fprintf(w, "Summary: %d match, %d mismatch, %d missing\n",
		summary.Matching, summary.Mismatched, summary.MissingInA+summary.MissingInB)
}
