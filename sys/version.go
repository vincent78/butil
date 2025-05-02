package sys

import (
	"fmt"
	"gitee.com/vincent78/gcutil/model"
	"os"
	"runtime"
)

var (
	VersionMajor = 0          // Major version component of the current release
	VersionMinor = 0          // Minor version component of the current release
	VersionPatch = 1          // Patch version component of the current release
	VersionMeta  = "snapshot" // Version metadata to append to the version string   stable   snapshot
)

type VersionModel struct {
	Identify  string `json:"id"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	Arch      string `json:"arch"`
	GoVersion string `json:"goVersion"`
	GoOS      string `json:"goos"`
	GoPath    string `json:"goPath"`
	GoRoot    string `json:"goRoot"`
}

func Version() VersionModel {
	obj := VersionModel{
		Identify:  model.Identify,
		Version:   VersionWithMeta(),
		Arch:      runtime.GOARCH,
		GoVersion: runtime.Version(),
		GoOS:      runtime.GOOS,
		GoPath:    os.Getenv("GOPATH"),
		GoRoot:    runtime.GOROOT(),
	}

	if model.GitCommit != "" {
		obj.Commit = model.GitCommit
	}
	if model.BuildDate != "" {
		obj.BuildDate = model.BuildDate
	}
	return obj
}

var VersionStr = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}

// VersionWithMeta holds the textual version string including the metadata.
var VersionWithMeta = func() string {
	v := VersionStr()
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}

func VersionWithCommit(gitCommit, buildDate string) string {
	vsn := VersionWithMeta()
	if len(gitCommit) >= 7 {
		vsn += "-" + gitCommit[:7]
	} else if len(gitCommit) > 0 {
		vsn += "-" + gitCommit
	}
	if buildDate != "" {
		vsn += "-" + buildDate
	}
	return vsn
}
