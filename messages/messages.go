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

func ModeShift(mode string) string {
	return fmt.Sprintf("               %v %v %v", note(), color.LightCyan(mode), note())
}

func heart() string {
	return color.Red("♥")
}

func note() string {
	return color.LightCyan("♪")
}
