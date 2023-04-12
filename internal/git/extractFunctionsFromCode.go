package git

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/dev-zipida-com/guruda/internal/chatgpt"
	"github.com/dev-zipida-com/guruda/internal/prompts"
	"github.com/dev-zipida-com/guruda/internal/structures"
)

func extract (code string, extension string) []structures.ExtractedObject {
	var objects []structures.ExtractedObject
	prompt := prompts.ExtractObjects(code)

	res, err := chatgpt.GetResponse(prompt)
	if err != nil {
		log.Fatal(err)
	}
	
	var resJson string
	err := json.Unmarshal([]byte(res), &resJson)
	if err != nil {
		log.Fatal(err)
	}
	
	for key, value := range resJson {
		var object = structures.ExtractedObject{
			Name: ,
			Code: ,
			FileRoute: "",
			Description: "",
			References: , // 이 객체가 참조하고 있는 함수, 클래스 등의 이름
		}

		objects = append(objects, object)
	}

	return objects
}

func ExtractFuctionsFromCode(paths map[string]string) {
	for path, code := range paths {
		slash := strings.Split(path, "/")
		extension := slash[len(slash)-1]

		objects := extract(code, extension)

	}

}