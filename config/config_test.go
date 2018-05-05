package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfigFilename = fixturePath("valid")
var invalidConfigFilename = fixturePath("non-existent")

func fixturePath(filename string) string {
	x, _ := os.Getwd()
	return fmt.Sprintf("%v.json", filepath.Join(x, "fixtures", filename))
}

func TestConfigFilename(t *testing.T) {
	assert := assert.New(t)

	expected := Dropbox.configPath(".dgl.json")
	received := configFilename
	assert.Equal(expected, received)
}

func TestExists(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(true, Dropbox.Exists(validConfigFilename))
	assert.Equal(false, Dropbox.Exists(invalidConfigFilename))
}

func TestDropboxConfig(t *testing.T) {
	assert := assert.New(t)

	Dropbox.load(validConfigFilename)

	assert.Equal("~/Dropbox", Dropbox.FullPath)
	assert.Equal("gifs/", Dropbox.GifDir)
}
