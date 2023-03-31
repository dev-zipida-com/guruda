package pkg

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v50/github"
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

func GetContent(url string) string {

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

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

func FetchProjects(username, reponame, path string) (*github.RepositoryContent, []*github.RepositoryContent, error) {
	client := github.NewClient(nil)
	// opt := &github.RepositoryContentGetOptions{}
	fileContent, directoryContent, _, err := client.Repositories.GetContents(context.Background(), username, reponame, path, nil)
	return fileContent, directoryContent, err
}
