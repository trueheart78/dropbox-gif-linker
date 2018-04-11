package messages

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type sampleError struct {
	s string
}

func (e sampleError) Error() string {
	return e.s
}

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
	expected := "\x1b[0;96m♪ mode shifted to url ♪\x1b[0m"
	assert.Equal(expected, received)

	received = ModeShift("md")
	expected = "\x1b[0;96m♪ mode shifted to md ♪\x1b[0m"
	assert.Equal(expected, received)
}

func TestAwaitingInput(t *testing.T) {
	assert := assert.New(t)

	received := AwaitingInput("url")
	expected := "\x1b[0;95mWaiting for input...\x1b[0m               \x1b[0;96m♪\x1b[0m \x1b[0;96murl\x1b[0m \x1b[0;96m♪\x1b[0m"

	assert.Equal(expected, received)

	received = AwaitingInput("md")
	expected = "\x1b[0;95mWaiting for input...\x1b[0m               \x1b[0;96m♪\x1b[0m \x1b[0;96mmd\x1b[0m \x1b[0;96m♪\x1b[0m"

	assert.Equal(expected, received)
}

func TestCurrentMode(t *testing.T) {
	assert := assert.New(t)

	received := CurrentMode("url")
	expected := "               \x1b[0;96m♪\x1b[0m \x1b[0;96murl\x1b[0m \x1b[0;96m♪\x1b[0m"

	assert.Equal(expected, received)
}

func TestInputError(t *testing.T) {
	assert := assert.New(t)

	err := sampleError{"sample error"}
	received := InputError(err)
	expected := "\x1b[0;31mError reading input: sample error\x1b[0m"

	assert.Equal(expected, received)
}
