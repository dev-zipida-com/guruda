package git

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os/exec"
	"strings"
)

// 함수의 내용을 추출하는 인터페이스
type FunctionExtractor interface {
	ExtractFunctionCode(code string) ([]string, error)
}

// Go 언어에서 함수의 내용을 추출하는 구조체
type GoFunctionExtractor struct{}

// 함수의 전체 코드를 추출하는 함수
func (extractor GoFunctionExtractor) ExtractFunctionCode(code string) ([]string, error) {
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

// Python 언어에서 함수의 내용을 추출하는 구조체
type PythonFunctionExtractor struct{}

// Python 언어에서 함수의 내용을 추출하는 메서드
func (extractor PythonFunctionExtractor) ExtractFunctionCode(code string) ([]string, error) {
	// 함수 정의를 찾기 위한 정규표현식
	cmd := exec.Command("python3", "./python/extractFunction.py", code)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	functions := strings.Split(string(output), "', '")
	for i, function := range functions {
		functions[i] = strings.Trim(function, "['")
		functions[i] = strings.TrimSuffix(functions[i], "']")
	}

	return functions, nil
}

// 추출기에 따라 함수의 내용을 추출하는 함수
func ExtractFunctionCodes(extractor FunctionExtractor, code string) ([]string, error) {
	return extractor.ExtractFunctionCode(code)
}
