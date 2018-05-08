package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

// Client for the Dropbox API interactions
type Client struct {
	Host    string
	Version int
	Config  clientConfig
}

type clientConfig interface {
	FullPath() string
	GifsPath() string
	Token() string
	Valid() bool
}

type existingPayload struct {
	RelativePath string `json:"path"`
}

type creationPayload struct {
	RelativePath string         `json:"path"`
	Settings     settingPayload `json:"settings"`
}

type settingPayload struct {
	Visibility string `json:"requested_visibility"`
}

type existsResponse struct {
	Links   []Link `json:"links"`
	HasMore bool   `json:"has_more"`
}

// Link is the data that is provided from the Dropbox API
type Link struct {
	Tag            string          `json:".tag"`
	URL            string          `json:"url"`
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Path           string          `json:"path_lower"`
	Permissions    LinkPermissions `json:"link_permissions"`
	ClientModified string          `json:"client_modified"`
	ServerModified string          `json:"server_modified"`
	Revision       string          `json:"rev"`
	FileSize       int             `json:"size"`
}

// LinkPermissions are the permissions that Dropbox as assigned
type LinkPermissions struct {
	ResolvedVisibility  LinkTag `json:"resolved_visibility"`
	RequestedVisibility LinkTag `json:"requested_visibility"`
	CanRevoke           bool    `json:"can_revoke"`
}

// LinkTag is a tag assgined to a link
type LinkTag struct {
	Tag string `json:".tag"`
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

// NewClient creates a new Client for interacting with Dropbox
func NewClient(config clientConfig) (c Client) {
	c.Host = "https://api.dropboxapi.com"
	c.Version = 2
	c.Config = config
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

// GifsPath provides the full dropbox & gifs path
func (c Config) GifsPath() string {
	ok, _ := c.Valid()
	if ok {
		return c.GifDir
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

// Valid returns whether the config is valid
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

func (c Client) valid() bool {
	if !c.Config.Valid() || c.Host == "" || c.Version == 0 {
		return false
	}
	return true
}

func (c Client) basicRequest(fullURL string, payload bytes.Buffer) (result *http.Response, err error) {
	request, err := http.NewRequest(http.MethodPost, fullURL, &payload)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.Config.Token()))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Dropbox Gif Linker")
	return http.DefaultClient.Do(request)
}

func (c Client) exists(filename string) (link Link, err error) {
	if !c.valid() {
		err = errors.New("client is not valid")
		return
	}

	payload := c.existingPayload(filename)
	fullURL := c.existingURL()
	result, err := c.basicRequest(fullURL, payload)

	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("dropbox returned a %d", result.StatusCode)
		return
	}

	var rawBody []byte
	var exists existsResponse
	rawBody, err = ioutil.ReadAll(result.Body)
	defer result.Body.Close()
	if err == nil {
		json.Unmarshal(rawBody, &exists)
		if len(exists.Links) > 0 {
			for _, l := range exists.Links {
				if l.Path == filename {
					link = l
					return
				}
			}
			err = fmt.Errorf("no existing link for %v", filename)
		} else {
			err = fmt.Errorf("no existing link for %v", filename)
		}
	}
	return
}

func (c Client) create(filename string) (link Link, err error) {
	if !c.valid() {
		err = errors.New("client is not valid")
		return
	}

	payload := c.creationPayload(filename)
	fullURL := c.creationURL()
	result, err := c.basicRequest(fullURL, payload)

	if err != nil {
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("dropbox returned a %d", result.StatusCode)
		return
	}

	var rawBody []byte
	rawBody, err = ioutil.ReadAll(result.Body)
	defer result.Body.Close()
	if err == nil {
		json.Unmarshal(rawBody, &link)
	}
	return
}

// DirectLink returns the embeddable string
// From: https://www.dropbox.com/s/eqoo012hoa0wq7k/taylor%20bat%20focused.gif?dl=0
// To:   https://dl.dropboxusercontent.com/s/eqoo012hoa0wq7k/taylor%20bat%20focused.gif
func (e Link) DirectLink() string {
	u, err := url.Parse(e.URL)
	if err != nil {
		panic(err)
	}
	// change the host to point directly to the content
	u.Host = "dl.dropboxusercontent.com"
	// remove the dl=0 query
	u.RawQuery = ""
	return u.String()
}

func (c Client) fixFilename(filename string) string {
	if !strings.HasPrefix(filename, c.Config.GifsPath()) {
		return filepath.Join(c.Config.GifsPath(), filename)
	}
	return filename
}

func (c Client) existingPayload(filename string) (buf bytes.Buffer) {
	//fmt.Println("before:", filename)
	//filename = c.fixFilename(filename)
	//fmt.Println("after:", filename)
	payload := existingPayload{filename}
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		fmt.Printf("There was an error encoding the json. err = %s", err)
	}
	return
}

func (c Client) creationPayload(filename string) (buf bytes.Buffer) {
	//filename = c.fixFilename(filename)
	payload := creationPayload{filename, c.settingPayload()}
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		fmt.Printf("There was an error encoding the json. err = %s", err)
	}
	return
}

func (c Client) settingPayload() settingPayload {
	return settingPayload{"public"}
}

func (c Client) creationURL() string {
	u := c.apiURL()
	u.Path = c.creationPath()
	return u.String()
}

func (c Client) existingURL() string {
	u := c.apiURL()
	u.Path = c.existingPath()
	return u.String()
}

func (c Client) apiURL() *url.URL {
	u, err := url.Parse(c.Host)
	if err != nil {
		panic(err)
	}
	return u
}

func (c Client) creationPath() string {
	return fmt.Sprintf("%d/sharing/create_shared_link_with_settings", c.Version)
}

func (c Client) existingPath() string {
	return fmt.Sprintf("%d/sharing/list_shared_links", c.Version)
}
