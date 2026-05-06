// Package schema loads and evaluates a schema definition file against
// parsed .env key/value maps.
//
// A schema file contains one rule per line:
//
//	# comment
//	!REQUIRED_KEY          – key must be present
//	OPTIONAL_KEY           – key may be absent
//	!PORT=^[0-9]+$         – required and must match the regexp
//	LOG_LEVEL=^(debug|info|warn|error)$
//
// Use LoadFile to parse a schema file, then Schema.Check to validate
// an env map and collect Violation entries.
package schema
