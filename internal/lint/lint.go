// Package lint provides heuristic checks on parsed .env file contents,
// flagging common mistakes such as whitespace in values, suspicious characters,
// and keys that shadow well-known environment variables.
package lint

import (
	"fmt"
	"strings"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warning Severity = "warning"
	Error   Severity = "error"
)

// Finding represents a single lint issue found in a .env file.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// wellKnownKeys are system-level keys that should not be overridden carelessly.
var wellKnownKeys = map[string]bool{
	"PATH": true, "HOME": true, "USER": true, "SHELL": true,
	"LANG": true, "PWD": true, "TERM": true,
}

// Check runs all lint rules against the provided key-value map and returns
// any findings. The map is typically produced by parser.ParseFile.
func Check(env map[string]string) []Finding {
	var findings []Finding

	for key, value := range env {
		if strings.TrimSpace(value) != value {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "value has leading or trailing whitespace",
				Severity: Warning,
			})
		}

		if strings.ContainsAny(value, "\t") {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "value contains a tab character",
				Severity: Warning,
			})
		}

		if value == "" {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "value is empty",
				Severity: Warning,
			})
		}

		if wellKnownKeys[strings.ToUpper(key)] {
			findings = append(findings, Finding{
				Key:      key,
				Message:  "key shadows a well-known system environment variable",
				Severity: Error,
			})
		}
	}

	return findings
}
