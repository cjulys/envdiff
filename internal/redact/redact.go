// Package redact provides utilities for masking sensitive values
// in .env file comparisons before displaying or exporting results.
package redact

import "strings"

// DefaultPatterns contains common key substrings that indicate sensitive data.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

const masked = "***REDACTED***"

// Redactor holds configuration for redacting sensitive values.
type Redactor struct {
	patterns []string
}

// New returns a Redactor using the provided patterns.
// If patterns is empty, DefaultPatterns are used.
func New(patterns []string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultPatterns
	}
	upper := make([]string, len(patterns))
	for i, p := range patterns {
		upper[i] = strings.ToUpper(p)
	}
	return &Redactor{patterns: upper}
}

// IsSensitive returns true if the key matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked constant if the key is sensitive,
// otherwise it returns the original value unchanged.
func (r *Redactor) MaskValue(key, value string) string {
	if r.IsSensitive(key) {
		return masked
	}
	return value
}

// MaskedValue is the string substituted for sensitive values.
func MaskedValue() string {
	return masked
}
