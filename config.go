/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package blog

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

// add() sets another key-value pair.
func (al tAguments) add(aKey string, aValue string) {
	al[aKey] = aValue
} // add()

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
	fmt.Fprint(os.Stderr, "\nMost options can be set in an INI file to keep he commandline short ;-)\n\nWith all file- and directory-names make sure that they're readable, and at\nleast the 'post' folder must be writeable for the user running this\nprogram to store the postings.\n\n")
} // ShowHelp()

// iniWalker() is an internal helper used to add all INI file
// key-value pairs to the global `AppArguments` list.
func iniWalker(aSect, aKey, aVal string) {
	// Since we're only using the `Default` section we can
	// ignore the `aSect` argument here.
	AppArguments.add(aKey, aVal)
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
		s, _ := filepath.Abs("./css/")
		data.AddSectionKey("", "css", s)
		s, _ = filepath.Abs("./img/")
		data.AddSectionKey("", "img", s)
		data.AddSectionKey("", "inifile", defIniFile)
		// s, _ = filepath.Abs("./intl.ini")
		// data.AddSectionKey("", "intl", s)
		s, _ = filepath.Abs("./js/")
		data.AddSectionKey("", "js", s)
		data.AddSectionKey("", "lang", "de")
		data.AddSectionKey("", "listen", "127.0.0.1")
		data.AddSectionKey("", "logfile", "")
		data.AddSectionKey("", "passfile", "")
		data.AddSectionKey("", "port", "8181")
		s, _ = filepath.Abs("./postings/")
		data.AddSectionKey("", "postdir", s)
		data.AddSectionKey("", "realm", "")
		s, _ = filepath.Abs("./static/")
		data.AddSectionKey("", "static", s)
		s, _ = filepath.Abs("./views/")
		data.AddSectionKey("", "tpldir", s)
	}
	defaults := data.GetSection("")

	s, _ := defaults.AsString("css")
	cssStr, _ := filepath.Abs(s)
	flag.StringVar(&cssStr, "css", cssStr,
		"<dirName> the directory with CSS file(s)\n")

	s, _ = defaults.AsString("img")
	imgStr, _ := filepath.Abs(s)
	flag.StringVar(&imgStr, "img", imgStr,
		"<dirName> the directory with images\n")

	/*
		s, _ = defaults.AsString("intl")
		intlStr, _ := filepath.Abs(s)
		flag.StringVar(&intlStr, "intl", intlStr,
			"<fileName> the path/filename of the localisation file\n")
	*/

	iniStr, _ := defaults.AsString("inifile")
	flag.StringVar(&iniStr, "ini", iniStr,
		"<fileName> the path/filename of the INI file\n")

	s, _ = defaults.AsString("js")
	jsStr, _ := filepath.Abs(s)
	flag.StringVar(&jsStr, "js", jsStr,
		"<dirName> the directory with JavaScript\n")

	langStr, _ := defaults.AsString("lang")
	flag.StringVar(&langStr, "lang", langStr,
		"(optional) the default language to use\n")

	listenStr, _ := defaults.AsString("listen")
	flag.StringVar(&listenStr, "listen", listenStr,
		"the host's IP to listen at\n")

	logStr, _ := defaults.AsString("logfile")
	flag.StringVar(&logStr, "log", logStr,
		"(optional) name of the logfile to write to\n")

	/*
		ndBool := false
		flag.BoolVar(&ndBool, "nd", ndBool,
			"(optional) no daemon: whether daemonise the program")
	*/

	portInt, _ := defaults.AsInt("port")
	flag.IntVar(&portInt, "port", portInt,
		"<portNumber> the IP port to listen to")
	portStr := fmt.Sprintf("%d", portInt)

	s, _ = defaults.AsString("postdir")
	postStr, _ := filepath.Abs(s)
	flag.StringVar(&postStr, "post", postStr,
		"<dirName> the directory used for storing the postings\n")

	paBool := false
	flag.BoolVar(&paBool, "pa", paBool,
		"(optional) posting add: whether to write a posting from the commandline")

	pfStr := ""
	flag.StringVar(&pfStr, "pf", pfStr,
		"<fileName> (optional) posting file: name of a file to add as new posting")

	realStr, _ := defaults.AsString("realm")
	flag.StringVar(&realStr, "realm", realStr,
		"(optional) <hostName> name of host/domain to secure by BasicAuth\n")

	s, _ = defaults.AsString("static")
	stcStr, _ := filepath.Abs(s)
	flag.StringVar(&stcStr, "static", stcStr,
		"<dirName> the directory with static files\n")

	s, _ = defaults.AsString("tpldir")
	tplStr, _ := filepath.Abs(s)
	flag.StringVar(&tplStr, "tpl", tplStr,
		"<dirName> directory with page templates\n")

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
		"(optional) user list: whether to show all users in the password file")

	uuStr := ""
	flag.StringVar(&uuStr, "uu", uuStr,
		"<userName> (optional) user update: update a username from the password file")

	flag.Usage = ShowHelp
	flag.Parse()

	cmdIniFile, _ := filepath.Abs(iniStr)
	if cmdIniFile != defIniFile {
		data := ini.NewSections()
		if data, err := data.Load(defIniFile); nil == err {
			data.Walk(iniWalker)
		}
	}
	AppArguments.add("inifile", cmdIniFile)

	if 0 < len(cssStr) {
		cssStr, _ = filepath.Abs(cssStr)
	}
	AppArguments.add("css", cssStr)

	if 0 < len(imgStr) {
		imgStr, _ = filepath.Abs(imgStr)
	}
	AppArguments.add("img", imgStr)

	/*
		if 0 == len(intlStr) {
			intlStr, _ = filepath.Abs(intlStr)
		}
			AppArguments.add("intl", intlStr)
	*/

	if 0 < len(jsStr) {
		jsStr, _ = filepath.Abs(jsStr)
	}
	AppArguments.add("js", jsStr)

	if 0 == len(langStr) {
		langStr = "en"
	}
	AppArguments.add("lang", langStr)

	if "0" == listenStr {
		listenStr = ""
	}
	AppArguments.add("listen", listenStr)

	if 0 < len(logStr) {
		logStr, _ = filepath.Abs(logStr)
		AppArguments.add("logfile", logStr)
	}

	/*
		if ndBool {
			s = fmt.Sprintf("%v", ndBool)
			AppArguments.add("nd", s)
		}
	*/

	portStr = fmt.Sprintf("%d", portInt)
	AppArguments.add("port", portStr)

	if paBool {
		s = fmt.Sprintf("%v", paBool)
		AppArguments.add("pa", s)
	}

	if 0 < len(pfStr) {
		pfStr, _ = filepath.Abs(pfStr)
		AppArguments.add("pf", pfStr)
	}

	if 0 < len(postStr) {
		postStr, _ = filepath.Abs(postStr)
	}
	AppArguments.add("postdir", postStr)

	if 0 < len(stcStr) {
		stcStr, _ = filepath.Abs(stcStr)
	}
	AppArguments.add("static", stcStr)

	if 0 < len(tplStr) {
		tplStr, _ = filepath.Abs(tplStr)
	}
	AppArguments.add("tpldir", tplStr)

	if 0 < len(uaStr) {
		AppArguments.add("ua", uaStr)
	}

	if 0 < len(ucStr) {
		AppArguments.add("uc", ucStr)
	}

	if 0 < len(udStr) {
		AppArguments.add("ud", udStr)
	}

	if 0 < len(ufStr) {
		ufStr, _ = filepath.Abs(ufStr)
		AppArguments.add("uf", ufStr)
		// w/o password file there's no BasicAuth
		if 0 < len(realStr) {
			AppArguments.add("real", realStr)
		}
	}

	if ulBool {
		s = fmt.Sprintf("%v", ulBool)
		AppArguments.add("ul", s)
	}

	if 0 < len(uuStr) {
		AppArguments.add("uu", uuStr)
	}

} // initArguments()

func init() {
	initArguments()
} // init()

/* _EoF_ */
