package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/user/envdiff/internal/diff"
)

// Format represents a supported export format.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// JSONRecord is the structure used when exporting results to JSON.
type JSONRecord struct {
	Key    string `json:"key"`
	Status string `json:"status"`
	ValueA string `json:"value_a,omitempty"`
	ValueB string `json:"value_b,omitempty"`
}

// WriteJSON writes diff results as a JSON array to w.
func WriteJSON(w io.Writer, results []diff.Result, verbose bool) error {
	records := toRecords(results, verbose)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

// WriteCSV writes diff results as CSV lines to w.
func WriteCSV(w io.Writer, results []diff.Result, verbose bool) error {
	if _, err := fmt.Fprintln(w, "key,status,value_a,value_b"); err != nil {
		return err
	}
	sorted := sortedResults(results)
	for _, r := range sorted {
		va, vb := "", ""
		if verbose {
			va = r.ValueA
			vb = r.ValueB
		}
		if _, err := fmt.Fprintf(w, "%s,%s,%s,%s\n", r.Key, r.Status, va, vb); err != nil {
			return err
		}
	}
	return nil
}

func toRecords(results []diff.Result, verbose bool) []JSONRecord {
	sorted := sortedResults(results)
	out := make([]JSONRecord, 0, len(sorted))
	for _, r := range sorted {
		rec := JSONRecord{Key: r.Key, Status: string(r.Status)}
		if verbose {
			rec.ValueA = r.ValueA
			rec.ValueB = r.ValueB
		}
		out = append(out, rec)
	}
	return out
}

func sortedResults(results []diff.Result) []diff.Result {
	copy := append([]diff.Result(nil), results...)
	sort.Slice(copy, func(i, j int) bool { return copy[i].Key < copy[j].Key })
	return copy
}
