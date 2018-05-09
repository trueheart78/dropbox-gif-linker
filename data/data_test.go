package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var h = NewHandler()

func TestDataClean(t *testing.T) {
	var data string
	var err error

	for _, d := range dirtyData() {
		data, err = h.Clean(d)

		assert.Equal(t, "/path/to/file name.gif", data)
		assert.Nil(t, err)
	}

	badGif := "/sample/not a gif"
	_, err = h.Clean(badGif)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("not a gif [%v]", badGif), err.Error())

	badGif = "/sample/gif.gif gif.gif"
	_, err = h.Clean(badGif)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("multiple gifs detected in %v", badGif), err.Error())
}

func TestDataIsGif(t *testing.T) {
	assert.True(t, h.isGif("sample.gif"))
	assert.True(t, h.isGif("sample.GIF"))
	assert.False(t, h.isGif("sample.gifk"))
	assert.False(t, h.isGif("sample.gf"))
	assert.False(t, h.isGif("sample.if"))
	assert.False(t, h.isGif("sample.jpg"))
	assert.False(t, h.isGif("samplegif"))
}

func TestDataHasApostrophes(t *testing.T) {
	assert.True(t, h.hasApostrophes("'sample'"))
	assert.False(t, h.hasApostrophes("sample'"))
	assert.False(t, h.hasApostrophes("'sample"))
	assert.False(t, h.hasApostrophes("\"sample\""))
}

func TestDataHasQuotes(t *testing.T) {
	assert.True(t, h.hasQuotes("\"sample\""))
	assert.False(t, h.hasQuotes("sample\""))
	assert.False(t, h.hasQuotes("\"sample"))
	assert.False(t, h.hasQuotes("'sample'"))
}

func dirtyData() []string {
	data := make([]string, 0)
	data = append(data, "/path/to/file name.gif")
	data = append(data, "/path/to/file\\ name.gif")
	data = append(data, "\"/path/to/file\\ name.gif\"")
	data = append(data, "'/path/to/file\\ name.gif'")
	data = append(data, "'/path/to/file\\ name.gif' \n")
	data = append(data, "\t'/path/to/file\\ name.gif' \n")
	return data
}
