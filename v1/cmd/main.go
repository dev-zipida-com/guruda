package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func main() {
    ctx := context.Background()
    ts := oauth2.StaticTokenSource(
        &oauth2.Token{AccessToken: "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"},
    )
    tc := oauth2.NewClient(ctx, ts)

    client := github.NewClient(tc)
    // use the client to make API calls
    functionName := "getDetectionImage"
	owner := "dev-zipida-com"
	repoName := "facematch"
	lang := "ts"
	
	query := fmt.Sprintf("%s in:file language:%s repo:%s/%s", functionName, lang, owner, repoName)
	opts := &github.SearchOptions{Sort: "indexed", Order: "desc"}
	result, _, err := client.Search.Code(ctx, query, opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, c := range result.CodeResults {
		fmt.Println(*c.Path)
	}

}