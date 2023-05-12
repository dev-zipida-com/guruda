package git

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dev-zipida-com/guruda/internal/chatgpt"
	"github.com/dev-zipida-com/guruda/internal/prompts"
	"github.com/dev-zipida-com/guruda/internal/structures"
)

func ExtractModules(paths map[string]string) (map[string][]structures.Module, error){
	m := make(map[string][]structures.Module)

	for path, code := range paths {
		splitedByALine := strings.Split(code, "\n")

		numLines := 100
		if len(splitedByALine) < numLines {
			numLines = len(splitedByALine)
		}

		newCode := strings.Join(splitedByALine[:numLines], "\n")

		chatgptResponse, err := chatgpt.GetResponse(prompts.ExtractModules(newCode))
		if err != nil {
			log.Println(err)
		}

		var modules []structures.Module
		err = json.Unmarshal([]byte(chatgptResponse), &modules)
		if err != nil {
			log.Println(err)
		}

		m[path] = modules
	}

	return m, nil
}