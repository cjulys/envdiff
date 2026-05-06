// Package env provides variable interpolation for .env file values.
//
// A Resolver can expand $VAR and ${VAR} references found inside values using
// a caller-supplied map of key/value pairs and, optionally, the host process
// environment as a fallback.
//
// Variable references that cannot be resolved are replaced with an empty
// string. When the fallback flag is set to true, the host process environment
// (os.Getenv) is consulted after the supplied map before falling back to
// an empty string.
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
//
// Example with OS environment fallback:
//
//	r := env.NewResolver(map[string]string{
//		"HOST": "localhost",
//	}, true)
//
//	// If PORT is not in the map, os.Getenv("PORT") is tried next.
//	expanded := r.Expand("postgres://${HOST}:${PORT}/mydb")
package env
