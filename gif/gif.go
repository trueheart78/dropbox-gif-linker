package gif

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	// the adapter for sqlite3
	_ "github.com/mattn/go-sqlite3"
	homedir "github.com/mitchellh/go-homedir"
)

var db *sql.DB
var databasePath string
var connected bool

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

// Record of a gif
type Record struct {
	ID           int64
	BaseName     string
	Directory    string
	FileSize     int
	Checksum     string
	SharedLinkID string
	CreatedAt    string
	UpdatedAt    string
	SharedLink   RecordSharedLink
}

// RecordSharedLink details for a gif record
type RecordSharedLink struct {
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

// Count returns the number of gifs cached in the database
func Count() int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM gifs").Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

// Find looks up a record by ID
func Find(id int) (record Record, err error) {
	fmt.Printf("[debug] finding record %d\n", id)
	err = errors.New("construction")
	return
}

// FindByMD5Checksum looks up a record by the md5 checksum
func FindByMD5Checksum(checksum string) (record Record, err error) {
	fmt.Printf("[debug] finding md5 %v\n", checksum)
	err = errors.New("construction")
	return
}

// FindByFilename looks up the record by filename
func FindByFilename(shortFilename string) (record Record, err error) {
	fmt.Printf("[debug] finding by filename %v\n", shortFilename)
	err = errors.New("construction")
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

	stmt, err := tx.Prepare("UPDATE gifs SET basename = ?, directory = ?, size = ?, md5 = ?, shared_link_id = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	var affected, affected2 int64
	var u sql.Result
	u, err = stmt.Exec(r.BaseName, r.Directory, r.FileSize, r.Checksum, r.SharedLinkID, dateString, r.ID)
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

	stmt, err := tx.Prepare("INSERT INTO gifs (basename, directory, size, md5, shared_link_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	var id int64
	var u sql.Result
	u, err = stmt.Exec(r.BaseName, r.Directory, r.FileSize, r.Checksum, r.SharedLinkID, dateString, dateString)
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
	_, err2 = stmt2.Exec(r.SharedLink.ID, id, r.SharedLink.RemotePath, 1, dateString, dateString)
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
	r.SharedLink.Count = 1
	r.SharedLink.CreatedAt = dateString
	r.SharedLink.UpdatedAt = dateString

	ok = true
	return
}

// Increment updates the Count value in memory and in the db
func (r *Record) Increment() (ok bool, err error) {
	if r.ID > 0 {
		r.SharedLink.Count++
	}
	ok, err = r.Save()
	return
}

// String returns a string formatted-Record
func (r Record) String() string {
	return fmt.Sprintf("[%d] [%v] %v (%v) [used: %d]", r.ID, r.Directory, r.BaseName, humanize.Bytes(uint64(r.FileSize)), r.SharedLink.Count)
}

// URL returns a publicly-accessible url
func (r Record) URL() string {
	u, err := url.Parse("https://dl.dropboxusercontent.com")
	if err != nil {
		return ""
	}
	u.Path = filepath.Join(r.SharedLink.RemotePath, url.QueryEscape(r.BaseName))
	return u.String()
}

// Markdown returns a publicly-accessible markdown-based url
func (r Record) Markdown() string {
	return fmt.Sprintf("![%v](%v)", r.BaseName, r.URL())
}

// Init queues up the database connection
func Init() (ok bool, err error) {
	var structure string
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
