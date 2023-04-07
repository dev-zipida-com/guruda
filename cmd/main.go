package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dev-zipida-com/guruda/pkg"
	"github.com/google/go-github/v50/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func main() {
	var username, reponame string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)
	fmt.Print("Enter GitHub repository name: ")
	fmt.Scanf("%s", &reponame)

	token := os.Getenv("GITHUB_TOKEN")

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	paths, err := pkg.FetchContentsRecursively(client, username, reponame, "")
	if err != nil {
		log.Fatal("FetchContentsRecursively error: ", err)
		return
	}

	fmt.Println("paths: ", len(paths))

}