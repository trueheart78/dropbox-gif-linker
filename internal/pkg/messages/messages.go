package messages

import (
	"fmt"

	"github.com/bclicn/color"
)

// Welcome returns a properly formatted greeting
func Welcome(version string) string {
	return fmt.Sprintf("%v %v %v", cheer(), color.Blue("Welcome to Dropbox Gif Listener v"+version), cheer())
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
	return color.Blue(fmt.Sprintf("🎵 mode shifted to %v 🎵", mode))
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
	return fmt.Sprintf("%v %v %v", cheer(), color.Red(text), cheer())
}

// Sad returns a properly formatted line of disappointed text
func Sad(text string) string {
	return fmt.Sprintf("%v %v %v", skull(), color.Red(text), skull())
}

// Error returns a properly formatted line of error-focused text
func Error(text string, err error) string {
	return Sad(fmt.Sprintf("%v: %v", text, err.Error()))
}

// Info returns a properly formatted line of info-focused text
func Info(text string) string {
	return color.Blue(fmt.Sprintf("🤘🏽 %v 🤘🏽", text))
}

func spacing() string {
	return "             "
}

func heart() string {
	return "💖"
}

func cheer() string {
	return "🎉"
}

func note() string {
	return "🎵"
}

func skull() string {
	return "☠️ "
}
