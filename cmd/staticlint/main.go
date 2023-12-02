package main

import (
	"fmt"
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
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"path/filepath"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name: "staticlint",
	Doc:  "проверяет наличие прямого вызова os.Exit в функции main пакета main",
	Run:  run,
}

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
									fmt.Println()
									pass.Reportf(selExp.Sel.NamePos, "нельзя использовать прямой вызов os.Exit в функции main")

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
func main() {
	analyzers := []*analysis.Analyzer{
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
	analyzers = append(analyzers, simple.Analyzers[0].Analyzer)
	analyzers = append(analyzers, stylecheck.Analyzers[0].Analyzer)
	analyzers = append(analyzers, quickfix.Analyzers[0].Analyzer)

	multichecker.Main(
		analyzers...,
	)

}
