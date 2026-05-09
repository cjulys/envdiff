package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildAnnotationResult(key string, status Status) Result {
	r := Result{Key: key, Status: status}
	if status == StatusMismatch {
		r.ValueA = "aaa"
		r.ValueB = "bbb"
	}
	return r
}

func TestAnnotate_ReturnsOnePerResult(t *testing.T) {
	results := []Result{
		buildAnnotationResult("FOO", StatusMatch),
		buildAnnotationResult("BAR", StatusMismatch),
	}
	anns := Annotate(results)
	if len(anns) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(anns))
	}
}

func TestAnnotate_SortedByKey(t *testing.T) {
	results := []Result{
		buildAnnotationResult("ZEBRA", StatusMatch),
		buildAnnotationResult("ALPHA", StatusMissingInB),
	}
	anns := Annotate(results)
	if anns[0].Key != "ALPHA" || anns[1].Key != "ZEBRA" {
		t.Errorf("expected sorted order, got %s, %s", anns[0].Key, anns[1].Key)
	}
}

func TestAnnotate_MatchNote(t *testing.T) {
	anns := Annotate([]Result{buildAnnotationResult("X", StatusMatch)})
	if !strings.Contains(anns[0].Note, "match") {
		t.Errorf("expected match note, got: %s", anns[0].Note)
	}
}

func TestAnnotate_MismatchNote(t *testing.T) {
	anns := Annotate([]Result{buildAnnotationResult("X", StatusMismatch)})
	if !strings.Contains(anns[0].Note, "differs") {
		t.Errorf("expected differs note, got: %s", anns[0].Note)
	}
}

func TestAnnotate_MissingInANote(t *testing.T) {
	anns := Annotate([]Result{buildAnnotationResult("X", StatusMissingInA)})
	if !strings.Contains(anns[0].Note, "missing in first") {
		t.Errorf("expected missing-in-first note, got: %s", anns[0].Note)
	}
}

func TestAnnotate_MissingInBNote(t *testing.T) {
	anns := Annotate([]Result{buildAnnotationResult("X", StatusMissingInB)})
	if !strings.Contains(anns[0].Note, "missing in second") {
		t.Errorf("expected missing-in-second note, got: %s", anns[0].Note)
	}
}

func TestWriteAnnotations_OnlyProblems_HidesMatch(t *testing.T) {
	results := []Result{
		buildAnnotationResult("OK_KEY", StatusMatch),
		buildAnnotationResult("BAD_KEY", StatusMismatch),
	}
	var buf bytes.Buffer
	WriteAnnotations(&buf, Annotate(results), true)
	out := buf.String()
	if strings.Contains(out, "OK_KEY") {
		t.Error("expected OK_KEY to be hidden when onlyProblems=true")
	}
	if !strings.Contains(out, "BAD_KEY") {
		t.Error("expected BAD_KEY to appear in output")
	}
}

func TestWriteAnnotations_AllShownWhenNotFiltered(t *testing.T) {
	results := []Result{
		buildAnnotationResult("OK_KEY", StatusMatch),
		buildAnnotationResult("BAD_KEY", StatusMissingInB),
	}
	var buf bytes.Buffer
	WriteAnnotations(&buf, Annotate(results), false)
	out := buf.String()
	if !strings.Contains(out, "OK_KEY") {
		t.Error("expected OK_KEY to appear when onlyProblems=false")
	}
	if !strings.Contains(out, "BAD_KEY") {
		t.Error("expected BAD_KEY to appear in output")
	}
}

func TestWriteAnnotations_PrefixFormat(t *testing.T) {
	results := []Result{buildAnnotationResult("K", StatusMismatch)}
	var buf bytes.Buffer
	WriteAnnotations(&buf, Annotate(results), false)
	if !strings.HasPrefix(buf.String(), "[DIFF]") {
		t.Errorf("expected [DIFF] prefix, got: %s", buf.String())
	}
}
