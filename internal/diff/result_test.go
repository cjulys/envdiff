package diff_test

import (
	"testing"

	"github.com/envdiff/internal/diff"
)

func TestResult_IsProblem(t *testing.T) {
	cases := []struct {
		name   string
		r      diff.Result
		want   bool
	}{
		{"match is not a problem", diff.Result{Status: diff.Match}, false},
		{"mismatch is a problem", diff.Result{Status: diff.Mismatch}, true},
		{"missing is a problem", diff.Result{Status: diff.Missing}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.r.IsProblem(); got != tc.want {
				t.Errorf("IsProblem() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResult_MissingInA(t *testing.T) {
	r := diff.Result{Status: diff.Missing, ValueA: "", ValueB: "something"}
	if !r.MissingInA() {
		t.Error("expected MissingInA() to be true")
	}
	r2 := diff.Result{Status: diff.Missing, ValueA: "something", ValueB: ""}
	if r2.MissingInA() {
		t.Error("expected MissingInA() to be false when ValueA is set")
	}
}

func TestResult_MissingInB(t *testing.T) {
	r := diff.Result{Status: diff.Missing, ValueA: "something", ValueB: ""}
	if !r.MissingInB() {
		t.Error("expected MissingInB() to be true")
	}
	r2 := diff.Result{Status: diff.Missing, ValueA: "", ValueB: "something"}
	if r2.MissingInB() {
		t.Error("expected MissingInB() to be false when ValueB is set")
	}
}
