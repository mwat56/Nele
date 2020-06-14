/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
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
//
//	`aMarkdown` The text to store as a new posting.
func addMarkdown(aMarkdown []byte) (int, error) {
	return NewPosting("").Set(aMarkdown).Store()
} // addMarkdown()

// AddConsolePost reads data from `StdIn` and saves it as a new posting,
// returning the number of bytes written and a possible I/O error.
func AddConsolePost() (int, error) {
	markdown, err := bufio.NewReader(os.Stdin).ReadBytes(0x03)
	if (nil != err) && (io.EOF != err) {
		return 0, err
	}

	return addMarkdown(markdown)
} // AddConsolePost()

// AddFilePost reads `aFilename` and adds it as a new posting,
// returning the number of bytes written and a possible I/O error.
//
//	`aFilename` The text file to add as a new posting.
func AddFilePost(aFilename string) (int, error) {
	markdown, err := ioutil.ReadFile(aFilename) /* #nosec G304 */
	if nil != err {
		return 0, err
	}

	return addMarkdown(markdown)
} // AddFilePost()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// UserAdd reads a password for `aUser` from the commandline
// and adds it to `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` The username to add to the password file.
//	`aFilename` The name of the password file to use.
func UserAdd(aUser, aFilename string) {
	passlist.AddUser(aUser, aFilename)
} // UserAdd()

// UserCheck reads a password for `aUser` from the commandline
// and compares it with the one stored in `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` The username to check in the password file.
//	`aFilename` The name of the password file to use.
func UserCheck(aUser, aFilename string) {
	passlist.CheckUser(aUser, aFilename)
} // UserCheck()

// UserDelete removes the entry for `aUser` from the password
// list `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aUser` The username to remove from the password file.
//	`aFilename` The name of the password file to use.
func UserDelete(aUser, aFilename string) {
	passlist.DeleteUser(aUser, aFilename)
} // UserDelete()

// UserList reads `aFilename` and lists all users stored in there.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
//	`aFilename` The name of the password file to use.
func UserList(aFilename string) {
	passlist.ListUsers(aFilename)
} // UserList()

// UserUpdate reads a password for `aUser` from the commandline
// and updates the entry in the password list `aFilename`.
//
// NOTE: This function does not return but terminates the program
// with error code `0` (zero) if successful, or `1` (one) otherwise.
//
// `aUser` The username to remove from the password file.
//
// `aFilename` The name of the password file to use.
func UserUpdate(aUser, aFilename string) {
	passlist.UpdateUser(aUser, aFilename)
} // UserUpdate()

/* _EoF_ */
