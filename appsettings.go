/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"fmt"
	"strings"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `TAppArgs` is a collection of commandline arguments and INI values.
	TAppArgs struct {
		AccessLog     string // (optional) name of page access logfile
		Addr          string // listen address ("1.2.3.4:5678")
		BlogName      string // name/description of this blog
		CertKey       string // TLS certificate key
		CertPem       string // private TLS certificate
		DataDir       string // base directory of application's data
		delWhitespace bool   // remove whitespace from generated pages
		Dump          bool   // Debug: dump this structure to `StdOut`
		ErrorLog      string // (optional) name of page error logfile
		GZip          bool   // send compressed data to remote browser
		HashFile      string // file of hashtag/mention database
		// Intl       string // path/filename of the localisation file
		Lang     string // default GUI language
		listen   string // IP of host to listen at
		LogStack bool   // log stack trace in case of errors

		MaxFileSize int64  // max. upload file size
		mfs         string // max. upload file size

		Name string // name of the actual program

		PageLength  uint   // the number of postings to show per page
		persistence string // either `db`, `fs`, or `tee`.`
		PostAdd     bool   // whether to write a posting from commandline
		PostFile    string // name of file to post
		port        int    // port to listen to
		Realm       string // host/domain to secure by BasicAuth
		Screenshot  bool   // whether to use page screenshots or not
		Theme       string // `dark` or `light` display theme
		UserAdd     string // username to add to password list
		UserCheck   string // username to check in password list
		UserDelete  string // username to delete from password list
		UserFile    string // (optional) name of page access logfile
		UserList    bool   // print out a list of current users
		UserUpdate  string // username to update in password list
	}
)

var (
	// `AppArgs` holds the commandline arguments and INI values combined.
	//
	// This structure should be considered `R/O` after it was set up
	// by a call to `InitConfig()`.
	AppArgs TAppArgs
)

// `String()` implements the `Stringer` interface returning a (pretty
// printed) string representation of the current `TAppArgs` instance.
//
// NOTE: This method is meant mostly for debugging purposes.
//
// Returns:
//   - `string`: The string representation of the current app configuration.
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

/* _EoF_ */
