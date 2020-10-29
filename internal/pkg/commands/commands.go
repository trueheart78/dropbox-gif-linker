package commands

import (
	"fmt"
	"strings"
)

var exitCommands = [4]string{"exit", "e", "quit", "q"}
var urlCommands = [2]string{"url", "u"}
var markdownCommands = [2]string{"md", "m"}
var bbcodeCommands = [2]string{"bbcode", "b"}
var deleteCommands = [2]string{"delete", "del"}
var configCommands = [2]string{"config", "details"}
var countCommands = [2]string{"count", "gifs"}
var helpCommands = [2]string{"help", "?"}
var versionCommands = [2]string{"version", "v"}
var taylorCommands = [4]string{"taylor", "taylorswift", "taylor swift", "swiftie"}

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

// BbcodeMode returns true if the input is a bbcode mode command
func BbcodeMode(input string) bool {
	return supported(input, bbcodeCommands[:])
}

// Delete returns true if the input is a delete command
func Delete(input string) (exists bool) {
	return supported(input, deleteCommands[:])
}

// Help returns true if the input is a help command
func Help(input string) bool {
	return supported(input, helpCommands[:])
}

// Config returns true if the input is a config command
func Config(input string) bool {
	return supported(input, configCommands[:])
}

// Count returns true if the input is a count command
func Count(input string) bool {
	return supported(input, countCommands[:])
}

// Version returns true if the input is a version command
func Version(input string) bool {
	return supported(input, versionCommands[:])
}

// Taylor returns true if the input is a <3 Taylor Swift <3 command
func Taylor(input string) bool {
	return supported(input, taylorCommands[:])
}

// Any returns true if the input is in any of the commands
func Any(input string) bool {
	var all []string
	for _, v := range urlCommands {
		all = append(all, v)
	}
	for _, v := range markdownCommands {
		all = append(all, v)
	}
	for _, v := range bbcodeCommands {
		all = append(all, v)
	}
	for _, v := range deleteCommands {
		all = append(all, v)
	}
	for _, v := range helpCommands {
		all = append(all, v)
	}
	for _, v := range versionCommands {
		all = append(all, v)
	}
	for _, v := range countCommands {
		all = append(all, v)
	}
	for _, v := range exitCommands {
		all = append(all, v)
	}
	for _, v := range configCommands {
		all = append(all, v)
	}
	for _, v := range taylorCommands {
		all = append(all, v)
	}

	return supported(input, all[:])
}

// HelpOutput outputs the entries for each command
func HelpOutput() string {
	output := "Supported Commands:\n"
	output += fmt.Sprintf(" %v - Shift to URL Mode\n", strings.Join(urlCommands[:], ", "))
	output += fmt.Sprintf(" %v - Shift to Markdown Mode\n", strings.Join(markdownCommands[:], ", "))
	output += fmt.Sprintf(" %v - Shift to BBCode Mode\n", strings.Join(bbcodeCommands[:], ", "))
	output += fmt.Sprintf(" %v - Delete Last Record\n", strings.Join(deleteCommands[:], ", "))
	output += fmt.Sprintf(" %v - Database Record Count\n", strings.Join(countCommands[:], ", "))
	output += fmt.Sprintf(" %v - Loaded Configuration\n", strings.Join(configCommands[:], ", "))
	output += fmt.Sprintf(" %v - Version Details\n", strings.Join(versionCommands[:], ", "))
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
