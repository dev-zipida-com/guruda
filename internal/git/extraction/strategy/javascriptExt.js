function getFunctionsStartingIndex(functionName) {
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

    let index = fileContents.search(functionDeclarationRegex);
    if (index !== -1) {
        return console.log(index);
    }

    index = fileContents.search(functionExpressionRegex);
    if (index !== -1) {
        return console.log(index);
    }

    index = fileContents.search(arrowFunctionRegex);
    if (index !== -1) {
        return console.log(index);
    }

    index = fileContents.search(classMethodRegex);
    if (index !== -1) {
        return console.log(index);
    }

    index = fileContents.search(staticClassMethodRegex);
    if (index !== -1) {
        return console.log(index);
    }

    console.log(`The function ${functionName} could not be found`);
    return;
}

function getFunctionsCode(startingIndex, fileContents) {
    let functionCode = "";
    let openBrackets = 0;
    let closeBrackets = 0;
    let i = startingIndex;

    while (openBrackets !== closeBrackets || openBrackets === 0) {
        if (fileContents[i] === "{") {
            openBrackets++;
        } else if (fileContents[i] === "}") {
            closeBrackets++;
        }
        functionCode += fileContents[i];
        i++;
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
        if (functionName === "function") continue;
        functionDeclarations.push(match[1]);
    }

    while ((match = regex.functionExpression.exec(code))) {
        const functionName = match[1] || match[2] || match[3];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push(functionName);
        }
    }

    while ((match = regex.arrowFunction.exec(code))) {
        const functionName = match[1] || match[4];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push(functionName);
        }
    }

    while ((match = regex.classMethod.exec(code))) {
        const functionName = match[1];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push(functionName);
        }
    }

    while ((match = regex.staticClassMethod.exec(code))) {
        const functionName = match[1];
        if (functionName && !functionDeclarations.includes(functionName)) {
            if (functionName === "function") continue;
            functionDeclarations.push(functionName);
        }
    }

    return {
        modules: modules,
        functionDeclarations: functionDeclarations,
    };
}
