package messages

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWelcome(t *testing.T) {
	assert := assert.New(t)

	expected := "\x1b[0;31m♥\x1b[0m \x1b[0;96mWelcome to Dropbox Gif Listener\x1b[0m \x1b[0;31m♥\x1b[0m"
	assert.Equal(expected, Welcome())
}

func TestGoodbye(t *testing.T) {
	assert := assert.New(t)

	expected := "\x1b[0;31m♥\x1b[0m \x1b[0;96mGoodbye\x1b[0m \x1b[0;31m♥\x1b[0m"
	assert.Equal(expected, Goodbye())
}

func TestModeShift(t *testing.T) {
	assert := assert.New(t)

	received := ModeShift("url")
	expected := "               \x1b[0;96m♪\x1b[0m \x1b[0;96murl\x1b[0m \x1b[0;96m♪\x1b[0m"
	assert.Equal(expected, received)

	received = ModeShift("md")
	expected = "               \x1b[0;96m♪\x1b[0m \x1b[0;96mmd\x1b[0m \x1b[0;96m♪\x1b[0m"
	assert.Equal(expected, received)
}
