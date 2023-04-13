package git

import (
	"github.com/dev-zipida-com/guruda/internal/structures"
)

func GetDependencies(paths map[string]string, modules map[string][]structures.Module) (ExtractedObject structures.ExtractedObject) {
	f := make(map[string]structures.ExtractedObject)

	for path, code := range paths {
		objects := modules[path]

	}
}