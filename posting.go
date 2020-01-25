/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

/*
 * This file provides article/posting related functions and methods.
 */

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/mwat56/apachelogger"
	// bf "github.com/russross/blackfriday/v2"
	bf "gopkg.in/russross/blackfriday.v2"
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// RegEx to correct redundant markup created by 'bf';
	// see `mdToHTML()`
	poPreCodeRE1 = regexp.MustCompile(`(?s)\s*(<pre>)<code>(.*?)\s*</code>(</pre>)\s*`)

	poPreCodeRE2 = regexp.MustCompile(`(?s)\s*(<pre)><code (class="language-\w+")>(.*?)\s*</code>(</pre>)\s*`)

	// Instead of creating this objects with every call to `mdToHTML()`
	// we use some prepared global instances.
	bfExtensions = bf.WithExtensions(
		bf.Autolink |
			bf.BackslashLineBreak |
			bf.DefinitionLists |
			bf.FencedCode |
			bf.Footnotes |
			bf.HeadingIDs |
			bf.NoIntraEmphasis |
			bf.SpaceHeadings |
			bf.Strikethrough |
			bf.Tables)

	bfRenderer = bf.WithRenderer(
		bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.FootnoteReturnLinks |
				bf.Smartypants |
				bf.SmartypantsFractions |
				bf.SmartypantsDashes |
				bf.SmartypantsLatexDashes,
		}),
	)

	bfPre = []byte("</pre>")

	bfPreCode = []byte("<pre><code ")
)

// `mdToHTML()` converts the `aMarkdown` data returning HTML data.
//
//	`aMarkdown` The raw Markdown text to convert.
func mdToHTML(aMarkdown []byte) (rHTML []byte) {
	rHTML = bytes.TrimSpace(bf.Run(aMarkdown, bfRenderer, bfExtensions))

	// Testing for PRE first makes this implementation twice as fast
	// if there's no PRE in the generated HTML and about the same
	// speed if there actually is a PRE part.
	if i := bytes.Index(rHTML, bfPre); 0 > i {
		return // no need for further RegEx execution
	}

	// return handlePreCode(rHTML)
	rHTML = poPreCodeRE1.ReplaceAll(rHTML, []byte("$1\n$2\n$3"))

	if i := bytes.Index(rHTML, bfPreCode); 0 > i {
		return // no need for the second RegEx execution
	}
	return poPreCodeRE2.ReplaceAll(rHTML, []byte("$1 $2>\n$3\n$4"))
} // mdToHTML()

// `newID()` returns an article ID based on `aTime` in hexadecimal notation.
//
//	`aTime` is the time to be returned in hexadecimal notation.
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
// It's assumed that no more than one ID per nanosecond is required.
func NewID() string {
	return newID(time.Now())
} // NewID()

// `timeID()` returns a posting's date/time represented by `aID`.
//
//	`aID` is a posting's ID as returned by `newID()`.
func timeID(aID string) (rTime time.Time) {
	if i64, err := strconv.ParseInt(aID, 16, 64); nil == err {
		return time.Unix(0, i64)
	}

	return
} // timeID()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// `poPostingBaseDirectory` is the base directory for storing articles.
	//
	// This variable's value must be set initially before creating any
	// `TPosting` or `TPostList` instances.
	// After that it should be considered `read/only`.
	// Its default value is `./postings`.
	//
	//	see `PostingBaseDirectory()`, `SetPostingBaseDirectory()`
	poPostingBaseDirectory = func() string {
		dir, _ := filepath.Abs(`./postings`)
		return dir
	}()

	// Cache of last/current posting count.
	// see `PostingCount()`, `delFile()`, `goCount()`, `TPosting.Store()`
	µCountCache uint32
)

// PostingBaseDirectory returns the base directory used for
// storing articles/postings.
func PostingBaseDirectory() string {
	return poPostingBaseDirectory
} // PostingBaseDirectory()

// PostingCount returns the number of postings currently available.
//
// In case of I/O errors the return value will be `-1`.
func PostingCount() (rCount int) {
	if result := atomic.LoadUint32(&µCountCache); 0 < result {
		return int(result)
	}
	// Apparently there's no current value ready so we compute a new one.
	// Instead of doing this in _one_ glob we actually do it in two
	// thus trading memory consumption with processing time and so
	// we're being better prepared for huge amounts of postings.
	dirnames, err := filepath.Glob(poPostingBaseDirectory + `/*`)
	if nil != err {
		return // we can't recover from this :-(
	}
	for _, mdName := range dirnames {
		if filesnames, err := filepath.Glob(mdName + `/*.md`); nil == err {
			rCount += len(filesnames)
		}
	}
	atomic.StoreUint32(&µCountCache, uint32(rCount))

	return
} // PostingCount()

// SetPostingBaseDirectory sets the base directory used for
// storing articles/postings.
//
//	`aBaseDir` The base directory to use for storing articles/postings.
func SetPostingBaseDirectory(aBaseDir string) {
	dir, err := filepath.Abs(aBaseDir)
	if nil == err {
		poPostingBaseDirectory = dir
	}
} // SetPostingBaseDirectory()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// TPosting is a single article/posting to be used by a template.
	TPosting struct {
		id       string // hex. representation of date/time
		markdown []byte // (article-/file-)contents in Markdown markup
	}
)

// NewPosting returns a new posting structure with an empty article text.
//
//	`aID` if an empty string the `NewID()` function is called
// to provide a new article ID.
func NewPosting(aID string) *TPosting {
	if 0 == len(aID) {
		aID = NewID()
	}

	return &TPosting{id: aID}
} // NewPosting()

// After reports whether this posting is younger than the one
// identified by `aID`.
//
//	`aID` is the ID of another posting to compare.
func (p *TPosting) After(aID string) bool {
	return (p.id > aID)
} // After()

// Before reports whether this posting is older than the one
// identified by `aID`.
//
//	`aID` is the ID of another posting to compare.
func (p *TPosting) Before(aID string) bool {
	return (p.id < aID)
} // Before()

// Clear resets the internal fields to their respective zero values.
//
// This method does NOT remove the file (if any) associated with this
// posting/article; for that call the `Delete()` method.
func (p *TPosting) Clear() *TPosting {
	// p.markdown = []byte(``) // doesn't pass the equality test :-((
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

// `delFile()` removes `aFileName` from the filesystem
// returning a possible I/O error.
//
// A non-existing file is not considered an error here.
//
//	`aFileName` The name of the file to delete.
func delFile(aFileName *string) error {
	err := os.Remove(*aFileName)
	if nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
	}
	atomic.StoreUint32(&µCountCache, 0) // invalidate count cache

	return err
} // delFile()

// Delete removes the posting/article from the filesystem
// returning a possible I/O error.
//
// This method does NOT empty the markdown text of the object;
// for that call the `Clear()` method.
func (p *TPosting) Delete() error {
	fName := p.PathFileName()
	if 0 == len(fName) {
		return nil
	}

	return delFile(&fName)
} // Delete()

// Equal reports whether this posting is of the same time as `aID`.
//
//	`aID` The ID of the posting to compare with this one.
func (p *TPosting) Equal(aID string) bool {
	return timeID(p.id).Equal(timeID(aID))
} // Equal()

// Exists returns whether there is a file with more than zero bytes.
func (p *TPosting) Exists() bool {
	fi, err := os.Stat(pathname(p.id))
	if (nil != err) || fi.IsDir() {
		return false
	}

	return (0 < fi.Size())
} // Exists()

// ID returns the article's identifier.
//
// The identifier is based on the article's creation time
// and given in hexadecimal notation.
//
// This method allows the template to validate and use
// the placeholder `.ID`
func (p *TPosting) ID() string {
	return p.id
} // ID()

// Len returns the current length of the posting's Markdown text.
//
// If the markup is not already in memory this methods calls
// `TPosting.Load()` to read the text data from the filesystem.
func (p *TPosting) Len() int {
	if result := len(p.markdown); 0 < result {
		return result
	}
	if err := p.Load(); nil != err {
		apachelogger.Err("TPosting.Len()",
			fmt.Sprintf("TPosting.Load('%s'): %v", p.id, err))
	}

	return len(p.markdown)
} // Len()

// Load reads the Markdown from disk, returning a possible I/O error.
func (p *TPosting) Load() error {
	fName := p.PathFileName()
	if 0 == len(fName) {
		return fmt.Errorf("No filename for posting '%s'", p.id)
	}
	if _, err := os.Stat(fName); nil != err {
		return err // probably ENOENT
	}

	bs, err := ioutil.ReadFile(fName) // #nosec G304
	if nil != err {
		return err
	}
	if p.markdown = bytes.TrimSpace(bs); nil == p.markdown {
		// `bytes.TrimSpace()` returns `nil` instead of an empty slice
		p.markdown = []byte(``)
	}

	return nil
} // Load()

// `makeDir()` creates the directory for storing the article
// returning the article's path/file-name but w/o filename extension.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
func (p *TPosting) makeDir() (string, error) {
	fMode := os.ModeDir | 0775

	// Using the ID's first three characters leads to
	// directories worth about 52 days of data.
	// We need the year to guard against ID overflows.
	dir := fmt.Sprintf(`%04d%s`, timeID(p.id).Year(), p.id[:3])
	dirname := path.Join(poPostingBaseDirectory, dir)
	if err := os.MkdirAll(filepath.FromSlash(dirname), fMode); nil != err {
		return "", err
	}

	return path.Join(dirname, p.id), nil
} // makeDir()

// Markdown returns the Markdown of this article.
//
// If the markup is not already in memory this methods calls
// `TPosting.Load()` to read the text data from the filesystem.
func (p *TPosting) Markdown() []byte {
	if 0 < len(p.markdown) {
		// that's the easy path …
		return p.markdown
	}

	if err := p.Load(); nil != err {
		apachelogger.Err("TPosting.Markdown()",
			fmt.Sprintf("TPosting.Load('%s'): %v", p.id, err))
	}

	return p.markdown
} // Markdown()

// `pathname()` returns the complete article path-/filename.
func pathname(aID string) string {
	if 0 == len(aID) {
		return ""
	}

	// Using the ID's first three characters leads to
	// directories worth about 52 days of data.
	// We need the year to guard against ID overflows.
	dir := fmt.Sprintf(`%04d%s`, timeID(aID).Year(), aID[:3])

	return path.Join(poPostingBaseDirectory, dir, aID+`.md`)
} // pathname()

// PathFileName returns the article's complete path-/filename.
func (p *TPosting) PathFileName() string {
	return pathname(p.id)
} // PathFileName()

// Post returns the article's HTML markup.
func (p *TPosting) Post() template.HTML {
	// make sure we have the most recent version:
	p.Markdown()

	return template.HTML(MarkupTags(mdToHTML(p.markdown))) // #nosec G203
} // Post()

// Set assigns the article's Markdown text.
//
//	`aMarkdown` is the actual Markdown text of the article to assign.
func (p *TPosting) Set(aMarkdown []byte) *TPosting {
	if 0 < len(aMarkdown) {
		if p.markdown = bytes.TrimSpace(aMarkdown); nil == p.markdown {
			p.markdown = []byte(``)
		}
	} else {
		p.markdown = []byte(``)
	}

	return p
} // Set()

// Store writes the article's Markdown to disk returning
// the number of bytes written and a possible I/O error.
//
// The file is created on disk with mode `0640` (`-rw-r-----`).
func (p *TPosting) Store() (int64, error) {
	if _, err := p.makeDir(); nil != err {
		// without an appropriate directory we can't save anything …
		return 0, err
	}
	if 0 == len(p.markdown) {
		return 0, p.Delete()
	}

	fName := p.PathFileName()
	if 0 == len(fName) {
		return 0, fmt.Errorf("No filename for posting '%s'", p.id)
	}
	if err := ioutil.WriteFile(fName, p.markdown, 0640); nil != err {
		return 0, err
	}

	fi, err := os.Stat(fName)
	if nil != err {
		return 0, err
	}
	atomic.StoreUint32(&µCountCache, 0) // invalidate count cache

	return fi.Size(), nil
} // Store()

// String returns a stringified version of the posting object.
//
// Note: This is mainly for debugging purposes and has no real life use.
func (p *TPosting) String() string {
	return p.id + `: [[` + string(p.markdown) + `]]`
} // String()

// Time returns the posting's date/time.
func (p *TPosting) Time() time.Time {
	return timeID(p.id)
} // Time()

/*
func updatePostDirs() {
	dirnames, err := filepath.Glob(poPostingBaseDirectory + "/*")
	if nil != err {
		return
	}
	for _, dirname := range dirnames {
		filenames, err := filepath.Glob(dirname + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filenames) {
			continue // skip empty directory
		}
		for _, fName := range filenames {
			pName := strings.TrimPrefix(fName, dirname+"/")
			id := pName[:len(pName)-3]
			newName := pathname(id)

			dir := fmt.Sprintf("%04d%s", timeID(id).Year(), id[:3])
			dirname := path.Join(poPostingBaseDirectory, dir)
			os.MkdirAll(filepath.FromSlash(dirname), os.ModeDir|0775)
			os.Rename(fName, newName)
		}
	}
} // updatePostDirs()
*/

/* _EoF_ */
