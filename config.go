/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1005 - capitalisation wanted
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
	"time"

	"github.com/mwat56/ini"
	"github.com/mwat56/pageview"
)

type (
	// tArguments is the list structure for the cmdline argument values
	// merged with the key-value pairs from the INI file.
	tArguments struct {
		ini.TSection // embedded INI section
	}
)

var (
	// AppArguments is the list for the cmdline arguments and INI values.
	AppArguments tArguments

	// RegEx to match a size value (xxx)
	cfKmgRE = regexp.MustCompile(`(?i)\s*(\d+)\s*([bgkm]+)?`)
)

// `set()` adds/sets another key-value pair.
//
// If `aValue` is empty then `aKey` gets removed.
func (al *tArguments) set(aKey, aValue string) {
	if 0 < len(aValue) {
		al.AddKey(aKey, aValue)
	} else {
		al.RemoveKey(aKey)
	}
} // set()

// Get returns the value associated with `aKey` and `nil` if found,
// or an empty string and an error.
//
//	`aKey` The key to lookup in the list.
func (al *tArguments) Get(aKey string) (string, error) {
	if result, ok := al.AsString(aKey); ok {
		return result, nil
	}

	return "", fmt.Errorf("Missing config value: %s", aKey)
} // Get()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `absolute()` returns `aDir` as an absolute path.
func absolute(aBaseDir, aDir string) string {
	if 0 == len(aDir) {
		return aDir
	}
	if '/' == aDir[0] {
		s, _ := filepath.Abs(aDir)
		return s
	}

	s, _ := filepath.Abs(filepath.Join(aBaseDir, aDir))
	return s
} // absolute()

// `kmg2Num()` returns a 'KB|MB|GB` string as an integer.
func kmg2Num(aString string) (rInt int) {
	matches := cfKmgRE.FindStringSubmatch(aString)
	if 2 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all Atoi() errors.
		rInt, _ = strconv.Atoi(matches[1])
		switch strings.ToLower(matches[2]) {
		case "", "b":
			return
		case "kb":
			rInt *= 1024
		case "mb":
			rInt *= 1024 * 1024
		case "gb":
			rInt *= 1024 * 1024 * 1024
		}
	}

	return
} // kmg2Num()

// `readIniData()` returns the config values read from INI file(s).
//
//	The steps here are:
//	(1) read the local `./.nele.ini`,
//	(2) read the global `/etc/nele.ini`,
//	(3) read the user-local `~/.nele.ini`,
//	(4) read the user-local `~/.config/nele.ini`,
//	(5) read the `-ini` commandline argument.
func readIniData() {
	// (1) ./
	fName, _ := filepath.Abs("./nele.ini")
	ini1, err := ini.New(fName)
	if nil == err {
		ini1.AddSectionKey("", "iniFile", fName)
	}

	// (2) /etc/
	fName = "/etc/nele.ini"
	if ini2, err2 := ini.New(fName); nil == err2 {
		ini1.Merge(ini2)
		ini1.AddSectionKey("", "iniFile", fName)
	}

	// (3) ~user/
	fName, _ = os.UserHomeDir()
	if 0 < len(fName) {
		fName, _ = filepath.Abs(filepath.Join(fName, ".nele.ini"))
		if ini2, err2 := ini.New(fName); nil == err2 {
			ini1.Merge(ini2)
			ini1.AddSectionKey("", "iniFile", fName)
		}
	}

	// (4) ~/.config/
	if confDir, err2 := os.UserConfigDir(); nil == err2 {
		fName, _ = filepath.Abs(filepath.Join(confDir, "nele.ini"))
		if ini2, err2 := ini.New(fName); nil == err2 {
			ini1.Merge(ini2)
			ini1.AddSectionKey("", "iniFile", fName)
		}
	}

	// (5) cmdline
	aLen := len(os.Args)
	for i := 1; i < aLen; i++ {
		if `-ini` == os.Args[i] {
			//XXX Note that this works only if `-ini` and
			// filename are two separate arguments. It will
			// fail if it's given in the form `-ini=filename`.
			i++
			if i < aLen {
				fName, _ = filepath.Abs(os.Args[i])
				if ini2, _ := ini.New(fName); nil == err {
					ini1.Merge(ini2)
					ini1.AddSectionKey("", "iniFile", fName)
				}
			}
			break
		}
	}

	AppArguments = tArguments{*ini1.GetSection("")}
} // readIniData()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

/*
func init() {
	// // see: https://github.com/microsoft/vscode-go/issues/2734
	// testing.Init() // workaround for Go 1.13
	InitConfig()
} // init()
*/

// InitConfig reads the commandline arguments into a list
// structure merging it with key-value pairs read from INI file(s).
//
// The steps here are:
// (1) read the INI file(s),
// (2) merge the commandline arguments the INI values
// into the global `AppArguments` variable.
//
// This function is meant to be called first thing in the program's `main()`.
func InitConfig() {
	readIniData()

	blogName, _ := AppArguments.Get("blogName")
	flag.StringVar(&blogName, "blogName", blogName,
		"Name of this Blog (shown on every page)\n")

	s, _ := AppArguments.Get("dataDir")
	dataDir, _ := filepath.Abs(s)
	flag.StringVar(&dataDir, "dataDir", dataDir,
		"<dirName> the directory with CSS, IMG, JS, POSTINGS, STATIC, VIEWS sub-directories\n")

	s, _ = AppArguments.Get("accessLog")
	accessLog := absolute(dataDir, s)
	flag.StringVar(&accessLog, "accessLog", accessLog,
		"<filename> Name of the access logfile to write to\n")

	s, _ = AppArguments.Get("certKey")
	certKey := absolute(dataDir, s)
	flag.StringVar(&certKey, "certKey", certKey,
		"<fileName> the name of the TLS certificate's private key\n")

	s, _ = AppArguments.Get("certPem")
	certPem := absolute(dataDir, s)
	flag.StringVar(&certPem, "certPem", certPem,
		"<fileName> the name of the TLS certificate PEM\n")

	s, _ = AppArguments.Get("errorLog")
	errorLog := absolute(dataDir, s)
	flag.StringVar(&errorLog, "errorLog", errorLog,
		"<filename> Name of the error logfile to write to\n")

	gzip, _ := AppArguments.AsBool("gzip")
	flag.BoolVar(&gzip, "gzip", gzip,
		"(optional) use gzip compression for server responses")

	s, _ = AppArguments.Get("hashFile")
	hashFile := absolute(dataDir, s)
	flag.StringVar(&hashFile, "hashFile", hashFile,
		"<fileName> (optional) the name of a file storing #hashtags and @mentions\n")

	iniFile, _ := AppArguments.Get("iniFile")
	flag.StringVar(&iniFile, "ini", iniFile,
		"<fileName> the path/filename of the INI file to use\n")

	language, _ := AppArguments.Get("lang")
	flag.StringVar(&language, "lang", language,
		"(optional) the default language to use ")

	listenStr, _ := AppArguments.Get("listen")
	flag.StringVar(&listenStr, "listen", listenStr,
		"the host's IP to listen at ")

	logStack, _ := AppArguments.AsBool("logStack")
	flag.BoolVar(&logStack, "logStack", logStack,
		"<boolean> Log a stack trace for recovered runtime errors ")

	maxFileSize, _ := AppArguments.Get("maxfilesize")
	flag.StringVar(&maxFileSize, "maxfilesize", maxFileSize,
		"max. accepted size of uploaded files")

	pageView, _ := AppArguments.AsBool("pageView")
	flag.BoolVar(&pageView, "pageView", pageView,
		"(optional) use page preview images for links")

	portInt, _ := AppArguments.AsInt("port")
	flag.IntVar(&portInt, "port", portInt,
		"<portNumber> the IP port to listen to ")

	postAdd := false
	flag.BoolVar(&postAdd, "pa", postAdd,
		"(optional) posting add: write a posting from the commandline")

	postFile := ""
	flag.StringVar(&postFile, "pf", postFile,
		"<fileName> (optional) post file: name of a file to add as new posting")

	realmStr, _ := AppArguments.Get("realm")
	flag.StringVar(&realmStr, "realm", realmStr,
		"(optional) <hostName> name of host/domain to secure by BasicAuth\n")

	themeStr, _ := AppArguments.Get("theme")
	flag.StringVar(&themeStr, "theme", themeStr,
		"<name> the display theme to use ('light' or 'dark')\n")

	userAdd := ""
	flag.StringVar(&userAdd, "ua", userAdd,
		"<userName> (optional) user add: add a username to the password file")

	userChange := ""
	flag.StringVar(&userChange, "uc", userChange,
		"<userName> (optional) user check: check a username in the password file")

	s, _ = AppArguments.Get("passFile")
	userFile := absolute(dataDir, s)
	flag.StringVar(&userFile, "uf", userFile,
		"<fileName> (optional) user passwords file storing user/passwords for BasicAuth\n")

	userList := false
	flag.BoolVar(&userList, "ul", userList,
		"(optional) user list: show all users in the password file")

	userName := ""
	flag.StringVar(&userName, "ud", userName,
		"<userName> (optional) user delete: remove a username from the password file")

	userUpdate := ""
	flag.StringVar(&userUpdate, "uu", userUpdate,
		"<userName> (optional) user update: update a username in the password file")

	flag.Usage = ShowHelp
	flag.Parse() // // // // // // // // // // // // // // // // // // //

	if 0 == len(blogName) {
		blogName = time.Now().Format("2006:01:02:15:04:05")
	}
	AppArguments.set("blogName", blogName)

	if 0 < len(dataDir) {
		dataDir, _ = filepath.Abs(dataDir)
	}
	if f, err := os.Stat(dataDir); nil != err {
		log.Fatalf("datadir == %s` problem: %v", dataDir, err)
	} else if !f.IsDir() {
		log.Fatalf("Error: Not a directory `%s`", dataDir)
	}
	AppArguments.set("dataDir", dataDir)

	// `postingBaseDirectory` defined in `posting.go`:
	SetPostingBaseDirectory(filepath.Join(dataDir, "./postings"))

	if 0 < len(accessLog) {
		accessLog = absolute(dataDir, accessLog)
	}
	AppArguments.set("accessLog", accessLog)

	if 0 < len(certKey) {
		certKey = absolute(dataDir, certKey)
		if fi, err := os.Stat(certKey); (nil != err) || (0 == fi.Size()) {
			certKey = ""
		}
	}
	AppArguments.set("certKey", certKey)

	if 0 < len(certPem) {
		certPem = absolute(dataDir, certPem)
		if fi, err := os.Stat(certPem); (nil != err) || (0 == fi.Size()) {
			certPem = ""
		}
	}
	AppArguments.set("certPem", certPem)

	if 0 < len(errorLog) {
		errorLog = absolute(dataDir, errorLog)
	}
	AppArguments.set("errorLog", errorLog)

	if gzip {
		AppArguments.set("gzip", "true")
	} else {
		AppArguments.set("gzip", "")
	}

	if 0 < len(hashFile) {
		hashFile = absolute(dataDir, hashFile)
	}
	AppArguments.set("hashFile", hashFile)

	if 0 == len(language) {
		language = "en"
	}
	AppArguments.set("lang", language)

	if "0" == listenStr {
		listenStr = ""
	}
	AppArguments.set("listen", listenStr)

	if logStack {
		AppArguments.set("logStack", "true")
	} else {
		AppArguments.set("logStack", "")
	}

	if 0 == len(maxFileSize) {
		maxFileSize = "10485760" // 10 MB
	} else {
		maxFileSize = fmt.Sprintf("%d", kmg2Num(maxFileSize))
	}
	AppArguments.set("mfs", maxFileSize)

	if pageView {
		_ = pageview.SetImageDirectory(absolute(dataDir, "img"))
		pageview.SetImageFileType(`png`)
		pageview.SetJavaScript(false)
		pageview.SetMaxAge(0)
		// pageview.SetUserAgent(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/78.0.3904.108 Chrome/78.0.3904.108 Safari/537.36`)
		// Doesn't work with Facebook:
		pageview.SetUserAgent(`Lynx/2.8.9dev.16 libwww-FM/2.14 SSL-MM/1.4.1 GNUTLS/3.5.17`)
		// pageview.SetUserAgent(`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36 OPR/66.0.3515.72`)
		AppArguments.set("pageView", "true")
	} else {
		AppArguments.set("pageView", "")
	}

	AppArguments.set("port", fmt.Sprintf("%d", portInt))

	if postAdd {
		AppArguments.set("pa", "true")
	} else {
		AppArguments.set("pa", "")
	}

	if 0 < len(postFile) {
		postFile = absolute(dataDir, postFile)
	}
	AppArguments.set("pf", postFile)
	AppArguments.set("realm", realmStr)
	AppArguments.set("theme", themeStr)
	AppArguments.set("ua", userAdd)
	AppArguments.set("uc", userChange)
	AppArguments.set("ud", userName)

	if 0 < len(userFile) {
		userFile = absolute(dataDir, userFile)
	}
	AppArguments.set("uf", userFile)

	if userList {
		AppArguments.set("ul", "true")
	} else {
		AppArguments.set("ul", "")
	}

	AppArguments.set("uu", userUpdate)
} // InitConfig()

// ShowHelp lists the commandline options to `Stderr`.
func ShowHelp() {
	fmt.Fprintf(os.Stderr, "\n  Usage: %s [OPTIONS]\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Most options can be set in an INI file to keep the command-line short ;-)")
} // ShowHelp()

/* _EoF_ */
