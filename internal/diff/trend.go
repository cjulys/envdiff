package diff

import (
	"fmt"
	"io"
	"sort"
	"time"
)

// TrendPoint represents the problem count at a specific point in time.
type TrendPoint struct {
	Label     string
	Timestamp time.Time
	Problems  int
	Total     int
}

// Trend holds an ordered sequence of TrendPoints.
type Trend struct {
	Points []TrendPoint
}

// AddRun appends a new data point derived from a labeled set of results.
func (t *Trend) AddRun(label string, ts time.Time, results []Result) {
	total := len(results)
	problems := 0
	for _, r := range results {
		if r.IsProblem() {
			problems++
		}
	}
	t.Points = append(t.Points, TrendPoint{
		Label:     label,
		Timestamp: ts,
		Problems:  problems,
		Total:     total,
	})
}

// Sorted returns a copy of the trend points sorted by timestamp ascending.
func (t *Trend) Sorted() []TrendPoint {
	cp := make([]TrendPoint, len(t.Points))
	copy(cp, t.Points)
	sort.Slice(cp, func(i, j int) bool {
		return cp[i].Timestamp.Before(cp[j].Timestamp)
	})
	return cp
}

// Direction returns +1 if problems increased, -1 if decreased, 0 if unchanged.
func (t *Trend) Direction() int {
	pts := t.Sorted()
	if len(pts) < 2 {
		return 0
	}
	first := pts[0].Problems
	last := pts[len(pts)-1].Problems
	switch {
	case last > first:
		return 1
	case last < first:
		return -1
	default:
		return 0
	}
}

// WriteTrend writes a simple ASCII trend table to w.
func WriteTrend(w io.Writer, t *Trend) {
	pts := t.Sorted()
	if len(pts) == 0 {
		fmt.Fprintln(w, "no trend data available")
		return
	}
	fmt.Fprintf(w, "%-20s  %-10s  %8s  %8s\n", "label", "date", "problems", "total")
	fmt.Fprintf(w, "%s\n", repeatChar('-', 54))
	for _, p := range pts {
		fmt.Fprintf(w, "%-20s  %-10s  %8d  %8d\n",
			p.Label,
			p.Timestamp.Format("2006-01-02"),
			p.Problems,
			p.Total,
		)
	}
	dir := t.Direction()
	arrow := "→"
	if dir > 0 {
		arrow = "↑ worsening"
	} else if dir < 0 {
		arrow = "↓ improving"
	}
	fmt.Fprintf(w, "\ntrend: %s\n", arrow)
}
