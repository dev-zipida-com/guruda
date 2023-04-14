package git

import (
	"fmt"
	"os/exec"
)

func executeAnotherLanguageFile(language, functionPath, content string) string {

	cmd := exec.Command(language, functionPath, content)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	return string(output)
}
