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
	assert.Equal(t, 1.2, Current)
}

func TestReleaseCandidate(t *testing.T) {
	assert.Equal(t, 0, ReleaseCandidate)
}

func TestFullVersion(t *testing.T) {
	var rc string
	if ReleaseCandidate > 0 {
		rc = fmt.Sprintf("-rc%d", ReleaseCandidate)
	}
	expected := fmt.Sprintf("%v version %.1f%v %v/%v", Library, Current, rc, runtime.GOOS, runtime.GOARCH)
	assert.Equal(t, expected, Full())
}
