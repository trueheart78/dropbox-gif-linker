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

type testConfig struct {
	fullPath string
	gifDir   string
	apiToken string
	valid    bool
}

func (t testConfig) FullPath() string {
	return t.fullPath
}
func (t testConfig) GifsPath() string {
	return t.gifDir
}
func (t testConfig) Token() string {
	return t.apiToken
}
func (t testConfig) Valid() bool {
	return t.valid
}
func (t testConfig) Environment() string {
	return "test"
}
func (t testConfig) DatabasePath() string {
	if t.Environment() != "" {
		return fmt.Sprintf("%v/gifs-%v.sqlite3.db", filepath.Join(t.FullPath(), ".gifs"), t.Environment())
	}
	return fmt.Sprintf("%v/gifs.sqlite3.db", filepath.Join(t.FullPath(), ".gifs"))
}
func (t testConfig) LoadedPath() string {
	return ""
}

var missingFile = "/gifs/def.gif"
var existingFile = "/gifs/taylor swift/excited/file name 1.gif"
var host = "https://example-api.com"
var version = 3
var client = Client{
	Host:    host,
	Version: version,
}
var apiToken = "xxx"
var fullPath = "/my/path/to/dropbox/"
var gifDir = "/gifs"
var validConfig = testConfig{fullPath, gifDir, apiToken, true}
var invalidConfig = testConfig{fullPath, gifDir, apiToken, false}

func fixturePath(filename string) string {
	workingDir, _ := os.Getwd()
	return fmt.Sprintf("%v.json", filepath.Join(workingDir, "fixtures", filename))
}

func TestConfigPath(t *testing.T) {
	fullConfigPath := configPath(".dgl.json")

	assert.True(t, strings.HasSuffix(fullConfigPath, configFilename))
	assert.False(t, strings.HasPrefix(fullConfigPath, configFilename))
}

func TestConfigExists(t *testing.T) {
	assert.True(t, configExists(validConfigFilename))
	assert.True(t, configExists(emptyConfigFilename))
	assert.False(t, configExists(missingConfigFilename))
}

func TestConfigLoad(t *testing.T) {
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

func TestConfigDatabasePath(t *testing.T) {
	assert := assert.New(t)

	// valid config
	d := Config{}
	d.load(validConfigFilename)

	dbPath := filepath.Join(d.FullPath(), ".gifs", "gifs.sqlite3.db")
	assert.Equal(dbPath, d.DatabasePath())
}

func TestConfigLoadedPath(t *testing.T) {
	// valid config
	d := Config{}
	d.load(validConfigFilename)

	assert.Equal(t, validConfigFilename, d.LoadedPath())
}

func TestConfigGifDirFix(t *testing.T) {
	d := Config{GifDir: "example/"}

	assert.Equal(t, "example/", d.GifDir)
	d.gifDirFix()
	assert.Equal(t, "/example/", d.GifDir)
}

func TestConfigValidate(t *testing.T) {
	assert := assert.New(t)

	d, derr := createFromConfig(validConfigFilename)
	ok, err := d.validate()
	assert.Nil(derr)
	assert.True(d.Valid())
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(invalidPathConfigFilename)
	ok, err = d.validate()
	assert.False(d.Valid())
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_path should be \"/Dropbox/\" instead of \"Dropbox/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	ok, err = d.validate()
	assert.False(d.Valid())
	assert.False(ok)
	assert.NotNil(err)
	assert.Equal("the dropbox_gif_dir should be \"/gifs/\" instead of \"gifs/\"", err.Error())

	d, derr = createFromConfig(invalidDirConfigFilename)
	d.gifDirFix()
	ok, err = d.validate()
	assert.True(d.Valid())
	assert.True(ok)
	assert.Nil(err)

	d, derr = createFromConfig(emptyConfigFilename)
	ok, err = d.validate()
	assert.False(d.Valid())
	assert.False(ok)
	assert.NotNil(err)
}

func TestClientTruncate(t *testing.T) {
	c := newClient(validConfig)

	originalFilename := filepath.Join(validConfig.FullPath(), "example", "sample.gif")
	truncatedFilename, err := c.Truncate(originalFilename)

	assert.Equal(t, "example/sample.gif", truncatedFilename)
	assert.Nil(t, err)

	_, err = c.Truncate("invalid/dropbox/path.gif")
	assert.Equal(t, fmt.Sprintf("filepath does not contain the dropbox path [%v]", fullPath), err.Error())
	assert.NotNil(t, err)
}

func TestClientFixFilename(t *testing.T) {
	c := newClient(validConfig)

	originalFilename := "sample.gif"
	fixedFilename := c.fixFilename(originalFilename)
	// makes sure the basic file is in the gifs dropbox path
	assert.Equal(t, "/gifs/sample.gif", fixedFilename)

	originalFilename = "sample/hello/sample.gif"
	fixedFilename = c.fixFilename(originalFilename)
	// makes sure the full file is in the gifs dropbox path
	assert.Equal(t, "/gifs/sample/hello/sample.gif", fixedFilename)

	originalFilename = "/example/sample/hello/sample.gif"
	fixedFilename = c.fixFilename(originalFilename)
	// makes sure the full file is in the gifs dropbox path
	assert.Equal(t, "/gifs/example/sample/hello/sample.gif", fixedFilename)

	originalFilename = "/gifs/sample/hello/sample.gif"
	fixedFilename = c.fixFilename(originalFilename)
	// doesn't change the original filename
	assert.Equal(t, originalFilename, fixedFilename)

	originalFilename = "gifs/sample/hello/sample.gif"
	fixedFilename = c.fixFilename(originalFilename)
	// double /gifs when gifs does not have a leading slash
	assert.Equal(t, "/gifs/gifs/sample/hello/sample.gif", fixedFilename)
}

func TestClientCreateLink(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL

	// TODO: write some tests!
	// 1. when it exists
	// 2. when it does not exist and creation succeeds
	// 3. when it does not exist and creation fails
	// 4. when the dropbox api returns a 409 on exist check
}

func TestClientCreationWithInvalidAuthServer(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL
	_, err := c.create(missingFile)
	assert.Equal(t, "dropbox returned a 400", err.Error())
}

func TestClientExistsWithInvalidAuthServer(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubInvalidAuth()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, "dropbox returned a 400", err.Error())
}

func TestClientCreationSuccess(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubCreationSuccess(existingFile)
	c.Host = apiStub.URL
	url, err := c.create(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
	assert.Equal(t, fmt.Sprintf("![%v](%v)", url.Name, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif"), url.Markdown())
	assert.Equal(t, "/s/DROPBOX_HASH", url.RemotePath())
	assert.Equal(t, "/taylor swift/excited", url.Directory())
	assert.Equal(t, "DROPBOX_ID", url.DropboxID())
}

func TestClientCreationExists(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubCreationExists()
	c.Host = apiStub.URL
	_, err := c.create(existingFile)
	assert.Equal(t, "dropbox returned a 409", err.Error())
}

func TestClientCreationFailure(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubCreationFailure()
	c.Host = apiStub.URL
	_, err := c.create(missingFile)
	assert.Equal(t, "dropbox returned a 409", err.Error())
}

func TestClientExistsWithNoLinks(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubUnshared()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, fmt.Sprintf("no existing link for %v", missingFile), err.Error())
}

func TestClientExistsWithLinks(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubShared(existingFile)
	c.Host = apiStub.URL
	url, err := c.exists(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
}

func TestClientExistsWithMultipleLinks(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubSharedMultiple(existingFile)
	c.Host = apiStub.URL
	url, err := c.exists(existingFile)
	assert.Nil(t, err)
	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/file+name+1.gif", url.DirectLink())
}

func TestClientExistsWithMultipleLinksNoMatch(t *testing.T) {
	c := newClient(validConfig)
	apiStub := stubSharedMultipleNoMatch()
	c.Host = apiStub.URL
	_, err := c.exists(missingFile)
	assert.Equal(t, fmt.Sprintf("no existing link for %v", missingFile), err.Error())
}

func TestNewClient(t *testing.T) {
	c := newClient(validConfig)
	assert.Equal(t, "https://api.dropboxapi.com", c.Host)
	assert.Equal(t, 2, c.Version)
}

func TestClientValidClient(t *testing.T) {
	c := Client{Config: invalidConfig}

	assert.False(t, c.valid())

	c.Host = "https://www.sample.com"
	assert.False(t, c.valid())

	c.Version = 1
	assert.False(t, c.valid())

	c.Config = validConfig
	assert.True(t, c.valid())
}

func TestClientExistingPayload(t *testing.T) {
	data := client.existingPayload(missingFile)
	json := fmt.Sprintf("{\"path\":\"%v\"}\n", missingFile)
	assert.Equal(t, json, data.String())
}

func TestClientCreationPayload(t *testing.T) {
	data := client.creationPayload(missingFile)
	json := fmt.Sprintf("{\"path\":\"%v\",\"settings\":{\"requested_visibility\":\"public\"}}\n", missingFile)
	assert.Equal(t, json, data.String())
}

func TestClientExistingURL(t *testing.T) {
	url := fmt.Sprintf("%v/%d/%v", host, version, "sharing/list_shared_links")
	assert.Equal(t, url, client.existingURL())
}

func TestClientCreationURL(t *testing.T) {
	url := fmt.Sprintf("%v/%d/%v", host, version, "sharing/create_shared_link_with_settings")
	assert.Equal(t, url, client.creationURL())
}

func TestClientCreationPath(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%d/sharing/create_shared_link_with_settings", client.Version), client.creationPath())
}

func TestClientExistingPath(t *testing.T) {
	assert.Equal(t, fmt.Sprintf("%d/sharing/list_shared_links", client.Version), client.existingPath())
}

func TestClientApiURL(t *testing.T) {
	assert.Equal(t, host, client.apiURL().String())
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"links\": [], \"has_more\": false}"))
	}))
}

func stubShared(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftExistingResponse(filePath)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func stubSharedMultiple(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftExistingResponseMultiple(filePath)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	}))
}

func stubSharedMultipleNoMatch() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(existingResponseMultipleNoMatch()))
	}))
}

func stubCreationFailure() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(creationFailsResponse()))
	}))
}

func stubCreationExists() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(creationExistsResponse()))
	}))
}

func stubCreationSuccess(filePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := craftCreationResponse(filePath)
		w.Header().Set("Content-Type", "application/json")
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
