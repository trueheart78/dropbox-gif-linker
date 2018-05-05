package dropbox

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = Client{}

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
	assert.Equal(t, "https://api.dropboxapi.com/2/sharing/list_shared_links", client.existingURL())
}

func TestCreationURL(t *testing.T) {
	assert.Equal(t, "https://api.dropboxapi.com/2/sharing/create_shared_link_with_settings", client.creationURL())
}

func TestCreationPath(t *testing.T) {
	assert.Equal(t, "sharing/create_shared_link_with_settings", client.creationPath())
}

func TestExistingPath(t *testing.T) {
	assert.Equal(t, "sharing/list_shared_links", client.existingPath())
}

func TestApiURL(t *testing.T) {
	assert.Equal(t, "https://api.dropboxapi.com/", client.apiURL().String())
}

func TestApiVersion(t *testing.T) {
	assert.Equal(t, 2, client.apiVersion())
}

func TestNil(t *testing.T) {
	assert := assert.New(t)
	// assert for nil - good for error-checking
	assert.Nil(nil)
	assert.NotNil(1)
}
