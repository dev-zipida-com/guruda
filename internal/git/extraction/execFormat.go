package git

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/dev-zipida-com/guruda/internal/structures"
)

func executeAnotherLanguageFile(path, language string) structures.Response {
	cmd := exec.Command(language, path)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	var res structures.Response
	err = json.Unmarshal(output, &res)
	if err != nil {
		fmt.Println(err)
	}

	return res
}
