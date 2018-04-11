package main

import (
	"fmt"
	"github.com/trueheart78/dropbox-gif-linker/commands"
)

func goodbyeMessage() string {
	return "Love you!"
}

func init() {
	fmt.Println("Dropbox Gif Listener")
}

func main() {
	var input string
	for {
		fmt.Scanln(&input)

		if commands.Exit(input) {
			fmt.Println(goodbyeMessage())
			break
		}
	}
}
