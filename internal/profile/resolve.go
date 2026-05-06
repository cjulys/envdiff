package profile

import "fmt"

// ResolveArgs takes a slice of arguments that may be either profile names or
// raw file paths, resolves any known profile names to their constituent files,
// and returns the flat list of file paths to compare.
//
// If a name matches a profile in the registry its files are expanded in-place.
// Unrecognised names are treated as literal file paths.
func ResolveArgs(r *Registry, args []string) ([]string, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("at least two files or profiles are required")
	}

	out := make([]string, 0, len(args))
	for _, arg := range args {
		if p, ok := r.Get(arg); ok {
			if len(p.Files) == 0 {
				return nil, fmt.Errorf("profile %q has no files defined", arg)
			}
			out = append(out, p.Files...)
		} else {
			out = append(out, arg)
		}
	}

	if len(out) < 2 {
		return nil, fmt.Errorf("resolved argument list must contain at least two file paths")
	}
	return out, nil
}

// MustGet returns the profile with the given name or panics.
// Intended for use in tests and CLI scaffolding where the name is known valid.
func MustGet(r *Registry, name string) Profile {
	p, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("profile %q not found", name))
	}
	return p
}
