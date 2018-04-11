package commands

var commands = make(map[string]string)

func init() {
	commands["exit"] = "exit"
	commands[":exit"] = "exit"
	commands["ex"] = "exit"
	commands["e"] = "exit"
	commands[":e"] = "exit"
	commands["quit"] = "exit"
	commands["q"] = "exit"
	commands[":quit"] = "exit"
	commands[":q"] = "exit"

	commands["url"] = "url"
	commands["u"] = "url"
	commands[":url"] = "url"
	commands[":u"] = "url"

	commands["markdown"] = "md"
	commands["md"] = "md"
	commands["m"] = "md"
	commands[":md"] = "md"
	commands[":m"] = "md"
}

func Exit(input string) bool {
	if commands[input] == "exit" {
		return true
	}
	return false
}

func UrlMode(input string) bool {
	if commands[input] == "url" {
		return true
	}
	return false
}

func MarkdownMode(input string) bool {
	if commands[input] == "md" {
		return true
	}
	return false
}
