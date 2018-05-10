package commands

import (
	"fmt"
	"strings"
)

var exitCommands = [5]string{"exit", "ex", "e", "quit", "q"}
var urlCommands = [2]string{"url", "u"}
var markdownCommands = [3]string{"markdown", "md", "m"}
var configCommands = [4]string{"config", "conf", "cfg", "c"}
var helpCommands = [4]string{"help", "he", "h", "?"}

// Exit returns true if the input is an exit command
func Exit(input string) (exists bool) {
	return supported(input, exitCommands[:])
}

// URLMode returns true if the input is a url mode command
func URLMode(input string) (exists bool) {
	return supported(input, urlCommands[:])
}

// MarkdownMode returns true if the input is a markdown mode command
func MarkdownMode(input string) bool {
	return supported(input, markdownCommands[:])
}

// Help returns true if the input is a help command
func Help(input string) bool {
	return supported(input, helpCommands[:])
}

// Config returns true if the input is a help command
func Config(input string) bool {
	return supported(input, configCommands[:])
}

// HelpOutput outputs the entries for each command
func HelpOutput() string {
	output := "Supported Commands:\n"
	output += fmt.Sprintf(" %v - Shift to URL Mode\n", strings.Join(urlCommands[:], ", "))
	output += fmt.Sprintf(" %v - Shift to Markdown Mode\n", strings.Join(markdownCommands[:], ", "))
	output += fmt.Sprintf(" %v - Loaded Configuration\n", strings.Join(configCommands[:], ", "))
	output += fmt.Sprintf(" %v - Exit Program\n", strings.Join(exitCommands[:], ", "))
	output += fmt.Sprintf(" %v - Help (This Menu)\n", strings.Join(helpCommands[:], ", "))

	return output
}

// returns whether the passed input (or a variant) exists in the commands slice
func supported(input string, commands []string) (exists bool) {
	if strings.HasPrefix(input, ":") {
		input = strings.Replace(input, ":", "", 1)
	}
	for _, k := range commands {
		if input == k {
			exists = true
			break
		}
	}
	return
}
