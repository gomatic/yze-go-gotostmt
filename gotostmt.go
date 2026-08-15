// Package gotostmt provides a go/analysis analyzer that forbids the goto
// statement, per the gomatic Go standard that control flow stays structured.
//
// The rule, stated so each clause can be broken separately:
//
//   - EVERY goto statement in the analyzed package is reported — not the first
//     per function, not the first per file, not the first per pass.
//   - BOTH of goto's directions are reported: the backward jump to an earlier
//     label, and the forward jump to a later one.
//   - NOTHING is exempt. The verdict reads one thing, the statement's token, so
//     no path, directory, filename or package name changes it: a goto in a
//     _test.go file, under internal/, or in a file named for a generator is
//     reported exactly like any other.
//   - break, continue and fallthrough are never reported, in their bare or
//     their labelled form.
//
// Why goto and not the other branch statements, which are jumps too and are
// the same *ast.BranchStmt: their targets are fixed by the grammar. A labelled
// break must name an enclosing for, switch or select and a labelled continue
// an enclosing for — the compiler rejects anything else with "invalid break
// label" / "invalid continue label" — so either can only transfer control to a
// boundary of a statement the reader is already inside, and the flow stays
// reducible. A goto names an arbitrary label anywhere in its function, in
// either direction and out of any nesting depth. That is the difference the
// rule is about, and it is why respelling a backward jump as a labelled
// continue, or a forward jump that already leaves an enclosing loop as a
// labelled break, is not an evasion of this rule but the restructuring it asks
// for.
//
// The boundary of that, stated rather than left to be found: a forward jump
// with NO enclosing loop can be silenced by wrapping the region in a
// `for { … break label }` that never iterates. That costs one wrapper, leaves
// the same jump, and changes nothing about the design, while the honest remedy
// — a defer, or an early return from an extracted helper — costs more. So for
// that one shape evading is cheaper than complying. This analyzer does not see
// it, deliberately: a loop that cannot iterate is a property of the loop and
// not of the branch statement, and a rule about fake loops is a different rule.
//
// The remedy depends on the direction, and both directions have one:
//
//   - a backward jump to an earlier label is a for statement;
//   - a forward jump that leaves an enclosing loop is a labelled break;
//   - a forward jump to a cleanup or exit label is a defer, or an early return
//     from an extracted helper.
//
// Scope. This analyzer declares no TestScope and applies no filter of its own,
// so the standalone binary in cmd/ — usable as a `go vet -vettool` — reports
// every goto it is shown, generated files included. Under the yze runner,
// go-yze drops every diagnostic in a file whose head carries the standard
// generated-file marker (go-yze's generated.go holds the pattern; it is not
// spelled out here, for the reason below), so in that configuration the rule is
// effectively scoped to hand-written code. That scoping is the FRAMEWORK's and
// not this rule's, and it is load-bearing: the fleet's entire goto population
// lives in generated parsers, where no remedy above can be taken at any price.
//
// The marker is described here and never quoted, because quoting it in this
// comment turns this file off. golangci-lint runs with `generated: lax`
// (.golangci.yaml), which classifies a file as generated when its leading
// comments merely CONTAIN the marker's words — so a doc comment that quotes
// the string, as this one did, drops every golangci-lint finding in the file
// while go-yze, whose pattern is anchored to a whole line, keeps reporting.
// Measured: with the string quoted a planted `unused` violation in this file
// gives "0 issues"; with the same string reworded it is reported. Half the
// gate goes silent, nothing says so, and the two halves disagree about what
// this file even is.
package gotostmt

import (
	"go/ast"
	"go/token"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// message names a remedy for each of goto's two directions, because the two
// take different ones and the analyzer does not know which it is looking at.
const message = "goto is not permitted; replace a backward jump with a for statement, " +
	"and a forward jump with a labelled break, a defer, or an early return from an extracted helper"

// Analyzer reports every use of the goto statement.
var Analyzer = &analysis.Analyzer{
	Name:     "gotostmt",
	Doc:      "reports use of the goto statement, which the gomatic Go standard forbids",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "gotostmt",
	Categories: []goyze.Category{"patterns"},
	URL:        "https://docs.gomatic.dev/yze/gotostmt",
	Analyzer:   Analyzer,
}

// run reports every goto statement in the analyzed package.
func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder([]ast.Node{(*ast.BranchStmt)(nil)}, func(n ast.Node) {
		if stmt := n.(*ast.BranchStmt); stmt.Tok == token.GOTO {
			pass.Reportf(stmt.Pos(), message)
		}
	})
	return nil, nil
}
