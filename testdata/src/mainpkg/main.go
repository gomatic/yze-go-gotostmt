// Command mainpkg pins the PACKAGE dimension at "main": a command's code is
// judged like a library's. Nothing here is importable, which is the reasoning
// an exemption would be built on, and it has no bearing on whether the control
// flow is readable.
package main

// limit bounds the summation.
type limit int

func main() {
	i := limit(0)
loop:
	if i < 10 {
		i++
		goto loop // want `^goto is not permitted; replace a backward jump with a for statement, and a forward jump with a labelled break, a defer, or an early return from an extracted helper$`
	}
	_ = i
}
