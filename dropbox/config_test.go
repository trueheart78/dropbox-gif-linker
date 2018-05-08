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
var invalidPathConfigFilename = fixturePath("invalid_path")
var invalidDirConfigFilename = fixturePath("invalid_dir")
var emptyConfigFilename = fixturePath("empty")
var missingConfigFilename = fixturePath("missing")

var dropbox = Config{}

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
	d := Config{}
	d.load(validConfigFilename)

	assert.Equal("~/Dropbox", d.DropboxPath)
	assert.Equal("/gifs", d.GifDir)
	assert.Equal("API_TOKEN", d.APIToken)
	assert.Equal(validConfigFilename, d.Path)
	assert.True(d.Loaded)

	// empty config
	d = Config{}
	d.load(emptyConfigFilename)

	assert.Equal("", d.DropboxPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(emptyConfigFilename, d.Path)
	assert.True(d.Loaded)

	// missing config
	d = Config{}
	d.load(missingConfigFilename)

	assert.Equal("", d.DropboxPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(missingConfigFilename, d.Path)
	assert.False(d.Loaded)
}

func TestGifDirFix(t *testing.T) {
	assert := assert.New(t)

	d := Config{GifDir: "example/"}

	assert.Equal("example/", d.GifDir)
	d.gifDirFix()
	assert.Equal("/example/", d.GifDir)
}

func TestValidConfig(t *testing.T) {
	assert := assert.New(t)

	d, derr := createFromConfig(validConfigFilename)
	ok, err := d.Valid()
	assert.Nil(derr)
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(invalidPathConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_path should be \"/Dropbox/\" instead of \"Dropbox/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_gif_dir should be \"/gifs/\" instead of \"gifs/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	d.gifDirFix()
	ok, err = d.Valid()
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(emptyConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
}
