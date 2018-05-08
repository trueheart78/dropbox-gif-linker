package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNil description
func TestDataIsGif(t *testing.T) {
	h := NewHandler()

	assert.True(t, h.isGif("sample.gif"))
	assert.True(t, h.isGif("sample.GIF"))
	assert.False(t, h.isGif("sample.gifk"))
	assert.False(t, h.isGif("sample.gf"))
	assert.False(t, h.isGif("sample.if"))
	assert.False(t, h.isGif("sample.jpg"))
	assert.False(t, h.isGif("samplegif"))
}
