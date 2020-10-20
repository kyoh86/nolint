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
		noLinter.file = f
		noLinter.commentMap = nil
		ast.Walk(noLinter, f)
	}
	return noLinter, nil
}

type NoLinter struct {
	fset       *token.FileSet
	file       *ast.File
	commentMap ast.CommentMap
	marks      map[mark]bool
}

func (n *NoLinter) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return n
	}

	// Ensure we initialized the comment map if we don't have one.
	if n.commentMap == nil {
		n.commentMap = ast.NewCommentMap(n.fset, node, n.file.Comments)
	}

	// No need to match comments with comments.
	switch node.(type) {
	case *ast.Comment, *ast.CommentGroup:
		return n
	}

	// Update the comment map to target the current node we're visiting.
	n.commentMap.Update(nil, node)

	// Get the comments related to this specific node. If there are no comments
	// there can't be any nolint directive.
	commentsForNode, ok := n.commentMap[node]
	if !ok {
		return n
	}

	// If we would have multiple comment groups for the node we only care for
	// the last one since we want the directive to be added right above our
	// node.
	lastCommentForNode := commentsForNode[len(commentsForNode)-1]
	n.regMarks(lastCommentForNode, node)

	return n
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

func (n *NoLinter) regMarks(comment *ast.CommentGroup, node ast.Node) {
	var (
		lastCommentIdx int
		allCommentText = comment.Text()
		commentLines   = strings.Split(allCommentText, "\n")
		commentLineNo  = n.fset.Position(comment.End()).Line
		nodeLineNo     = n.fset.Position(node.Pos()).Line
	)

	if commentLineNo != nodeLineNo && commentLineNo != nodeLineNo-1 {
		// The `nolint` comment does not end on the same line or the line above
		// the node..
		return
	}

	switch {
	case len(commentLines) == 1:
		lastCommentIdx = len(commentLines) - 1
	case len(commentLines) >= 2:
		lastCommentIdx = len(commentLines) - 2
	}

	lastLineOfComment := commentLines[lastCommentIdx]
	for _, block := range strings.Split(lastLineOfComment, "//") {
		block := strings.TrimSpace(block)
		if block == "nolint" {
			n.regMark(node, "")
		}
		if strings.HasPrefix(block, "nolint:") {
			for _, category := range strings.Split(strings.TrimPrefix(block, "nolint:"), ",") {
				n.regMark(node, strings.TrimSpace(category))
			}
		}
	}
}
