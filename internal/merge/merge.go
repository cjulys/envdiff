// Package merge provides functionality to merge multiple .env file maps
// into a unified key set, useful for generating template or baseline files.
package merge

import "sort"

// Result holds the merged output of multiple env maps.
type Result struct {
	// Keys is the sorted union of all keys across all sources.
	Keys []string
	// Sources maps each key to the list of source file names that define it.
	Sources map[string][]string
	// Values maps each key to the value from the first source that defines it.
	Values map[string]string
}

// Merge combines multiple named env maps into a single Result.
// The names slice should correspond 1-to-1 with the maps slice.
func Merge(names []string, maps []map[string]string) Result {
	keySet := make(map[string]struct{})
	sources := make(map[string][]string)
	values := make(map[string]string)

	for i, m := range maps {
		name := ""
		if i < len(names) {
			name = names[i]
		}
		for k, v := range m {
			keySet[k] = struct{}{}
			sources[k] = append(sources[k], name)
			if _, seen := values[k]; !seen {
				values[k] = v
			}
		}
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return Result{
		Keys:    keys,
		Sources: sources,
		Values:  values,
	}
}

// UniqueKeys returns keys that appear in exactly one of the provided maps.
func UniqueKeys(names []string, maps []map[string]string) map[string]string {
	r := Merge(names, maps)
	unique := make(map[string]string)
	for _, k := range r.Keys {
		if len(r.Sources[k]) == 1 {
			unique[k] = r.Sources[k][0]
		}
	}
	return unique
}
