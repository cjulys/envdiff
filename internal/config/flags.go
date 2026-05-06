package config

import (
	"flag"
	"strings"
)

// multiFlag is a flag.Value that accumulates repeated string flags.
type multiFlag []string

func (m *multiFlag) String() string  { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

// ParseFlags reads command-line flags into a Config and returns the
// remaining positional arguments (the file paths).
//
// Usage:
//
//	envdiff [flags] file1.env file2.env [...]
func ParseFlags(args []string) (*Config, []string, error) {
	cfg := Default()

	fs := flag.NewFlagSet("envdiff", flag.ContinueOnError)

	var (
		format     string
		filterKeys multiFlag
		filterStat multiFlag
	)

	fs.StringVar(&format, "format", "pretty", "Output format: pretty | markdown | json")
	fs.BoolVar(&cfg.Verbose, "verbose", false, "Show key values in output")
	fs.BoolVar(&cfg.Color, "color", true, "Enable ANSI colors (pretty format only)")
	fs.BoolVar(&cfg.ExitOnDiff, "exit-on-diff", false, "Exit with code 1 when differences are found")
	fs.StringVar(&cfg.FilterPrefix, "prefix", "", "Only compare keys with this prefix")
	fs.Var(&filterKeys, "key", "Restrict comparison to this key (repeatable)")
	fs.Var(&filterStat, "status", "Filter output by status: match|mismatch|missing (repeatable)")

	if err := fs.Parse(args); err != nil {
		return nil, nil, err
	}

	cfg.Format = OutputFormat(format)
	cfg.FilterKeys = []string(filterKeys)
	cfg.FilterStatus = []string(filterStat)
	cfg.Files = fs.Args()

	if err := cfg.Validate(); err != nil {
		return nil, nil, err
	}

	return cfg, cfg.Files, nil
}
