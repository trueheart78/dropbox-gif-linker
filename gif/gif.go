package gif

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	// the adapter for sqlite3
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

var db *sql.DB
var databasePath string

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
	ID           int64
	BaseName     string
	Directory    string
	Size         int
	SharedLinkID string
	CreatedAt    string
	UpdatedAt    string
	SharedLink   SharedLink
}

// SharedLink details for a gif record
type SharedLink struct {
	ID         string
	GifID      int64
	RemotePath string
	Count      int
	CreatedAt  string
	UpdatedAt  string
}

// matches the time format used by sqlite3: "2018-02-19 13:56:25.741308"
func dbTime() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000000")
}

// Find looks up a record by ID
func Find(id int) (record Record, err error) {

	return
}

// Save captures the record to the database
func (r *Record) Save() (ok bool, err error) {
	if r.ID == 0 {
		return r.Create()
	}
	return r.Update()
}

// Update makes the record more up-to-date
func (r *Record) Update() (ok bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	dateString := dbTime()

	stmt, err := tx.Prepare("UPDATE gifs SET basename = ?, directory = ?, size = ?, shared_link_id = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	var affected, affected2 int64
	var u sql.Result
	u, err = stmt.Exec(r.BaseName, r.Directory, r.Size, r.SharedLinkID, dateString, r.ID)
	if err != nil {
		return
	}

	affected, err = u.RowsAffected()
	if err != nil {
		return
	}
	stmt2, err2 := tx.Prepare("UPDATE shared_links SET gif_id = ?, remote_path = ?, count = ?, updated_at = ? WHERE id = ?")
	if err2 != nil {
		err = err2
		return
	}
	defer stmt2.Close() // danger!

	_, err2 = stmt2.Exec(r.ID, r.SharedLink.RemotePath, r.SharedLink.Count, dateString, r.SharedLink.ID)
	if err2 != nil {
		err = err2
		return
	}
	affected2, err2 = u.RowsAffected()
	if err2 != nil {
		err = err2
		return
	}

	err2 = tx.Commit()
	if err2 != nil {
		err = err2
		return
	}

	if affected > 0 {
		r.UpdatedAt = dateString
	}

	if affected2 > 0 {
		r.SharedLink.UpdatedAt = dateString
	}

	ok = true
	return
}

// Create makes a new record in the database
func (r *Record) Create() (ok bool, err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()
	dateString := dbTime()

	stmt, err := tx.Prepare("INSERT INTO gifs (basename, directory, size, shared_link_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	var id int64
	var u sql.Result
	u, err = stmt.Exec(r.BaseName, r.Directory, r.Size, r.SharedLinkID, dateString, dateString)
	if err != nil {
		return
	}
	stmt2, err2 := tx.Prepare("INSERT INTO shared_links (id, gif_id, remote_path, count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err2 != nil {
		err = err2
		return
	}
	defer stmt2.Close() // danger!

	id, err2 = u.LastInsertId()
	if err2 != nil {
		err = err2
		return
	}
	_, err2 = stmt2.Exec(r.SharedLink.ID, id, r.SharedLink.RemotePath, r.SharedLink.Count, dateString, dateString)
	if err2 != nil {
		err = err2
		return
	}
	err2 = tx.Commit()
	if err2 != nil {
		err = err2
		return
	}

	r.ID = id
	r.CreatedAt = dateString
	r.UpdatedAt = dateString
	r.SharedLink.CreatedAt = dateString
	r.SharedLink.UpdatedAt = dateString

	ok = true
	return
}

// Init queues up the database connection
func Init() (ok bool, err error) {
	var structure string
	if databasePath == "" {
		err = errors.New("no database path set")
		return
	}
	// connect to the database
	_, err = Connect()
	if err != nil {
		return
	}
	// load the structure sql
	structure, err = structureStatement()
	if err != nil {
		return
	}
	// create the structure
	_, err = db.Exec(structure)
	if err != nil {
		return
	}
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
	db, err = sql.Open("sqlite3", databasePath)
	if err != nil {
		return
	}
	ok = true
	return
}

// Disconnect from the database connection
func Disconnect() {
	if db != nil {
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

func structureStatement() (data string, err error) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	schemaPath := filepath.Join(wd, "../", "db", "schema.sql")

	file, err := os.Open(schemaPath)
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
