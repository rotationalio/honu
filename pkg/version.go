package pkg

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 1
	VersionMinor         = 0
	VersionPatch         = 0
	VersionReleaseLevel  = "alpha"
	VersionReleaseNumber = 1
)

// Set the GitVersion via -ldflags="-X 'github.com/rotationalio/honu/pkg.GitVersion=$(git rev-parse --short HEAD)'"
var GitVersion string

// Set the BuildDate via -ldflags="-X github.com/rotationalio/honu/pkg.BuildDate=YYYY-MM-DD"
var BuildDate string

// Version returns the semantic version for the current build.
func Version() string {
	versionCore := fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			versionCore = fmt.Sprintf("%s-%s.%d", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		} else {
			versionCore = fmt.Sprintf("%s-%s", versionCore, VersionReleaseLevel)
		}
	}

	if GitVersion != "" {
		if BuildDate != "" {
			versionCore = fmt.Sprintf("%s (revision %s built on %s)", versionCore, GitVersion, BuildDate)
		} else {
			versionCore = fmt.Sprintf("%s (%s)", versionCore, GitVersion)
		}
	}

	return versionCore
}
