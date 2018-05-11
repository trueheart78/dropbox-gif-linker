package gif

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func dbPath() string {
	workingDir, _ := os.Getwd()
	return filepath.Join(workingDir, "..", "db", "test.sqlite3.db")
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

func generateRecord(id int64, sharedID string) (r Record) {
	if id > 0 {
		r.ID = id
		r.CreatedAt = dbTime()
		r.UpdatedAt = dbTime()
	}
	r.BaseName = "swiftie life 'the best' - 02.gif"
	r.Directory = "/taylor swift"
	r.Checksum = fmt.Sprintf("%vabcdefghijklmnopqrstuvxyz", sharedID)
	r.FileSize = 3456
	r.SharedLinkID = sharedID
	r.SharedLink = RecordSharedLink{
		ID:         sharedID,
		GifID:      r.ID,
		RemotePath: "s/DROPBOX_HASH",
		Count:      1,
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
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
	SetDatabasePath("./missing.sqlite3.db")
	ok, err := removeDatabase()
	assert.False(t, ok)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("remove %v: no such file or directory", "./missing.sqlite3.db"), err.Error())

	SetDatabasePath("./sample.sqlite3.db")
	ok, err = Init()
	assert.True(t, ok)
	assert.Nil(t, err)

	ok, err = removeDatabase()
	assert.True(t, ok)
	assert.Nil(t, err)
}

func TestGifSave(t *testing.T) {
	setUp()

	recordOne := generateRecord(0, "abcd")
	recordTwo := generateRecord(0, "efgh")

	ok, err := recordOne.Save()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), recordOne.ID)

	ok, err = recordTwo.Save()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), recordTwo.ID)

	tearDown()
}

func TestGifCreate(t *testing.T) {
	setUp()

	recordOne := generateRecord(0, "abcd")
	recordTwo := generateRecord(0, "efgh")
	recordThree := generateRecord(0, "efgh")

	ok, err := recordOne.Create()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), recordOne.ID)

	ok, err = recordTwo.Create()
	assert.True(t, ok)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), recordTwo.ID)

	// we allow a single type of error
	ok, err = recordThree.Create()
	assert.True(t, ok)
	assert.Nil(t, err)

	tearDown()
}

func TestGifUpdate(t *testing.T) {
	setUp()

	recordOne := generateRecord(0, "abcd")
	ok, err := recordOne.Create()
	assert.True(t, ok)
	assert.Nil(t, err)

	oldTime := recordOne.UpdatedAt
	oldTime2 := recordOne.SharedLink.UpdatedAt
	ok, err = recordOne.Update()
	assert.Nil(t, err)
	assert.NotEqual(t, recordOne.UpdatedAt, oldTime, "gif update times should differ")
	assert.NotEqual(t, recordOne.SharedLink.UpdatedAt, oldTime2, "shared_link update times should differ")

	tearDown()
}

func TestGifSaveExtended(t *testing.T) {
	setUp()

	recordOne := generateRecord(0, "abcd")
	ok, err := recordOne.Save()
	assert.True(t, ok)
	assert.Nil(t, err)

	oldTime := recordOne.UpdatedAt
	oldTime2 := recordOne.SharedLink.UpdatedAt
	ok, err = recordOne.Save()
	assert.Nil(t, err)
	assert.NotEqual(t, recordOne.UpdatedAt, oldTime, "gif update times should differ")
	assert.NotEqual(t, recordOne.SharedLink.UpdatedAt, oldTime2, "shared_link update times should differ")

	tearDown()
}

func TestGifCount(t *testing.T) {
	setUp()

	count := Count()
	assert.Equal(t, 0, count)

	record := generateRecord(0, "abcd")
	record.Create()

	count = Count()
	assert.Equal(t, 1, count)

	record = generateRecord(0, "wxyz")
	record.Create()

	count = Count()
	assert.Equal(t, 2, count)

	tearDown()
}

func TestGifRecordIncrement(t *testing.T) {
	setUp()

	record := generateRecord(0, "swift")
	_, err := record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 1, record.SharedLink.Count)

	_, err = record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 2, record.SharedLink.Count)

	_, err = record.Increment()
	assert.Nil(t, err)
	assert.Equal(t, 3, record.SharedLink.Count)

	tearDown()
}

func TestFind(t *testing.T) {
	setUp()

	record := generateRecord(0, "swift")
	record.Save()

	recordTwo, err := Find(record.ID)
	assert.Nil(t, err)
	assert.Equal(t, record.ID, recordTwo.ID)
	assert.Equal(t, record.SharedLinkID, recordTwo.SharedLink.ID)

	_, err = Find(1989)
	assert.NotNil(t, err)
	assert.Equal(t, "no gif with ID 1989", err.Error())

	tearDown()
}

func TestFindByMD5Checksum(t *testing.T) {
	setUp()

	record := generateRecord(0, "swift")
	record.Save()

	recordTwo, err := FindByMD5Checksum(record.Checksum)
	assert.Nil(t, err)
	assert.Equal(t, record.ID, recordTwo.ID)
	assert.Equal(t, record.SharedLinkID, recordTwo.SharedLink.ID)

	_, err = FindByMD5Checksum("unused-checksum")
	assert.NotNil(t, err)
	assert.Equal(t, "no gif with checksum unused-checksum", err.Error())

	tearDown()
}

func TestFindByFilename(t *testing.T) {
	setUp()

	record := generateRecord(0, "swift")
	record.Save()

	filename := filepath.Join(record.Directory, record.BaseName)
	recordTwo, err := FindByFilename(filename)
	assert.Nil(t, err)
	assert.Equal(t, record.ID, recordTwo.ID)

	filename = "/miss swift/no gif for swift.gif"
	_, err = FindByFilename(filename)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("no gif with that directory/basename [%v]", filename), err.Error())

	tearDown()
}

func TestRecordString(t *testing.T) {
	record := generateRecord(1989, "swift")

	assert.Equal(t, "[1989] [taylor swift] swiftie life 'the best' - 02.gif (3.5 kB) [used: 1]", record.String())
}

func TestRecordTags(t *testing.T) {
	record := Record{}

	record.Directory = "/taylor swift/love story"
	assert.Equal(t, "taylor swift, love story", record.Tags())

	record.Directory = "/swift/love"
	assert.Equal(t, "swift, love", record.Tags())

	record.Directory = "/taylor swift"
	assert.Equal(t, "taylor swift", record.Tags())
}

func TestRecordURL(t *testing.T) {
	record := generateRecord(1989, "swift")

	assert.Equal(t, "https://dl.dropboxusercontent.com/s/DROPBOX_HASH/swiftie+life+%2527the+best%2527+-+02.gif", record.URL())
}

func TestRecordMarkdown(t *testing.T) {
	record := generateRecord(1989, "swift")

	assert.Equal(t, fmt.Sprintf("![%v](%v)", record.BaseName, record.URL()), record.Markdown())
}
