package diff

// Status represents the comparison outcome for a single key.
type Status string

const (
	// Match indicates both files have the key with identical values.
	Match Status = "match"
	// Mismatch indicates both files have the key but with different values.
	Mismatch Status = "mismatch"
	// Missing indicates the key is absent in one of the files.
	Missing Status = "missing"
)

// Result holds the comparison outcome for a single environment variable key.
type Result struct {
	Key     string
	Status  Status
	ValueA  string
	ValueB  string
	FileA   string
	FileB   string
}

// IsProblem returns true if the result represents a non-matching state.
func (r Result) IsProblem() bool {
	return r.Status != Match
}

// MissingInA returns true when the key exists only in file B.
func (r Result) MissingInA() bool {
	return r.Status == Missing && r.ValueA == ""
}

// MissingInB returns true when the key exists only in file A.
func (r Result) MissingInB() bool {
	return r.Status == Missing && r.ValueB == ""
}
