package messages

import (
	"fmt"

	"github.com/bclicn/color"
)

// Welcome returns a properly formatted greeting
func Welcome(version float64) string {
	versionOutput := color.LightCyan(fmt.Sprintf("v%.2f", version))
	return fmt.Sprintf("%v %v %v %v", heart(), color.LightCyan("Welcome to Dropbox Gif Listener"), versionOutput, heart())
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

// LinkTextOld returns a properly formatted link
func LinkTextOld(text string) string {
	return color.White(fmt.Sprintf("%v", text))
}

// LinkTextNew returns a properly formatted link
func LinkTextNew(text string) string {
	return color.LightYellow(fmt.Sprintf("%v", text))
}

// Help returns a properly formatted line of help text
func Help(text string) string {
	return color.LightYellow(fmt.Sprintf("%v", text))
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
