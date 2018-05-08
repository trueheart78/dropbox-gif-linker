package dropbox

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

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
	d, err = createFromConfig(fullConfig)
	if err != nil {
		return
	}
	d.gifDirFix()
	ok, _ := d.Valid()
	if !ok {
		err = fmt.Errorf("please validate the %v file. See README for details", fullConfig)
	}
	return
}

func (c *Config) gifDirFix() {
	if c.GifDir != "" && !strings.HasPrefix(c.GifDir, "/") {
		c.GifDir = fmt.Sprintf("/%v", c.GifDir)
	}
}

// FullPath provides the full dropbox & gifs path
func (c Config) FullPath() string {
	ok, _ := c.Valid()
	if ok {
		return filepath.Join(c.DropboxPath, c.GifDir)
	}
	return ""
}

// Token returns the api token for use in API calls
func (c Config) Token() string {
	ok, _ := c.Valid()
	if ok {
		return c.APIToken
	}
	return ""
}

func (c Config) Valid() (ok bool, err error) {
	if !c.Loaded {
		err = errors.New("the config has yet to be loaded")
		return
	}
	if c.DropboxPath == "" || c.GifDir == "" || c.APIToken == "" {
		err = errors.New("the config is incomplete")
		return
	}
	if !strings.HasPrefix(c.DropboxPath, "~/") && !strings.HasPrefix(c.DropboxPath, "/") {
		err = fmt.Errorf("the dropbox_path should be \"/%v\" instead of \"%v\"", c.DropboxPath, c.DropboxPath)
		return
	}
	if !strings.HasPrefix(c.GifDir, "/") {
		err = fmt.Errorf("the dropbox_gif_dir should be \"/%v\" instead of \"%v\"", c.GifDir, c.GifDir)
		return
	}
	ok = true
	return
}

func (c *Config) load(configFilename string) (ok bool, err error) {
	var raw []byte
	raw, err = ioutil.ReadFile(configFilename)
	if err != nil {
		c.Path = configFilename
		return
	}
	json.Unmarshal(raw, c)
	ok = true
	c.Path = configFilename
	if configExists(configFilename) {
		c.Loaded = true
	}
	return
}

func createFromConfig(configFilename string) (dropbox Config, err error) {
	_, err = dropbox.load(configFilename)
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
