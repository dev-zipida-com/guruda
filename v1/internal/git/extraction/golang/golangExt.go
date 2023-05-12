package strategy

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/dev-zipida-com/guruda/internal/structures"
)
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ExtractImportedModules(code string, fPath string) []structures.Module {
    var result []structures.Module
	var myModules []string
	var calls []string

    fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.AllErrors|parser.ParseComments)
    if err != nil {
        panic(err)
    }

    for _, importSpec := range file.Imports {
		module := strings.Trim(importSpec.Path.Value, "\"")
		if importSpec.Name != nil {
			module = importSpec.Name.Name
		}
		myModules = append(myModules, module)
	}

	ast.Inspect(file, func(node ast.Node) bool {
		switch x := node.(type) {
		case *ast.CallExpr:
			switch f := x.Fun.(type) {
			case *ast.SelectorExpr:
				if ident, ok := f.X.(*ast.Ident); ok {
					module := ident.Name
					function := f.Sel.Name
					calls = append(calls, fmt.Sprintf("%s.%s", module, function))
				}
			}
		}
		return true
	})

	for _, m := range myModules {
		cr := []string{}
		mSlashed := strings.Split(m, "/")
		mr := mSlashed[len(mSlashed)-1]

		for _, c := range calls {
			cModule := strings.Split(c, ".")[0]
			cFunc := strings.Split(c, ".")[1]
			if cModule == mr {
				if !contains(cr, cFunc) {
					cr = append(cr, cFunc)
				}
			}
		}

		result = append(result, structures.Module{
			Name: m,
			Imported: cr,
		})
	}
	
    return result
}
