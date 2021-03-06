package gifkv

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	bolt "github.com/coreos/bbolt"
	humanize "github.com/dustin/go-humanize"
	homedir "github.com/mitchellh/go-homedir"
)

var db *bolt.DB
var databasePath string
var connected bool
var bucketName = "gifs"
var dropboxBaseURL = "https://dl.dropboxusercontent.com"

// SetDatabasePath sets the db path
func SetDatabasePath(filePath string) (ok bool, err error) {
	filePath, err = homedir.Expand(filePath)
	if err != nil {
		return
	}
	databasePath = filePath
	ok = true
	return
}

// GetDatabasePath returns the db path
func GetDatabasePath() string {
	return databasePath
}

func databaseDir() string {
	if databasePath == "" {
		return ""
	}
	return strings.Replace(databasePath, filepath.Base(databasePath), "", 1)
}

func resetDatabasePath() {
	databasePath = ""
}

// Record of a dropbox-linked gif
type Record struct {
	ID           string `json:"checksum"`
	BaseName     string `json:"base_name"`
	Directory    string `json:"directory"`
	FileSize     int    `json:"file_size"`
	SharedLinkID string `json:"shared_link_id"`
	RemotePath   string `json:"remote_path"`
	persisted    bool
}

// Count returns the number of gifs cached in the database
func Count() int {
	var s bolt.BucketStats
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		s = b.Stats()
		return nil
	})
	return s.KeyN
}

// Find looks up a record by checksum
func Find(checksum string) (record Record, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		v := b.Get([]byte(checksum))
		if v != nil {
			json.Unmarshal(v, &record)
			record.persisted = true
		} else {
			return fmt.Errorf("Unable to find id \"%s\"", checksum)
		}
		return nil
	})
	return
}

// Save captures the record to the database
func (r *Record) Save() (bool, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Put([]byte(r.ID), r.json())
		return err
	})
	if err != nil {
		return false, err
	}
	r.persisted = true
	return true, nil
}

// Delete removes the record from the database
func (r *Record) Delete() (bool, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		err := b.Delete([]byte(r.ID))
		return err
	})
	if err != nil {
		return false, err
	}
	r.persisted = false
	return true, nil
}

// RemoteOK checks to see if a persisted record returns a 200 status code
func (r Record) RemoteOK() (bool, error) {
	if r.URL() == "" {
		return false, errors.New("empty url")
	}
	resp, err := http.Get(r.URL())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, nil
}

func (r Record) json() []byte {
	data, _ := json.Marshal(r)
	return data
}

// String returns a string formatted-Record
func (r Record) String() string {
	return fmt.Sprintf("[%v] %v (%v)", r.Tags(), r.BaseName, humanize.Bytes(uint64(r.FileSize)))
}

//Persisted returns whether the record is saved in the database
func (r Record) Persisted() bool {
	return r.persisted
}

// Tags break down the directory
func (r Record) Tags() string {
	tags := strings.Replace(r.Directory, string(os.PathSeparator), "", 1)
	tags = strings.Replace(tags, string(os.PathSeparator), ", ", -1)
	return tags
}

// URL returns a publicly-accessible url
func (r Record) URL() string {
	u, err := url.Parse(dropboxBaseURL)
	if err != nil {
		return ""
	}
	u.Path = filepath.Join(r.RemotePath, url.QueryEscape(r.BaseName))
	return u.String()
}

// Markdown returns a publicly-accessible markdown-based url
func (r Record) Markdown() string {
	return fmt.Sprintf("![%v](%v)", r.BaseName, r.URL())
}

// BBCode returns a publicly-accessible bbcode-based url
func (r Record) BBCode() string {
	return fmt.Sprintf("[img]%v[/img]", r.URL())
}

// Init queues up the database connection
func Init() (ok bool, err error) {
	if databasePath == "" {
		err = errors.New("no database path set")
		return
	}
	if _, err := os.Stat(databaseDir()); os.IsNotExist(err) {
		os.MkdirAll(databaseDir(), os.ModePerm)
	}
	// connect to the database
	_, err = Connect()
	if err != nil {
		return
	}
	// initiate the bucket
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
	defer db.Close()
	ok = true
	return
}

// Connect to the database
func Connect() (ok bool, err error) {
	if databasePath == "" {
		err = errors.New("no database path set")
		return
	}
	db, err = bolt.Open(databasePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		connected = false
		return
	}
	connected = true
	ok = true
	return
}

// Disconnect from the database connection
func Disconnect() {
	if db != nil && connected {
		db.Close()
	}
}

func removeDatabase() (ok bool, err error) {
	if databasePath == "" {
		err = errors.New("no database path set")
		return
	}
	err = os.Remove(databasePath)
	if err != nil {
		return
	}
	ok = true
	return
}
