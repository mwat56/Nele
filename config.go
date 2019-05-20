/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package blog

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ini "github.com/mwat56/go-ini"
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

// ShowHelp lists the commandline options to `Stderr`.
func ShowHelp() {
	fmt.Fprintf(os.Stderr, "\n  Usage: %s [OPTIONS]\n\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\nMost options can be set in an INI file to keep he commandline short ;-)\n")
} // ShowHelp()

// `iniWalker()` is an internal helper used to set all INI file
// key-value pairs to the global `AppArguments` list.
func iniWalker(aSect, aKey, aVal string) {
	// Since we're only using the `Default` section we can
	// ignore the `aSect` argument here.
	AppArguments.set(aKey, aVal)
} // iniWalker()

// `initArguments()` reads the commandline arguments into a list
// structure merging it with key-value pairs read from an INI file.
//
// The steps here are:
// (1) a hard-coded INI filename ("./blog.ini") is used to (try to)
// read an INI file into a local `data` variable.
// (2) The commandline arguments are read/parsed into the global
// `AppArguments` variable.
// (3) If the commandline arguments named another INI filename that
// INI file is read as well and merged with the old data values.
// (4) Finally the commandline arguments are merged with the INI values
// into the global `AppArguments` variable.
func initArguments() {
	AppArguments = make(tAguments)

	defIniFile, _ := filepath.Abs("./blog.ini")
	data, err := ini.LoadFile(defIniFile)
	if nil == err {
		data.AddSectionKey("", "inifile", defIniFile)
		data.Walk(iniWalker)
	} else {
		data = ini.NewSections()
		data.AddSectionKey("", "certKey", "")
		data.AddSectionKey("", "certPem", "")
		s, _ := filepath.Abs("./")
		data.AddSectionKey("", "datadir", s)
		s, _ = filepath.Abs("./hashfile.db")
		data.AddSectionKey("", "hashfile", s)
		data.AddSectionKey("", "inifile", defIniFile)
		// s, _ = filepath.Abs("./intl.ini")
		// data.AddSectionKey("", "intl", s)
		s, _ = filepath.Abs("./js/")
		data.AddSectionKey("", "js", s)
		data.AddSectionKey("", "lang", "de")
		data.AddSectionKey("", "listen", "127.0.0.1")
		data.AddSectionKey("", "logfile", "")
		data.AddSectionKey("", "maxfilesize", "10MB")
		data.AddSectionKey("", "passfile", "")
		data.AddSectionKey("", "port", "8181")
		data.AddSectionKey("", "realm", "")
		data.AddSectionKey("", "theme", "light")
	}
	defaults := data.GetSection("")

	s, _ := defaults.AsString("certKey")
	ckStr, _ := filepath.Abs(s)
	flag.StringVar(&ckStr, "certKey", ckStr,
		"<fileName> the name of the TLS certificate key\n")

	s, _ = defaults.AsString("certPem")
	cpStr, _ := filepath.Abs(s)
	flag.StringVar(&cpStr, "certPem", cpStr,
		"<fileName> the name of the TLS certificate PEM\n")

	s, _ = defaults.AsString("datadir")
	dataStr, _ := filepath.Abs(s)
	flag.StringVar(&dataStr, "datadir", dataStr,
		"<dirName> the directory with CSS, IMG, JS, POSTINGS, STATIC, VIEWS sub-directories\n")

	s, _ = defaults.AsString("hashfile")
	hashStr, _ := filepath.Abs(s)
	flag.StringVar(&hashStr, "hashfile", hashStr,
		"<fileName> (optional) the name of a file storing #hashtags and @mentions\n")

	/*
		s, _ = defaults.AsString("intl")
		intlStr, _ := filepath.Abs(s)
		flag.StringVar(&intlStr, "intl", intlStr,
			"<fileName> the path/filename of the localisation file\n")
	*/

	iniStr := ""
	flag.StringVar(&iniStr, "ini", iniStr,
		"<fileName> the path/filename of the INI file\n")

	langStr, _ := defaults.AsString("lang")
	flag.StringVar(&langStr, "lang", langStr,
		"(optional) the default language to use ")

	listenStr, _ := defaults.AsString("listen")
	flag.StringVar(&listenStr, "listen", listenStr,
		"the host's IP to listen at ")

	logStr, _ := defaults.AsString("logfile")
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
	ufStr, _ := filepath.Abs(s)
	flag.StringVar(&ufStr, "uf", ufStr,
		"<fileName> (optional) user passwords file storing user/passwords for BasicAuth\n")

	ulBool := false
	flag.BoolVar(&ulBool, "ul", ulBool,
		"(optional) user list: show all users in the password file")

	uuStr := ""
	flag.StringVar(&uuStr, "uu", uuStr,
		"<userName> (optional) user update: update a username in the password file")

	flag.Usage = ShowHelp
	flag.Parse()

	cmdIniFile, _ := filepath.Abs(iniStr)
	if cmdIniFile != defIniFile {
		data := ini.NewSections()
		if data, err := data.Load(defIniFile); nil == err {
			data.Walk(iniWalker)
		}
	}
	AppArguments.set("inifile", cmdIniFile)

	if 0 < len(ckStr) {
		ckStr, _ = filepath.Abs(ckStr)
		AppArguments.set("certKey", ckStr)
	}

	if 0 < len(cpStr) {
		cpStr, _ = filepath.Abs(cpStr)
		AppArguments.set("certPem", cpStr)
	}

	if 0 < len(dataStr) {
		dataStr, _ = filepath.Abs(dataStr)
	}
	if f, err := os.Stat(dataStr); nil != err {
		log.Fatalf("`%s` problem: %v", dataStr, err)
	} else if !f.IsDir() {
		log.Fatalf("Error: Not a directory `%s`", dataStr)
	}
	AppArguments.set("datadir", dataStr)
	// defined in `posting.go`:
	postingBaseDirectory = filepath.Join(dataStr, "./postings")

	if 0 < len(hashStr) {
		hashStr, _ = filepath.Abs(hashStr)
		AppArguments.set("hashfile", hashStr)
	}

	/*
		if 0 == len(intlStr) {
			intlStr, _ = filepath.Abs(intlStr)
		}
			AppArguments.set("intl", intlStr)
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
		logStr, _ = filepath.Abs(logStr)
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
		pfStr, _ = filepath.Abs(pfStr)
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
		ufStr, _ = filepath.Abs(ufStr)
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
} // initArguments()

func init() {
	initArguments()
} // init()

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

/* _EoF_ */
