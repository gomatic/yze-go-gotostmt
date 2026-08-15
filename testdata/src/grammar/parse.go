// Package grammar sits on a path containing "grammar", which is where 1552 of
// the fleet's 1802 gotos live. The four other path fixtures pin a last segment
// or one directory name; this one pins the segment that carries most of the
// real population, and it is hand-written with no generated marker.
package grammar

// limit bounds the summation.
type limit int

// Parse loops with a goto.
func Parse(bound limit) limit {
	i := limit(0)
loop:
	if i < bound {
		i++
		goto loop // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	return i
}
