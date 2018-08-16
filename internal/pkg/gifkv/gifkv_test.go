package gifkv

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dbPath() string {
	workingDir, _ := os.Getwd()
	return filepath.Join(workingDir, "..", "db", "test.boltdb.db")
}

func initDbPath() {
	SetDatabasePath(dbPath())
}

func setUp() {
	initDbPath()
	removeDatabase()
	Init()
	Connect()
}

func tearDown() {
	Disconnect()
	removeDatabase()
}

func generateRecord(checksum string, sharedID string) (r Record) {
	if checksum != "" {
		r.ID = checksum
	} else {
		r.ID = "random-checksum"

	}
	r.BaseName = "swiftie life 'the best' - 02.gif"
	r.Directory = "/taylor swift"
	r.FileSize = 3456
	r.SharedLinkID = sharedID
	r.RemotePath = "s/DROPBOX_HASH"
	r.Count = 1
	return
}

func TestSetDatabasePath(t *testing.T) {
	assert.Equal(t, "", GetDatabasePath())

	SetDatabasePath(dbPath())
	assert.Equal(t, dbPath(), GetDatabasePath())
}

func TestInit(t *testing.T) {
	resetDatabasePath()
	ok, err := Init()

	assert.NotNil(t, err)
	assert.Equal(t, "no database path set", err.Error())

	initDbPath()
	ok, err = Init()

	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestRemoveDatabase(t *testing.T) {
	SetDatabasePath("./missing.boltdb.db")
	ok, err := removeDatabase()
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("remove %v: no such file or directory", "./missing.boltdb.db"), err.Error())

	SetDatabasePath("./sample.boltdb.db")
	ok, err = Init()
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = removeDatabase()
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestGifSave(t *testing.T) {
	setUp()

	recordOne := generateRecord("checksum-a", "abcd")
	recordTwo := generateRecord("checksum-b", "efgh")

	ok, err := recordOne.Save()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "checksum-a", recordOne.ID)

	ok, err = recordTwo.Save()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, "checksum-b", recordTwo.ID)

	tearDown()
}

func TestGifCount(t *testing.T) {
	setUp()

	count := Count()
	assert.Equal(t, 0, count)

	record := generateRecord("checksum-a", "abcd")
	record.Save()

	count = Count()
	assert.Equal(t, 1, count)

	record = generateRecord("checksum-b", "wxyz")
	record.Save()

	count = Count()
	assert.Equal(t, 2, count)

	tearDown()
}

func TestGifRecordIncrement(t *testing.T) {
	setUp()

	record := generateRecord("checksum-a", "swift")
	_, err := record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 2, record.Count)

	_, err = record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 3, record.Count)

	_, err = record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 4, record.Count)

	tearDown()
}

func TestGifFind(t *testing.T) {
	setUp()

	record := generateRecord("checksum-a", "swift")
	record.Save()

	recordTwo, err := Find(record.ID)
	assert.Nil(t, err)
	assert.Equal(t, record.ID, recordTwo.ID)
	assert.Equal(t, record.SharedLinkID, recordTwo.SharedLinkID)

	_, err = Find("1989")
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to find id \"1989\"", err.Error())

	tearDown()
}

func TestGifRecordString(t *testing.T) {
	record := generateRecord("1989", "swift")

	// assert.Equal(t, "[taylor swift] swiftie life 'the best' - 02.gif (3.5 kB) [used: 1]", record.String())
	assert.Equal(t, "[taylor swift] swiftie life 'the best' - 02.gif (3.5 kB)", record.String())
}

func TestGifRecordTags(t *testing.T) {
	record := Record{}

	record.Directory = "/taylor swift/love story"
	assert.Equal(t, "taylor swift, love story", record.Tags())

	record.Directory = "/swift/love"
	assert.Equal(t, "swift, love", record.Tags())

	record.Directory = "/taylor swift"
	assert.Equal(t, "taylor swift", record.Tags())
}

func TestGifRecordURL(t *testing.T) {
	record := generateRecord("1989", "swift")

	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/swiftie+life+%2527the+best%2527+-+02.gif", record.URL())
}

func TestGifRecordMarkdown(t *testing.T) {
	record := generateRecord("1989", "swift")

	assert.Equal(t, fmt.Sprintf("![%v](%v)", record.BaseName, record.URL()), record.Markdown())
}
