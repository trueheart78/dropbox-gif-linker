package dropbox

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var missingFile = "gifs/def.gif"
var host = "https://example-api.com"
var version = 3
var client = Client{
	Host:    host,
	Version: version,
}

func TestExistsWithInvalidAuthServer(t *testing.T) {
	c := NewClient()
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, "dropbox returned a 400", err.Error())
}

func TestExistsWithNoLinks(t *testing.T) {
	c := NewClient()
	apiStub := stubUnshared()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, fmt.Sprintf("no existing link for %v", missingFile), err.Error())
}

func TestExistsWithLinks(t *testing.T) {
	c := NewClient()
	apiStub := stubShared("/gifs/file name 1.gif")
	c.Host = apiStub.URL
	url, err := c.exists("/gifs/file name 1.gif")
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/dropbox-hash/file+name+1.gif", url.DirectLink())
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	assert.Equal(t, "https://api.dropboxapi.com", c.Host)
	assert.Equal(t, 2, c.Version)
}

func TestValidClient(t *testing.T) {
	c := Client{}

	assert.False(t, c.valid())

	c.Host = "https://www.sample.com"
	assert.False(t, c.valid())

	c.Version = 1
	assert.True(t, c.valid())
}

func TestExistingPayload(t *testing.T) {
	data := client.existingPayload(missingFile)
	json := fmt.Sprintf("{\"path\":\"%v\"}\n", missingFile)
	assert.Equal(t, json, data.String())
}

func TestCreationPayload(t *testing.T) {
	data := client.creationPayload(missingFile)
	json := fmt.Sprintf("{\"path\":\"%v\",\"settings\":{\"requested_visibility\":\"public\"}}\n", missingFile)
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
	assert.Equal(t, fmt.Sprintf("%d/sharing/create_shared_link_with_settings", client.Version), client.creationPath())
}

func TestExistingPath(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%d/sharing/list_shared_links", client.Version), client.existingPath())
}

func TestApiURL(t *testing.T) {
	assert.Equal(t, host, client.apiURL().String())
}

func TestNil(t *testing.T) {
	assert := assert.New(t)
	// assert for nil - good for error-checking
	assert.Nil(nil)
	assert.NotNil(1)
}

func stubInvalidAuth() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := fmt.Sprintf("Error in call to API function \"%v\": The given OAuth 2 access token is malformed.", r.RequestURI)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(resp))
	}))
}

func stubUnshared() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"links\": [], \"has_more\": false}"))
	}))
}

func stubShared(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftExistingResponse(filePath)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func stubCreationFailure() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(creationFailsResponse()))
	}))
}

func stubCreationExists() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(creationExistsResponse()))
	}))
}

func stubCreationSuccess(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftCreationResponse(filePath)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func craftExistingResponse(filePath string) string {
	basename := path.Base(filePath)
	basenameEscaped := url.QueryEscape(basename)
	response := existingResponse()
	response = strings.Replace(response, "URL_BASENAME", basenameEscaped, 1)
	response = strings.Replace(response, "PATH_LOWER", filePath, 1)
	return response
}

func craftCreationResponse(filePath string) string {
	basename := path.Base(filePath)
	basenameEscaped := url.QueryEscape(basename)
	response := creationValidResponse()
	response = strings.Replace(response, "URL_BASENAME", basenameEscaped, 1)
	response = strings.Replace(response, "PATH_LOWER", filePath, 1)
	return response
}

//returns a 200 - OK
func existingResponse() string {
	return `
	{
		"links": [
		{
			".tag": "file",
			"url": "https://www.dropbox.com/s/dropbox-hash/URL_BASENAME?dl=0",
			"id": "id:DROPBOX_ID",
			"name": "file name 1.gif",
			"path_lower": "/gifs/examples/funny/file name 1.gif",
			"link_permissions": {
				"resolved_visibility": {
					".tag": "public"
				},
				"requested_visibility": {
					".tag": "public"
				},
				"can_revoke": true
			},
			"client_modified": "2017-09-01T15:37:19Z",
			"server_modified": "2017-12-01T16:37:11Z",
			"rev": "5d050301f24e",
			"size": 2078402
		}
		],
		"has_more": false
	}
	`
}

// returns a 409 - Conflict
func creationExistsResponse() string {
	return `
	{
	   "error_summary": "shared_link_already_exists/..",
	   "error": {
			".tag": "shared_link_already_exists"
	   }
	}
	`
}

// returns a 409 - Conflict
func creationFailsResponse() string {
	return `
	{
    "error_summary": "path/not_found/...",
    "error": {
        ".tag": "path",
        "path": {
            ".tag": "not_found"
			}
		}
	}
	`
}

// returns a 200 - OK
func creationValidResponse() string {
	return `
	{
		".tag": "file",
		"url": "https://www.dropbox.com/s/dropbox-hash/URL_BASENAME?dl=0",
		"id": "id:DROPBOX_ID",
		"name": "file name 1.gif",
		"path_lower": "/gifs/examples/funny/file name 1.gif",
		"link_permissions": {
			"resolved_visibility": {
				".tag": "public"
			},
			"requested_visibility": {
				".tag": "public"
			},
			"can_revoke": true
		},
		"client_modified": "2017-09-01T15:37:19Z",
		"server_modified": "2017-12-01T16:37:11Z",
		"rev": "5d050301f24e",
		"size": 2078402
	}
	`
}
