package diff

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// PivotTable renders a matrix of keys vs environment file labels,
// showing each key's value (or status) per environment.
type PivotTable struct {
	Labels  []string
	Rows    []PivotRow
}

// PivotRow holds a single key and its per-environment values.
type PivotRow struct {
	Key    string
	Cells  map[string]string // label -> display value
	Status string            // "match", "mismatch", "missing"
}

// BuildPivot constructs a PivotTable from a slice of Results and their source labels.
// labels must correspond to the Sources used in each Result.
func BuildPivot(results []Result, labels []string, verbose bool) PivotTable {
	type keyLabel struct {
		key   string
		label string
		val   string
	}

	cellMap := map[string]map[string]string{}
	statusMap := map[string]string{}

	for _, r := range results {
		key := r.Key
		if _, ok := cellMap[key]; !ok {
			cellMap[key] = map[string]string{}
		}

		displayA := "(missing)"
		displayB := "(missing)"

		if r.Status != StatusMissingInA {
			if verbose {
				displayA = r.ValueA
			} else {
				displayA = "(set)"
			}
		}
		if r.Status != StatusMissingInB {
			if verbose {
				displayB = r.ValueB
			} else {
				displayB = "(set)"
			}
		}

		if len(labels) >= 1 {
			cellMap[key][labels[0]] = displayA
		}
		if len(labels) >= 2 {
			cellMap[key][labels[1]] = displayB
		}

		statusMap[key] = string(r.Status)
	}

	keys := make([]string, 0, len(cellMap))
	for k := range cellMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	rows := make([]PivotRow, 0, len(keys))
	for _, k := range keys {
		rows = append(rows, PivotRow{
			Key:    k,
			Cells:  cellMap[k],
			Status: statusMap[k],
		})
	}

	return PivotTable{Labels: labels, Rows: rows}
}

// Write renders the PivotTable as a tab-separated table to w.
func (p PivotTable) Write(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// Header row
	fmt.Fprintf(tw, "KEY\t")
	for _, l := range p.Labels {
		fmt.Fprintf(tw, "%s\t", l)
	}
	fmt.Fprintf(tw, "STATUS\n")

	for _, row := range p.Rows {
		fmt.Fprintf(tw, "%s\t", row.Key)
		for _, l := range p.Labels {
			v := row.Cells[l]
			if v == "" {
				v = "(missing)"
			}
			fmt.Fprintf(tw, "%s\t", v)
		}
		fmt.Fprintf(tw, "%s\n", row.Status)
	}

	tw.Flush()
}
