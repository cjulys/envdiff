package diff

import (
	"fmt"
	"io"
)

// RadarOptions controls the behaviour of RunRadar.
type RadarOptions struct {
	// Out is the writer for the radar report. Defaults to os.Stdout when nil.
	Out io.Writer
}

// DefaultRadarOptions returns sensible defaults for RunRadar.
func DefaultRadarOptions() RadarOptions {
	return RadarOptions{}
}

// RunRadar builds and writes a radar report for the supplied named environments.
// It returns true when at least one environment has problems, allowing callers
// to propagate a non-zero exit code.
func RunRadar(envs map[string][]Result, opts RadarOptions) bool {
	report := BuildRadar(envs)

	out := opts.Out
	if out == nil {
		panic("RadarOptions.Out must not be nil")
	}

	WriteRadar(out, report)

	for _, e := range report.Entries {
		if e.Problems > 0 {
			fmt.Fprintf(out, "\n%d environment(s) have differences.\n", problemEnvCount(report))
			return true
		}
	}
	fmt.Fprintln(out, "\nAll environments are in sync.")
	return false
}

func problemEnvCount(r RadarReport) int {
	count := 0
	for _, e := range r.Entries {
		if e.Problems > 0 {
			count++
		}
	}
	return count
}
