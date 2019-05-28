/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

/*
 * This file provides article/posting related functions and methods.
 */

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newID()` returns an article ID based on `aTime` in hexadecimal notation.
//
// `aTime` is the time to be returned in hexadecimal notation.
//
// Internal function to allow for unit testing.
// The `timeID()` function reverses this computation.
func newID(aTime time.Time) string {
	id := fmt.Sprintf("%x", aTime.UnixNano())
	if 16 > len(id) {
		return strings.Repeat("0", 16-len(id)) + id
	}

	return id
} // newID()

// NewID returns a new article ID.
// It is based on the current date/time and given in hexadecimal notation.
// It's assumend that no more than one ID per nanosecond is needed.
func NewID() string {
	return newID(time.Now())
} // NewID()

// `timeID()` returns a posting's date/time represented by `aID`.
//
// `aID` is a posting's ID as returned by `newID()`.
func timeID(aID string) (rTime time.Time) {
	if i64, err := strconv.ParseInt(aID, 16, 64); nil == err {
		return time.Unix(0, i64)
	}

	return
} // timeID()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// `postingBaseDirectory` is the base directory for storing the articles.
	//
	// This variable's value _must_ be set initially before creating any
	// `TPosting` or `TPostList` instances.
	// After that it should be considered `read/only`.
	postingBaseDirectory = "./postings"
)

// PostingBaseDirectory returns the base directory used for
// storing articles/postings.
func PostingBaseDirectory() string {
	return postingBaseDirectory
} // PostingBaseDirectory()

// SetPostingBaseDirectory set the base directory used for
// storing articles/postings.
//
// `aBaseDir` is the base directory to use for storing articles/postings.
func SetPostingBaseDirectory(aBaseDir string) (rErr error) {
	postingBaseDirectory, rErr = filepath.Abs(aBaseDir)
	return
} // SetPostingBaseDirectory()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TPosting is a single /article/posting to be injected into a template.
	TPosting struct {
		id       string // hex. representation of date/time
		markdown []byte // (article-/file-)contents in Markdown markup
	}
)

// NewPosting returns a new posting structure with an empty article text.
func NewPosting() *TPosting {
	return newPosting("")
} // NewPosting()

// newPosting() is the core function of `NewPost()` (for testing purposes).
//
// `aID` if an empty string the `NewID()` function is called
// to provide a new article ID.
func newPosting(aID string) *TPosting {
	if 0 == len(aID) {
		aID = NewID()
	}
	result := TPosting{id: aID}

	return &result
} // newPosting()

// After reports whether this posting is younger than the one
// identified by `aID`.
//
// `aID` is the ID of another posting to compare.
func (p *TPosting) After(aID string) bool {
	return (p.id > aID)
} // After()

// Before reports whether this posting is older than the one
// identified by `aID`.
//
// `aID` is the ID of another posting to compare.
func (p *TPosting) Before(aID string) bool {
	return (p.id < aID)
} // Before()

// Clear resets the internal fields to their respective zero values.
//
// This method does NOT remove the file (if any) associated with this
// posting/article; for that call the `Delete()` method.
func (p *TPosting) Clear() *TPosting {
	var bs []byte

	p.markdown = bs

	return p
} // Clear()

// clone() returns a copy of this posting/article.
func (p *TPosting) clone() *TPosting {
	return &TPosting{
		id:       p.id,
		markdown: p.markdown,
	}
} // clone()

// Date returns the posting's date as a formatted string (`yyy-mm-dd`).
func (p *TPosting) Date() string {
	y, m, d := timeID(p.id).Date()

	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
} // Date()

// delFile() removes `aFileName` from the filesystem
// returning a possible I/O error.
//
// A non-existing file is not considered an error here.
//
// `aFileName` is the name of the file to delFile.
func delFile(aFileName *string) error {
	err := os.Remove(*aFileName)
	if nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
	}

	return err
} // delFile()

// Delete removes the posting/article from the filesystem
// returning a possible I/O error.
//
// This method does NOT empty the internal fields of the object;
// for that call the `Clear()` method.
func (p *TPosting) Delete() error {
	filepathname := p.PathFileName()

	return delFile(&filepathname)
} // Delete()

// Equal reports whether this posting is of the same time as `aID`.
func (p *TPosting) Equal(aID string) bool {
	return timeID(p.id).Equal(timeID(aID))
} // Equal()

// Exists returns whether there is a file with more than zero bytes.
func (p *TPosting) Exists() bool {
	fi, err := os.Stat(pathname(p.id))
	if nil != err {
		return false
	}
	if (0 == fi.Size()) || fi.IsDir() {
		return false
	}

	return true
} // Exists()

// ID returns the article's identifier.
//
// The identifier is based on the article's creation time
// and given in hexadecimal notation.
//
// This methods allows the template to validate and use
// the placeholder `.ID`
func (p *TPosting) ID() string {
	return p.id
} // ID()

// Len returns the current length of the posting's Markdown text.
func (p *TPosting) Len() int {
	return len(p.markdown)
} // Len()

// Load reads the Markdown from disk, returning a possible I/O error.
func (p *TPosting) Load() error {
	var err error
	filepathname := p.PathFileName()
	if _, err = os.Stat(filepathname); nil != err {
		return err // probably ENOENT
	}

	var bs []byte
	if bs, err = ioutil.ReadFile(filepathname); nil != err {
		return err
	}
	p.markdown = []byte(strings.TrimSpace(string(bs)))

	return nil
} // Load()

// makeDir() creates the directory for storing the article
// returning the article's path/file-name but w/o filename extension.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
func (p *TPosting) makeDir() (string, error) {
	fmode := os.ModeDir | 0775

	// We need the year to guard against ID overflows.
	y, _, _ := timeID(p.id).Date()
	// Using the aID's first three characters leads to
	// directories worth about 52 days of data.
	dir := fmt.Sprintf("%04d%s", y, p.id[:3])
	dirname := path.Join(postingBaseDirectory, dir)
	if err := os.MkdirAll(filepath.FromSlash(dirname), fmode); nil != err {
		return "", err
	}

	return path.Join(dirname, p.id), nil
} // makeDir()

// Markdown returns the Markdown of this article.
//
// If the markup is not already in memory this methods tries
// to read the required data from the filesystem.
//
// In case of I/O errors while accessing the file an empty
// byte slice is returned.
func (p *TPosting) Markdown() []byte {
	if 0 < len(p.markdown) {
		// that's the easy path ...
		return p.markdown
	}
	var err error

	// now we have to check the filesystem
	fName := p.PathFileName()
	if _, err = os.Stat(fName); nil != err {
		return p.markdown // return empty slice
	}

	var bs []byte
	if bs, err = ioutil.ReadFile(fName); nil != err {
		return p.markdown // return empty slice
	}
	p.markdown = []byte(strings.TrimSpace(string(bs)))

	return p.markdown
} // Markdown()

// `pathname()` returns the complete article path-/filename.
func pathname(aID string) string {
	// We need the year to guard against ID overflows.
	y, _, _ := timeID(aID).Date()

	// Using the aID's first three characters leads to
	// directories worth about 52 days of data.
	dir := fmt.Sprintf("%04d%s", y, aID[:3])

	return path.Join(postingBaseDirectory, dir, aID+".md")
} // pathname()

// PathFileName returns the article's complete path-/filename.
func (p *TPosting) PathFileName() string {
	return pathname(p.id)
} // PathFileName()

// Post returns the article's HTML markup.
func (p *TPosting) Post() template.HTML {
	return cachedHTML(p)
	/*
		// make sure we have the most recent version:
		p.Markdown()

		return template.HTML(MDtoHTML(p.markdown))
	*/
} // Post()

// Set assigns the article's Markdown text.
//
// `aMarkdown` is the actual Markdown text of the article to assign.
func (p *TPosting) Set(aMarkdown []byte) *TPosting {
	if 0 < len(aMarkdown) {
		p.markdown = []byte(strings.TrimSpace(string(aMarkdown)))
	} else {
		p.Load()
	}

	return p
} // Set()

// Store writes the article's Markdown to disk returning
// the number of bytes written and a possible I/O error.
//
// The file is created on disk with mode `0644` (`-rw-r--r--`).
func (p *TPosting) Store() (int64, error) {
	var result int64
	var err error

	if _, err = p.makeDir(); nil != err {
		// without an appropriate directory we can't save anything ...
		return result, err
	}
	if 0 == len(p.markdown) {
		p.Load()
		if 0 == len(p.markdown) {
			return result, fmt.Errorf("Markdown '%s' is empty", p.id)
		}
	}
	fName := p.PathFileName()
	if err = ioutil.WriteFile(fName, p.markdown, 0644); nil != err {
		return result, err
	}

	if fi, err := os.Stat(fName); nil == err {
		result = fi.Size()
	}

	return result, err
} // Store()

// Time returns the posting's date/time.
func (p *TPosting) Time() time.Time {
	return timeID(p.id)
} // Time()

/*
func updatePostDirs() {
	dirnames, err := filepath.Glob(postingBaseDirectory + "/*")
	if nil != err {
		return
	}
	for _, dirname := range dirnames {
		filesnames, err := filepath.Glob(dirname + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filesnames) {
			continue // skip empty directory
		}
		for _, fName := range filesnames {
			pName := strings.TrimPrefix(fName, dirname+"/")
			id := pName[:len(pName)-3]
			newName := pathname(id)

			y, _, _ := timeID(id).Date()
			dir := fmt.Sprintf("%04d%s", y, id[:3])
			dirname := path.Join(postingBaseDirectory, dir)
			os.MkdirAll(filepath.FromSlash(dirname), os.ModeDir|0775)
			os.Rename(fName, newName)
		}
	}
} // updatePostDirs()
*/

/* _EoF_ */
