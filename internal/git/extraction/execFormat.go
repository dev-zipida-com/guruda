package git

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/dev-zipida-com/guruda/internal/structures"
)

func executeAnotherLanguageFile(language, functionPath, content string) structures.Response {

	cmd := exec.Command(language, functionPath, content)

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
