package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedCommands(t *testing.T) {
	assert.Equal(t, [5]string{"exit", "ex", "e", "quit", "q"}, exitCommands)
	assert.Equal(t, [2]string{"url", "u"}, urlCommands)
	assert.Equal(t, [3]string{"markdown", "md", "m"}, markdownCommands)
	assert.Equal(t, [4]string{"help", "he", "h", "?"}, helpCommands)
}

func TestExit(t *testing.T) {
	assert := assert.New(t)
	assert.True(Exit("e"))
	assert.True(Exit("ex"))
	assert.True(Exit("exit"))
	assert.True(Exit(":exit"))
	assert.True(Exit(":e"))
	assert.True(Exit("q"))
	assert.True(Exit("quit"))
	assert.True(Exit(":quit"))
	assert.True(Exit(":q"))

	assert.False(Exit("url"))
}

func TestURLMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(URLMode("url"))
	assert.True(URLMode("u"))
	assert.True(URLMode(":url"))
	assert.True(URLMode(":u"))

	assert.False(URLMode("markdown"))
	assert.False(URLMode("exit"))
	assert.False(URLMode("config"))
	assert.False(URLMode("help"))
}

func TestMarkdownMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(MarkdownMode("markdown"))
	assert.True(MarkdownMode("md"))
	assert.True(MarkdownMode("m"))
	assert.True(MarkdownMode(":md"))
	assert.True(MarkdownMode(":m"))

	assert.False(MarkdownMode("url"))
	assert.False(MarkdownMode("exit"))
	assert.False(MarkdownMode("config"))
	assert.False(MarkdownMode("help"))
}

func TestHelp(t *testing.T) {
	assert := assert.New(t)

	assert.True(Help("help"))
	assert.True(Help("he"))
	assert.True(Help("h"))
	assert.True(Help(":help"))
	assert.True(Help(":he"))
	assert.True(Help(":h"))

	assert.False(Help("url"))
	assert.False(Help("markdown"))
	assert.False(Help("exit"))
	assert.False(Help("config"))
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config("config"))
	assert.True(Config("conf"))
	assert.True(Config("cfg"))
	assert.True(Config("c"))

	assert.False(Config("url"))
	assert.False(Config("markdown"))
	assert.False(Config("help"))
	assert.False(Config("exit"))
}

func TestSupported(t *testing.T) {
	commands := [3]string{"no", "noway", "bogus"}

	assert.True(t, supported("no", commands[:]))
	assert.True(t, supported(":no", commands[:]))
	assert.True(t, supported("noway", commands[:]))
	assert.True(t, supported(":noway", commands[:]))
	assert.True(t, supported("bogus", commands[:]))
	assert.True(t, supported(":bogus", commands[:]))

	assert.False(t, supported("awesome", commands[:]))
	assert.False(t, supported(":awesome", commands[:]))
}

func TestHelpOutput(t *testing.T) {
	assert := assert.New(t)

	output := HelpOutput()
	assert.True(strings.HasPrefix(output, "Supported Commands"))
}
