package commands

var exitCommands = make(map[string]string)

func init() {
	exitCommands["exit"] = "exit"
	exitCommands[":exit"] = "exit"
	exitCommands["ex"] = "exit"
	exitCommands["e"] = "exit"
	exitCommands[":e"] = "exit"
	exitCommands["quit"] = "exit"
	exitCommands["q"] = "exit"
	exitCommands[":quit"] = "exit"
	exitCommands[":q"] = "exit"
}

func Exit(input string) bool {
	if exitCommands[input] == "exit" {
		return true
	}
	return false
}
