package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Annotation holds a human-readable note attached to a Result.
type Annotation struct {
	Key     string
	Status  Status
	Note    string
}

// Annotate produces a list of Annotations from a slice of Results.
// Each annotation provides a plain-English description of the diff outcome
// for that key, suitable for display in reports or CI output.
func Annotate(results []Result) []Annotation {
	annotations := make([]Annotation, 0, len(results))
	for _, r := range results {
		annotations = append(annotations, Annotation{
			Key:    r.Key,
			Status: r.Status,
			Note:   annotationNote(r),
		})
	}
	sort.Slice(annotations, func(i, j int) bool {
		return annotations[i].Key < annotations[j].Key
	})
	return annotations
}

func annotationNote(r Result) string {
	switch r.Status {
	case StatusMatch:
		return fmt.Sprintf("%s: values match across environments", r.Key)
	case StatusMismatch:
		return fmt.Sprintf("%s: value differs between environments", r.Key)
	case StatusMissingInA:
		return fmt.Sprintf("%s: present in second file but missing in first", r.Key)
	case StatusMissingInB:
		return fmt.Sprintf("%s: present in first file but missing in second", r.Key)
	default:
		return fmt.Sprintf("%s: unknown status", r.Key)
	}
}

// WriteAnnotations writes all annotations to w, one per line.
func WriteAnnotations(w io.Writer, annotations []Annotation, onlyProblems bool) {
	for _, a := range annotations {
		if onlyProblems && a.Status == StatusMatch {
			continue
		}
		fmt.Fprintln(w, formatAnnotation(a))
	}
}

func formatAnnotation(a Annotation) string {
	var prefix string
	switch a.Status {
	case StatusMatch:
		prefix = "[OK]     "
	case StatusMismatch:
		prefix = "[DIFF]   "
	case StatusMissingInA:
		prefix = "[MISS-A] "
	case StatusMissingInB:
		prefix = "[MISS-B] "
	default:
		prefix = "[?]      "
	}
	return prefix + strings.TrimSpace(a.Note)
}
