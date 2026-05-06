package config

import (
	"errors"
	"strings"
)

// OutputFormat controls how results are rendered.
type OutputFormat string

const (
	FormatPretty   OutputFormat = "pretty"
	FormatMarkdown OutputFormat = "markdown"
	FormatJSON     OutputFormat = "json"
)

// Config holds all runtime options for an envdiff run.
type Config struct {
	// Files is the ordered list of .env file paths to compare.
	Files []string

	// Format is the output format (pretty, markdown, json).
	Format OutputFormat

	// Verbose includes actual key values in the output.
	Verbose bool

	// Color enables ANSI color codes in pretty output.
	Color bool

	// FilterStatus limits output to specific result statuses.
	// Accepted values: "match", "mismatch", "missing".
	FilterStatus []string

	// FilterPrefix restricts comparison to keys with this prefix.
	FilterPrefix string

	// FilterKeys restricts comparison to this explicit set of keys.
	FilterKeys []string

	// ExitOnDiff causes the process to exit with code 1 if any diff is found.
	ExitOnDiff bool
}

// Validate returns an error if the Config is not usable.
func (c *Config) Validate() error {
	if len(c.Files) < 2 {
		return errors.New("at least two files must be provided for comparison")
	}

	switch c.Format {
	case FormatPretty, FormatMarkdown, FormatJSON:
		// valid
	case "":
		c.Format = FormatPretty
	default:
		return errors.New("unknown format: " + string(c.Format) +
			"; expected one of: pretty, markdown, json")
	}

	for _, s := range c.FilterStatus {
		switch strings.ToLower(s) {
		case "match", "mismatch", "missing":
			// valid
		default:
			return errors.New("unknown filter status: " + s)
		}
	}

	return nil
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Format:     FormatPretty,
		Color:      true,
		ExitOnDiff: false,
	}
}
