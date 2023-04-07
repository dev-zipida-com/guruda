package main

import (
	"fmt"
	"log"

	"github.com/dev-zipida-com/guruda/pkg"
	"github.com/google/go-github/v50/github"
)

func main() {
	var username, reponame string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)
	fmt.Print("Enter GitHub repository name: ")
	fmt.Scanf("%s", &reponame)

	client := github.NewClient(nil)

	_, err := pkg.FetchContentsRecursively(client, username, reponame, "")
	if err != nil {
		log.Fatal(err)
		return
	}

	
}