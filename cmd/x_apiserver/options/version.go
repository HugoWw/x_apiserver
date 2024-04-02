package options

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"runtime"
	"strconv"
)

var (
	// GitCommit is the commit of repo
	// sha1 from git, output of $(git rev-parse HEAD)
	GitCommit = "UNKNOWN"

	// APPVersion is the version of apiserver.
	APPVersion = "UNKNOWN"

	platform = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	goVersion = runtime.Version()

	compiler = runtime.Compiler

	//BuildDate in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	BuildDate = "1970-01-01T00:00:00Z"
)

type versionOptions struct {
	GitCommit  string `json:"gitCommit"`
	Date       string `json:"buildDate"`
	GoVersion  string `json:"goVersion"`
	Compiler   string `json:"compiler"`
	APPVersion string `json:"appVersion"`
	Platform   string `json:"platform"`
}

func newVersionOption() versionOptions {
	return versionOptions{
		GitCommit:  GitCommit,
		Date:       BuildDate,
		GoVersion:  goVersion,
		Compiler:   compiler,
		APPVersion: APPVersion,
		Platform:   platform,
	}
}

func formatVersion() string {
	vfo := newVersionOption()
	return fmt.Sprintf("gitCommit:%v\nbuildDate:%v\ngoVersion:%v\ncompiler:%v\nplatform:%v\n",
		vfo.GitCommit, vfo.Date, vfo.GoVersion, vfo.Compiler, vfo.Platform)
}

const (
	VersionFalse    versionValue = 0
	VersionTrue     versionValue = 1
	VersionRaw      versionValue = 2
	strRawVersion   string       = "raw"
	versionFlagName              = "version"
)

// The versionValue type implementation pflag.Value interface
type versionValue int

func (v *versionValue) Set(s string) error {

	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

func (v *versionValue) Type() string {
	return versionFlagName
}

var versionFlag = version(versionFlagName, VersionFalse, "Print version information and quit")

func versionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	pflag.Lookup(name).NoOptDefVal = "true"
}

func version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	versionVar(p, name, value, usage)
	return p
}

func VersionAddFlags(fs *pflag.FlagSet, name string) {
	fs.AddFlag(pflag.Lookup("version"))
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}

// PrintVersionAndExitIfRequested used output app version info
// if the -version flag was passed
func PrintVersionAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", formatVersion())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s:\n%s\n", "[X-APIServer]", formatVersion())
		os.Exit(0)
	}
}
