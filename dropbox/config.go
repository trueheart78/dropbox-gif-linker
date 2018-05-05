package dropbox

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var configFilename = ".dgl.json"

// Config is the object to be used when working with Client
type Config struct {
	DropboxPath string `json:"dropbox_path"`
	GifDir      string `json:"dropbox_gif_dir"`
	APIToken    string `json:"dropbox_api_token"`
	Path        string
	Loaded      bool
}

// NewConfig attempts to load an existing configuration
func NewConfig() (d Config, err error) {
	fullConfig := configPath(configFilename)
	d = createFromConfig(fullConfig)
	if !d.valid() {
		err = fmt.Errorf("please validate the %v file. See README for details", fullConfig)
	}
	return
}

// FullPath provides the full dropbox & gifs path
func (c Config) FullPath() string {
	if c.valid() {
		return filepath.Join(c.DropboxPath, c.GifDir)
	}
	return ""
}

func (c Config) valid() bool {
	if !c.Loaded || c.DropboxPath == "" || c.GifDir == "" || c.APIToken == "" {
		return false
	}
	return true
}

func (c *Config) load(configFilename string) (ok bool, err error) {
	var raw []byte
	raw, err = ioutil.ReadFile(configFilename)
	if err == nil {
		json.Unmarshal(raw, c)
		ok = true
	}
	c.Path = configFilename
	if configExists(configFilename) {
		c.Loaded = true
	}
	return
}

func createFromConfig(configFilename string) (dropbox Config) {
	_, err := dropbox.load(configFilename)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	return
}

func configExists(configFilename string) bool {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return false
	}
	return true
}

func configPath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, filename)
}
