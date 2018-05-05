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

// NewClient create a new Client for interacting with Dropbox
func NewClient() (c Client) {
	return c
}

func (c Client) exists(filename string) (ok bool, err error) {
	data := c.existingPayload(filename)
	request, err := http.NewRequest(http.MethodPost, c.existingURL(), &data)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Authorization", "API_TOKEN")
	request.Header.Set("Content-Type", "application/json")

	result, err := http.DefaultClient.Do(request)

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		err = fmt.Errorf("dropbox returned a %d", result.StatusCode)
		return
	}
	fmt.Println(result.StatusCode)
	fmt.Println(result)
	if err != nil {
		panic(err)
	}
	_, err = ioutil.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	result.Body.Close()
	return
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
	u.Path = fmt.Sprintf("%d/%v", c.apiVersion(), c.creationPath())
	return u.String()
}

func (c Client) existingURL() string {
	u := c.apiURL()
	u.Path = fmt.Sprintf("%d/%v", c.apiVersion(), c.existingPath())
	return u.String()
}

func (c Client) apiURL() *url.URL {
	u, err := url.Parse("https://api.dropboxapi.com/")
	if err != nil {
		panic(err)
	}
	return u
}

func (c Client) creationPath() string {
	return "sharing/create_shared_link_with_settings"
}

func (c Client) existingPath() string {
	return "sharing/list_shared_links"
}

func (c Client) apiVersion() int {
	return 2
}
