package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

var configFilename = ".dgl.json"

// Dropbox config
var Dropbox = config{}

type config struct {
	FullPath string `json:"dropbox_path"`
	GifDir   string `json:"dropbox_gif_dir"`
	APIToken string `json:"dropbox_api_token"`
}

type configInterface interface {
	configPath() string
	Exists() bool
}

func (c *config) load(configFilename string) {
	raw, err := ioutil.ReadFile(configFilename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	json.Unmarshal(raw, c)
}

func (c config) Exists(configFilename string) bool {
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c config) configPath(filename string) string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, filename)
}

func init() {
	configFilename = Dropbox.configPath(configFilename)
	if Dropbox.Exists(configFilename) {
		Dropbox.load(configFilename)
	} else {
		fmt.Printf("Please create the %v file. See README for details.\n", Dropbox.configPath(configFilename))
		os.Exit(1)
	}
}
