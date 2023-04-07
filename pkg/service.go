package pkg

import (
	"encoding/base64"
	"encoding/json"
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
	ProcessedContent string   `json:"processedContent"`
}

func GetContent(url string) (string, error) {

	method := "GET"
	token := os.Getenv("GITHUB_TOKEN")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "token "+token)

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
		importRegex = regexp.MustCompile(`^\s*import\s+(\w+|\.)+\s*$`)
	case ".java":
		importRegex = regexp.MustCompile(`^\s*import\s+(?:static\s+)?([\w\.]+)\s*;?$`)
	case ".ts", ".js":
		importRegex = regexp.MustCompile(`import\s+(?:.+\s+from\s+)?['"](.+)['"]`)
	default:
		return nil, fmt.Errorf("invalid file extension: %s", fileExtension)
	}

	if content == "" {
		return nil, fmt.Errorf("content not found")
	}

	matches := importRegex.FindAllStringSubmatch(content, -1)
	modules := make([]string, len(matches))
	for i, match := range matches {
		modules[i] = match[1]
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

func FetchContentsRecursively(client *github.Client, owner, repo, path string) ([]Box, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)

	if err != nil {
		return nil, err
	}

	var boxes = []Box{}

	for _, content := range directoryContent {
		if *content.Type == "dir" {
			subBoxes, err := FetchContentsRecursively(client, owner, repo, *content.Path)
			if err != nil {
				return nil, err
			}

			boxes = append(boxes, subBoxes...)
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
				myFilePath := *content.Path
				myModules, err := GetImportedModulesList(myFilePath, myContent)
				if err != nil {
					return nil, err
				}

				box := Box{
					Content:          myContent,
					FilePath:         myFilePath,
					Modules:          myModules,
					ProcessedContent: "",
				}
				fmt.Println(box.FilePath, box.Modules)
				boxes = append(boxes, box)
			}
		}
	}

	return boxes, nil
}
