package a_test

import "testing"

// TestExternalTestPackage pins the PACKAGE dimension at its other test value:
// package a_test, not package a. An exemption keyed on the "_test" package-name
// suffix rather than on the file name loses this finding.
func TestExternalTestPackage(t *testing.T) {
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
