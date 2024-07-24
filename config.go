/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mwat56/ini"
	"github.com/mwat56/screenshot"
	"github.com/mwat56/whitespace"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

var (
	// `iniValues` is the list for the cmdline arguments and INI
	// values used during program startup.
	iniValues *ini.TSection // embedded INI section // tArguments
)

// --------------------------------------------------------------------------

// `absolute()` returns `aDir` as an absolute path.
//
// If `aDir` is an empty string the current directory (`./`) gets returned.
//
// If `aDir` starts with a slash (`/`) it's returned after cleaning.
//
// If `aBaseDir` is an empty string the current directory (`./`) is used.
//
// Otherwise `aBaseDir` gets prepended to `aDir` and returned after cleaning.
//
// Parameters:
//   - `aBaseDir` The base directory to prepend to `aDir`.
//   - `aDir` The directory to make absolute.
func absolute(aBaseDir, aDir string) string {
	if 0 == len(aDir) {
		aDir, _ = filepath.Abs(`./`)
		return aDir
	}
	if '/' == aDir[0] {
		return filepath.Clean(aDir)
	}
	if 0 == len(aBaseDir) {
		aBaseDir, _ = filepath.Abs(`./`)
	}

	return filepath.Join(aBaseDir, aDir)
} // absolute()

// --------------------------------------------------------------------------

var (
	// RegEx to match a size value (xxx)
	cfKmgRE = regexp.MustCompile(`(?i)\s*(\d+)\s*([bgkm]+)?`)
)

// `kmg2Num()` returns a 'B|KB|MB|GB` string as an integer.
//
// Returns:
//   - `int64`: The integer value of `aString`.
func kmg2Num(aString string) (rInt int64) {
	matches := cfKmgRE.FindStringSubmatch(aString)
	if 2 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all ParseInt() errors.
		rInt, _ = strconv.ParseInt(matches[1], 10, 64)
		switch strings.ToLower(matches[2]) {
		case ``, `b`:
			return
		case `kb`:
			rInt *= 1024
		case `mb`:
			rInt *= 1024 * 1024
		case `gb`:
			rInt *= 1024 * 1024 * 1024
		}
	}

	return
} // kmg2Num()

// --------------------------------------------------------------------------

/*
func init() {
	// // see: https://github.com/microsoft/vscode-go/issues/2734
	// testing.Init() // workaround for Go 1.13
	InitConfig()
} // init()
*/

// `copyIniDataToAppArgs()` copies the commandline and INI values
// into the global `TAppArgs` instance.
func copyIniDataToAppArgs() {
	var ( // re-use variables
		err error
		fi  os.FileInfo
	)
	if 0 == len(AppArgs.BlogName) {
		// AppArgs.BlogName = time.Now().Format("2006:01:02:15:04:05")
		AppArgs.BlogName = `<! BlogName not configured !>`
	}

	if 0 == len(AppArgs.DataDir) {
		AppArgs.DataDir = `./`
	}
	AppArgs.DataDir, _ = filepath.Abs(AppArgs.DataDir)
	if fi, err = os.Stat(AppArgs.DataDir); nil != err {
		log.Fatalf("`dataDir` == %q problem: %v", AppArgs.DataDir, err)
	} else if !fi.IsDir() {
		log.Fatalf("Error: `dataDir` not a directory %q", AppArgs.DataDir)
	}
	// `postingBaseDirectory` defined in `posting.go`:
	SetPostingBaseDirectory(filepath.Join(AppArgs.DataDir, "./postings"))

	if 0 < len(AppArgs.AccessLog) {
		AppArgs.AccessLog = absolute(AppArgs.DataDir, AppArgs.AccessLog)
	}

	if 0 < len(AppArgs.CertKey) {
		AppArgs.CertKey = absolute(AppArgs.DataDir, AppArgs.CertKey)
		if fi, err = os.Stat(AppArgs.CertKey); (nil != err) || (0 >= fi.Size()) {
			AppArgs.CertKey = ``
		}
	}

	if 0 < len(AppArgs.CertPem) {
		AppArgs.CertPem = absolute(AppArgs.DataDir, AppArgs.CertPem)
		if fi, err = os.Stat(AppArgs.CertPem); (nil != err) || (0 >= fi.Size()) {
			AppArgs.CertPem = ``
		}
	}

	whitespace.UseRemoveWhitespace = AppArgs.delWhitespace

	if 0 < len(AppArgs.ErrorLog) {
		AppArgs.ErrorLog = absolute(AppArgs.DataDir, AppArgs.ErrorLog)
	}

	if 0 < len(AppArgs.HashFile) {
		AppArgs.HashFile = absolute(AppArgs.DataDir, AppArgs.HashFile)
	} else {
		log.Fatalln("Error: `hashFile` argument missing")
	}

	if 0 < len(AppArgs.Lang) {
		AppArgs.Lang = strings.ToLower(AppArgs.Lang)
	}
	switch AppArgs.Lang {
	case `de`, `en`:
		// accepted values

	default:
		AppArgs.Lang = `en`
	}

	if `0` == AppArgs.listen {
		AppArgs.listen = ``
	}
	if 0 >= AppArgs.port {
		AppArgs.port = 8181
	}
	// an empty `listen` value means: listen on all interfaces
	AppArgs.Addr = fmt.Sprintf("%s:%d", AppArgs.listen, AppArgs.port)

	if 0 == len(AppArgs.mfs) {
		AppArgs.MaxFileSize = 10485760 // 10 MB
	} else {
		AppArgs.MaxFileSize = kmg2Num(AppArgs.mfs)
	}

	if 0 < len(AppArgs.PostFile) {
		AppArgs.PostFile = absolute(AppArgs.DataDir, AppArgs.PostFile)
	}

	if 0 == len(AppArgs.Realm) {
		AppArgs.Realm = `My Blog`
	}

	if AppArgs.Screenshot {
		processScreenshotOptions()
	}

	if 0 < len(AppArgs.Theme) {
		AppArgs.Theme = strings.ToLower(AppArgs.Theme)
	}
	switch AppArgs.Theme {
	case `dark`, `light`:
		// accepted values

	default:
		AppArgs.Theme = `dark`
	}

	if 0 < len(AppArgs.UserFile) {
		AppArgs.UserFile = absolute(AppArgs.DataDir, AppArgs.UserFile)
	}

	if AppArgs.Dump {
		// Print out the arguments and terminate:
		log.Fatalf("runtime arguments:\n%s", AppArgs.String())
	}
} // copyIniDataToAppArgs()

// `InitConfig()` reads both the INI values and the commandline arguments.
//
// The steps here are:
//
// (1) read the INI file(s):
//
//	(a) read the local `./.nele.ini`,
//	(b) read the global `/etc/nele.ini`,
//	(c) read the user-local `~/.nele.ini`,
//	(d) read the user-local `~/.config/nele.ini`,
//
// (2) merge the commandline arguments with the INI values into
// the global `AppArgs` variable.
//
// This function is meant to be called first thing in the application's
// `main()` function.
func InitConfig() {
	// `InitConfig()` calls `flag.parse()` which in turn will cause
	// errors when run with `go test …`.

	const appName string = `nele`

	section, _ := ini.ReadIniData(appName)
	iniValues = section

	readCmdlineArgs()

	parseCmdlineArgs()

	copyIniDataToAppArgs()

	var persistence IPersistence
	switch AppArgs.persistence {
	case `db`:
		persistence = NewDBpersistence(AppArgs.Name)

	case `fs`:
		fallthrough

	default:
		persistence = NewFSpersistence()
	}
	SetPersistence(IPersistence(persistence))

} // InitConfig()

// `parseCmdlineArgs()` parses the actual commandline arguments.
func parseCmdlineArgs() {
	defer func() {
		// make sure a `panic` won't kill the program
		if err := recover(); nil != err {
			return
		}
	}()

	flag.CommandLine.Usage = ShowHelp
	_ = flag.CommandLine.Parse(os.Args[1:])
} // parseCmdlineArgs()

// `processScreenshotOptions()` handles the INI- and commandline options
// for the screenshot facility.
func processScreenshotOptions() {
	// Don't depend on defaults set by the package
	ssOptions := screenshot.Options()

	//TODO make this values configurable by INI and cmdline.

	ssOptions.AcceptOther = true
	ssOptions.CertErrors = false
	ssOptions.Cookies = false
	ssOptions.HostsAvoidJSfile = absolute(`./`, screenshot.HostsAvoidJS)
	ssOptions.HostsNeedJSfile = absolute(`./`, screenshot.HostsNeedJS)
	ssOptions.ImageAge = 0
	ssOptions.ImageDir = absolute(AppArgs.DataDir, `img`)
	ssOptions.ImageHeight = 800
	ssOptions.ImageOverwrite = false
	ssOptions.ImageQuality = 75
	ssOptions.ImageScale = 0
	ssOptions.ImageWidth = 800
	ssOptions.JavaScript = false
	ssOptions.MaxProcessTime = 32
	ssOptions.Mobile = false // !important
	ssOptions.Platform = screenshot.DefaultPlatform
	ssOptions.Scrollbars = false
	// We need an agent that is accepted by `Facebook`, `Twitter`,
	// and `YouTube` at least because still a lot of sites do some
	// sort of browser-sniffing (instead of capability checking).
	ssOptions.UserAgent = screenshot.DefaultAgent
	ssOptions.Do() // activate the settings
} // processScreenshotOptions()

// `readCmdlineArgs()` reads the commandline arguments into a list
// structure merging it with key-value pairs read from INI file(s).
func readCmdlineArgs() {
	var (
		ok bool
		s  string // temp. value
	)
	defer func() {
		// make sure a `panic` won't kill the program
		if err := recover(); nil != err {
			ok = false
		}
	}()

	if s, ok = iniValues.AsString(`dataDir`); (ok) && (0 < len(s)) {
		AppArgs.DataDir, _ = filepath.Abs(s)
	} else {
		AppArgs.DataDir, _ = filepath.Abs(`./`)
	}
	flag.CommandLine.StringVar(&AppArgs.DataDir, `dataDir`, AppArgs.DataDir,
		"<dirName> Directory with CSS, FONTS, IMG, SESSIONS, and VIEWS sub-directories\n")

	if AppArgs.BlogName, ok = iniValues.AsString(`blogName`); (!ok) || (0 == len(AppArgs.BlogName)) {
		AppArgs.BlogName = `<! BlogName not configured !>`
	}
	flag.CommandLine.StringVar(&AppArgs.BlogName, `blogName`, AppArgs.BlogName,
		"<string> Name of this Blog (shown on every page)\n")

	if s, ok = iniValues.AsString(`accessLog`); (ok) && (0 < len(s)) {
		AppArgs.AccessLog = absolute(AppArgs.DataDir, s)
	}
	flag.CommandLine.StringVar(&AppArgs.AccessLog, `accessLog`, AppArgs.AccessLog,
		"<filename> Name of the access logfile to write to\n")

	if s, ok = iniValues.AsString(`certKey`); (ok) && (0 < len(s)) {
		AppArgs.CertKey = absolute(AppArgs.DataDir, s)
	}
	flag.CommandLine.StringVar(&AppArgs.CertKey, `certKey`, AppArgs.CertKey,
		"<fileName> Name of the TLS certificate key\n")

	if s, ok = iniValues.AsString(`certPem`); (ok) && (0 < len(s)) {
		AppArgs.CertPem = absolute(AppArgs.DataDir, s)
	}
	flag.CommandLine.StringVar(&AppArgs.CertPem, `certPem`, AppArgs.CertPem,
		"<fileName> Name of the TLS certificate PEM\n")

	if AppArgs.delWhitespace, ok = iniValues.AsBool(`delWhitespace`); !ok {
		AppArgs.delWhitespace = true
	}
	flag.CommandLine.BoolVar(&AppArgs.delWhitespace, `delWhitespace`, AppArgs.delWhitespace,
		"<boolean> Delete superfluous whitespace in generated pages")

	// Debug aid:
	flag.CommandLine.BoolVar(&AppArgs.Dump, `d`, AppArgs.Dump, `dump`)

	if s, ok = iniValues.AsString(`errorLog`); (ok) && (0 < len(s)) {
		AppArgs.ErrorLog = absolute(AppArgs.DataDir, s)
	}
	flag.CommandLine.StringVar(&AppArgs.ErrorLog, `errorlog`, AppArgs.ErrorLog,
		"<filename> Name of the error logfile to write to\n")

	if AppArgs.GZip, ok = iniValues.AsBool(`gzip`); !ok {
		AppArgs.GZip = true
	}
	flag.CommandLine.BoolVar(&AppArgs.GZip, `gzip`, AppArgs.GZip,
		"<boolean> use gzip compression for server responses")

	if s, ok = iniValues.AsString(`hashFile`); (ok) && (0 < len(s)) {
		AppArgs.HashFile = absolute(AppArgs.DataDir, s)
	} else {
		AppArgs.HashFile = absolute(AppArgs.DataDir, `HashFile.db`)
	}
	flag.CommandLine.StringVar(&AppArgs.HashFile, `hashFile`, AppArgs.HashFile,
		"<fileName> Name of the file storing #hashtags and @mentions\n")

	iniFile, _ := iniValues.AsString(`iniFile`)
	flag.CommandLine.StringVar(&iniFile, `ini`, iniFile,
		"<fileName> the path/filename of the INI file to use\n")

	if AppArgs.Lang, ok = iniValues.AsString(`lang`); (!ok) || (0 == len(AppArgs.Lang)) {
		AppArgs.Lang = `en`
	}
	flag.CommandLine.StringVar(&AppArgs.Lang, `lang`, AppArgs.Lang,
		"<de|en> Default language to use ")

	if AppArgs.listen, ok = iniValues.AsString(`listen`); (!ok) || (0 == len(AppArgs.listen)) {
		AppArgs.listen = `127.0.0.1`
	}
	flag.CommandLine.StringVar(&AppArgs.listen, `listen`, AppArgs.listen,
		"<IP number> The host's IP to listen at ")

	AppArgs.LogStack, _ = iniValues.AsBool(`logStack`)
	flag.CommandLine.BoolVar(&AppArgs.LogStack, "lst", AppArgs.LogStack,
		"<boolean> Log a stack trace for recovered runtime errors ")

	if AppArgs.mfs, ok = iniValues.AsString(`maxfilesize`); ok && (0 < len(AppArgs.mfs)) {
		AppArgs.mfs = strings.ToLower(AppArgs.mfs)
	} else {
		AppArgs.mfs = `10485760` // 10 MB
	}
	flag.CommandLine.StringVar(&AppArgs.mfs, `mfs`, AppArgs.mfs,
		"<filesize> Max. accepted size of uploaded files")

	//
	//TODO: read persistence (db|fs|tee)
	//

	AppArgs.port, ok = iniValues.AsInt(`port`)
	if (!ok) || (0 == AppArgs.port) {
		AppArgs.port = 8181
	}
	flag.CommandLine.IntVar(&AppArgs.port, `port`, AppArgs.port,
		"<port number> The IP port to listen to ")

	flag.CommandLine.BoolVar(&AppArgs.PostAdd, `pa`, AppArgs.PostAdd,
		"<boolean> (optional) posting add: write a posting from the commandline")

	flag.CommandLine.StringVar(&AppArgs.PostFile, `pf`, AppArgs.PostFile,
		"<fileName> (optional) post file: name of a file to add as new posting")

	AppArgs.Screenshot, _ = iniValues.AsBool(`Screenshot`)
	flag.CommandLine.BoolVar(&AppArgs.Screenshot, `pv`, AppArgs.Screenshot,
		"<boolean> Use page preview/screenshot images for links")

	//TODO implement various screenshot options

	if AppArgs.Realm, ok = iniValues.AsString(`realm`); (!ok) || (0 == len(AppArgs.Realm)) {
		AppArgs.Realm = `My Blog`
	}
	flag.CommandLine.StringVar(&AppArgs.Realm, `realm`, AppArgs.Realm,
		"<hostName> Name of host/domain to secure by BasicAuth\n")

	if AppArgs.Theme, _ = iniValues.AsString(`theme`); 0 < len(AppArgs.Theme) {
		AppArgs.Theme = strings.ToLower(AppArgs.Theme)
	}
	switch AppArgs.Theme {
	case `dark`, `light`:
		// these are okay
	default:
		AppArgs.Theme = `dark`
	}
	flag.CommandLine.StringVar(&AppArgs.Theme, `theme`, AppArgs.Theme,
		"<name> The display theme to use ('light' or 'dark')\n")

	flag.CommandLine.StringVar(&AppArgs.UserAdd, `ua`, AppArgs.UserAdd,
		"<userName> User add: add a username to the password file")

	flag.CommandLine.StringVar(&AppArgs.UserCheck, `uc`, AppArgs.UserCheck,
		"<userName> User check: check a username in the password file")

	flag.CommandLine.StringVar(&AppArgs.UserDelete, `ud`, AppArgs.UserDelete,
		"<userName> User delete: remove a username from the password file")

	if s, ok = iniValues.AsString(`passFile`); ok && (0 < len(s)) {
		AppArgs.UserFile = absolute(AppArgs.DataDir, s)
	}
	flag.CommandLine.StringVar(&AppArgs.UserFile, `uf`, AppArgs.UserFile,
		"<fileName> Passwords file storing user/passwords for BasicAuth\n")

	flag.CommandLine.BoolVar(&AppArgs.UserList, `ul`, AppArgs.UserList,
		"<boolean> User list: show all users in the password file")

	flag.CommandLine.StringVar(&AppArgs.UserUpdate, `uu`, AppArgs.UserUpdate,
		"<userName> User update: update a username in the password file")

	iniValues.Clear() // release unneeded memory
} // readCmdlineArgs()

// ShowHelp lists the commandline options to `Stderr`.
func ShowHelp() {
	fmt.Fprintf(os.Stderr, "\n\tUsage: %s [OPTIONS]\n\n", os.Args[0])
	flag.CommandLine.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n\tMost options can be set in an INI file to keep the command-line short ;-)")
} // ShowHelp()

/* _EoF_ */
