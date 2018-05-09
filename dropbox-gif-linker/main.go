package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/trueheart78/dropbox-gif-linker/commands"
	"github.com/trueheart78/dropbox-gif-linker/data"
	"github.com/trueheart78/dropbox-gif-linker/dropbox"
	"github.com/trueheart78/dropbox-gif-linker/messages"
)

var dropboxClient dropbox.Client

func init() {
	var err error
	dropboxClient, err = dropbox.DefaultClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(messages.Welcome())
}

func main() {
	var link dropbox.Link
	var input, cleaned string
	var err error
	mode := "url"
	reader := bufio.NewReader(os.Stdin)
	handler := data.NewHandler()
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
			cleaned, err = handler.Clean(input)
			if err != nil {
				fmt.Printf("Woops! %v", err.Error())
				continue
			}
			fmt.Println(cleaned)
			link, err = dropboxClient.CreateLink(cleaned)
			if err != nil {
				fmt.Printf("Error creating link: %v", err.Error())
				continue
			}
			fmt.Println(messages.LinkText(link.DirectLink()))
		}
	}
}
