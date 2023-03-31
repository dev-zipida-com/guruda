package main

import (
	"fmt"

	"github.com/dev-zipida-com/pkg/service"
)

func main() {
	var username, reponame string
	fmt.Print("Enter GitHub username: ")
	fmt.Scanf("%s", &username)
	fmt.Print("Enter GitHub repository name: ")
	fmt.Scanf("%s", &reponame)

	_, directoryContent, err := service.FetchProjects(username, reponame, "/")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println()

	for _, content := range directoryContent {
		fmt.Println(content)
		fmt.Println(service.GetContent(*content.URL))
		fmt.Println()
	}
}
