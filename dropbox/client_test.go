package dropbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var host = "https://sample.com"
var version = 3
var client = Client{
	Host:    host,
	Version: version,
}

func skipTestThatThing(t *testing.T) {
	filename := "gifs/def.gif"
	ok, _ := client.exists(filename)
	assert.False(t, ok)
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.Equal(t, "https://api.dropboxapi.com", c.Host)
	assert.Equal(t, 2, c.Version)
}

func TestValid(t *testing.T) {
	c := Client{}

	assert.False(t, c.valid())

	c.Host = "https://www.sample.com"
	assert.False(t, c.valid())

	c.Version = 1
	assert.True(t, c.valid())
}

func TestExistingPayload(t *testing.T) {
	filename := "gifs/def.gif"
	data := client.existingPayload(filename)
	json := fmt.Sprintf("{\"path\":\"%v\"}\n", filename)
	assert.Equal(t, json, data.String())
}

func TestCreationPayload(t *testing.T) {
	filename := "gifs/def.gif"
	data := client.creationPayload(filename)
	json := fmt.Sprintf("{\"path\":\"%v\",\"settings\":{\"requested_visibility\":\"public\"}}\n", filename)
	assert.Equal(t, json, data.String())
}

func TestExistingURL(t *testing.T) {
	url := fmt.Sprintf("%v/%d/%v", host, version, "sharing/list_shared_links")
	assert.Equal(t, url, client.existingURL())
}

func TestCreationURL(t *testing.T) {
	url := fmt.Sprintf("%v/%d/%v", host, version, "sharing/create_shared_link_with_settings")
	assert.Equal(t, url, client.creationURL())
}

func TestCreationPath(t *testing.T) {
	assert.Equal(t, "sharing/create_shared_link_with_settings", client.creationPath())
}

func TestExistingPath(t *testing.T) {
	assert.Equal(t, "sharing/list_shared_links", client.existingPath())
}

func TestApiURL(t *testing.T) {
	assert.Equal(t, host, client.apiURL().String())
}

func TestApiVersion(t *testing.T) {
	assert.Equal(t, version, client.apiVersion())
}

func TestNil(t *testing.T) {
	assert := assert.New(t)
	// assert for nil - good for error-checking
	assert.Nil(nil)
	assert.NotNil(1)
}
