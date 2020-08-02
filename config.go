/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

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
	"github.com/mwat56/pageview"
	"github.com/mwat56/whitespace"
)

type (
	// TAppArgs Collection of commandline arguments and INI values.
	TAppArgs struct {
		AccessLog     string // (optional) name of page access logfile
		Addr          string // listen address ("1.2.3.4:5678")
		BlogName      string // name/description of this blog
		CertKey       string // TLS certificate key
		CertPem       string // private TLS certificate
		DataDir       string // base directory of application's data
		delWhitespace bool   // remove whitespace from generated pages
		dump          bool   // Debug: dump this structure to `StdOut`
		ErrorLog      string // (optional) name of page error logfile
		GZip          bool   // send compressed data to remote browser
		HashFile      string // file with hashtag/mention database
		// Intl       string // path/filename of the localisation file
		Lang     string // default GUI language
		listen   string // IP of host to listen at
		LogStack bool   // log stack trace in case of errors

		MaxFileSize int64  // max. upload file size
		mfs         string // max. upload file size

		PageView   bool   // wether to use page previews or not
		PostAdd    bool   // whether to write a posting from commandline
		PostFile   string // name of file to post
		port       int    // port to listen to
		Realm      string // host/domain to secure by BasicAuth
		Theme      string // `dark` or `light` display theme
		UserAdd    string // username to add to password list
		UserCheck  string // username to check in password list
		UserDelete string // username to delete from password list
		UserFile   string // (optional) name of page access logfile
		UserList   bool   // print out a list of current users
		UserUpdate string // username to update in password list
	}

	// tArguments is the list structure for the cmdline argument values
	// merged with the key-value pairs from the INI file.
	tArguments struct {
		ini.TSection // embedded INI section
	}
)

var (
	// AppArgs holds the commandline arguments and INI values.
	//
	// This structure should be considered R/O after it was
	// set up by a call to `InitConfig()`.
	AppArgs TAppArgs

	// iniValues is the list for the cmdline arguments and INI values.
	iniValues tArguments
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

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
//	`aBaseDir` The base directory to prepend to `aDir`.
//	`aDir` The directory to make absolute.
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

// String implements the `Stringer` interface returning a (pretty printed)
// string representation of the current `TAppArgs` instance.
//
// NOTE: This method is meant mostly for debugging purposes.
func (aa TAppArgs) String() string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				strings.Replace(
					fmt.Sprintf("%#v", aa),
					`, `, ",\n\t", -1),
				`{`, "{\n\t", -1),
			`}`, ",\n}", -1),
		`:`, ` : `, -1) //FIXME this affects property values as well!
} // String()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// RegEx to match a size value (xxx)
	cfKmgRE = regexp.MustCompile(`(?i)\s*(\d+)\s*([bgkm]+)?`)
)

// `kmg2Num()` returns a 'B|KB|MB|GB` string as an integer.
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

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

/*
func init() {
	// // see: https://github.com/microsoft/vscode-go/issues/2734
	// testing.Init() // workaround for Go 1.13
	InitConfig()
} // init()
*/

// InitConfig reads both the INI values and the commandline arguments.
//
// The steps here are:
//
// (1) read the INI file(s):
//	(a) read the local `./.nele.ini`,
//	(b) read the global `/etc/nele.ini`,
//	(c) read the user-local `~/.nele.ini`,
//	(d) read the user-local `~/.config/nele.ini`,
// (2) merge the commandline arguments with the INI values into
// the global `AppArgs` variable.
//
// This function is meant to be called first thing in the application's
// `main()` function.
func InitConfig() {
	flag.CommandLine = flag.NewFlagSet(`Kaliber`, flag.ExitOnError)
	iniValues = tArguments{*ini.ReadIniData(`nele`)}

	setAppArgs()
	parseAppArgs()
	readAppArgs()
} // InitConfig()

// `parseAppArgs()` parsed the actual commandline arguments.
func parseAppArgs() {
	flag.CommandLine.Usage = ShowHelp
	_ = flag.CommandLine.Parse(os.Args[1:])
} // parseAppArgs()

// `readAppArgs()` copies the commandline values into the `TAppArgs` instance.
func readAppArgs() {
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
		log.Fatalf("`dataDir` == `%s` problem: %v", AppArgs.DataDir, err)
	} else if !fi.IsDir() {
		log.Fatalf("Error: `dataDir` not a directory `%s`", AppArgs.DataDir)
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

	if AppArgs.PageView {
		_ = pageview.SetImageDirectory(absolute(AppArgs.DataDir, `img`))
		pageview.SetImageFileType(`png`)
		pageview.SetJavaScript(false)
		pageview.SetMaxAge(0)
		// pageview.SetUserAgent(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/78.0.3904.108 Chrome/78.0.3904.108 Safari/537.36`)
		// Doesn't work with Facebook:
		// pageview.SetUserAgent(`Lynx/2.8.9dev.16 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/3.5.17`)
		pageview.SetUserAgent(`Lynx/2.9.0dev.5 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/3.5.17`)
		// see: https://lynx.invisible-island.net/current/CHANGES.html
		// pageview.SetUserAgent(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.72`)
	}

	if 0 < len(AppArgs.PostFile) {
		AppArgs.PostFile = absolute(AppArgs.DataDir, AppArgs.PostFile)
	}

	if 0 == len(AppArgs.Realm) {
		AppArgs.Realm = `My Blog`
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

	if AppArgs.dump {
		// Print out the arguments and terminate:
		log.Fatalf("runtime arguments:\n%s", AppArgs.String())
	}

	flag.CommandLine = nil // free unneeded memory
} // readAppArgs()

// `setAppArgs()` reads the commandline arguments into a list
// structure merging it with key-value pairs read from INI file(s).
func setAppArgs() {
	var (
		ok bool
		s  string // temp. value
	)

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
	flag.CommandLine.BoolVar(&AppArgs.dump, `d`, AppArgs.dump, `dump`)

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

	if AppArgs.port, ok = iniValues.AsInt(`port`); (!ok) || (0 == AppArgs.port) {
		AppArgs.port = 8181
	}
	flag.CommandLine.IntVar(&AppArgs.port, `port`, AppArgs.port,
		"<port number> The IP port to listen to ")

	flag.CommandLine.BoolVar(&AppArgs.PostAdd, `pa`, AppArgs.PostAdd,
		"<boolean> (optional) posting add: write a posting from the commandline")

	flag.CommandLine.StringVar(&AppArgs.PostFile, `pf`, AppArgs.PostFile,
		"<fileName> (optional) post file: name of a file to add as new posting")

	AppArgs.PageView, _ = iniValues.AsBool(`pageView`)
	flag.CommandLine.BoolVar(&AppArgs.PageView, `pv`, AppArgs.PageView,
		"<boolean> Use page preview images for links")

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

	iniValues.Clear()           // release unneeded memory
	iniValues = tArguments{nil} // dito
} // setAppArgs()()

// ShowHelp lists the commandline options to `Stderr`.
func ShowHelp() {
	fmt.Fprintf(os.Stderr, "\n  Usage: %s [OPTIONS]\n\n", os.Args[0])
	flag.CommandLine.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Most options can be set in an INI file to keep the command-line short ;-)")
} // ShowHelp()

/* _EoF_ */
