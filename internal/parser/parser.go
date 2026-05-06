package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents the key-value pairs parsed from a .env file.
type EnvMap map[string]string

// ParseFile reads a .env file and returns a map of key-value pairs.
// It skips blank lines and comment lines (starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file %q: %w", path, err)
	}
	defer f.Close()

	env := make(EnvMap)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, lineNum, err)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading env file %q: %w", path, err)
	}

	return env, nil
}

// parseLine splits a single KEY=VALUE line into its components.
// Inline comments (# ...) after the value are stripped.
// Surrounding quotes on the value are removed.
func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid line %q: missing '='" , line)
	}

	key := strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", fmt.Errorf("invalid line %q: empty key", line)
	}

	// Keys must not contain spaces
	if strings.ContainsAny(key, " \t") {
		return "", "", fmt.Errorf("invalid line %q: key contains whitespace", line)
	}

	rawValue := strings.TrimSpace(parts[1])

	// Strip inline comment
	if idx := strings.Index(rawValue, " #"); idx != -1 {
		rawValue = strings.TrimSpace(rawValue[:idx])
	}

	// Unquote value if wrapped in matching quotes
	value := unquote(rawValue)

	return key, value, nil
}

// unquote removes surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
