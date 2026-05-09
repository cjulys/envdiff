package diff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// ContextLine represents a single line of context output around a diff result.
type ContextLine struct {
	Key    string
	Value  string
	Source string // file label or path
}

// ContextBlock groups a Result with surrounding key context from the same env map.
type ContextBlock struct {
	Result  Result
	Before  []ContextLine
	After   []ContextLine
}

// BuildContext returns ContextBlocks for all problem results, each annotated
// with up to n neighbouring keys (alphabetically) from the provided env maps.
func BuildContext(results []Result, envA, envB map[string]string, n int) []ContextBlock {
	keys := sortedKeys(envA, envB)
	index := buildKeyIndex(keys)

	var blocks []ContextBlock
	for _, r := range results {
		if !r.IsProblem() {
			continue
		}
		blocks = append(blocks, ContextBlock{
			Result: r,
			Before: neighbours(keys, index, r.Key, envA, envB, -n, 0),
			After:  neighbours(keys, index, r.Key, envA, envB, 1, n+1),
		})
	}
	return blocks
}

// WriteContext writes context blocks to w in a human-readable format.
func WriteContext(w io.Writer, blocks []ContextBlock, verbose bool) {
	for i, b := range blocks {
		if i > 0 {
			fmt.Fprintln(w, strings.Repeat("-", 40))
		}
		for _, cl := range b.Before {
			writeContextLine(w, " ", cl, verbose)
		}
		writeContextLine(w, ">", ContextLine{Key: b.Result.Key, Value: b.Result.ValueA}, verbose)
		for _, cl := range b.After {
			writeContextLine(w, " ", cl, verbose)
		}
	}
}

func writeContextLine(w io.Writer, prefix string, cl ContextLine, verbose bool) {
	if verbose {
		fmt.Fprintf(w, "%s %s=%s\n", prefix, cl.Key, cl.Value)
	} else {
		fmt.Fprintf(w, "%s %s\n", prefix, cl.Key)
	}
}

func sortedKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func buildKeyIndex(keys []string) map[string]int {
	m := make(map[string]int, len(keys))
	for i, k := range keys {
		m[k] = i
	}
	return m
}

func neighbours(keys []string, index map[string]int, key string, envA, envB map[string]string, from, to int) []ContextLine {
	idx, ok := index[key]
	if !ok {
		return nil
	}
	var lines []ContextLine
	start := idx + from
	end := idx + to
	if start < 0 {
		start = 0
	}
	if end > len(keys) {
		end = len(keys)
	}
	for i := start; i < end; i++ {
		k := keys[i]
		v := envA[k]
		if v == "" {
			v = envB[k]
		}
		lines = append(lines, ContextLine{Key: k, Value: v})
	}
	return lines
}
