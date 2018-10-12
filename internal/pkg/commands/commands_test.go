package commands

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupportedCommands(t *testing.T) {
	assert.Equal(t, [4]string{"exit", "e", "quit", "q"}, exitCommands)
	assert.Equal(t, [2]string{"url", "u"}, urlCommands)
	assert.Equal(t, [2]string{"md", "m"}, markdownCommands)
	assert.Equal(t, [2]string{"delete", "del"}, deleteCommands)
	assert.Equal(t, [2]string{"version", "v"}, versionCommands)
	assert.Equal(t, [2]string{"help", "?"}, helpCommands)
	assert.Equal(t, [2]string{"config", "details"}, configCommands)
	assert.Equal(t, [2]string{"count", "gifs"}, countCommands)
	assert.Equal(t, [4]string{"taylor", "taylorswift", "taylor swift", "swiftie"}, taylorCommands)
}

func TestExit(t *testing.T) {
	assert := assert.New(t)
	assert.True(Exit("e"))
	assert.True(Exit("exit"))
	assert.True(Exit(":exit"))
	assert.True(Exit(":e"))
	assert.True(Exit("q"))
	assert.True(Exit("quit"))
	assert.True(Exit(":quit"))
	assert.True(Exit(":q"))

	assert.False(Exit("url"))
	assert.False(Exit("md"))
	assert.False(Exit("delete"))
	assert.False(Exit("config"))
	assert.False(Exit("help"))
	assert.False(Exit("count"))
	assert.False(Exit("version"))
	assert.False(Exit("taylor"))
}

func TestURLMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(URLMode("url"))
	assert.True(URLMode("u"))
	assert.True(URLMode(":url"))
	assert.True(URLMode(":u"))

	assert.False(URLMode("md"))
	assert.False(URLMode("exit"))
	assert.False(URLMode("delete"))
	assert.False(URLMode("config"))
	assert.False(URLMode("help"))
	assert.False(URLMode("count"))
	assert.False(URLMode("version"))
	assert.False(URLMode("taylor"))
}

func TestMarkdownMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(MarkdownMode("md"))
	assert.True(MarkdownMode("m"))
	assert.True(MarkdownMode(":md"))
	assert.True(MarkdownMode(":m"))

	assert.False(MarkdownMode("url"))
	assert.False(MarkdownMode("exit"))
	assert.False(MarkdownMode("delete"))
	assert.False(MarkdownMode("config"))
	assert.False(MarkdownMode("help"))
	assert.False(MarkdownMode("count"))
	assert.False(MarkdownMode("version"))
	assert.False(MarkdownMode("taylor"))
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)

	assert.True(Delete("delete"))
	assert.True(Delete("del"))
	assert.True(Delete(":delete"))
	assert.True(Delete(":del"))

	assert.False(Delete("url"))
	assert.False(Delete("md"))
	assert.False(Delete("help"))
	assert.False(Delete("exit"))
	assert.False(Delete("config"))
	assert.False(Delete("count"))
	assert.False(Delete("version"))
	assert.False(Delete("taylor"))
}

func TestHelp(t *testing.T) {
	assert := assert.New(t)

	assert.True(Help("help"))
	assert.True(Help("?"))
	assert.True(Help(":help"))
	assert.True(Help(":?"))

	assert.False(Help("url"))
	assert.False(Help("md"))
	assert.False(Help("exit"))
	assert.False(Help("delete"))
	assert.False(Help("config"))
	assert.False(Help("count"))
	assert.False(Help("version"))
	assert.False(Help("taylor"))
}

func TestConfig(t *testing.T) {
	assert := assert.New(t)

	assert.True(Config("config"))
	assert.True(Config("details"))

	assert.False(Config("url"))
	assert.False(Config("md"))
	assert.False(Config("help"))
	assert.False(Config("delete"))
	assert.False(Config("exit"))
	assert.False(Config("count"))
	assert.False(Config("version"))
	assert.False(Config("taylor"))
}

func TestCount(t *testing.T) {
	assert := assert.New(t)

	assert.True(Count("count"))

	assert.False(Count("url"))
	assert.False(Count("md"))
	assert.False(Count("help"))
	assert.False(Count("delete"))
	assert.False(Count("exit"))
	assert.False(Count("config"))
	assert.False(Count("version"))
	assert.False(Count("taylor"))
}

func TestVersion(t *testing.T) {
	assert := assert.New(t)

	assert.True(Version("version"))

	assert.False(Version("url"))
	assert.False(Version("md"))
	assert.False(Version("help"))
	assert.False(Version("delete"))
	assert.False(Version("exit"))
	assert.False(Version("config"))
	assert.False(Version("count"))
	assert.False(Version("taylor"))
}

func TestTaylor(t *testing.T) {
	assert := assert.New(t)

	assert.True(Taylor("taylor"))

	assert.False(Taylor("url"))
	assert.False(Taylor("md"))
	assert.False(Taylor("help"))
	assert.False(Taylor("delete"))
	assert.False(Taylor("exit"))
	assert.False(Taylor("config"))
	assert.False(Taylor("count"))
	assert.False(Taylor("version"))
}

func TestAny(t *testing.T) {
	assert.True(t, Any("url"))
	assert.True(t, Any("md"))
	assert.True(t, Any("help"))
	assert.True(t, Any("delete"))
	assert.True(t, Any("exit"))
	assert.True(t, Any("config"))
	assert.True(t, Any("count"))
	assert.True(t, Any("version"))
	assert.True(t, Any("taylor"))
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
