package main

import (
	"context"
	"fmt"
	"math"
	"os"

	"github.com/google/go-github/v50/github"
	openai "github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2"
)

func GetResponse(sysMessage string, message string) (string, error) {
	client := openai.NewClient("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")

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: sysMessage,
				},
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

func splitCode(s string) []string {
	maxWordLen := 2400
	strLen := len(s)

	splitedText := []string{}
	pageNum := int(math.Ceil(float64(strLen)/float64(maxWordLen)))
	interval := 50

	startIndex := 0
	endIndex := 0

	if strLen >= maxWordLen {
		for i := 0; i < pageNum; i++ {
			if i == 0 {
				startIndex = 0
				endIndex = maxWordLen + interval
			} else if i == pageNum-1 {
				startIndex = i*maxWordLen - interval
				endIndex = strLen 
			} else {
				startIndex = i*maxWordLen - interval
				endIndex = (i+1)*maxWordLen + interval
			}
			splitedText = append(splitedText, s[startIndex:endIndex])
		}
	} else {
		fmt.Println("Text splition is done")
		return []string{s}
	}

	fmt.Println("Text splition is done")
	return splitedText
}

func main() {
	lang := "golang"
	prompt := `I want you to act as a code analyzer and generate only a one JSON object of all the functions, modules, and methods used in the given code. The code is written in ` + lang + ` and is a part of a larger program. The code is given below:
	
	`
	token := "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"

    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: token},
    )
    tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)

    owner := "dev-zipida-com"
    repo := "ogada"
    ref := "deploy"
    path := "/frontend/src/pages/components/Address_PR.js"

    file, _, _, err := client.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{
        Ref: ref,
    })

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
	content, err := file.GetContent()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

	code := content

	splitedText := splitCode(code)

	for _, text := range splitedText {
		res, err := GetResponse(prompt, text)
		if err != nil {
			fmt.Println("ERROR: ")
			fmt.Println(err)
		}

		fmt.Println(res)
	}
}