package diff

// Severity represents how serious a diff result is.
type Severity int

const (
	SeverityNone    Severity = iota // match
	SeverityLow                     // mismatch in non-sensitive key
	SeverityMedium                  // missing in one env
	SeverityHigh                    // missing in both or sensitive key mismatch
)

// String returns a human-readable label for the severity.
func (s Severity) String() string {
	switch s {
	case SeverityLow:
		return "low"
	case SeverityMedium:
		return "medium"
	case SeverityHigh:
		return "high"
	default:
		return "none"
	}
}

// sensitivePatterns are key substrings that elevate severity.
var sensitivePatterns = []string{
	"SECRET", "PASSWORD", "TOKEN", "KEY", "PASS", "CREDENTIAL",
}

// isSensitiveKey reports whether the key name suggests sensitive data.
func isSensitiveKey(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// Classify returns the Severity for a single Result.
func Classify(r Result) Severity {
	switch r.Status {
	case StatusMatch:
		return SeverityNone
	case StatusMissingInA, StatusMissingInB:
		if isSensitiveKey(r.Key) {
			return SeverityHigh
		}
		return SeverityMedium
	case StatusMismatch:
		if isSensitiveKey(r.Key) {
			return SeverityHigh
		}
		return SeverityLow
	default:
		return SeverityNone
	}
}

// ClassifyAll annotates each result with its severity and returns a map
// keyed by result key for quick lookup.
func ClassifyAll(results []Result) map[string]Severity {
	out := make(map[string]Severity, len(results))
	for _, r := range results {
		out[r.Key] = Classify(r)
	}
	return out
}
