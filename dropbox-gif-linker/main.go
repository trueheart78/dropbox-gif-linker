package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	humanize "github.com/dustin/go-humanize"
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

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
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

	clearScreen()
	fmt.Println(messages.Welcome(version.Current))
}

func convert(link dropbox.Link) (newGif gif.Record, err error) {
	if link == (dropbox.Link{}) {
		err = errors.New("invalid link")
		return
	}
	newGif.BaseName = link.Name
	newGif.Directory = link.Directory()
	newGif.FileSize = link.FileSize
	newGif.SharedLinkID = link.DropboxID()
	newGif.SharedLink = gif.RecordSharedLink{
		ID:         link.DropboxID(),
		RemotePath: link.RemotePath(),
		Count:      1,
	}
	return
}

func capture(gifRecord gif.Record) {
	if gifRecord != (gif.Record{}) {
		fmt.Println(messages.LinkTextNew(gifRecord.String()))
		if md() {
			fmt.Println(messages.LinkTextNew(gifRecord.Markdown()))
			clipboard.Write(gifRecord.Markdown())
		} else {
			fmt.Println(messages.LinkTextNew(gifRecord.URL()))
			clipboard.Write(gifRecord.URL())
		}
	}
}

func configMessage() string {
	config := "Current Config:\n"
	config += fmt.Sprintf("- Path:      %v\n", dropboxClient.Config.LoadedPath())
	config += fmt.Sprintf("- Gifs Path: %v\n", dropboxClient.Config.FullPath())
	config += fmt.Sprintf("- Db Path:   %v\n", dropboxClient.Config.DatabasePath())
	config += fmt.Sprintf("- Db Gifs:   %v\n", humanize.Comma(int64(gif.Count())))
	config += fmt.Sprintf("- Token:     %v\n", dropboxClient.Config.Token())
	return config
}

func helpMessage() string {
	return fmt.Sprintf("Usage: Drag and drop a single gif at a time.\n\n%v", commands.HelpOutput())
}

func main() {
	var link dropbox.Link
	var gifRecord, gifRecordCached gif.Record
	var input, cleaned string
	var err error
	var id int
	defer gif.Disconnect()
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
			capture(gifRecord)
		} else if commands.MarkdownMode(input) {
			mode = "md"
			fmt.Println(messages.ModeShift("md"))
			capture(gifRecord)
		} else if commands.Help(input) {
			fmt.Println(messages.Help(helpMessage()))
		} else if commands.Config(input) {
			fmt.Println(messages.Help(configMessage()))
		} else {
			gif.Connect()
			id, err = handler.ID(input)
			if err == nil {
				fmt.Printf("[Not Implemented]: Finding a gif for ID %d\n", id)
				continue
			} else {
				cleaned, err = handler.Clean(input)
				if err != nil {
					fmt.Printf("Woops! %v\n", err.Error())
					continue
				}
				// TODO: let's check with the Gifs table to see if it already exists!

				link, err = dropboxClient.CreateLink(cleaned)
				if err != nil {
					fmt.Printf("Error creating link: %v\n", err.Error())
					continue
				}
				if gifRecord != (gif.Record{}) {
					gifRecordCached = gifRecord
				}
				gifRecord, err = convert(link)
				if err != nil {
					gifRecord = gifRecordCached
					fmt.Printf("Error converting link: %v\n", err.Error())
					continue
				}
				capture(gifRecord)
			}
			gif.Disconnect()
		}
	}
}
