package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLibrary(t *testing.T) {
	assert.Equal(t, "dropbox-gif-linker", Library)
}

func TestCurrent(t *testing.T) {
	assert.Equal(t, "1.5.1", Current())
}

func TestFullVersion(t *testing.T) {
	var rc string
	if ReleaseCandidate > 0 {
		rc = fmt.Sprintf("-rc%d", ReleaseCandidate)
	}
	expected := fmt.Sprintf("%v version %d.%d.%d%v %v/%v", Library, Major, Minor, Patch, rc, runtime.GOOS, runtime.GOARCH)
	assert.Equal(t, expected, Full())
}
