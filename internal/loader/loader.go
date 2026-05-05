package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded .env file with its name and parsed key-value pairs.
type EnvFile struct {
	Name string
	Path string
	Vars map[string]string
}

// LoadFile reads and parses a single .env file from the given path.
func LoadFile(path string) (*EnvFile, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	vars, err := parser.ParseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	return &EnvFile{
		Name: filepath.Base(path),
		Path: path,
		Vars: vars,
	}, nil
}

// LoadFiles reads and parses multiple .env files, returning them in order.
// It collects all errors and returns a combined error if any files fail to load.
func LoadFiles(paths []string) ([]*EnvFile, error) {
	var files []*EnvFile
	var errs []string

	for _, p := range paths {
		f, err := LoadFile(p)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		files = append(files, f)
	}

	if len(errs) > 0 {
		return files, fmt.Errorf("errors loading files: %v", errs)
	}

	return files, nil
}
