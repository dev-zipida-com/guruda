package git

import (
	"context"
	"fmt"

	"github.com/google/go-github/v50/github"
)

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

func GetRepository(client *github.Client, owner, repo, path string) (map[string]string, error) {
	_, directoryContent, _, err := client.Repositories.GetContents(context.Background(), owner, repo, path, nil)
	if err != nil {
		return nil, err
	}

	paths := make(map[string]string)

	for _, content := range directoryContent {
		if *content.Type == "dir" {
			subPath, err := GetRepository(client, owner, repo, *content.Path)
			if err != nil {
				return nil, err
			}

			paths = mergeMaps(paths, subPath)

		} else {
			myContent, err := GetContent(*content.URL)
			if err != nil {
				fmt.Println(err)
			}

			myFilePath := *content.Path
			
			paths[myFilePath] = myContent
		}
	}

	return paths, nil
}