package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testError struct {
	s string
}

func (e testError) Error() string {
	return e.s
}

func TestWelcome(t *testing.T) {
	assert := assert.New(t)

	expected := "\x1b[0;31m♥\x1b[0m \x1b[0;34mWelcome to Dropbox Gif Listener\x1b[0m \x1b[0;34mv3.1-rc1\x1b[0m \x1b[0;31m♥\x1b[0m"
	assert.Equal(expected, Welcome(3.1, 1))
	expected = "\x1b[0;31m♥\x1b[0m \x1b[0;34mWelcome to Dropbox Gif Listener\x1b[0m \x1b[0;34mv0.8\x1b[0m \x1b[0;31m♥\x1b[0m"
	assert.Equal(expected, Welcome(0.8, 0))
}

func TestGoodbye(t *testing.T) {
	assert := assert.New(t)

	expected := "\x1b[0;31m♥\x1b[0m \x1b[0;34mGoodbye\x1b[0m \x1b[0;31m♥\x1b[0m"
	assert.Equal(expected, Goodbye())
}

func TestModeShift(t *testing.T) {
	assert := assert.New(t)

	received := ModeShift("url")
	expected := "\x1b[0;34m♪ mode shifted to url ♪\x1b[0m"
	assert.Equal(expected, received)

	received = ModeShift("md")
	expected = "\x1b[0;34m♪ mode shifted to md ♪\x1b[0m"
	assert.Equal(expected, received)
}

func TestAwaitingInput(t *testing.T) {
	assert := assert.New(t)

	received := AwaitingInput("url")
	expected := "\x1b[0;31m♥\x1b[0m \x1b[0;95mWaiting for input\x1b[0m \x1b[0;31m♥\x1b[0m             \x1b[0;31m♪\x1b[0m \x1b[0;34murl\x1b[0m \x1b[0;31m♪\x1b[0m"

	assert.Equal(expected, received)

	received = AwaitingInput("md")
	expected = "\x1b[0;31m♥\x1b[0m \x1b[0;95mWaiting for input\x1b[0m \x1b[0;31m♥\x1b[0m             \x1b[0;31m♪\x1b[0m \x1b[0;34mmd\x1b[0m \x1b[0;31m♪\x1b[0m"

	assert.Equal(expected, received)
}

func TestCurrentMode(t *testing.T) {
	assert := assert.New(t)

	received := CurrentMode("url")
	expected := "             \x1b[0;31m♪\x1b[0m \x1b[0;34murl\x1b[0m \x1b[0;31m♪\x1b[0m"

	assert.Equal(expected, received)
}

func TestInputError(t *testing.T) {
	assert := assert.New(t)

	err := testError{"sample error"}
	received := InputError(err)
	expected := "\x1b[0;31mError reading input: sample error\x1b[0m"

	assert.Equal(expected, received)
}

func TestLinkTextOld(t *testing.T) {
	assert := assert.New(t)

	received := LinkTextOld("sample.gif")
	expected := "\x1b[0;92msample.gif\x1b[0m"
	assert.Equal(expected, received)
}

func TestLinkTextNew(t *testing.T) {
	assert := assert.New(t)

	received := LinkTextNew("sample.gif")
	expected := "\x1b[0;32msample.gif\x1b[0m"
	assert.Equal(expected, received)
}
