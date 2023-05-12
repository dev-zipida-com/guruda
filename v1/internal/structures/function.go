package structures

type References struct {
	Name string
	Imported []string
	Path string
}

type Module struct {
	Name string
	Imported []string
}

type ExtractedObject struct {
	Code string
	Path string
	Description string
	DescriptionInShort string
	References References
}
