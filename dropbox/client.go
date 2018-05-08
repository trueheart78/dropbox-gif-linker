package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client for the Dropbox API interactions
type Client struct {
	Host    string
	Version int
	Config  clientConfig
}

type clientConfig interface {
	FullPath() string
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

// NewClient creates a new Client for interacting with Dropbox
func NewClient(config clientConfig) (c Client) {
	c.Host = "https://api.dropboxapi.com"
	c.Version = 2
	c.Config = config
	return
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
	request.Header.Set("Authorization", "API_TOKEN")
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

func (c Client) existingPayload(filename string) (buf bytes.Buffer) {
	payload := existingPayload{filename}
	err := json.NewEncoder(&buf).Encode(&payload)
	if err != nil {
		fmt.Printf("There was an error encoding the json. err = %s", err)
	}
	return
}

func (c Client) creationPayload(filename string) (buf bytes.Buffer) {
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
