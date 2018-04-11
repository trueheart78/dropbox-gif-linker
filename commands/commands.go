package commands

var exitCommands = make(map[string]bool)
var urlCommands = make(map[string]bool)
var markdownCommands = make(map[string]bool)

func init() {
	exitCommands["exit"] = true
	exitCommands[":exit"] = true
	exitCommands["ex"] = true
	exitCommands["e"] = true
	exitCommands[":e"] = true
	exitCommands["quit"] = true
	exitCommands["q"] = true
	exitCommands[":quit"] = true
	exitCommands[":q"] = true

	urlCommands["url"] = true
	urlCommands["u"] = true
	urlCommands[":url"] = true
	urlCommands[":u"] = true

	markdownCommands["markdown"] = true
	markdownCommands["md"] = true
	markdownCommands["m"] = true
	markdownCommands[":md"] = true
	markdownCommands[":m"] = true
}

func Exit(input string) bool {
	_, exist := exitCommands[input]
	return exist
}

func UrlMode(input string) bool {
	_, exist := urlCommands[input]
	return exist
}

func MarkdownMode(input string) bool {
	_, exist := markdownCommands[input]
	return exist
}
