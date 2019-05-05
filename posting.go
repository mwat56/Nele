/*
    Copyright © 2019  M.Watermann, 10247 Berlin, Germany
                All rights reserved
            EMail : <support@mwat.de>
*/

package blog

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
	"syscall"
	"time"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newID()` returns an article ID based on `aTime` in hexadecimal notation.
//
// `aTime` is the time to be returned in hexadecimal notation.
//
// Internal function to allow for unit testing.
// The `timeID()` function reverses this computation.
func newID(aTime time.Time) string {
	return fmt.Sprintf("%x", aTime.UnixNano())
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
		rTime = time.Unix(0, i64)
	}

	return
} // timeID()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TPosting is a single /article/posting to be injected into a template.
	TPosting struct {
		basedir  string // the base directory for storing the articles
		id       string // hex. representation of date/time
		markdown []byte // (article-/file-)contents in Markdown markup
	}

	// TPost is a simple structure for data access in templates
	TPost struct {
		Date string        // article date
		ID   string        // article ID: hex. representation of date/time
		Post template.HTML // article contents in HTML markup
	}
)

// NewPosting returns a new posting structure with an empty article text.
//
// `aBaseDir` is the base direcetory under which all postings are stored.
func NewPosting(aBaseDir string) *TPosting {
	return newPosting(aBaseDir, "")
} // NewPosting()

// newPosting() is the core function of `NewPost()` (for testing purposes).
//
// `aID` if an empty string the `NewID()` function is called
// to provide a new article ID.
func newPosting(aBaseDir, aID string) *TPosting {
	if 0 == len(aBaseDir) {
		aBaseDir = "./"
	}
	if 0 == len(aID) {
		aID = NewID()
	}
	bd, _ := filepath.Abs(aBaseDir)
	result := TPosting{basedir: bd, id: aID}

	return &result
} // newPosting()

// After reports whether this posting is younger than the one
// identified by `aID`.
//
// `aID` is the (timestamp base) ID of another posting to compare.
func (p *TPosting) After(aID string) bool {
	return timeID(p.id).After(timeID(aID))
} // After()

// Before reports whether this posting is older than the one
// identified by `aID`.
//
// `aID` is the (timestamp base) ID of another posting to compare.
func (p *TPosting) Before(aID string) bool {
	return timeID(p.id).Before(timeID(aID))
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
		basedir:  p.basedir,
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

// IsFile returns whether the posting is stored in the filesystem.
func (p *TPosting) IsFile() bool {
	f, err := os.Stat(p.PathFileName())
	if nil != err {
		return false
	}

	return (0 < f.Size())
} // IsFile()

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
	p.markdown = bs

	return nil
} // Load()

// makeDir() creates the directory for storing the article
// returning the article's path/file-name but w/o filename extension.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
func (p *TPosting) makeDir() (string, error) {
	fmode := os.ModeDir | 0775

	// Using the aID's first three characters leads to
	// directories worth about 52 days of data.
	dirname := path.Join(p.basedir, string(p.id[:3]))
	if err := os.MkdirAll(filepath.FromSlash(dirname), fmode); nil != err {
		return "", err
	}

	return path.Join(dirname, p.id), nil
} // makeDir()

// Markdown returns the Markdown of this article.
//
// If the markup is not already in memory this methods tries
// to read the required data from the filesystem.
func (p *TPosting) Markdown() []byte {
	if 0 < len(p.markdown) {
		// that's the easy path ...
		return p.markdown
	}

	// now we have to check the filesystem
	filepathname := p.PathFileName()
	if _, err := os.Stat(filepathname); nil != err {
		return p.markdown // return empty slice
	}
	p.markdown, _ = ioutil.ReadFile(filepathname)

	return p.markdown
} // Markdown()

func pathname(aBaseDir, aID string) string {
	// Using the aID's first three characters leads to
	// directories worth about 52 days of data.
	return path.Join(aBaseDir, string(aID[:3]), aID+".md")
} // pathname()

// PathFileName returns the article's path-/filename.
func (p *TPosting) PathFileName() string {
	return pathname(p.basedir, p.id)
} // PathFileName()

// Post returns the artile's HTML markup.
func (p *TPosting) Post() template.HTML {
	// make sure we have the most recent version:
	p.Markdown()

	return template.HTML(MDtoHTML(p.markdown))
} // Post()

// Posting returns this article suitable for use in templates.
//
// This method allows the template to validate
// and use the placeholders `Date`, `ID` and `Post`.
func (p *TPosting) Posting() *TPost {
	// make sure we have the most recent version:
	p.Markdown()

	return &TPost{
		Date: p.Date(),
		ID:   p.id,
		Post: p.Post(),
	}
} // Posting()

// Set assigns the article's Markdown text.
//
// `aMarkdown` is the actual Markdown text of the article to assign.
func (p *TPosting) Set(aMarkdown []byte) *TPosting {
	if 0 < len(aMarkdown) {
		p.markdown = aMarkdown
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
	filepathname := p.PathFileName()
	if err = ioutil.WriteFile(filepathname, p.markdown, 0644); nil != err {
		return result, err
	}

	if fi, err := os.Stat(filepathname); nil == err {
		result = fi.Size()
	}

	return result, err
} // Store()

// Time returns the posting's date/time.
func (p *TPosting) Time() time.Time {
	return timeID(p.id)
} // Time()

/* _EoF_ */
