package version

import (
	"fmt"
	"runtime"
)

const (
	// Library name
	Library = "dropbox-gif-linker"
	// Current version of the library
	Current = 0.8
)

// Full returns the full version string
func Full() string {
	return fmt.Sprintf("%v version %.2f %v/%v", Library, Current, runtime.GOOS, runtime.GOARCH)
}
