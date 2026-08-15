// Package deep sits under internal/, pinning that segment of the PATH
// dimension: an exemption keyed on "/internal/" loses this finding.
package deep

// Deep loops with a goto.
func Deep(limit int) int {
	i := 0
loop:
	if i < limit {
		i++
		goto loop // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	return i
}
