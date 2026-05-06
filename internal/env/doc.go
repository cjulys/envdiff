// Package env provides variable interpolation for .env file values.
//
// A Resolver can expand $VAR and ${VAR} references found inside values using
// a caller-supplied map of key/value pairs and, optionally, the host process
// environment as a fallback.
//
// Example:
//
//	r := env.NewResolver(map[string]string{
//		"HOST": "localhost",
//		"PORT": "5432",
//	}, false)
//
//	expanded := r.Expand("postgres://${HOST}:${PORT}/mydb")
//	// expanded == "postgres://localhost:5432/mydb"
package env
