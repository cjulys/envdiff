package diff_test

import (
	"os"

	"github.com/yourusername/envdiff/internal/diff"
)

// ExampleAnnotate demonstrates how to generate and print annotations
// for a set of diff results.
func ExampleAnnotate() {
	results := []diff.Result{
		{Key: "DB_HOST", Status: diff.StatusMatch, ValueA: "localhost", ValueB: "localhost"},
		{Key: "DB_PASS", Status: diff.StatusMismatch, ValueA: "secret", ValueB: "hunter2"},
		{Key: "NEW_FLAG", Status: diff.StatusMissingInA},
	}

	annotations := diff.Annotate(results)
	diff.WriteAnnotations(os.Stdout, annotations, false)
	// Output:
	// [OK]     DB_HOST: values match across environments
	// [DIFF]   DB_PASS: value differs between environments
	// [MISS-A] NEW_FLAG: present in second file but missing in first
}
