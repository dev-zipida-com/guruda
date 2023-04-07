package pkg

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/v50/github"

	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
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
	Content string `json:"content"`
	FilePath string `json:"filePath"`
	Modules []string `json:"modules"`
	ProcessedContent string `json:"processedContent"`
}

func GetContent(url string) string {
	// err := godotenv.Load("../.env")
	// if err != nil {
    //     log.Fatal("Error loading .env file")
    // }

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", "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")

	if err != nil {
		fmt.Println(err)
		return ""
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	var resContent Response

	err = json.NewDecoder(res.Body).Decode(&resContent)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(resContent.Content)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(decodedBytes)
}

func GetImportedModulesList(filePath string, content string) ([]string, error) {
	splitedNamebySlash := strings.Split(filePath, "/")
	fileName := splitedNamebySlash[len(splitedNamebySlash)-1]
	extensions := strings.Split(fileName, ".")
	extensionName := extensions[len(extensions)-1]
	var importRegex *regexp.Regexp

	if(extensionName == "go") {
		importRegex = regexp.MustCompile(`^\s*import\s+(?:.+\s+)?["'](.+)["']`)
	} else if (extensionName == "py") {
		importRegex = regexp.MustCompile(`^\s*import\s+(\w+|\.)+\s*$`)
	} else if (extensionName == "java") {
		importRegex = regexp.MustCompile(`^\s*import\s+(?:static\s+)?([\w\.]+)\s*;?$`)
	} else if (extensionName == "ts" || extensionName == "js") {
		importRegex = regexp.MustCompile(`import\s+(?:.+\s+from\s+)?['"](.+)['"]`)
	} else {
		return make([]string, 1), fmt.Errorf("invalid file extension")
	}

	matches := importRegex.FindAllStringSubmatch(content, -1)
    modules := make([]string, len(matches))
    for i, match := range matches {
        modules[i] = match[1]
    }

    return modules, nil
}

func getResponseFromChatGPT(message string) (string, error) {
	err := godotenv.Load("../.env")
	if err != nil {
        log.Fatal("Error loading .env file")
    }

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
				extensionName = extension[len(extension) - 1]
			}
			if extensionName == "js" || extensionName == "ts" || extensionName == "py" || extensionName == "go" || extensionName == "java" {
				myContent := GetContent(*content.URL)
				myFilePath := *content.Path
				myModules, err := GetImportedModulesList(myFilePath, myContent)
				if err != nil {
					return nil, err
				}

				box := Box{
					Content: myContent,
					FilePath: myFilePath,
					Modules: myModules,
					ProcessedContent: "",
				}
				
				boxes = append(boxes, box)
			}
		}
	}

	return boxes, nil
}
