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
	assert.Equal(t, 1.0, Current)
}

func TestReleaseCandidate(t *testing.T) {
	assert.Equal(t, 3, ReleaseCandidate)
}

func TestFullVersion(t *testing.T) {
	expected := fmt.Sprintf("%v version %.1f-rc%d %v/%v", Library, Current, ReleaseCandidate, runtime.GOOS, runtime.GOARCH)
	assert.Equal(t, expected, Full())
}
