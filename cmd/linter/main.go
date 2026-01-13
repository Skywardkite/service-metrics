package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
	"honnef.co/go/tools/go/ast/astutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "banpanic",
	Doc:  "reports usage of panic, log.Fatal and os.Exit outside main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {

		pkgName := file.Name.Name

		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			switch fun := call.Fun.(type) {

			// Проверяем на вызов panic(...)
			case *ast.Ident:
				if fun.Name == "panic" {
					pass.Reportf(call.Pos(), "usage of panic is forbidden")
				}

			// Проверяем на вызов log.Fatal(...) или os.Exit(...) вне функции main пакета main
			case *ast.SelectorExpr:
				pkgIdent, ok := fun.X.(*ast.Ident)
				if !ok {
					return true
				}

				if (pkgIdent.Name == "log" && fun.Sel.Name == "Fatal") ||
					(pkgIdent.Name == "os" && fun.Sel.Name == "Exit") {

					if pkgName != "main" || !isInMainFunc(file, call) {
						pass.Reportf(call.Pos(),
							"call to %s.%s is only allowed in main",
							pkgIdent.Name, fun.Sel.Name)
					}
				}
			}

			return true
		})
	}

	return nil, nil
}

func isInMainFunc(file *ast.File, call *ast.CallExpr) bool {
	path, _ := astutil.PathEnclosingInterval(file, call.Pos(), call.End())
	for _, node := range path {
		if fn, ok := node.(*ast.FuncDecl); ok {
			return fn.Name.Name == "main"
		}
	}
	return false
}

func main() {
	singlechecker.Main(Analyzer)
}
