// Package lint provides heuristic analysis of parsed .env file contents.
//
// It surfaces common authoring mistakes such as:
//   - Empty values that may indicate an unset placeholder
//   - Values with leading or trailing whitespace that could cause subtle bugs
//   - Values containing tab characters
//   - Keys that shadow well-known system environment variables (e.g. PATH, HOME)
//
// Usage:
//
//	env, _ := parser.ParseFile("staging.env")
//	findings := lint.Check(env)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
package lint
