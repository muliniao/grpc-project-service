package version

import (
	"fmt"
	"strconv"
	"strings"
)

// The BuildString is a variable that is supposed to be set during build-time.
// This can be done by using the -ldflags='-X "fitstation-hp/lib-fs-core-go/pkg/v1/version.BuildString=VALUE"' parameter on go build.
// For examples of this, check any the Makefile contents of the projects using this library.
//
// Expected format: "0.13.90 a722bdb 2018-01-09T22:32:37+01:00 go version go1.11 linux/amd64"
var BuildString string

// Version is an object containing the application's version information.
type Version struct {
	Major int // Mayor version (vX)
	Minor int // Minor version (v0.X)
	Patch int // Patch version (v0.0.X)

	GitRevision   string // Git commit hash.
	GitAuthorDate string // Git commit date + time.

	GoVersion string // Go compiler version.
	GoArch    string // Go GOOS/GOARCH for which this application is compiled.
}

// CurrentVersion uses the BuildString variable to generate a Version object.
func CurrentVersion() Version {
	fields := strings.Split(BuildString, " ")
	if len(fields) != 7 {
		return Version{}
	}

	verFields := strings.Split(fields[0], ".")
	if len(verFields) != 3 {
		return Version{}
	}

	major, err := strconv.Atoi(verFields[0])
	if err != nil {
		return Version{}
	}
	minor, err := strconv.Atoi(verFields[1])
	if err != nil {
		return Version{}
	}
	patch, err := strconv.Atoi(verFields[2])
	if err != nil {
		return Version{}
	}

	ver := Version{
		Major:         major,
		Minor:         minor,
		Patch:         patch,
		GitRevision:   fields[1],
		GitAuthorDate: fields[2],
		GoVersion:     fields[5],
		GoArch:        fields[6],
	}

	return ver
}

// ToString function for the Version.
// Format: "0.13.90 a722bdb 2018-01-09T22:32:37+01:00 go version go1.9 linux/amd64"
func (v Version) String() string {
	version := fmt.Sprintf(
		"%d.%d.%d",
		v.Major,
		v.Minor,
		v.Patch,
	)

	if v.GitRevision != "" {
		version = version + " " + v.GitRevision
	}

	if v.GitAuthorDate != "" {
		version = version + " " + v.GitAuthorDate
	}

	if v.GoVersion != "" {
		version = version + " go version " + v.GoVersion
	}

	if v.GoArch != "" {
		version = version + " " + v.GoArch
	}

	return version
}
