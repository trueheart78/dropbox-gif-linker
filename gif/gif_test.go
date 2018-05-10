package gif

import (
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
