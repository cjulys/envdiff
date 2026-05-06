// Package template generates .env.template files from one or more environment
// maps produced by the loader or merge packages.
//
// A template file contains all discovered keys with their values replaced by
// an empty string or a configurable placeholder, making it safe to commit to
// version control as documentation for required environment variables.
//
// Example usage:
//
//	envs, _ := loader.LoadFiles([]string{".env.production", ".env.staging"})
//	f, _ := os.Create(".env.template")
//	defer f.Close()
//	template.Generate(f, envs, template.DefaultOptions())
package template
