/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides functions to add postings from the commandline
 * and maintain the user/password list.
 */

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"

	"github.com/mwat56/passlist"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `addMarkdown()` saves `aMarkdown` as a new posting,
// returning the number of bytes written and a possible I/O error.
func addMarkdown(aMarkdown []byte) (int64, error) {
	return NewPosting("").Set(aMarkdown).Store()
} // addMarkdown()

// AddConsolePost reads data from `StdIn` and saves it as a new posting,
// returning the number of bytes written and a possible I/O error.
func AddConsolePost() (int64, error) {
	var (
		err      error
		markdown []byte
	)
	if markdown, err = bufio.NewReader(os.Stdin).ReadBytes(0x03); (nil != err) && (io.EOF != err) {
		return 0, err
	}

	return addMarkdown(markdown)
} // AddConsolePost()

// AddFilePost reads `aFilename` and adds it as a new posting,
// returning the number of bytes written and a possible I/O error.
func AddFilePost(aFilename string) (int64, error) {
	var (
		err      error
		markdown []byte
	)
	if markdown, err = ioutil.ReadFile(aFilename); /* #nosec G304 */ nil != err {
		return 0, err
	}

	return addMarkdown(markdown)
} // AddFilePost()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// AddUser reads a password for `aUser` from the commandline
// and adds it to `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` the username to add to the password file.
//
//	`aFilename` name of the password file to use.
func AddUser(aUser, aFilename string) {
	passlist.AddUser(aUser, aFilename)
} // AddUser()

// CheckUser reads a password for `aUser` from the commandline
// and compares it with the one stored in `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` the username to check in the password file.
//
//	`aFilename` name of the password file to use.
func CheckUser(aUser, aFilename string) {
	passlist.CheckUser(aUser, aFilename)
} // CheckUser()

// DeleteUser removes the entry for `aUser` from the password
// list `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` the username to remove from the password file.
//
//	`aFilename` name of the password file to use.
func DeleteUser(aUser, aFilename string) {
	passlist.DeleteUser(aUser, aFilename)
} // DeleteUser()

// ListUsers reads `aFilename` and lists all users stored in there.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
// `aFilename` name of the password file to use.
func ListUsers(aFilename string) {
	passlist.ListUsers(aFilename)
} // ListUsers()

// UpdateUser reads a password for `aUser` from the commandline
// and updates the entry in the password list `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
// `aUser` the username to remove from the password file.
//
// `aFilename` name of the password file to use.
func UpdateUser(aUser, aFilename string) {
	passlist.UpdateUser(aUser, aFilename)
} // UpdateUser()

/* _EoF_ */
