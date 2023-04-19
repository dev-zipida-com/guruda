package git

import (
	"github.com/dev-zipida-com/guruda/internal/structures"
)

// paths : path(페이지의 경로)와 code(그 페이지의 코드 전체)
// GetDependencies 함수는 Paths를 받아와서 function 사이의 ExtractedObject 를 정의 및 반환합니다.
// ExtractedObject는 함수의 코드(Code), 경로(Path), 설명(Description), 간단한 설명(DescriptionInShort), 참조(References)를 포함합니다.
// 참조는 이 함수가 호출하는 함수들을 가리킵니다.
// References는 이 함수가 참조하는 다른 함수의 이름(Name), 다른 함수가 참조하는 또 다른 함수의 이름(Imported), 다른 함수가 위치한 페이지의 상대 경로(Path)를 포함합니다. 만약 함수가 참조하는 다른 함수가 없다면 References는 nil이 됩니다.
// 만약 References가 nil이거나 모두 외부 함수를 참조하고 있다면(즉 내가 만든 함수를 하나도 참조하고 있지 않다면), 이 함수는 말단에 위치한 함수이므로, chatgpt를 호출하여 이 함수에 대한 자연어 설명을 가져옵니다.
// 이 함수는 내가 만든 함수를 참조하고 있는 함수들을 찾아내고, 그 함수들을 다시 이 함수에게 넘겨주는 방식으로 재귀적으로 호출됩니다.
func GetDependencies(paths map[string]string) (map[string]structures.ExtractedObject) {
	f := make(map[string]structures.ExtractedObject)


	
}
