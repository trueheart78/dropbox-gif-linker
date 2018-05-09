package messages

import (
	"fmt"

	"github.com/bclicn/color"
)

// Welcome returns a properly formatted greeting
func Welcome() string {
	return fmt.Sprintf("%v %v %v", heart(), color.LightCyan("Welcome to Dropbox Gif Listener"), heart())
}

// Goodbye returns a properly formatted goodbye
func Goodbye() string {
	return fmt.Sprintf("%v %v %v", heart(), color.LightCyan("Goodbye"), heart())
}

// AwaitingInput returns an informational message
func AwaitingInput(mode string) string {
	return fmt.Sprintf("%v%v", color.LightPurple("Waiting for input..."), CurrentMode(mode))
}

// CurrentMode returns the current mode
func CurrentMode(mode string) string {
	return fmt.Sprintf("%v%v %v %v", spacing(), note(), color.LightCyan(mode), note())
}

// ModeShift returns the mode shifted to
func ModeShift(mode string) string {
	return color.LightCyan(fmt.Sprintf("♪ mode shifted to %v ♪", mode))
}

// InputError returns a properly formatted error
func InputError(err error) string {
	return color.Red(fmt.Sprintf("Error reading input: %v", err.Error()))
}

// LinkText returns a properly formatted link
func LinkText(text string) string {
	return color.LightRed(fmt.Sprintf("%v", text))
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
