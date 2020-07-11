# nolint

Nolint will make the `go/analysis` linters
be able to ignore diagnostics with `//nolint` comment

[![Go Report Card](https://goreportcard.com/badge/github.com/kyoh86/nolint)](https://goreportcard.com/report/github.com/kyoh86/nolint)
[![Coverage Status](https://img.shields.io/codecov/c/github/kyoh86/nolint.svg)](https://codecov.io/gh/kyoh86/nolint)

## Install

```
go get github.com/kyoh86/nolint
```

## Usage

### For linter users

If you are using the linters which using the `nolint`,
you can ignore (like below) diagnostics that they reported.

```go
for _, p := []int{10, 11, 12} {
	t.Run("dummy", func(t *testing.T) {
		foo.Bar(&p) // nolint
	})
}
```

`// nolint` will be ignore all diagnostics in the line.
And you can specify categories which you want to ignore.

```go
// nolint:someCategory,anotherCategory
```

### For custom linter users

If you are using the linters with `go/analysis/xxxxchecker`,
linters can be wrapped like below.

```go
multichecker.Main(
	nolint.WrapAll(
		exportloopref.Analyzer,
		bodyclose.Analyzer,
		// ...
	),
)
```
　
Then, diagnostics will be able to be ignored with a comment.

```go
// nolint
```

### For linter creators

If you are creator of `go/analysis` linters,
use the `nolint.Analyzer` like below.

```go
var Analyzer = &analysis.Analyzer{
	Run:      run,
	Requires: []*analysis.Analyzer{nolint.Analyzer},
	// ...
}

func run(pass *analysis.Pass) (interface{}, error) {
	noLinter := pass.ResultOf[nolint.Analyzer].(*nolint.NoLinter)
	// ...
	if !noLinter.IgnoreNode(node, "someCategory") {
		pass.Report(analysis.Diagnostic{
			Category: "someCategory",
			// ...
		})
	}
	// ...
}
```

NOTE: Category will be used to specify which diagnostic should be ignored.

```go
// nolint:someCategory,anotherCategory
```

# LICENSE

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](http://www.opensource.org/licenses/MIT)

This is distributed under the [MIT License](http://www.opensource.org/licenses/MIT).
