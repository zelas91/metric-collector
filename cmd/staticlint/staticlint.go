// Package main staticlint checks for the use of a direct os.Exit call in the main function of the main package.
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"path/filepath"
	"strings"
)

var (
	analyzers []*analysis.Analyzer
	Analyzer  = &analysis.Analyzer{
		Name: "staticlint",
		Doc:  "проверяет наличие прямого вызова os.Exit в функции main пакета main",
		Run:  run,
	}
)

// shouldIgnoreFile checks whether the absolute path to the file contains the substring "go-build".
// If it contains, the file is ignored.
func shouldIgnoreFile(filename string) bool {
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return false
	}
	return strings.Contains(absPath, "go-build")
}
func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Package).Filename
		if !shouldIgnoreFile(filename) {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Name.Name != "main" {
					continue
				}
				ast.Inspect(fn.Body, func(node ast.Node) bool {
					if call, ok := node.(*ast.CallExpr); ok {
						if selExp, ok := call.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := selExp.X.(*ast.Ident); ok {
								if ident.Name == "os" && selExp.Sel.Name == "Exit" {
									pass.Reportf(selExp.Sel.NamePos, "os.Exit calls are prohibited in main()")

								}

							}
						}
					}
					return true
				})
			}
		}
	}

	return nil, nil
}
func init() {
	analyzers = []*analysis.Analyzer{
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		cgocall.Analyzer,
		copylock.Analyzer,
		deepequalerrors.Analyzer,
		errorsas.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		stdmethods.Analyzer,
		structtag.Analyzer,
		tests.Analyzer,
		unmarshal.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
	}

	analyzers = append(analyzers, Analyzer)
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}
	analyzers = append(analyzers, findLintAnalyzer(simple.Analyzers, "S1000"))
	analyzers = append(analyzers, findLintAnalyzer(stylecheck.Analyzers, "ST1001"))
	analyzers = append(analyzers, findLintAnalyzer(quickfix.Analyzers, "QF1004"))
}

func findLintAnalyzer(analyzers []*lint.Analyzer, name string) *analysis.Analyzer {
	for _, v := range analyzers {
		if v.Analyzer.Name == name {
			return v.Analyzer
		}
	}
	return nil
}
func main() {
	multichecker.Main(
		analyzers...,
	)

}
