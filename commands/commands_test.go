package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

func TestUrlMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(UrlMode("url"))
	assert.True(UrlMode("u"))
	assert.True(UrlMode(":url"))
	assert.True(UrlMode(":u"))

	assert.False(UrlMode("exit"))
}

func TestMarkdownMode(t *testing.T) {
	assert := assert.New(t)

	assert.True(MarkdownMode("markdown"))
	assert.True(MarkdownMode("md"))
	assert.True(MarkdownMode("m"))
	assert.True(MarkdownMode(":md"))
	assert.True(MarkdownMode(":m"))

	assert.False(MarkdownMode("exit"))
}
