import ast
import sys

def extract_imports(code_str):
    imports = []
    from_imports = []

    # 코드 문자열을 AST로 변환
    tree = ast.parse(code_str)

    # 모든 import 구문과 from-import 구문을 찾아서 imports와 from_imports 리스트에 추가
    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            imports.extend([alias.name for alias in node.names])
        elif isinstance(node, ast.ImportFrom):
            module_name = node.module
            function_names = [alias.name for alias in node.names]
            from_imports.append((module_name, function_names))

    # imports와 from_imports 리스트를 이용해 key-value 형태의 딕셔너리를 생성하여 반환
    import_dict = {}
    for module_name, function_names in from_imports:
        if function_names:
            for function_name in function_names:
                if module_name in import_dict:
                    import_dict[module_name].append(function_name)
                else:
                    import_dict[module_name] = [function_name]
        else:
            import_dict[module_name] = [""]

    for imported_module in imports:
        import_dict[imported_module] = [""]

    return import_dict

arg1 = sys.argv[1]
print(extract_imports(arg1))