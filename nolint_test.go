package nolint_test

import (
	"testing"

	"github.com/kyoh86/nolint"
	"golang.org/x/tools/go/analysis/analysistest"
	"golang.org/x/tools/go/analysis/passes/printf"
)

func TestA(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, nolint.Wrap(printf.Analyzer), "a")
}
