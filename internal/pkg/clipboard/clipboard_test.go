package clipboard

import (
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

var cachedData string

func TestClipboardWrite(t *testing.T) {
	//cache the data for later because it's touching our clipboard
	cacheClipboard()

	expected := "this data should be on the clipboard"

	// the method under test
	Write(expected)

	received, err := clipboard.ReadAll()

	assert.Equal(t, expected, received)
	assert.Nil(t, err)

	// restore to original state
	restoreClipboard()
}

func TestClipboardRead(t *testing.T) {
	//cache the data for later because it's touching our clipboard
	cacheClipboard()

	expected := "this data should be on the clipboard"
	clipboard.WriteAll(expected)

	// the method under test
	received, err := Read()

	assert.Equal(t, expected, received)
	assert.Nil(t, err)

	// restore to original state
	restoreClipboard()
}

func cacheClipboard() {
	cachedData, _ = clipboard.ReadAll()
	clipboard.WriteAll("")
}

func restoreClipboard() {
	clipboard.WriteAll(cachedData)
}
