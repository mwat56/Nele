/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/mwat56/apachelogger"
	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

/* Defined in `persistence.go`:
type (
	TPosting struct {
		id           uint64    // integer representation of date/time
		lastModified time.Time // file modification time
		markdown     []byte    // article contents in Markdown markup
		mtx          *sync.RWMutex
	}

	TPostList []TPosting

	IPersistence interface {
		Create(aPost *TPosting) (int, error)
		Read(aID uint64) (*TPosting, error)
		Update(aPost *TPosting) (int, error)
		Delete(aID uint64) error

		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		PostingCount() uint32
	}
)
*/

// --------------------------------------------------------------------------
// constructor function:

// `NewPosting()` returns a new posting structure with the given article text.
//
// If `aID` is zero, the current time is used to generate a unique ID.
//
// Parameters:
//   - `aID`: A uint64 representing the unique identifier of the posting.
//   - `aText`: A string containing the initial Markdown text of the article.
//
// Returns:
//   - `*TPosting`: A new `TPosting` instance.
func NewPosting(aID uint64, aText string) *TPosting {
	if 0 == aID {
		aID = uint64(time.Now().UnixNano())
	}

	return &TPosting{
		id:           aID,
		lastModified: time.Now(),
		markdown:     []byte(aText),
		mtx:          new(sync.RWMutex),
	}
} // NewPosting()

// --------------------------------------------------------------------------
// TPosting methods

// `After()` reports whether this posting is younger than the one
// identified by `aID`.
//
// Parameters:
//
//   - `aID` is the ID of another posting to compare.
//
// Returns:
func (p *TPosting) After(aID uint64) bool {
	return (p.id > aID)
} // After()

// `Before()` reports whether this posting is older than the one
// identified by `aID`.
//
// Parameters:
//
//   - `aID` is the ID of another posting to compare.
//
// Returns:
func (p *TPosting) Before(aID uint64) bool {
	return (p.id < aID)
} // Before()

// `ChangeID()` changes the ID of the current posting including the
// persistence layer.
//
// Parameters:
//
// Returns:
func (p *TPosting) ChangeID(aID uint64) error {
	oldID := p.id
	p.id = aID

	if poPersistence.Exists(aID) {
		if _, err := poPersistence.Update(p); nil != err {
			p.id = oldID
			return err
		}
	} else if _, err := poPersistence.Create(p); nil != err {
		p.id = oldID
		return err
	}

	if poPersistence.Exists(oldID) {
		if err := poPersistence.Delete(oldID); nil != err {
			return err
		}
	}

	return nil
} // ChangeID()

// `Clear()` resets the text field to its zero value.
//
// This method does NOT remove the file (if any) associated with this
// posting/article; for that call the `Delete()` method.
//
// Returns:
func (p *TPosting) Clear() *TPosting {
	if nil == p {
		return p
	}

	p.mtx.Lock()
	defer p.mtx.Unlock()

	p.markdown = nil

	return p
} // Clear()

// `clone()` returns a copy of this posting/article.
//
// Returns:
func (p *TPosting) clone() *TPosting {
	return &TPosting{
		id:           p.id,
		lastModified: p.lastModified,
		markdown:     bytes.TrimSpace(p.markdown),
		mtx:          new(sync.RWMutex),
	}
} // clone()

// `Date()` returns the posting's date as a formatted string (`yyyy-mm-dd`).
//
// Returns:
func (p *TPosting) Date() string {
	y, m, d := id2time(p.id).Date()

	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
} // Date()

// `Delete()` removes the posting/article from the persistence layer
// returning a possible I/O error.
//
// This method does NOT empty the markdown text of the object;
// for that call the `Clear()` method.
//
// Returns:
func (p *TPosting) Delete() error {
	return poPersistence.Delete(p.id)
} // Delete()

// `Equal()` reports whether this posting is of the same time as `aID`.
//
// Parameters:
//
//   - `aID` The ID of the posting to compare with this one.
//
// Returns:
func (p *TPosting) Equal(aID uint64) bool {
	return (p.id == aID)
} // Equal()

// // `dir()` returns the fully qualified name of the directory where this
// // posting is stored.
// //
// // Returns:
// func (p *TPosting) dir() string {
// 	fname := poPersistence.PathFileName(p.id)

// 	return path.Dir(fname)
// } // dir()

// `Exists()` returns whether there is a file with more than zero bytes.
//
// Returns:
func (p *TPosting) Exists() bool {
	return poPersistence.Exists(p.id)
} // Exists()

// `ID()` returns the article's identifier.
//
// This method allows the template to validate and use
// the placeholder `.ID`
//
// Returns:
func (p *TPosting) ID() uint64 {
	return p.id
} // ID()

// `IDstr()` returns the article's identifier in string format.
//
// The identifier is based on the article's creation time
// and given in hexadecimal notation.
//
// This method allows the template to validate and use
// the placeholder `.ID`
//
// Returns:
//   - `string`: The article's identifier in string format.
func (p *TPosting) IDstr() string {
	return id2str(p.id)
} // IDstr()

// `LastModified()` returns the last-modified date/time of the posting.
//
// Returns:
func (p *TPosting) LastModified() string {
	return p.lastModified.Format(time.RFC1123)
} // LastModified()

// `Len()` returns the current length of the posting's Markdown text.
//
// If the markup is not already in memory this methods calls
// `TPosting.Load()` to read the text data from the filesystem.
//
// Returns:
func (p *TPosting) Len() int {
	if nil == p {
		return 0
	}

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if result := len(p.markdown); 0 < result {
		return result
	}

	if err := p.Load(); nil != err {
		apachelogger.Err("TPosting.Len()",
			fmt.Sprintf("TPosting.Load('%s'): %v", p.IDstr(), err))
	}

	return len(p.markdown)
} // Len()

// `Load()` reads the Markdown from disk, returning a possible I/O error.
//
// Returns:
func (p *TPosting) Load() error {
	pp, err := poPersistence.Read(p.id)
	if nil != err {
		return err
	}
	p.lastModified = pp.lastModified
	p.markdown = pp.markdown

	return nil
} // Load()

// `makeDir()` creates the directory for storing the article
// returning the article's path/file-name but w/o filename extension.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
//
// Returns:
func (p *TPosting) makeDir() (string, error) {
	return "", errors.New("TPosting.makeDir() need redirect to persistence layer")
} // makeDir()

// `PathFileName()` returns the article's complete path-/filename.
//
// Returns:
func (p *TPosting) PathFileName() string {
	return poPersistence.PathFileName(p.id)
} // PathFileName()

// `Post()` returns the article's HTML markup.
//
// This method uses the `Markdown()` method to get the latest version of
// the article's Markdown text.
// It then converts this text to HTML using the `MDtoHTML()` function and
// wraps it with the necessary HTML tags using the `MarkupTags()` function.
//
// The resulting HTML is returned as a `template.HTML` value.
//
// Returns:
func (p *TPosting) Post() template.HTML {
	// make sure we have the most recent version:
	return template.HTML(MarkupTags(MDtoHTML(p.Markdown()))) // #nosec G203
} // Post()

// `Markdown()` returns the Markdown of this article.
//
// If the markup is not already in memory this methods calls
// `TPosting.Load()` to read the text data from the filesystem.
//
// Returns:
func (p *TPosting) Markdown() []byte {
	if nil == p {
		return []byte(``)
	}

	p.mtx.RLock()
	defer p.mtx.RUnlock()

	if 0 < len(p.markdown) { // that's the easy path …
		return p.markdown
	}

	if err := p.Load(); nil != err {
		apachelogger.Err("TPosting.Markdown()", fmt.Sprintf("TPosting.Load('%d'): %v", p.id, err))
	}

	return p.markdown
} // Markdown()

// `Set()` assigns the article's Markdown text.
//
// Parameters:
//
//   - `aMarkdown`: The actual Markdown text of the article to assign.
//
// Returns:
func (p *TPosting) Set(aMarkdown []byte) *TPosting {
	if nil == p {
		return p
	}
	p.mtx.Lock()
	defer p.mtx.Unlock()

	if 0 < len(aMarkdown) {
		if p.markdown = bytes.TrimSpace(aMarkdown); nil == p.markdown {
			p.markdown = []byte(``)
		}
	} else {
		p.markdown = []byte(``)
	}
	p.lastModified = time.Now()

	return p
} // Set()

// `Store()` writes the article's Markdown to disk returning
// the number of bytes written and a possible I/O error.
//
// The file is created on disk with mode `0640` (`-rw-r-----`).
//
// Returns:
func (p *TPosting) Store() (int, error) {
	if nil == p {
		return 0, se.Wrap(errors.New("nil point"), 1)
	}

	return poPersistence.Create(p)
} // Store()

// `String()` returns a stringified version of the posting object.
//
// Note: This is mainly for debugging purposes and has no real life use.
//
// Returns:
func (p *TPosting) String() (rStr string) {
	if nil == p {
		return
	}
	rStr = id2str(p.id) + `: [[` + string(p.Markdown()) + `]]`

	return
} // String()

// `Time()` returns the posting's date/time.
//
// Returns:
func (p *TPosting) Time() time.Time {
	return id2time(p.id)
} // Time()

/* _EoF_ */
