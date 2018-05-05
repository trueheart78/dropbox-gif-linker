package main

import (
	"bufio"
	"fmt"
	"github.com/trueheart78/dropbox-gif-linker/commands"
	"github.com/trueheart78/dropbox-gif-linker/messages"
	"os"
	"strings"
)

func init() {
	fmt.Println(messages.Welcome())
}

func main() {
	var input string
	mode := "url"
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(messages.AwaitingInput(mode))
		input, _ = reader.ReadString('\n')
		input = strings.Trim(strings.TrimSpace(input), "\"'")
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
			fmt.Printf("You entered: %v\n", input)
		}
	}
}
