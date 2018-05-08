package dropbox

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validConfigFilename = fixturePath("valid")
var invalidPathConfigFilename = fixturePath("invalid_path")
var invalidDirConfigFilename = fixturePath("invalid_dir")
var emptyConfigFilename = fixturePath("empty")
var missingConfigFilename = fixturePath("missing")

var dropbox = Config{}

type testConfig struct {
	fullPath string
	apiToken string
	valid    bool
}

func (t testConfig) FullPath() string {
	return t.fullPath
}
func (t testConfig) Token() string {
	return t.apiToken
}
func (t testConfig) Valid() bool {
	return t.valid
}

var missingFile = "gifs/def.gif"
var existingFile = "/gifs/file name 1.gif"
var host = "https://example-api.com"
var version = 3
var client = Client{
	Host:    host,
	Version: version,
}
var apiToken = "xxx"
var fullPath = "xxxx/xxx"
var validConfig = testConfig{fullPath, apiToken, true}
var invalidConfig = testConfig{fullPath, apiToken, false}

func fixturePath(filename string) string {
	workingDir, _ := os.Getwd()
	return fmt.Sprintf("%v.json", filepath.Join(workingDir, "fixtures", filename))
}

func TestConfigPath(t *testing.T) {
	assert := assert.New(t)

	fullConfigPath := configPath(".dgl.json")
	assert.True(strings.HasSuffix(fullConfigPath, configFilename))
	assert.False(strings.HasPrefix(fullConfigPath, configFilename))
}

func TestConfigExists(t *testing.T) {
	assert := assert.New(t)

	assert.True(configExists(validConfigFilename))
	assert.True(configExists(emptyConfigFilename))
	assert.False(configExists(missingConfigFilename))
}

func TestLoad(t *testing.T) {
	assert := assert.New(t)

	// valid config
	d := Config{}
	d.load(validConfigFilename)

	assert.Equal("~/Dropbox", d.DropboxPath)
	assert.Equal("/gifs", d.GifDir)
	assert.Equal("API_TOKEN", d.APIToken)
	assert.Equal(validConfigFilename, d.Path)
	assert.True(d.Loaded)

	// empty config
	d = Config{}
	d.load(emptyConfigFilename)

	assert.Equal("", d.DropboxPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(emptyConfigFilename, d.Path)
	assert.True(d.Loaded)

	// missing config
	d = Config{}
	d.load(missingConfigFilename)

	assert.Equal("", d.DropboxPath)
	assert.Equal("", d.GifDir)
	assert.Equal("", d.APIToken)
	assert.Equal(missingConfigFilename, d.Path)
	assert.False(d.Loaded)
}

func TestGifDirFix(t *testing.T) {
	assert := assert.New(t)

	d := Config{GifDir: "example/"}

	assert.Equal("example/", d.GifDir)
	d.gifDirFix()
	assert.Equal("/example/", d.GifDir)
}

func TestValidConfig(t *testing.T) {
	assert := assert.New(t)

	d, derr := createFromConfig(validConfigFilename)
	ok, err := d.Valid()
	assert.Nil(derr)
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(invalidPathConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_path should be \"/Dropbox/\" instead of \"Dropbox/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_gif_dir should be \"/gifs/\" instead of \"gifs/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	d.gifDirFix()
	ok, err = d.Valid()
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(emptyConfigFilename)
	ok, err = d.Valid()
	assert.False(ok)
	assert.NotNil(err)
}

func TestCreationWithInvalidAuthServer(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL
	_, err := c.create(missingFile)
	assert.Equal(t, "dropbox returned a 400", err.Error())
}

func TestExistsWithInvalidAuthServer(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, "dropbox returned a 400", err.Error())
}

func TestCreationSuccess(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubCreationSuccess(existingFile)
	c.Host = apiStub.URL
	url, err := c.create(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
}

func TestCreationExists(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubCreationExists()
	c.Host = apiStub.URL
	_, err := c.create(existingFile)
	assert.Equal(t, "dropbox returned a 409", err.Error())
}

func TestCreationFailure(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubCreationFailure()
	c.Host = apiStub.URL
	_, err := c.create(missingFile)
	assert.Equal(t, "dropbox returned a 409", err.Error())
}

func TestExistsWithNoLinks(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubUnshared()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, fmt.Sprintf("no existing link for %v", missingFile), err.Error())
}

func TestExistsWithLinks(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubShared(existingFile)
	c.Host = apiStub.URL
	url, err := c.exists(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
}

func TestExistsWithMultipleLinks(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubSharedMultiple(existingFile)
	c.Host = apiStub.URL
	url, err := c.exists(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
}

func TestExistsWithMultipleLinksNoMatch(t *testing.T) {
	c := NewClient(validConfig)
	apiStub := stubSharedMultipleNoMatch()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, fmt.Sprintf("no existing link for %v", missingFile), err.Error())
}

func TestNewClient(t *testing.T) {
	c := NewClient(validConfig)
	assert.Equal(t, "https://api.dropboxapi.com", c.Host)
	assert.Equal(t, 2, c.Version)
}

func TestValidClient(t *testing.T) {
	c := Client{Config: invalidConfig}

	assert.False(t, c.valid())

	c.Host = "https://www.sample.com"
	assert.False(t, c.valid())

	c.Version = 1
	assert.False(t, c.valid())

	c.Config = validConfig
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

func stubSharedMultiple(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftExistingResponseMultiple(filePath)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func stubSharedMultipleNoMatch() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(existingResponseMultipleNoMatch()))
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
	return craftResponse(existingResponse(), filePath)
}

func craftExistingResponseMultiple(filePath string) string {
	return craftResponse(existingResponseMultiple(), filePath)
}
func craftCreationResponse(filePath string) string {
	return craftResponse(creationValidResponse(), filePath)
}

func craftResponse(input string, filePath string) string {
	basename := path.Base(filePath)
	basenameEscaped := url.QueryEscape(basename)
	response := input
	response = strings.Replace(response, "RAW_BASENAME", basename, 1)
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
			"url": "https://www.dropbox.com/s/DROPBOX_HASH/URL_BASENAME?dl=0",
			"id": "id:DROPBOX_ID",
			"name": "RAW_BASENAME",
			"path_lower": "PATH_LOWER",
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
func existingResponseMultiple() string {
	return `
	{
    "links": [
		{
            ".tag": "folder",
            "url": "https://www.dropbox.com/sh/skm3nfeqlxb1ekw/AAB6wMgH5yJjFCOvbiRRHVZqa?dl=0",
            "id": "id:DROPBOX_ID_OTHER",
            "name": "gifs",
            "path_lower": "/gifs",
            "link_permissions": {
                "resolved_visibility": {
                    ".tag": "public"
                },
                "requested_visibility": {
                    ".tag": "public"
                },
                "can_revoke": true
            }
        },
        {
            ".tag": "file",
			"url": "https://www.dropbox.com/s/DROPBOX_HASH/URL_BASENAME?dl=0",
			"id": "id:DROPBOX_ID",
			"name": "RAW_BASENAME",
			"path_lower": "PATH_LOWER",
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
        },
        {
            ".tag": "folder",
            "url": "https://www.dropbox.com/sh/skm3nfeqlxb1ekw/AAB6wMgH5yJjFCOvbiRRHVZqa?dl=0",
            "id": "id:DROPBOX_ID_OTHER",
            "name": "gifs",
            "path_lower": "/gifs",
            "link_permissions": {
                "resolved_visibility": {
                    ".tag": "public"
                },
                "requested_visibility": {
                    ".tag": "public"
                },
                "can_revoke": true
            }
        }
    ],
    "has_more": false
	}
	`
}

func existingResponseMultipleNoMatch() string {
	return `
	{
    "links": [
		{
            ".tag": "folder",
            "url": "https://www.dropbox.com/sh/skm3nfeqlxb1ekw/AAB6wMgH5yJjFCOvbiRRHVZqa?dl=0",
            "id": "id:DROPBOX_ID_OTHER",
            "name": "gifs",
            "path_lower": "/gifs",
            "link_permissions": {
                "resolved_visibility": {
                    ".tag": "public"
                },
                "requested_visibility": {
                    ".tag": "public"
                },
                "can_revoke": true
            }
        },
        {
            ".tag": "folder",
            "url": "https://www.dropbox.com/sh/skm3nfeqlxb1ekw/AAB6wMgH5yJjFCOvbiRRHVZqa?dl=0",
            "id": "id:DROPBOX_ID_OTHER",
            "name": "gifs",
            "path_lower": "/gifs",
            "link_permissions": {
                "resolved_visibility": {
                    ".tag": "public"
                },
                "requested_visibility": {
                    ".tag": "public"
                },
                "can_revoke": true
            }
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
		"url": "https://www.dropbox.com/s/DROPBOX_HASH/URL_BASENAME?dl=0",
		"id": "id:DROPBOX_ID",
		"name": "RAW_BASENAME",
		"path_lower": "PATH_LOWER",
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
