package honu

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 0
	VersionMinor         = 4
	VersionPatch         = 0
	VersionReleaseLevel  = "beta"
	VersionReleaseNumber = 12
)

// Set the GitVersion via -ldflags="-X 'github.com/rotationalio/honu.GitVersion=$(git rev-parse --short HEAD)'"
var GitVersion string

// Set the BuildDate via -ldflags="-X github.com/rotationalio/honu.BuildDate=YYYY-MM-DD"
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
