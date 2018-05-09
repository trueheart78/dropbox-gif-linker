package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	assert.Equal(t, 0.6, Current)
}

func TestFullVersion(t *testing.T) {
	expected := fmt.Sprintf("%v version %.2f %v/%v", Library, Current, runtime.GOOS, runtime.GOARCH)
	assert.Equal(t, expected, Full())
}
