package models

const name = "surveysvc"

// These variables are assigned values during the build process using the -ldflags="-X ..." linker option.
var version = "Not set"
var origin = "Not set"
var commit = "Not set"
var branch = "Not set"
var built = "Not set"

// Version represents the version of a service and other useful metadata.
type Version struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Origin  string `json:"origin"`
	Commit  string `json:"commit"`
	Branch  string `json:"branch"`
	Built   string `json:"built"`
}

// NewVersion returns a Version.
func NewVersion() Version {
	return Version{Name: name, Version: version, Origin: origin, Commit: commit, Branch: branch, Built: built}
}
