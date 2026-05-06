// Package validate provides checks for .env file key conventions,
// such as naming rules and duplicate key detection.
package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// validKeyPattern matches conventional env var names: uppercase letters, digits, underscores.
var validKeyPattern = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Issue represents a single validation problem found in an env map.
type Issue struct {
	Key     string
	Message string
}

func (i Issue) String() string {
	return fmt.Sprintf("%s: %s", i.Key, i.Message)
}

// CheckKeys validates the keys of an env map and returns any issues found.
// It checks for:
//   - Non-conventional key names (must be UPPER_SNAKE_CASE)
//   - Empty keys
//   - Keys containing whitespace
func CheckKeys(env map[string]string) []Issue {
	var issues []Issue

	for k := range env {
		if k == "" {
			issues = append(issues, Issue{Key: k, Message: "empty key"})
			continue
		}
		if strings.ContainsAny(k, " \t") {
			issues = append(issues, Issue{Key: k, Message: "key contains whitespace"})
			continue
		}
		if !validKeyPattern.MatchString(k) {
			issues = append(issues, Issue{
				Key:     k,
				Message: fmt.Sprintf("key %q does not match UPPER_SNAKE_CASE convention", k),
			})
		}
	}

	return issues
}

// CheckDuplicates detects duplicate keys in a raw list of key=value lines.
// The standard env map cannot hold duplicates, so this operates on raw lines.
func CheckDuplicates(lines []string) []Issue {
	seen := make(map[string]int)
	var issues []Issue

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		seen[key]++
		if seen[key] == 2 {
			issues = append(issues, Issue{Key: key, Message: "duplicate key detected"})
		}
	}

	return issues
}
