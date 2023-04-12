package pkg

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"go/parser"
	"go/token"
	"net/http"
	"path/filepath"

	"github.com/google/go-github/v50/github"

	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Response struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Sha         string `json:"sha"`
	Size        int    `json:"size"`
	Url         string `json:"url"`
	HtmlUrl     string `json:"html_url"`
	GitUrl      string `json:"git_url"`
	DownloadUrl string `json:"download_url"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Links       struct {
		Self string `json:"self"`
		Git  string `json:"git"`
		Html string `json:"html"`
	} `json:"_links"`
}

type Box struct {
	Content          string   `json:"content"`
	FilePath         string   `json:"filePath"`
	Modules          []string `json:"modules"`
	Funcs			[]string `json:"funcs"`
	ProcessedContent string   `json:"processedContent"`
	ProcessedCodeContent  string   `json:"processedCodeContent"`
}

func GetContent(url string) (string, error) {
	method := "GET"
	token := os.Getenv("GITHUB_TOKEN")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "token " + token)

	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var resContent Response

	err = json.NewDecoder(res.Body).Decode(&resContent)
	if err != nil {
		return "", err
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(resContent.Content)
	if err != nil {
		return "", err
	}

	if string(decodedBytes) == "" {
		return "", fmt.Errorf("Content Not Found")
	}

	return string(decodedBytes), nil
}

func GetImportedModulesListInGolang(content string) ([]string, error) {
	// 문자열을 파일로 변환
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	// import 문에서 사용된 패키지 추출
	var packages []string
	for _, i := range file.Imports {
		path := i.Path.Value[1 : len(i.Path.Value)-1]
		packages = append(packages, path)
	}

	// 추출된 패키지 리스트 출력
	return packages, nil
}

func GetImportedModulesList(filePath string, content string) ([]string, error) {
	var importRegex *regexp.Regexp
	fileExtension := filepath.Ext(filePath)

	switch fileExtension {
	case ".go":
		return GetImportedModulesListInGolang(content)
	case ".py":
		importRegex = regexp.MustCompile(`import\s+(?:(?:(?:[\w.]+)\s+as\s+\w+\s*,\s*)*(?:[\w.]+)\s*?)`)
	case ".java":
		importRegex = regexp.MustCompile(`import\s+(?:(?:(?:static\s+)?[\w.*]+(?:\s+as\s+\w+)?\s*,\s*)*(?:static\s+)?[\w.*]+\s*(?:as\s+\w+)?);`)
	case ".ts", ".js":
		importRegex = regexp.MustCompile(`import\s+(?:(?:(?:{.*?})|\S+)\s+from\s+)?['"](?P<path>@?[^'"]+)['"](?:;)?`)
	default:
		return nil, fmt.Errorf("invalid file extension: %s", fileExtension)
	}

	if content == "" {
		return nil, fmt.Errorf("content not found")
	}

	matches := importRegex.FindAllStringSubmatch(content, -1)
	modules := make([]string, len(matches))
	for i, match := range matches {
		cleanPath := filepath.Clean(match[1])
		relPathExpr := regexp.MustCompile(`^[.~$@].*?[/\\]?`)
		cleanPath = relPathExpr.ReplaceAllString(cleanPath, "")

		modules[i] = cleanPath
	}

	return modules, nil
}

func getResponseFromChatGPT(message string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			Stream:      false,
			Temperature: 0.5,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func mergeMaps[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	merged := make(map[K]V)
	for key, value := range m1 {
		merged[key] = value
	}
	for key, value := range m2 {
		merged[key] = value
	}
	return merged
}

func FetchContentsRecursively(client *github.Client, owner, repo, path string) (map[string]Box, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	everyFilePathsList := []string{}
	var packageName string
	
	if err != nil {
		return nil, err
	}

	paths := make(map[string]Box)

	for _, content := range directoryContent {
		if *content.Type == "dir" {
			subPath, err := FetchContentsRecursively(client, owner, repo, *content.Path)
			if err != nil {
				return nil, err
			}

			paths = mergeMaps(paths, subPath)

		} else {
			extension := strings.Split(*content.Path, ".")
			var extensionName string

			if len(extension) == 1 {
				continue
			} else {
				extensionName = extension[len(extension)-1]
			}

			if extensionName == "js" || extensionName == "ts" || extensionName == "py" || extensionName == "go" || extensionName == "java" {
				myContent, err := GetContent(*content.URL)
				if err != nil {
					fmt.Println(err)
				}

				if strings.Contains(*content.Path, "go.mod") {	
					packageName = strings.Split(strings.Split(myContent, "module ")[1], "\n")[0]
				}

				myFilePath := *content.Path
				everyFilePathsList = append(everyFilePathsList, myFilePath)

				myModules, err := GetImportedModulesList(myFilePath, myContent)
				if err != nil {
					return nil, err
				}

				myFuncs := getFuncs(myModules, myContent)

				box := Box{
					Content:          myContent,
					FilePath:         myFilePath,
					Modules:          myModules,
					Funcs:            myFuncs,
					ProcessedContent: "",
					ProcessedCodeContent: "",
				}

				paths[myFilePath] = box
			}
		}
	}

	pathToModule, moduleToFunc := setPaths(paths)

	functionsExplanation, err := getFunctionsMeaning(everyFilePathsList, pathToModule, moduleToFunc)
	if err != nil {
		return nil, err
	}

	return paths, nil
}
func getPackageNameFromGoFile(directoryContent []*github.RepositoryContent) (string, error) {
	for _, content := range directoryContent {
		if strings.Contains(*content.Path, "go.mod") {
			mod, err := GetContent(*content.URL)
			if err != nil {
				fmt.Println(err)
			}

			// get the string started with "module"
			packageName := strings.Split(strings.Split(mod, "module ")[1], "\n")[0]
			return packageName, nil
		}
	}

	return "", errors.New("package name not found")
}
func isInternalModule(module string, everyFilePathsList []string) bool {
	for _, fp := range everyFilePathsList {
		if strings.Contains(fp, module) {
			return true
		}
	}
	return false
}

func getFunctionsMeaning(everyFilePathsList []string, pathToModule map[string][]string , moduleToFunc map[string][]string) (map[string]string, error) {
	functionsExplanation := map[string]string{}

	for _, fp := range everyFilePathsList {
		modulesListOfFileUses := pathToModule[fp]
		for _, module := range modulesListOfFileUses {
			functionsListOfModuleUses := moduleToFunc[module]
			for _, function := range functionsListOfModuleUses {
				if functionsExplanation[function] == "" {
					// if function is originated inner context of the file stream, or its one of the inner functions, then ask chatgpt for explanation
					// else, continue the loof.
					if isInternalModule(module, everyFilePathsList) {
						var message string
						response, err := getResponseFromChatGPT(message)
						functionsExplanation[function] = response

						if err != nil {
							fmt.Println(err)
						}
						
					} else {
						continue
					}
				}
			}
			
		}
	}

	return functionsExplanation, nil
}

func setPaths(paths map[string]Box) (map[string][]string, map[string][]string){
	pathToModule := map[string][]string{}
	moduleToFunc := map[string][]string{}

	for _, box := range paths {
		filePath := box.FilePath
		modules := box.Modules
		funcs := box.Funcs
		
		pathToModule[filePath] = modules

		for _, module := range modules {
			if _, ok := moduleToFunc[module]; !ok {
				moduleToFunc[module] = []string{}
			}
			
			for _, f := range funcs {
				moduleNameInFunc := strings.Split(f, ".")[0]
				if strings.Contains(moduleNameInFunc, module) {
					moduleToFunc[module] = append(moduleToFunc[module], f)
				}
			}
		}
	}

	return pathToModule, moduleToFunc
}

func getFuncs(modules []string, content string) []string {
	var funcs []string
	
	for _, module := range modules {
		slash := strings.Split(module, "/")
		moduleName := slash[len(slash)-1]

		funcRegex := regexp.MustCompile(`\b` + moduleName + `\.[a-zA-Z0-9_]+\b`)
		matches := funcRegex.FindAllString(content, -1)

		for _, m := range matches {
			for _, f := range funcs {
				if f == m {
					continue
				}
				funcs = append(funcs, m)
			}
		}
	}

	return funcs
}

// func getTree(paths map[string]Box) []Box {
// 	var tree []string

// 	for filePath, box := range paths {
// 		slash := strings.Split(filePath, "/")
// 		if len(slash) == 0 {
// 		   continue
// 		}

// 		content := box.Content
// 		modules := box.Modules
// 		funcs := box.Funcs

		



// 	}

	

// 	return tree
	
// }