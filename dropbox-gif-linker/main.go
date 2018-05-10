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
	"github.com/trueheart78/dropbox-gif-linker/gif"
	"github.com/trueheart78/dropbox-gif-linker/messages"
	"github.com/trueheart78/dropbox-gif-linker/version"
)

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
	if len(os.Args) >= 2 {
		if os.Args[1] == "version" {
			fmt.Println(version.Full())
			os.Exit(0)
		}

		if os.Args[1] == "md" || os.Args[1] == "markdown" {
			mode = "md"
		}

		if os.Args[1] == "url" {
			mode = "url"
		}
	}

	dropboxClient, err = dropbox.DefaultClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gif.SetDatabasePath(dropboxClient.Config.DatabasePath())
	gif.Init()

	fmt.Println(messages.Welcome(version.Current))
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

func configMessage() string {
	config := "Current Config:\n"
	config += fmt.Sprintf("- Path:      %v\n", dropboxClient.Config.LoadedPath())
	config += fmt.Sprintf("- Gifs Path: %v\n", dropboxClient.Config.FullPath())
	config += fmt.Sprintf("- Database:  %v\n", dropboxClient.Config.DatabasePath())
	config += fmt.Sprintf("- Token:     %v\n", dropboxClient.Config.Token())
	return config
}

func helpMessage() string {
	return fmt.Sprintf("Usage: Drag and drop a single gif at a time.\n\n%v", commands.HelpOutput())
}

func main() {
	var link, cachedLink dropbox.Link
	var input, cleaned string
	var err error
	var id int
	defer gif.Close()
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
			fmt.Println(messages.Help(helpMessage()))
		} else if commands.Config(input) {
			fmt.Println(messages.Help(configMessage()))
		} else {
			gif.Connect()
			id, err = handler.ID(input)
			if err == nil {
				fmt.Printf("Finding a gif for ID %d\n", id)
				continue
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
			gif.Close()
		}
	}
}
