package diff

import (
	"fmt"
	"io"
	"strings"
)

// WriteScore writes a human-readable score card to w.
func WriteScore(w io.Writer, s Score, label string) {
	if label == "" {
		label = "overall"
	}

	bar := buildBar(s.Value, 20)

	fmt.Fprintf(w, "Env Health Score (%s)\n", label)
	fmt.Fprintf(w, "  Score : %.1f / 100  [%s]  Grade: %s\n", s.Value, bar, s.Grade)
	fmt.Fprintf(w, "  Keys  : %d total, %d ok, %d problem(s)\n",
		s.Total, s.Total-s.Problems, s.Problems)
}

// buildBar returns a simple ASCII progress bar of width w filled proportionally to pct.
func buildBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(pct / 100.0 * float64(width))
	empty := width - filled
	return strings.Repeat("█", filled) + strings.Repeat("░", empty)
}
