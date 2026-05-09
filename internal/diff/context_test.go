package diff

import (
	"bytes"
	"strings"
	"testing"
)

func buildContextResult(key, va, vb string, status Status) Result {
	return Result{Key: key, ValueA: va, ValueB: vb, Status: status}
}

func TestBuildContext_OnlyProblems(t *testing.T) {
	envA := map[string]string{"ALPHA": "1", "BETA": "2", "GAMMA": "3"}
	envB := map[string]string{"ALPHA": "1", "BETA": "9", "GAMMA": "3"}
	results := []Result{
		buildContextResult("ALPHA", "1", "1", StatusMatch),
		buildContextResult("BETA", "2", "9", StatusMismatch),
		buildContextResult("GAMMA", "3", "3", StatusMatch),
	}

	blocks := BuildContext(results, envA, envB, 1)
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	if blocks[0].Result.Key != "BETA" {
		t.Errorf("expected BETA, got %s", blocks[0].Result.Key)
	}
}

func TestBuildContext_BeforeAndAfterNeighbours(t *testing.T) {
	envA := map[string]string{"A": "1", "B": "2", "C": "3", "D": "4", "E": "5"}
	envB := map[string]string{"A": "1", "B": "2", "C": "99", "D": "4", "E": "5"}
	results := []Result{
		buildContextResult("C", "3", "99", StatusMismatch),
	}

	blocks := BuildContext(results, envA, envB, 2)
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block")
	}
	if len(blocks[0].Before) != 2 {
		t.Errorf("expected 2 before lines, got %d", len(blocks[0].Before))
	}
	if len(blocks[0].After) != 2 {
		t.Errorf("expected 2 after lines, got %d", len(blocks[0].After))
	}
	if blocks[0].Before[0].Key != "A" || blocks[0].Before[1].Key != "B" {
		t.Errorf("unexpected before keys: %+v", blocks[0].Before)
	}
	if blocks[0].After[0].Key != "D" || blocks[0].After[1].Key != "E" {
		t.Errorf("unexpected after keys: %+v", blocks[0].After)
	}
}

func TestBuildContext_NoProblems_ReturnsEmpty(t *testing.T) {
	envA := map[string]string{"X": "1"}
	envB := map[string]string{"X": "1"}
	results := []Result{
		buildContextResult("X", "1", "1", StatusMatch),
	}
	blocks := BuildContext(results, envA, envB, 1)
	if len(blocks) != 0 {
		t.Errorf("expected no blocks, got %d", len(blocks))
	}
}

func TestWriteContext_VerboseShowsValues(t *testing.T) {
	blocks := []ContextBlock{
		{
			Result: buildContextResult("SECRET", "abc", "xyz", StatusMismatch),
			Before: []ContextLine{{Key: "ALPHA", Value: "1"}},
			After:  []ContextLine{{Key: "ZETA", Value: "9"}},
		},
	}
	var buf bytes.Buffer
	WriteContext(&buf, blocks, true)
	out := buf.String()
	if !strings.Contains(out, "ALPHA=1") {
		t.Errorf("expected ALPHA=1 in verbose output, got: %s", out)
	}
	if !strings.Contains(out, "SECRET=abc") {
		t.Errorf("expected SECRET=abc in verbose output, got: %s", out)
	}
}

func TestWriteContext_NonVerboseHidesValues(t *testing.T) {
	blocks := []ContextBlock{
		{
			Result: buildContextResult("SECRET", "abc", "xyz", StatusMismatch),
			Before: []ContextLine{{Key: "ALPHA", Value: "1"}},
		},
	}
	var buf bytes.Buffer
	WriteContext(&buf, blocks, false)
	out := buf.String()
	if strings.Contains(out, "=abc") {
		t.Errorf("should not expose values in non-verbose mode, got: %s", out)
	}
	if !strings.Contains(out, "ALPHA") {
		t.Errorf("expected ALPHA key in output, got: %s", out)
	}
}

func TestWriteContext_SeparatorBetweenBlocks(t *testing.T) {
	blocks := []ContextBlock{
		{Result: buildContextResult("A", "1", "2", StatusMismatch)},
		{Result: buildContextResult("B", "", "3", StatusMissingInA)},
	}
	var buf bytes.Buffer
	WriteContext(&buf, blocks, false)
	out := buf.String()
	if !strings.Contains(out, "---") {
		t.Errorf("expected separator between blocks, got: %s", out)
	}
}
