package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client for the Dropbox API interactions
type Client struct {
	Host    string
	Version int
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
	Links   []existingLink `json:"links"`
	HasMore bool           `json:"has_more"`
}

type existingLink struct {
	Tag            string          `json:".tag"`
	URL            string          `json:"url"`
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Path           string          `json:"path_lower"`
	Permissions    linkPermissions `json:"link_permissions"`
	ClientModified string          `json:"client_modified"`
	ServerModified string          `json:"server_modified"`
	Revision       string          `json:"rev"`
	FileSize       int             `json:"size"`
}

type linkPermissions struct {
	ResolvedVisibility  linkTag `json:"resolved_visibility"`
	RequestedVisibility linkTag `json:"requested_visibility"`
	CanRevoke           bool    `json:"can_revoke"`
}

type linkTag struct {
	Tag string `json:".tag"`
}

// NewClient creates a new Client for interacting with Dropbox
func NewClient() (c Client) {
	c.Host = "https://api.dropboxapi.com"
	c.Version = 2
	return
}

func (c Client) valid() bool {
	if c.Host == "" || c.Version == 0 {
		return false
	}
	return true
}

func (c Client) exists(filename string) (ok bool, url string, err error) {
	if !c.valid() {
		return
	}
	data := c.existingPayload(filename)
	request, err := http.NewRequest(http.MethodPost, c.existingURL(), &data)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "API_TOKEN")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Dropbox Gif Linker")
	result, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if result.StatusCode != http.StatusOK {
		err = fmt.Errorf("dropbox returned a %d", result.StatusCode)
		return
	}

	var rawBody []byte
	var exists existsResponse
	rawBody, err = ioutil.ReadAll(result.Body)
	if err == nil {
		json.Unmarshal(rawBody, &exists)
		if len(exists.Links) > 0 {
			ok = true
			url = exists.Links[0].directLink()
		}
	}
	result.Body.Close()
	return
}

// changes the host and removes the RawQuery to provide the correct URL
// https://dl.dropboxusercontent.com/s/eqoo012hoa0wq7k/taylor%20bat%20focused.gif
// https://www.dropbox.com/s/eqoo012hoa0wq7k/taylor%20bat%20focused.gif?dl=0
func (e existingLink) directLink() string {
	u, err := url.Parse(e.URL)
	if err != nil {
		panic(err)
	}
	u.Host = "www.dropbox.com"
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
