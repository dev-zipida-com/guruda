import ast
import sys

# 모든 함수를 추출하여 리스트에 담는 함수
def extract_all_functions(code):
    # 파이썬 코드를 AST로 변환하는 함수
    module = ast.parse(code)

    # 모든 함수의 위치 정보를 찾아서 함수의 내용을 추출하는 함수
    function_codes = []
    for node in module.body:
        if isinstance(node, ast.FunctionDef):
            start = node.lineno
            end = node.body[-1].lineno
            function_code = ""
            for i, line in enumerate(code.split('\n')):
                if i >= start-1 and i < end:
                    function_code += line + '\n'
            function_codes.append(function_code)

    if len(function_codes) == 0:
        raise ValueError("No functions found in code")

    return function_codes


arg1 = sys.argv[1]
print(extract_all_functions(arg1))