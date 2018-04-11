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
}
