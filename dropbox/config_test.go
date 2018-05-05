package dropbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfigFilename = fixturePath("valid")
var emptyConfigFilename = fixturePath("empty")
var missingConfigFilename = fixturePath("missing")

var dropbox = DropboxConfig{}

func fixturePath(filename string) string {
	workingDir, _ := os.Getwd()
	return fmt.Sprintf("%v.json", filepath.Join(workingDir, "fixtures", filename))
}

func TestConfigPath(t *testing.T) {
	assert := assert.New(t)

	fullConfigPath := configPath(".dgl.json")
	assert.True(strings.HasSuffix(fullConfigPath, configFilename))
	assert.False(strings.HasPrefix(fullConfigPath, configFilename))
}

func TestConfigExists(t *testing.T) {
	assert := assert.New(t)

	assert.True(configExists(validConfigFilename))
	assert.True(configExists(emptyConfigFilename))
	assert.False(configExists(missingConfigFilename))
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	// valid config
	d := DropboxConfig{}
	d.load(validConfigFilename)

	assert.Equal("~/Dropbox", d.FullPath)
	assert.Equal("gifs/", d.GifDir)
	assert.Equal("API_TOKEN", d.APIToken)
	assert.Equal(validConfigFilename, d.ConfigPath)
	assert.True(d.ConfigLoaded)

	// empty config
	d = DropboxConfig{}
	d.load(emptyConfigFilename)

	assert.Equal("", d.FullPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(emptyConfigFilename, d.ConfigPath)
	assert.True(d.ConfigLoaded)

	// missing config
	d = DropboxConfig{}
	d.load(missingConfigFilename)

	assert.Equal("", d.FullPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(missingConfigFilename, d.ConfigPath)
	assert.False(d.ConfigLoaded)
}

func TestValid(t *testing.T) {
	assert := assert.New(t)

	d := createFromConfig(validConfigFilename)
	assert.True(d.valid())

	d = createFromConfig(emptyConfigFilename)
	assert.False(d.valid())
}
