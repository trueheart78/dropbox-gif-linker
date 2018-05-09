package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/trueheart78/dropbox-gif-linker/clipboard"
	"github.com/trueheart78/dropbox-gif-linker/commands"
	"github.com/trueheart78/dropbox-gif-linker/data"
	"github.com/trueheart78/dropbox-gif-linker/dropbox"
	"github.com/trueheart78/dropbox-gif-linker/messages"
)

var version = "0.6"
var dropboxClient dropbox.Client
var mode = "url"

func url() bool {
	return mode == "url"
}

func md() bool {
	return mode == "md"
}

func init() {
	var err error
	if len(os.Args) >= 2 && os.Args[1] == "version" {
		fmt.Printf("Version %v\n", version)
		os.Exit(0)
	}

	dropboxClient, err = dropbox.DefaultClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(messages.Welcome())
}

func capture(link dropbox.Link) {
	if link != (dropbox.Link{}) {
		fmt.Println(mode)
		if md() {
			clipboard.Write(link.Markdown())
			fmt.Println(messages.LinkTextNew(link.Markdown()))
		} else {
			clipboard.Write(link.DirectLink())
			fmt.Println(messages.LinkTextNew(link.DirectLink()))
		}
	}
	// fmt.Println(messages.LinkTextOld(link.DirectLink()))
}

func main() {
	var link, cachedLink dropbox.Link
	var input, cleaned, help string
	var err error
	reader := bufio.NewReader(os.Stdin)
	handler := data.NewHandler()
	for {
		fmt.Println(messages.AwaitingInput(mode))
		input, _ = reader.ReadString('\n')
		input = strings.Trim(strings.TrimSpace(input), "\"'")
		if commands.Exit(input) {
			fmt.Println(messages.Goodbye())
			break
		} else if commands.URLMode(input) {
			mode = "url"
			fmt.Println(messages.ModeShift("url"))
			capture(link)
		} else if commands.MarkdownMode(input) {
			mode = "md"
			fmt.Println(messages.ModeShift("md"))
			capture(link)
		} else if commands.Help(input) {
			if help == "" {
				help = fmt.Sprintf("Usage: Drag and drop a single gif at a time.\n\n%v", commands.HelpOutput())
			}
			fmt.Println(messages.Help(help))
		} else {
			cleaned, err = handler.Clean(input)
			if err != nil {
				fmt.Printf("Woops! %v\n", err.Error())
				continue
			}
			cachedLink = link
			link, err = dropboxClient.CreateLink(cleaned)
			if err != nil {
				link = cachedLink
				fmt.Printf("Error creating link: %v\n", err.Error())
				continue
			}
			capture(link)
		}
	}
}
