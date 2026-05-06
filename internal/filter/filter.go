package filter

import (
	"strings"

	"github.com/envdiff/internal/diff"
)

// Options controls which diff results are included after filtering.
type Options struct {
	OnlyMissing  bool
	OnlyMismatch bool
	Prefix       string
	Keys         []string
}

// Apply returns a filtered subset of results based on the provided Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	keySet := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[strings.ToUpper(k)] = struct{}{}
	}

	var out []diff.Result
	for _, r := range results {
		if opts.OnlyMissing && r.Status != diff.Missing {
			continue
		}
		if opts.OnlyMismatch && r.Status != diff.Mismatch {
			continue
		}
		if opts.Prefix != "" && !strings.HasPrefix(r.Key, strings.ToUpper(opts.Prefix)) {
			continue
		}
		if len(keySet) > 0 {
			if _, ok := keySet[r.Key]; !ok {
				continue
			}
		}
		out = append(out, r)
	}
	return out
}
