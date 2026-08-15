// Package gen holds a HAND-WRITTEN file named _parser.go and carrying no
// generated marker, pinning the last segment of the PATH dimension: an
// exemption keyed on the _parser.go suffix — the one that takes the fleet's
// whole population to zero — loses this finding.
package gen

// ParserLoop loops with a goto.
func ParserLoop(limit int) int {
	i := 0
loop:
	if i < limit {
		i++
		goto loop // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	return i
}
