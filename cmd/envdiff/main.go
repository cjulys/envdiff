package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/loader"
)

func main() {
	quiet := flag.Bool("quiet", false, "suppress detailed output, only print summary")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: envdiff [flags] <file1.env> <file2.env>\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	files, err := loader.LoadFiles(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if len(files) != 2 {
		fmt.Fprintf(os.Stderr, "error: could not load both files\n")
		os.Exit(1)
	}

	results := diff.Compare(files[0].Vars, files[1].Vars)
	summary := diff.Summarize(results)

	if !*quiet {
		if err := diff.WriteReport(os.Stdout, files[0].Name, files[1].Name, results); err != nil {
			fmt.Fprintf(os.Stderr, "error writing report: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\nSummary: %d match, %d mismatch, %d missing in %s, %d missing in %s\n",
		summary.Matching,
		summary.Mismatched,
		summary.MissingInB,
		files[1].Name,
		summary.MissingInA,
		files[0].Name,
	)

	if summary.Mismatched > 0 || summary.MissingInA > 0 || summary.MissingInB > 0 {
		os.Exit(1)
	}
}
