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
	mode := "url"
	for {
		fmt.Println(messages.AwaitingInput(mode))
		fmt.Scanln(&input)

		if commands.Exit(input) {
			fmt.Println(messages.Goodbye())
			break
		} else if commands.UrlMode(input) {
			mode = "url"
			fmt.Println(messages.ModeShift("url"))
		} else if commands.MarkdownMode(input) {
			mode = "md"
			fmt.Println(messages.ModeShift("md"))
		} else {
			fmt.Println("// parse content")
		}
	}
}
