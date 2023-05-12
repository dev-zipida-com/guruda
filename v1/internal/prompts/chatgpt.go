package prompts

func ExtractObjects(code string) string {
	prompt := `
		Extract the Functions, Variables, Classes, Decorators name, and Modules name that used in the code below "---". 
		Follow there conditions.
		1. Answer me in a json format. No need additional explanations of you, and only give answers that extract the names of the object.
		2. Your json type response is composed of "Functions", "Modules". "Functions" is an array of objects like functions, variables, classes, and decorators used in the code. "Modules" is an array of strings that composed of imported modules in the code.
		3. In case of getting the "Functions", describe a function with it's name, arguments, and return values. for example, {"Functions": [{"name": "create", "arguments": [{"name": "createEventFilesDto", "type":"createEventFilesDto"}], "return": [{"name": "createdEventFiles", type: "Promise<EventFiles>"}]}, ... ], ...}.
		4. In case of getting the "Modules", describe a module with it's name and a imported functions, classes, or decorators with its name too. for example, {"Modules": [{"name": "@nestjs/common", import:["Injectable"]}, {"name": "mongoose", import:["FilterQuery", "Model"]}, ...], ...}.
		---

	`
	return prompt + code
}

func ExtractModules(code string) string {
	prompt := `
		Extract the Modules name that used in the code below "---". 
		Follow there conditions.
			1. Answer me in a json format. No need additional explanations of you, and only give answers that extract the names of the object.
			2. Your json type response is composed of "Modules". "Modules" is an array of strings that composed of imported modules in the code.
			3. In case of getting the "Modules", describe a module with it's name and a imported functions, classes, or decorators with its name too. for example, {"Modules": [{"name": "@nestjs/common", import:["Injectable"]}, {"name": "mongoose", import:["FilterQuery", "Model"]}, ...], ...}.
		---

	`
	return prompt + code
}

func AskRefactoredCode(explanationByNL string, extension string) string {
	prompt := `
		write the code with ` + extension + ` language. following below my explanation.
	
	`

	return prompt + explanationByNL
}

func AskReadMeContents(explanationByNL string) string {
	prompt := `
		write the README.md contents with markdown style. following below my explanation.

	`

	return prompt + explanationByNL
}

func AskFunctionsMeaning(code string, explanationByNL []string) string {
	var expInOneString string
	for _, explanation := range explanationByNL {
		expInOneString += explanation + " "
	}

	prompt := `
		Explain the code below "---" with an English.
		Follow there conditions.
			1. 
			2. No additional explanations of you, and only give answers that convert code to natural language or natural language to code.
			3. Answer me in a json format.
			4. There are some explanations about the functions that used in the code. You can use the explanations as a reference. explanations are below.
				` + expInOneString + `
		---
		
	`
	
	return prompt + code
}

func AskFunctionsMeaningInShort(code string, explanationByNL []string) string {
	var expInOneString string
	for _, explanation := range explanationByNL {
		expInOneString += explanation + " "
	}

	prompt := `
		Explain the code below "---" with an English.
		Follow there conditions.
			1. Explain the code in short. Maybe 1~2 sentences be right.
			2. No additional explanations of you, and only give answers that convert code to natural language or natural language to code.
			3. Answer me in a json format.
			4. There are some explanations about the functions that used in the code. You can use the explanations as a reference. explanations are below.
				` + expInOneString + `
		---
		
	`
	
	return prompt + code
}