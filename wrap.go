package nolint

import (
	"golang.org/x/tools/go/analysis"
)

// WrapAll other Analyzers to ignore diagnostics marked by `//nolint`.
func WrapAll(analyzer ...*analysis.Analyzer) []*analysis.Analyzer {
	var wrapped []*analysis.Analyzer
	for _, a := range analyzer {
		wrapped = append(wrapped, Wrap(a))
	}
	return wrapped
}

// Wrap other Analyzer to ignore diagnostics marked by `//nolint`.
func Wrap(analyzer *analysis.Analyzer) *analysis.Analyzer {
	wrapped := *analyzer
	wrapped.Requires = append(wrapped.Requires, Analyzer)
	wrapped.Run = func(pass *analysis.Pass) (interface{}, error) {
		noLinter := pass.ResultOf[Analyzer].(*NoLinter)
		pseudo := *pass
		report := pseudo.Report
		pseudo.Report = func(diagnostic analysis.Diagnostic) {
			if noLinter.IgnoreDiagnostic(diagnostic) {
				return
			}
			report(diagnostic)
		}
		pseudo.Analyzer = analyzer
		return analyzer.Run(&pseudo)
	}
	return &wrapped
}

// WrapFunc will wrap an Analyzer to ignore diagnostics
// with a function.
func WrapFunc(
	analyzer *analysis.Analyzer,
	ignore func(analysis.Diagnostic) bool,
) *analysis.Analyzer {
	wrapped := *analyzer
	wrapped.Run = func(pass *analysis.Pass) (interface{}, error) {
		pseudo := *pass
		report := pseudo.Report
		pseudo.Report = func(diagnostic analysis.Diagnostic) {
			if !ignore(diagnostic) {
				report(diagnostic)
			}
		}
		pseudo.Analyzer = analyzer
		return wrapped.Run(&pseudo)
	}
	return &wrapped
}
