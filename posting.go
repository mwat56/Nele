/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/mwat56/apachelogger"
	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

/* Defined in `persistence.go`:
type (
	TPosting struct {
		id           uint64    // integer representation of date/time
		lastModified time.Time // file modification time
		markdown     []byte    // article contents in Markdown markup
	}

	TPostList []TPosting

	TWalkFunc func(aID uint64) error

	IPersistence interface {
		Create(aPost *TPosting) (int, error)
		Read(aID uint64) (*TPosting, error)
		Update(aPost *TPosting) (int, error)
		Delete(aID uint64) error

		Count() uint32
		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		Rename(aOldID, aNewID uint64) error
		Walk(aWalkFunc TWalkFunc) error
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
	}
} // NewPosting()

// --------------------------------------------------------------------------
// TPosting methods

// `After()` reports whether this posting is younger than the one
// identified by `aID`.
//
// Parameters:
//   - `aID`: The ID of another posting to compare.
//
// Returns:
//   - `bool`: Whether the current posting is older than `aID`.
func (p *TPosting) After(aID uint64) bool {
	return (p.id > aID)
} // After()

// `Before()` reports whether this posting is older than the one
// identified by `aID`.
//
// Parameters:
//   - `aID`: The ID of another posting to compare.
//
// Returns:
//   - `bool`: Whether the current posting is younger than `aID`.
func (p *TPosting) Before(aID uint64) bool {
	return (p.id < aID)
} // Before()

// `ChangeID()` changes the ID of the current posting including the
// persistence layer.
//
// Note: This method is provided for rare cases when a posting's ID
// has to be changed.
//
// Parameters:
//   - `aID`: is the ID of another posting to compare.
//
// Returns:
//   - `error`: A possible error during processing the request.
func (p *TPosting) ChangeID(aID uint64) error {
	oldID := p.id
	p.id = aID

	if err := poPersistence.Rename(oldID, aID); nil != err {
		p.id = oldID

		return err
	}

	return nil
} // ChangeID()

// `Clear()` resets the text field to its zero value.
//
// This method does NOT remove the file (if any) associated with this
// posting/article; for that call the `Delete()` method.
//
// Returns:
//   - `*TPosting`: The current posting without any text.
func (p *TPosting) Clear() *TPosting {
	if nil == p {
		return p
	}

	p.markdown = nil
	p.lastModified = time.Now()

	return p
} // Clear()

// `clone()` returns a copy of this posting/article.
//
// Returns:
//   - `*TPosting`: A clone of the current posting.
func (p *TPosting) clone() *TPosting {
	return &TPosting{
		id:           p.id,
		lastModified: p.lastModified,
		markdown:     bytes.TrimSpace(p.markdown),
	}
} // clone()

// `Date()` returns the posting's date as a formatted string (`yyyy-mm-dd`).
//
// Returns:
//   - `string`: The posting's creation date in string format.
func (p *TPosting) Date() string {
	y, m, d := id2time(p.id).Date()

	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
} // Date()

// `Delete()` removes the posting/article from the persistence layer
// returning a possible I/O error.
//
// This method does NOT empty the markdown text of the object;
// for that call the `[Clear]` method.
//
// Returns:
//   - `error`: A possible error during processing the request.
func (p *TPosting) Delete() error {
	return poPersistence.Delete(p.id)
} // Delete()

// `Equal()` reports whether this posting is of the same time as `aID`.
//
// Parameters:
//   - `aID`: The ID of the posting to compare with this one.
//
// Returns:
//   - `bool`: `true` if the posting is of the same as `aID`.
func (p *TPosting) Equal(aID uint64) bool {
	return (p.id == aID)
} // Equal()

// `Exists()` returns whether there is a file with more than zero bytes.
//
// Returns:
//   - `bool`: `true` if the posting exists in the persistence layer.
func (p *TPosting) Exists() bool {
	return poPersistence.Exists(p.id)
} // Exists()

// `ID()` returns the article's identifier.
//
// Returns:
//   - `uint64`: The posting's unique ID.
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
// The format of the returned string would be like
// "Mon, 02 Jan 2006 15:04:05 MST"
//
// Returns:
//   - `string`: The article's last modified date.
func (p *TPosting) LastModified() string {
	return p.lastModified.Format(time.RFC1123)
} // LastModified()

// `Len()` returns the current length of the posting's Markdown text.
//
// If the markup is not already in memory this method calls
// `[Load]` to read the text data from the persistence layer.
//
// Returns:
//   - `int`: The current length of the posting's text.
func (p *TPosting) Len() int {
	if nil == p {
		return 0
	}

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
//   - `error`: A possible error during processing the request.
func (p *TPosting) Load() error {
	pp, err := poPersistence.Read(p.id)
	if nil != err {
		return err
	}
	p.lastModified = pp.lastModified
	p.markdown = pp.markdown

	return nil
} // Load()

// `PathFileName()` returns the article's complete path-/filename.
//
// NOTE: The actual definition of the path-/filename depends on
// the implementation of this interface. In a file-based persistence
// layer it would be a `/path/directory/filename` string. However,
// in a database-based persistence layer it would be the `/path/file`
// of the database file.
//
// Returns:
//   - `string`: The current posting's path-/filename.
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
//   - `template.HTML`: The HTML markup of the current posting's text.
func (p *TPosting) Post() template.HTML {
	// make sure we have the most recent version:
	return template.HTML(MarkupTags(MDtoHTML(p.Markdown()))) // #nosec G203
} // Post()

// `Markdown()` returns the Markdown of this article.
//
// If the markdown is not already in memory this methods calls
// `[Load]` to read the text data from the persistence layer.
//
// Returns:
//   - `[]byte`: The current posting's text.
func (p *TPosting) Markdown() []byte {
	if nil == p {
		return []byte(``)
	}

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
//   - `aMarkdown`: The actual Markdown text to assign.
//
// Returns:
//   - `*TPosting`: The updated posting.
func (p *TPosting) Set(aMarkdown []byte) *TPosting {
	if nil == p {
		return p
	}

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
// The actual storing is delegated to the persistence layer.
//
// Returns:
//   - `int`: The number of bytes written.
//   - 'error`: A possible error, or `nil` on success.
func (p *TPosting) Store() (int, error) {
	if nil == p {
		return 0, se.Wrap(errors.New("nil point"), 1)
	}

	return poPersistence.Create(p)
} // Store()

// `String()` returns a stringified version of the posting instance.
//
// Note: This is mainly for debugging purposes and has no real life use.
//
// Returns:
//   - `string`: The stringified version of the current posting.
func (p *TPosting) String() string {
	return id2str(p.id) + `: [[` + string(p.Markdown()) + `]]`
} // String()

// `Time()` returns the posting's date/time.
//
// Returns:
//   - `time.Time`: The current posting's creation time.
func (p *TPosting) Time() time.Time {
	return id2time(p.id)
} // Time()

/* _EoF_ */
