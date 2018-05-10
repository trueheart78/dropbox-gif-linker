package gif

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// the adapter for sqlite3
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

var db *sql.DB
var databasePath string

// var initialized bool

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

func resetDatabasePath() {
	databasePath = ""
}

// Record of a gif
type Record struct {
	ID           int
	BaseName     string
	Directory    string
	SharedLinkID string
	CreatedAt    string
	UpdatedAt    string
	SharedLink   SharedLink
}

// SharedLink details for a gif record
type SharedLink struct {
	ID         string
	GifID      int
	RemotePath string
	Count      int
	CreatedAt  string
	UpdatedAt  string
}

// Find looks up a record by ID
func Find(id int) (record Record, err error) {

	return
}

// Save captures the record to the database
func (r Record) Save() (ok bool, err error) {
	if r.ID == 0 {
		return r.Create()
	}
	return
}

// Create makes a new record in the database
func (r Record) Create() (ok bool, err error) {

	return
}

// Init queues up the database connection
func Init() (ok bool, err error) {
	var structure string
	if databasePath == "" {
		err = errors.New("no database path set")
		return
	}
	fmt.Println(databasePath)
	os.Remove(databasePath)
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return
	}
	structure, err = structureStatement()
	if err != nil {
		return
	}

	_, err = db.Exec(structure)
	if err != nil {
		return
	}

	ok = true
	defer db.Close()
	return
}

// Close the database connection
func Close() {
	if db != nil {
		db.Close()
	}
}

func structureStatement() (data string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	schema := filepath.Join(wd, "../", "db", "schema.sql")

	file, err := os.Open(schema)
	if err != nil {
		return
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	data = string(bytes)
	return
}
