package a

import "testing"

// TestGotoInATestFileIsReported pins the PATH dimension: a _test.go file is
// judged like any other, so a path-keyed exemption on "_test" loses a finding.
func TestGotoInATestFileIsReported(t *testing.T) {
	n := 0
retry:
	if n < 1 {
		n++
		goto retry // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	if n != 1 {
		t.Fatal("unreachable")
	}
}
