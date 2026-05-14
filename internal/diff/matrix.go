package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// MatrixCell holds the comparison result between two environment files.
type MatrixCell struct {
	RowLabel string
	ColLabel string
	Matched  int
	Missing  int
	Mismatch int
	Total    int
}

// Score returns a 0–100 health score for the cell.
func (c MatrixCell) Score() int {
	if c.Total == 0 {
		return 100
	}
	return (c.Matched * 100) / c.Total
}

// Matrix is a pairwise comparison table across N environments.
type Matrix struct {
	Labels []string
	Cells  map[string]MatrixCell // key: "rowLabel:colLabel"
}

// cellKey builds a canonical map key.
func cellKey(row, col string) string {
	return row + ":" + col
}

// BuildMatrix computes a pairwise diff matrix from a map of label→results.
// Each entry in runs is the Compare output for that labelled environment
// against a shared reference set of keys.
func BuildMatrix(runs map[string][]Result) Matrix {
	labels := make([]string, 0, len(runs))
	for l := range runs {
		labels = append(labels, l)
	}
	sort.Strings(labels)

	cells := make(map[string]MatrixCell)

	for _, row := range labels {
		for _, col := range labels {
			if row == col {
				continue
			}
			// Build value maps from each run's results.
			rowMap := resultsToMap(runs[row])
			colMap := resultsToMap(runs[col])
			results := Compare(rowMap, colMap)

			cell := MatrixCell{RowLabel: row, ColLabel: col, Total: len(results)}
			for _, r := range results {
				switch r.Status {
				case StatusMatch:
					cell.Matched++
				case StatusMismatch:
					cell.Mismatch++
				case StatusMissingInA, StatusMissingInB:
					cell.Missing++
				}
			}
			cells[cellKey(row, col)] = cell
		}
	}

	return Matrix{Labels: labels, Cells: cells}
}

// resultsToMap converts a []Result into the key→value map expected by Compare.
func resultsToMap(results []Result) map[string]string {
	m := make(map[string]string, len(results))
	for _, r := range results {
		if r.ValueA != "" {
			m[r.Key] = r.ValueA
		} else {
			m[r.Key] = r.ValueB
		}
	}
	return m
}

// WriteMatrix renders the matrix as an ASCII table to w.
func WriteMatrix(w io.Writer, m Matrix) {
	if len(m.Labels) == 0 {
		fmt.Fprintln(w, "no environments to compare")
		return
	}

	colWidth := 8
	header := fmt.Sprintf("%-12s", "")
	for _, col := range m.Labels {
		header += fmt.Sprintf(" %-*s", colWidth, truncate(col, colWidth))
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, strings.Repeat("-", len(header)))

	for _, row := range m.Labels {
		line := fmt.Sprintf("%-12s", truncate(row, 12))
		for _, col := range m.Labels {
			if row == col {
				line += fmt.Sprintf(" %-*s", colWidth, "  --  ")
				continue
			}
			cell := m.Cells[cellKey(row, col)]
			line += fmt.Sprintf(" %-*s", colWidth, fmt.Sprintf("%d%%", cell.Score()))
		}
		fmt.Fprintln(w, line)
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
