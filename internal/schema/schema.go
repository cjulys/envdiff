// Package schema provides validation of .env files against a schema definition
// that declares required keys, optional keys, and expected value patterns.
package schema

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// KeyRule describes the constraints for a single key.
type KeyRule struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp // nil means any value is accepted
}

// Schema holds all rules loaded from a schema file.
type Schema struct {
	Rules []KeyRule
	ruleMap map[string]*KeyRule
}

// Violation represents a single schema violation found during Check.
type Violation struct {
	Key     string
	Message string
}

// LoadFile parses a schema definition file.
// Each non-blank, non-comment line has the form:
//
//	[!]KEY[=PATTERN]
//
// A leading '!' marks the key as required. PATTERN is a Go regexp.
func LoadFile(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("schema: open %s: %w", path, err)
	}
	defer f.Close()

	s := &Schema{ruleMap: make(map[string]*KeyRule)}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		rule, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("schema: parse %q: %w", line, err)
		}
		s.Rules = append(s.Rules, rule)
		s.ruleMap[rule.Key] = &s.Rules[len(s.Rules)-1]
	}
	return s, scanner.Err()
}

func parseLine(line string) (KeyRule, error) {
	rule := KeyRule{}
	if strings.HasPrefix(line, "!") {
		rule.Required = true
		line = line[1:]
	}
	parts := strings.SplitN(line, "=", 2)
	rule.Key = strings.TrimSpace(parts[0])
	if len(parts) == 2 && parts[1] != "" {
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return rule, err
		}
		rule.Pattern = re
	}
	return rule, nil
}

// Check validates env key/value pairs against the schema.
// It returns a list of violations (missing required keys, pattern mismatches).
func (s *Schema) Check(env map[string]string) []Violation {
	var violations []Violation
	for _, rule := range s.Rules {
		val, exists := env[rule.Key]
		if rule.Required && !exists {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			continue
		}
		if exists && rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			violations = append(violations, Violation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value %q does not match pattern %s", val, rule.Pattern),
			})
		}
	}
	return violations
}
