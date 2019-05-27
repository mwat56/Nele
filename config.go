/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package nele

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mwat56/ini"
)

type (
	// tAguments is the list structure for the cmdline argument values
	// merged with the key-value pairs from the INI file.
	tAguments map[string]string
)

var (
	// AppArguments is the list for the cmdline arguments and INI values.
	AppArguments tAguments
)

// `set()` adds/sets another key-value pair.
func (al tAguments) set(aKey string, aValue string) {
	al[aKey] = aValue
} // set()

// Get returns the value associated with `aKey` and `nil` if found,
// or an empty string and an error.
func (al tAguments) Get(aKey string) (string, error) {
	if result, ok := al[aKey]; ok {
		return result, nil
	}

	return "", fmt.Errorf("Missing config value: %s", aKey)
} // Get()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `absolute()` return `aDir` as an absolute path
func absolute(aBaseDir, aDir string) string {
	if 0 == len(aDir) {
		return aDir
	}
	if '/' == aDir[0] {
		s, _ := filepath.Abs(aDir)
		return s
	}

	return filepath.Join(aBaseDir, aDir)
} // absolute()

// `iniData()` returns the config values read from INI file(s).
// The steps here are:
// (1) read the global `/etc/nele.ini`,
// (2) read the user-local `~/.nele.ini`,
// (3) read the `-ini` commandline argument.
func iniData() *ini.TSection {
	// (1) /etc/
	fName := "/etc/nele.ini"
	ini1, err := ini.New(fName)
	if nil == err {
		ini1.AddSectionKey("", "inifile", fName)
	}

	// (2) ~user
	fName = ""
	if usr, err := user.Current(); nil != err {
		fName = os.Getenv("HOME")
	} else {
		fName = usr.HomeDir
	}
	if 0 < len(fName) {
		fName = filepath.Join(fName, ".nele.ini")
		if ini2, err := ini.New(fName); nil == err {
			ini1.Merge(ini2)
			ini1.AddSectionKey("", "inifile", fName)
		}
	}

	// (3) cmdline
	fName = ""
	aLen := len(os.Args)
	for i := 1; i < aLen; i++ {
		if `-ini` == os.Args[i] {
			if i < aLen {
				fName, _ = filepath.Abs(os.Args[i+1])
			}
			break
		}
	}
	if 0 < len(fName) {
		if ini2, _ := ini.New(fName); nil == err {
			ini1.Merge(ini2)
			ini1.AddSectionKey("", "inifile", fName)
		}
	}

	return ini1.GetSection("")
} // iniData()

func init() {
	initArguments()
} // init()

// `initArguments()` reads the commandline arguments into a list
// structure merging it with key-value pairs read from an INI file.
//
// The steps here are:
// (1) read the INI file(s),
// (2) merge the commandline arguments the INI values
// into the global `AppArguments` variable.
func initArguments() {
	AppArguments = make(tAguments)
	defaults := iniData()
	// fmt.Printf("initArguments(defaults): %v", defaults) //FIXME REMOVE

	bnStr, _ := defaults.AsString("blogname")
	flag.StringVar(&bnStr, "blogname", bnStr,
		"Name of this Blog (shown on every page)\n")

	s, _ := defaults.AsString("datadir")
	dataStr, _ := filepath.Abs(s)
	flag.StringVar(&dataStr, "datadir", dataStr,
		"<dirName> the directory with CSS, IMG, JS, POSTINGS, STATIC, VIEWS sub-directories\n")

	s, _ = defaults.AsString("certKey")
	ckStr := absolute(dataStr, s)
	if fi, err := os.Stat(ckStr); (nil != err) || (0 >= fi.Size()) {
		ckStr = ""
	}
	flag.StringVar(&ckStr, "certKey", ckStr,
		"<fileName> the name of the TLS certificate key\n")

	s, _ = defaults.AsString("certPem")
	cpStr := absolute(dataStr, s)
	if fi, err := os.Stat(cpStr); (nil != err) || (0 >= fi.Size()) {
		cpStr = ""
	}
	flag.StringVar(&cpStr, "certPem", cpStr,
		"<fileName> the name of the TLS certificate PEM\n")

	s, _ = defaults.AsString("hashfile")
	hashStr := absolute(dataStr, s)
	if fi, err := os.Stat(hashStr); (nil != err) || (0 >= fi.Size()) {
		hashStr = ""
	}
	flag.StringVar(&hashStr, "hashfile", hashStr,
		"<fileName> (optional) the name of a file storing #hashtags and @mentions\n")

	/*
		s, _ = defaults.AsString("intl")
		intlStr := absolute(dataStr, s)
		flag.StringVar(&intlStr, "intl", intlStr,
			"<fileName> the path/filename of the localisation file\n")
	*/

	iniStr, _ := defaults.AsString("inifile")
	if fi, err := os.Stat(iniStr); (nil != err) || (0 >= fi.Size()) {
		iniStr = ""
	}
	flag.StringVar(&iniStr, "ini", iniStr,
		"<fileName> the path/filename of the INI file\n")

	langStr, _ := defaults.AsString("lang")
	flag.StringVar(&langStr, "lang", langStr,
		"(optional) the default language to use ")

	listenStr, _ := defaults.AsString("listen")
	flag.StringVar(&listenStr, "listen", listenStr,
		"the host's IP to listen at ")

	s, _ = defaults.AsString("logfile")
	logStr := absolute(dataStr, s)
	if fi, err := os.Stat(logStr); (nil != err) || (0 >= fi.Size()) {
		logStr = ""
	}
	flag.StringVar(&logStr, "log", logStr,
		"(optional) name of the logfile to write to\n")

	mfsStr, _ := defaults.AsString("maxfilesize")
	flag.StringVar(&mfsStr, "maxfilesize", mfsStr,
		"max. accepted size of uploaded files")

	/*
		ndBool := false
		flag.BoolVar(&ndBool, "nd", ndBool,
			"(optional) no daemon: whether to not daemonise the program")
	*/

	portInt, _ := defaults.AsInt("port")
	flag.IntVar(&portInt, "port", portInt,
		"<portNumber> the IP port to listen to ")
	portStr := fmt.Sprintf("%d", portInt)

	paBool := false
	flag.BoolVar(&paBool, "pa", paBool,
		"(optional) posting add: write a posting from the commandline")

	pfStr := ""
	flag.StringVar(&pfStr, "pf", pfStr,
		"<fileName> (optional) post file: name of a file to add as new posting")

	realStr, _ := defaults.AsString("realm")
	flag.StringVar(&realStr, "realm", realStr,
		"(optional) <hostName> name of host/domain to secure by BasicAuth\n")

	themStr, _ := defaults.AsString("theme")
	flag.StringVar(&themStr, "theme", themStr,
		"<name> the display theme to use ('light' or 'dark')\n")

	uaStr := ""
	flag.StringVar(&uaStr, "ua", uaStr,
		"<userName> (optional) user add: add a username to the password file")

	ucStr := ""
	flag.StringVar(&ucStr, "uc", ucStr,
		"<userName> (optional) user check: check a username in the password file")

	udStr := ""
	flag.StringVar(&udStr, "ud", udStr,
		"<userName> (optional) user delete: remove a username from the password file")

	s, _ = defaults.AsString("passfile")
	ufStr := absolute(dataStr, s)
	if fi, err := os.Stat(ufStr); (nil != err) || (0 >= fi.Size()) {
		ufStr = ""
	}
	flag.StringVar(&ufStr, "uf", ufStr,
		"<fileName> (optional) user passwords file storing user/passwords for BasicAuth\n")

	ulBool := false
	flag.BoolVar(&ulBool, "ul", ulBool,
		"(optional) user list: show all users in the password file")

	uuStr := ""
	flag.StringVar(&uuStr, "uu", uuStr,
		"<userName> (optional) user update: update a username in the password file")

	flag.Usage = ShowHelp
	flag.Parse() // // // // // // // // // // // // // // // // // // //

	if 0 == len(bnStr) {
		bnStr = time.Now().Format("2006:01:02:15:04:05")
	}
	AppArguments.set("blogname", bnStr)

	if 0 < len(dataStr) {
		dataStr, _ = filepath.Abs(dataStr)
	}
	if f, err := os.Stat(dataStr); nil != err {
		log.Fatalf("datadir == %s` problem: %v", dataStr, err)
	} else if !f.IsDir() {
		log.Fatalf("Error: Not a directory `%s`", dataStr)
	}
	AppArguments.set("datadir", dataStr)

	// `postingBaseDirectory` defined in `posting.go`:
	postingBaseDirectory = filepath.Join(dataStr, "./postings")

	// AppArguments.set("certKey", "")
	if 0 < len(ckStr) {
		ckStr = absolute(dataStr, ckStr)
		if fi, err := os.Stat(ckStr); (nil == err) && (0 < fi.Size()) {
			AppArguments.set("certKey", ckStr)
		}
	}

	if 0 < len(cpStr) {
		cpStr = absolute(dataStr, cpStr)
		if fi, err := os.Stat(cpStr); (nil == err) && (0 < fi.Size()) {
			AppArguments.set("certPem", cpStr)
		}
	}

	if 0 < len(hashStr) {
		hashStr = absolute(dataStr, hashStr)
		AppArguments.set("hashfile", hashStr)
	}

	/*
		if 0 <len(intlStr) {
			intlStr = absolute(dataStr, intlStr)
			AppArguments.set("intl", intlStr)
		}
	*/

	if 0 == len(langStr) {
		langStr = "en"
	}
	AppArguments.set("lang", langStr)

	if "0" == listenStr {
		listenStr = ""
	}
	AppArguments.set("listen", listenStr)

	if 0 < len(logStr) {
		logStr = absolute(dataStr, logStr)
		AppArguments.set("logfile", logStr)
	}

	if 0 == len(mfsStr) {
		mfsStr = "10485760" // 10 MB
	} else {
		mfs := kmg2Num(mfsStr)
		mfsStr = fmt.Sprintf("%d", mfs)
	}
	AppArguments.set("mfs", mfsStr)

	/*
		if ndBool {
			s = fmt.Sprintf("%v", ndBool)
			AppArguments.set("nd", s)
		}
	*/

	portStr = fmt.Sprintf("%d", portInt)
	AppArguments.set("port", portStr)

	if paBool {
		s = fmt.Sprintf("%v", paBool)
		AppArguments.set("pa", s)
	}

	if 0 < len(pfStr) {
		pfStr = absolute(dataStr, pfStr)
		AppArguments.set("pf", pfStr)
	}

	if 0 < len(themStr) {
		AppArguments.set("theme", themStr)
	}

	if 0 < len(uaStr) {
		AppArguments.set("ua", uaStr)
	}

	if 0 < len(ucStr) {
		AppArguments.set("uc", ucStr)
	}

	if 0 < len(udStr) {
		AppArguments.set("ud", udStr)
	}

	if 0 < len(ufStr) {
		ufStr = absolute(dataStr, ufStr)
		AppArguments.set("uf", ufStr)

		// w/o password file there's no BasicAuth
		if 0 < len(realStr) {
			AppArguments.set("real", realStr)
		}
	}

	if ulBool {
		s = fmt.Sprintf("%v", ulBool)
		AppArguments.set("ul", s)
	}

	if 0 < len(uuStr) {
		AppArguments.set("uu", uuStr)
	}
	// fmt.Printf("initArguments(AppArguments): %v", AppArguments) //FIXME REMOVE
} // initArguments()

// `iniWalker()` is an internal helper used to set all INI file
// key-value pairs to the global `AppArguments` list.
func iniWalker(aSect, aKey, aVal string) {
	// Since we're only using the `Default` section we can
	// ignore the `aSect` argument here.
	AppArguments.set(aKey, aVal)
} // iniWalker()

var (
	kmgRE = regexp.MustCompile(`(?i)\s*(\d+)\s*([bgkm]+)?`)
)

// `kmg2Num()` returns a 'KB|MB|GB` string as an integer.
func kmg2Num(aString string) (rInt int) {
	matches := kmgRE.FindStringSubmatch(aString)
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

// ShowHelp lists the commandline options to `Stderr`.
func ShowHelp() {
	fmt.Fprintf(os.Stderr, "\n  Usage: %s [OPTIONS]\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\n  Most options can be set in an INI file to keep the command-line short ;-)")
} // ShowHelp()

/* _EoF_ */
