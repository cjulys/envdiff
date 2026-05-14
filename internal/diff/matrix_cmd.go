package diff

import (
	"fmt"
	"io"
)

// MatrixOptions controls the behaviour of RunMatrix.
type MatrixOptions struct {
	// Verbose includes per-cell problem details in output.
	Verbose bool
	// MinScore causes RunMatrix to return true (problems found) when any
	// cell's health score drops below this threshold (0–100).
	MinScore int
}

// DefaultMatrixOptions returns sensible defaults for RunMatrix.
func DefaultMatrixOptions() MatrixOptions {
	return MatrixOptions{
		Verbose:  false,
		MinScore: 80,
	}
}

// RunMatrix builds and writes a pairwise comparison matrix from the provided
// labelled result sets.  It returns true when at least one cell's score falls
// below opts.MinScore, signalling that the caller should exit non-zero.
func RunMatrix(w io.Writer, runs map[string][]Result, opts MatrixOptions) bool {
	if len(runs) == 0 {
		fmt.Fprintln(w, "no environments provided")
		return false
	}

	m := BuildMatrix(runs)
	WriteMatrix(w, m)

	if opts.Verbose {
		fmt.Fprintln(w, "")
		writeMatrixDetail(w, m)
	}

	for _, cell := range m.Cells {
		if cell.Score() < opts.MinScore {
			return true
		}
	}
	return false
}

// writeMatrixDetail prints per-cell problem summaries when verbose mode is on.
func writeMatrixDetail(w io.Writer, m Matrix) {
	fmt.Fprintln(w, "Cell details:")
	for _, row := range m.Labels {
		for _, col := range m.Labels {
			if row == col {
				continue
			}
			cell := m.Cells[cellKey(row, col)]
			fmt.Fprintf(w, "  %-10s → %-10s  matched=%d missing=%d mismatch=%d score=%d%%\n",
				cell.RowLabel, cell.ColLabel,
				cell.Matched, cell.Missing, cell.Mismatch,
				cell.Score(),
			)
		}
	}
}
