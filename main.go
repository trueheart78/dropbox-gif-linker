package main

import (
	"fmt"
	"github.com/trueheart78/dropbox-gif-linker/commands"
	"github.com/trueheart78/dropbox-gif-linker/messages"
)

func goodbyeMessage() string {
	return "Love you!"
}

func init() {
	fmt.Println(messages.Welcome())
}

func main() {
	var input string
	for {
		fmt.Scanln(&input)

		if commands.Exit(input) {
			fmt.Println(messages.Goodbye())
			break
		} else if commands.UrlMode(input) {
			fmt.Println(messages.ModeShift("url"))
		} else if commands.MarkdownMode(input) {
			fmt.Println(messages.ModeShift("md"))
		} else {
			fmt.Println("// parse content")
		}
	}
}
