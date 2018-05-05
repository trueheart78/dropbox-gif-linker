package dropbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var client = Client{}

func TestApiURL(t *testing.T) {
	assert.Equal(t, "https://api.dropboxapi.com/2", client.apiURL())
}

func TestCreationPath(t *testing.T) {
	assert.Equal(t, "sharing/create_shared_link_with_settings", client.creationPath())
}

func TestExistingPath(t *testing.T) {
	assert.Equal(t, "sharing/list_shared_links", client.existingPath())
}

func TestNil(t *testing.T) {
	assert := assert.New(t)
	// assert for nil - good for error-checking
	assert.Nil(nil)
	assert.NotNil(1)
}
