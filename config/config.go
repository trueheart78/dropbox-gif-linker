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

type DropboxConfig struct {
	FullPath     string `json:"dropbox_path"`
	GifDir       string `json:"dropbox_gif_dir"`
	APIToken     string `json:"dropbox_api_token"`
	ConfigPath   string
	ConfigLoaded bool
}

// LoadConfig attempts to load an existing configuration
func LoadConfig() (d DropboxConfig, err error) {
	fullConfig := configPath(configFilename)
	d = createFromConfig(fullConfig)
	if !d.valid() {
		err = fmt.Errorf("please validate the %v file. See README for details", fullConfig)
	}
	return
}

func (c DropboxConfig) valid() bool {
	if !c.ConfigLoaded || c.FullPath == "" || c.GifDir == "" || c.APIToken == "" {
		return false
	}
	return true
}

func (c *DropboxConfig) load(configFilename string) (ok bool, err error) {
	var raw []byte
	raw, err = ioutil.ReadFile(configFilename)
	if err == nil {
		json.Unmarshal(raw, c)
		ok = true
	}
	c.ConfigPath = configFilename
	if configExists(configFilename) {
		c.ConfigLoaded = true
	}
	return
}

func createFromConfig(configFilename string) (dropbox DropboxConfig) {
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
