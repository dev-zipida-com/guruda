package git

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// 함수의 전체 코드를 추출하는 함수
func extractFunctionCode(code string) ([]string, error) {
	// Go 언어 코드를 파싱하여 AST를 생성하는 함수
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", code, 0)
	if err != nil {
		return nil, err
	}

	// 함수의 위치 정보를 찾아서 함수의 전체 코드를 추출하는 함수
	var functionCode []string
	for _, decl := range node.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			start := function.Pos()
			end := function.End()
			functionCode = append(functionCode, code[start-1:end-1])
		}
	}

	return functionCode, nil
}
