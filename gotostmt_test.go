package gotostmt_test

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"

	gotostmt "github.com/gomatic/yze-go-gotostmt"
)

// packages are every fixture package: the rule's own shapes in "a", and three
// packages that exist only to vary the FILE PATH — under internal/, named for a
// mock, and named _parser.go — so an exemption keyed on any of those loses a
// finding here rather than passing unseen.
var packages = []string{"a", "gen", "grammar", "internal/deep", "mainpkg", "mocks"}

// wantMessage is the diagnostic's contract, transcribed from the package doc
// comment's remedy list rather than from a run: a backward jump takes a for
// statement, a forward one takes a labelled break, a defer, or an early return.
// Asserting it here is what stops the printed instruction from being rewritten
// — into an instruction to disable the rule, say — while every gate stays green.
const wantMessage = "goto is not permitted; replace a backward jump with a for statement, " +
	"and a forward jump with a labelled break, a defer, or an early return from an extracted helper"

// wantSites is every goto in the fixtures, at the column of the GOTO KEYWORD
// rather than of its label, each read off the fixture source. The corpus
// Finding is {rule, path, line} and can see neither the message nor the column,
// so docs/s03.md assigns both to this test.
var wantSites = []string{
	"a/a.go:9:3",                 // backward jump, the loop shape
	"a/a.go:19:3",                // forward jump, the skip-ahead shape
	"a/a.go:31:3",                // first of two gotos in one function
	"a/a.go:36:3",                // second of two gotos in the same function
	"a/a_test.go:12:3",           // a _test.go file, package a
	"a/a_x_test.go:13:3",         // the same directory, package a_test
	"a/second.go:7:3",            // a second file of the same package
	"gen/hand_parser.go:13:3",    // a hand-written file named _parser.go
	"grammar/parse.go:16:3",      // the path segment carrying most of the real population
	"internal/deep/deep.go:11:3", // under internal/
	"mainpkg/main.go:15:3",       // a main package is judged like any other
	"mocks/mock_thing.go:11:3",   // named for a mock
}

func TestEveryGotoIsReportedAtItsKeywordWithTheWholeMessage(t *testing.T) {
	dir := analysistest.TestData()
	results := analysistest.Run(t, dir, gotostmt.Analyzer, packages...)
	require.NotEmpty(t, results)

	root := filepath.Join(dir, "src") + string(filepath.Separator)
	sites := map[string]struct{}{}
	for _, result := range results {
		require.NoError(t, result.Err)
		for _, diagnostic := range result.Diagnostics {
			assert.Equal(t, wantMessage, diagnostic.Message,
				"the diagnostic's whole message is the contract, not its first five words")

			// The file's identity is read from the FileSet's own name, never
			// from Position().Filename, which a //line directive in the judged
			// file can rewrite.
			file := result.Pass.Fset.File(diagnostic.Pos)
			require.NotNil(t, file)
			position := result.Pass.Fset.Position(diagnostic.Pos)
			name := filepath.ToSlash(strings.TrimPrefix(file.Name(), root))
			sites[fmt.Sprintf("%s:%d:%d", name, position.Line, position.Column)] = struct{}{}
		}
	}

	got := make([]string, 0, len(sites))
	for site := range sites {
		got = append(got, site)
	}
	sort.Strings(got)
	want := append([]string(nil), wantSites...)
	sort.Strings(want)
	assert.Equal(t, want, got,
		"every goto in every fixture package is reported, at the goto keyword")
}

func TestRegistrationIsWellFormed(t *testing.T) {
	assert.NoError(t, gotostmt.Registration.Validate())
	assert.Equal(t, "yze/gotostmt", gotostmt.Registration.RuleID())
	assert.Same(t, gotostmt.Analyzer, gotostmt.Registration.Analyzer)
}
