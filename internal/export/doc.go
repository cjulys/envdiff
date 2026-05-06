// Package export provides serialisation of diff results into machine-readable
// formats such as JSON and CSV.
//
// Usage:
//
//	results := diff.Compare(envA, envB)
//	err := export.Write(os.Stdout, results, export.Options{
//		Format:  export.FormatJSON,
//		Verbose: true,
//	})
//
// Supported formats are FormatJSON ("json") and FormatCSV ("csv").
// Use IsSupported to validate a format string before constructing Options.
package export
