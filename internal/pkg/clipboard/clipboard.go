package clipboard

import (
	"github.com/atotto/clipboard"
)

// Read returns the data on the clipboard
func Read() (data string, err error) {
	data, err = clipboard.ReadAll()
	return
}

// Write writes data to the clipboard
func Write(output string) (ok bool, err error) {
	err = clipboard.WriteAll(output)
	if err != nil {
		ok = true
	}
	return
}
