package version

import (
	"fmt"
	"runtime"
)

const (
	// Library name
	Library = "dropbox-gif-linker"
	// Major version
	Major = 1
	// Minor version
	Minor = 5
	// Patch version
	Patch = 1
	// ReleaseCandidate version of the library
	ReleaseCandidate = 0
)

// Current returns the semver version
func Current() string {
	var rc string
	if ReleaseCandidate > 0 {
		rc = fmt.Sprintf("-rc%d", ReleaseCandidate)
	}
	return fmt.Sprintf("%d.%d.%d%v", Major, Minor, Patch, rc)
}

// Full returns the full version string
func Full() string {
	return fmt.Sprintf("%v version %v %v/%v", Library, Current(), runtime.GOOS, runtime.GOARCH)
}
