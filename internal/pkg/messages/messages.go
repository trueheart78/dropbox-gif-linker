package messages

import (
	"fmt"

	"github.com/bclicn/color"
)

// Welcome returns a properly formatted greeting
func Welcome(version float64, releaseCandidate int) string {
	var rc string
	if releaseCandidate > 0 {
		rc = fmt.Sprintf("-rc%d", releaseCandidate)
	}
	versionOutput := color.Blue(fmt.Sprintf("v%.1f%v", version, rc))
	return fmt.Sprintf("%v %v %v %v", heart(), color.Blue("Welcome to Dropbox Gif Listener"), versionOutput, heart())
}

// Goodbye returns a properly formatted goodbye
func Goodbye() string {
	return fmt.Sprintf("%v %v %v", heart(), color.Blue("Goodbye"), heart())
}

// AwaitingInput returns an informational message
func AwaitingInput(mode string) string {
	return fmt.Sprintf("%v %v %v%v", heart(), color.LightPurple("Waiting for input"), heart(), CurrentMode(mode))
}

// CurrentMode returns the current mode
func CurrentMode(mode string) string {
	return fmt.Sprintf("%v%v %v %v", spacing(), note(), color.Blue(mode), note())
}

// ModeShift returns the mode shifted to
func ModeShift(mode string) string {
	return color.Blue(fmt.Sprintf("♪ mode shifted to %v ♪", mode))
}

// InputError returns a properly formatted error
func InputError(err error) string {
	return color.Red(fmt.Sprintf("Error reading input: %v", err.Error()))
}

// LinkTextOld returns a properly formatted link
func LinkTextOld(text string) string {
	return color.LightGreen(fmt.Sprintf("%v", text))
}

// LinkTextNew returns a properly formatted link
func LinkTextNew(text string) string {
	return color.Green(fmt.Sprintf("%v", text))
}

// Help returns a properly formatted line of help text
func Help(text string) string {
	return color.Green(fmt.Sprintf("%v\n", text))
}

// Happy returns a properly formatted line of disappointed text
func Happy(text string) string {
	return fmt.Sprintf("%v %v %v", heart(), color.Red(text), heart())
}

// Sad returns a properly formatted line of disappointed text
func Sad(text string) string {
	return fmt.Sprintf("%v %v %v", skull(), color.Red(text), skull())
}

func spacing() string {
	return "             "
}

func heart() string {
	return color.Red("💖")
}

func note() string {
	return color.Red("♪")
}

func skull() string {
	return color.Red("☠️ ")
}
