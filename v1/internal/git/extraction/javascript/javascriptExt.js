function getFunc(functionName, fileContents) {
    const functionDeclarationRegex = new RegExp(`function ${functionName}\\(`);
    const functionExpressionRegex = new RegExp(
        `const ${functionName} = function\\(`
    );
    const arrowFunctionRegex = new RegExp(
        `const ${functionName} = \\(.*\\) =>`
    );
    const classMethodRegex = new RegExp(
        `(${functionName}\\(.*\\)\\s*{)|(\\s*${functionName}\\(.*\\)\\s*{)`
    );
    const staticClassMethodRegex = new RegExp(
        `static\\s+${functionName}\\(.*\\)\\s*{`
    );

    const regexpes = [
        functionDeclarationRegex,
        functionExpressionRegex,
        arrowFunctionRegex,
        classMethodRegex,
        staticClassMethodRegex,
    ];

    const index = regexpes.findIndex(
        (regex) => fileContents.search(regex) !== -1
    );

    if (index !== -1) {
        return getFunctionsCode(
            fileContents.search(regexpes[index]),
            fileContents
        );
    } else {
        return "Found No Function";
    }
}

function getFunctionsCode(startingIndex, fileContents) {
    let functionCode = "";
    let openBraces = 0;
    let isFunc = false;
    let lastBrace = "";

    for (let i = startingIndex; i < fileContents.length; i++) {
        const char = fileContents[i];

        if (char === "{" || char === "(") {
            openBraces++;
            lastBrace = char;
            if (isFunc === false) {
                isFunc = true;
            }
        } else if (char === "}" || char === ")") {
            openBraces--;
            lastBrace = char;
        }

        functionCode += char;

        if (isFunc && openBraces === 0 && lastBrace !== ")") {
            break;
        }
    }

    return functionCode;
}

function extractModules(code) {
    const statements = code.match(
        /(import\s+(?:([\w*{}\n\r\t, ]+?)\s+from\s+)?['"](.+?)['"]|require\s*\(\s*['"][^'"]+['"]\s*\))/gm
    );

    const modules = {};

    statements.forEach((statement) => {
        const moduleMatch = statement.match(
            /(import\s+(?:([\w*{}\n\r\t, ]+?)\s+from\s+)?['"](.+?)['"]|require\s*\(\s*['"]([^'"]+)['"]\s*\))/
        );
        if (moduleMatch) {
            const [, , functions, moduleName, functionName] = moduleMatch;
            if (!modules[moduleName]) {
                modules[moduleName] = [];
            }
            if (functionName) {
                modules[moduleName].push(functionName);
            } else if (functions) {
                const functionList = functions
                    .split(",")
                    .map((func) => func.trim());
                modules[moduleName].push(...functionList);
            }
        }
    });

    const regex = {
        functionDeclaration: /function\s+(\w+)\s*\(/g,
        functionExpression:
            /const\s+(\w+)\s*=\s*function\s*\(|const\s+(\w+)\s*=\s*\(|(\w+)\s*=\s*\(/g,
        arrowFunction:
            /const\s+(\w+)\s*=\s*\((\w+(,\s*\w+)*)?\)\s*=>|(\w+)\s*=\s*\((\w+(,\s*\w+)*)?\)\s*=>/g,
        classMethod: /\s+(\w+)\s*\(\w*(,\s*\w+)*\)\s*{/g,
        staticClassMethod: /static\s+(\w+)\s*\(\w*(,\s*\w+)*\)\s*{/g,
    };

    const functionDeclarations = [];

    let match;

    while ((match = regex.functionDeclaration.exec(code))) {
        const functionName = match[1];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push([
                functionName,
                getFunc(functionName, code),
            ]);
        }
    }

    while ((match = regex.functionExpression.exec(code))) {
        const functionName = match[1] || match[2] || match[3];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push([
                functionName,
                getFunc(functionName, code),
            ]);
        }
    }

    while ((match = regex.arrowFunction.exec(code))) {
        const functionName = match[1] || match[4];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push([
                functionName,
                getFunc(functionName, code),
            ]);
        }
    }

    while ((match = regex.classMethod.exec(code))) {
        const functionName = match[1];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push([
                functionName,
                getFunc(functionName, code),
            ]);
        }
    }

    while ((match = regex.staticClassMethod.exec(code))) {
        const functionName = match[1];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push([
                functionName,
                getFunc(functionName, code),
            ]);
        }
    }

    return {
        modules: modules,
        functionDeclarations: functionDeclarations,
    };
}
