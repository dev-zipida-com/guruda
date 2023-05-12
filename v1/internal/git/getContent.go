package git

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dev-zipida-com/guruda/internal/structures"
)

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

	var resContent structures.Response

	err = json.NewDecoder(res.Body).Decode(&resContent)
	if err != nil {
		return "", err
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(resContent.Content)
	if err != nil {
		return "", err
	}

	if string(decodedBytes) == "" {
		return "", fmt.Errorf("content not found")
	}

	return string(decodedBytes), nil
}