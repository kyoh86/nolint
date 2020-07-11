package nolint

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:             "nolint",
	Doc:              "make the go/analysis checkers be able to ignore diagnostics with \"Nolint\" comment",
	Run:              run,
	RunDespiteErrors: true,
	ResultType:       reflect.TypeOf(new(NoLinter)),
	// FactTypes []Fact
}

func run(pass *analysis.Pass) (interface{}, error) {
	noLinter := &NoLinter{
		fset:  pass.Fset,
		marks: map[mark]bool{},
	}
	for _, f := range pass.Files {
		for _, c := range f.Comments {
			noLinter.regMarks(c)
		}
	}
	return noLinter, nil
}

type NoLinter struct {
	fset  *token.FileSet
	marks map[mark]bool
}

func (n *NoLinter) IgnoreDiagnostic(diagnostic analysis.Diagnostic) bool {
	for _, category := range []string{diagnostic.Category, ""} {
		if n.marks[n.genMark(n.fset, diagnostic.Pos, category)] {
			return true
		}
	}
	return false
}

func (n *NoLinter) IgnoreNode(r analysis.Range, category string) bool {
	for _, category := range []string{category, ""} {
		if n.marks[n.genMark(n.fset, r.Pos(), category)] {
			return true
		}
	}
	return false
}

type mark struct {
	Filename string
	Line     int
	Category string
}

func (n *NoLinter) genMark(fset *token.FileSet, pos token.Pos, category string) mark {
	p := fset.Position(pos)
	return mark{
		Filename: p.Filename,
		Line:     p.Line,
		Category: category,
	}
}

func (n *NoLinter) regMark(r analysis.Range, category string) {
	category = strings.TrimSpace(category)
	n.marks[n.genMark(n.fset, r.Pos(), category)] = true
}

func (n *NoLinter) regMarks(comment *ast.CommentGroup) {
	text := strings.TrimSpace(comment.Text())
	if text == "nolint" {
		n.regMark(comment, "")
	}
	if strings.HasPrefix(text, "nolint:") {
		for _, category := range strings.Split(strings.TrimPrefix(text, "nolint:"), ",") {
			n.regMark(comment, category)
		}
	}
}
