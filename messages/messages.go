package messages

import (
	"fmt"
	"github.com/bclicn/color"
)

func Welcome() string {
	return fmt.Sprintf("%v %v %v", heart(), color.LightCyan("Welcome to Dropbox Gif Listener"), heart())
}

func Goodbye() string {
	return fmt.Sprintf("%v %v %v", heart(), color.LightCyan("Goodbye"), heart())
}

func AwaitingInput(mode string) string {
	return fmt.Sprintf("%v%v", color.LightPurple("Waiting for input..."), CurrentMode(mode))
}

func CurrentMode(mode string) string {
	return fmt.Sprintf("%v%v %v %v", spacing(), note(), color.LightCyan(mode), note())
}

func ModeShift(mode string) string {
	return color.LightCyan(fmt.Sprintf("♪ mode shifted to %v ♪", mode))
}

func InputError(err error) string {
	return color.Red(fmt.Sprintf("Error reading input: %v", err))
}

func spacing() string {
	return "               "
}

func heart() string {
	return color.Red("♥")
}

func note() string {
	return color.LightCyan("♪")
}
