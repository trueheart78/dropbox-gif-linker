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
	"github.com/trueheart78/dropbox-gif-linker/taylor"
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

func handleFirstArg(argument string) {
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

func init() {
	var err error
	if len(os.Args) >= 2 {
		handleFirstArg(os.Args[1])
	}

	dropboxClient, err = dropbox.DefaultClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	gif.SetDatabasePath(dropboxClient.Config.DatabasePath())
	_, err = gif.Init()
	if err != nil {
		fmt.Printf("Error initiating database: %v", err.Error())
		os.Exit(1)
	}

	clearScreen()
	fmt.Println(messages.Welcome(version.Current, version.ReleaseCandidate))
}

func convert(link dropbox.Link, checksum string) (newGif gif.Record, err error) {
	if link == (dropbox.Link{}) {
		err = errors.New("invalid link")
		return
	}
	newGif.BaseName = link.Name
	newGif.Directory = link.Directory()
	if !strings.HasPrefix(newGif.Directory, string(os.PathSeparator)) {
		newGif.Directory = fmt.Sprintf("%v%v", string(os.PathSeparator), newGif.Directory)
	}
	newGif.FileSize = link.FileSize
	newGif.Checksum = checksum
	newGif.SharedLinkID = link.DropboxID()
	newGif.SharedLink = gif.RecordSharedLink{
		ID:         link.DropboxID(),
		RemotePath: link.RemotePath(),
		Count:      1,
	}
	return
}

func capture(gifRecord gif.Record, increment bool) {
	if gifRecord != (gif.Record{}) {
		if increment {
			gifRecord.Increment()
		}
		fmt.Println(messages.LinkTextNew(gifRecord.String()))
		if md() {
			fmt.Println(messages.LinkTextNew(gifRecord.Markdown()))
			clipboard.Write(gifRecord.Markdown())
		} else {
			fmt.Println(messages.LinkTextNew(gifRecord.URL()))
			clipboard.Write(gifRecord.URL())
		}
		fmt.Println("")
	}
}

func configMessage() string {
	config := "Current Config:\n"
	config += fmt.Sprintf("- Path:      %v\n", dropboxClient.Config.LoadedPath())
	config += fmt.Sprintf("- Gifs Path: %v\n", dropboxClient.Config.FullPath())
	config += fmt.Sprintf("- Db Path:   %v\n", dropboxClient.Config.DatabasePath())
	config += fmt.Sprintf("- Db Gifs:   %v\n", humanize.Comma(int64(gif.Count())))
	config += fmt.Sprintf("- Token:     %v", dropboxClient.Config.Token())
	return config
}

func helpMessage() string {
	return fmt.Sprintf("Usage: Drag and drop a single gif at a time.\n\n%v", commands.HelpOutput())
}

func handleCommand(input string, gifRecord gif.Record) bool {
	if commands.Exit(input) {
		fmt.Println(messages.Goodbye())
		return false
	} else if commands.URLMode(input) {
		mode = "url"
		fmt.Println(messages.ModeShift("url"))
		capture(gifRecord, false)
	} else if commands.MarkdownMode(input) {
		mode = "md"
		fmt.Println(messages.ModeShift("md"))
		capture(gifRecord, false)
	} else if commands.Help(input) {
		fmt.Println(messages.Help(helpMessage()))
	} else if commands.Config(input) {
		fmt.Println(messages.Help(configMessage()))
	} else if commands.Count(input) {
		fmt.Println(messages.Help(humanize.Comma(int64(gif.Count())) + " total"))
	} else if commands.Version(input) {
		fmt.Println(messages.Help(version.Full()))
	} else if commands.Taylor(input) {
		fmt.Println(taylor.HeadShot())
	}
	return true
}

func main() {
	var link dropbox.Link
	var gifRecord, gifRecordCached gif.Record
	var input, cleaned, md5checksum, shortFilename string
	var err error
	var id int
	var continueOn bool
	defer gif.Disconnect()
	reader := bufio.NewReader(os.Stdin)
	handler := data.NewHandler()
	for {
		gif.Disconnect() // make sure we're always disconnected while awaiting input
		fmt.Println(messages.AwaitingInput(mode))
		input, _ = reader.ReadString('\n')
		input = strings.Trim(strings.TrimSpace(input), "\"'")
		gif.Connect()
		if commands.Any(input) {
			continueOn = handleCommand(input, gifRecord)
			if !continueOn {
				break
			}
		} else {
			// cache the gifRecord, then reset it
			if gifRecord != (gif.Record{}) {
				gifRecordCached = gifRecord
			}
			gifRecord = gif.Record{}

			id, err = handler.ID(input)
			if err == nil {
				gifRecord, err = gif.Find(int64(id))
				if err == nil {
					capture(gifRecord, true)
					continue
				}
			}

			cleaned, err = handler.Clean(input)
			if err != nil {
				fmt.Printf("Woops! %v\n", err.Error())
				continue
			}

			md5checksum, err = handler.MD5Checksum(cleaned)
			if err == nil {
				gifRecord, err = gif.FindByMD5Checksum(md5checksum)
				if err == nil {
					capture(gifRecord, true)
					continue
				}
			}
			shortFilename, err = dropboxClient.Truncate(cleaned)
			if err != nil {
				fmt.Printf("Error truncating filename: %v\n", err.Error())
				continue
			}
			gifRecord, err = gif.FindByFilename(shortFilename)
			if err == nil {
				capture(gifRecord, true)
				continue
			} else {
				link, err = dropboxClient.CreateLink(cleaned)
				if err != nil {
					fmt.Printf("Error creating link: %v\n", err.Error())
					continue
				}
				gifRecord, err = convert(link, md5checksum)
				if err != nil {
					gifRecord = gifRecordCached
					fmt.Printf("Error converting link: %v\n", err.Error())
					continue
				}
				_, err := gifRecord.Save()
				if err != nil {
					gifRecord = gifRecordCached
					fmt.Printf("Error saving gif: %v\n", err.Error())
					continue
				}

				capture(gifRecord, false)
			}
		}
	}
}
