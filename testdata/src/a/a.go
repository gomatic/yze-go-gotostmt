package a

// withGoto jumps BACKWARD to an earlier label: the loop shape.
func withGoto() int {
	i := 0
loop:
	if i < 10 {
		i++
		goto loop // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	return i
}

// withForwardGoto jumps FORWARD to a later label: the skip-ahead/cleanup shape,
// which no for statement expresses and which the backward fixture cannot pin.
func withForwardGoto(bad bool) string {
	out := "start"
	if bad {
		goto done // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	out = "work"
done:
	return out
}

// withTwoGotos holds two gotos in ONE function, so an implementation reporting
// at most once per function, file or pass loses a finding here.
func withTwoGotos(a, b bool) int {
	n := 0
	if a {
		goto one // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	n++
one:
	if b {
		goto two // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	n += 2
two:
	return n
}

// withBreak uses break (a non-goto branch statement) and must NOT be flagged.
func withBreak() {
	for {
		break
	}
}

// withContinue uses continue (a non-goto branch statement) and must NOT be flagged.
func withContinue() {
	for i := 0; i < 3; i++ {
		continue
	}
}

// withFallthrough uses fallthrough (a non-goto branch statement) and must NOT be flagged.
func withFallthrough(n int) int {
	switch n {
	case 0:
		fallthrough
	case 1:
		return 1
	}
	return 0
}

// withLabeledBreak uses a labeled break and must NOT be flagged: its target is
// fixed by the grammar to a boundary of the enclosing for, so the jump is
// bounded rather than arbitrary.
func withLabeledBreak() {
loop:
	for {
		break loop
	}
}

// withLabeledContinue uses a labeled continue and must NOT be flagged, for the
// same reason: continue may only name an ENCLOSING for.
func withLabeledContinue(n int) int {
	total := 0
outer:
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > i {
				continue outer
			}
			total++
		}
	}
	return total
}
