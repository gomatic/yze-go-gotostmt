package a

// inSecondFile is a goto in a SECOND file of the same package, so an
// implementation reporting once per pass loses this finding.
func inSecondFile(x bool) int {
	if x {
		goto out // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	return 1
out:
	return 2
}
