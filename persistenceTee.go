/*
Copyright © 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"fmt"
	"sync"

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

		Count() int
		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		Rename(aOldID, aNewID uint64) error
		Search(aText string, aOffset, aLimit uint) (*TPostList, error)
		Walk(aWalkFunc TWalkFunc) error
	}
)
*/

type (
	// `TeePersistence` is a `IPersistence` implementation that uses
	// two different`IPersistence` implementations to load/store postings.
	TeePersistence struct {
		first  IPersistence
		second IPersistence
	}
)

// --------------------------------------------------------------------------

func init() {
	// ensure proper interface implementation
	var (
		_ IPersistence = TeePersistence{}
		_ IPersistence = (*TeePersistence)(nil)
	)
} // init()

// --------------------------------------------------------------------------
// constructor function

// `NewFSpersistence()` creates a new instance of `TeePersistence`.
//
// It does not take any parameters.
//
// Returns:
//   - `*TeePersistence`: A persistence instance instance.
func NewTeePersistence(aFirst, aSecond IPersistence) *TeePersistence {
	return &TeePersistence{
		first:  aFirst,
		second: aSecond,
	}
} // NewTeePersistence()

// --------------------------------------------------------------------------
// TeePersistence methods

// `Count()` returns the number of postings currently available.
//
// Returns:
//   - `int`: The number of available postings, or `0` in case of I/O errors.
func (tp TeePersistence) Count() int {
	var (
		r, r1, r2 int
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		r1 = tp.first.Count()
	}()

	go func() {
		defer wg.Done()
		r2 = tp.second.Count()
	}()

	wg.Wait()
	if r1 >= r2 {
		r = r1
	} else {
		r = r2
	}

	return r
} // Count()

// `Create()` creates a new posting in the filesystem.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
//   - `aPost`: The `TPosting` instance containing the article's data.
//
// Returns:
//   - `int`: The number of bytes written to the storage.
//   - 'error`:` A possible error, or `nil` on success.
func (tp TeePersistence) Create(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, se.Wrap(ErrEmptyPosting, 1)
	}
	var (
		e, e1, e2 error
		i, i1, i2 int
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		i1, e1 = tp.first.Create(aPost)
	}()

	go func() {
		defer wg.Done()
		i2, e2 = tp.second.Create(aPost)
	}()

	wg.Wait()
	if nil != e1 {
		e = e1
		i = i1
	} else if nil != e2 {
		e = e2
		i = i2
	}
	//
	//TODO: try to UPDATE the one failing from the successful one
	//

	return i, e
} // Create()

// `Delete()` removes the posting/article from the storage
// and returns a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to delete.
//
// Returns:
//   - 'error`: A possible I/O error, or `nil` on success.
func (tp TeePersistence) Delete(aID uint64) error {
	var (
		e, e1 error
		wg    sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		e1 = tp.first.Delete(aID)
	}()

	go func() {
		defer wg.Done()
		// we can ignore a failure on the second storage
		// because the post was meant to be deleted anyway
		_ = tp.second.Delete(aID)
	}()

	wg.Wait()
	if nil != e1 {
		e = e1
	}

	return e
} // Delete()

// `Exists()` checks if a posting with the given ID exists in the storage.
//
// It returns a boolean value indicating whether the post exists.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to check.
//
// Returns:
//   - `bool`: `true` if the post exists, `false` otherwise.
func (tp TeePersistence) Exists(aID uint64) bool {
	var (
		b, b1, b2 bool
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		b1 = tp.first.Exists(aID)
	}()

	go func() {
		defer wg.Done()
		b2 = tp.second.Exists(aID)
	}()

	wg.Wait()
	if !b1 {
		b = b1
	} else if !b2 {
		b = b2
		//
		//TODO: try to CREATE the missing posting
		//
	}

	return b
} // Exists()

// `PathFileName()` returns the posting's complete path-/filename.
//
// The returned path-/filename is in the format:
//
//	"1:<file/path_of_first_storage>|2:<file/path_of_second_storage>"
//
// Parameters:
//   - `aID`: The unique identifier of the posting to handle.
//
// Returns:
//   - `string`: The path-/filename associated with `aID`.
func (tp TeePersistence) PathFileName(aID uint64) string {
	f1 := tp.first.PathFileName(aID)
	f2 := tp.second.PathFileName(aID)

	return fmt.Sprintf("1:%q|2:%q", f1, f2)
} // PathFileName()

// `Read()` reads the posting from disk, returning a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to be read.
//
// Returns:
//   - `*TPosting`: The `TPosting` instance containing the article's data, or `nil` if the file does not exist.
//   - 'error`: A possible I/O error, or `nil` on success.
func (tp TeePersistence) Read(aID uint64) (*TPosting, error) {
	var (
		e, e1, e2 error
		p, p1, p2 *TPosting
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		p1, e1 = tp.first.Read(aID)
	}()

	go func() {
		defer wg.Done()
		p2, e2 = tp.second.Read(aID)
	}()

	wg.Wait()
	if nil == e1 {
		p = p1
		//
		//TODO: check whether to UPDATE the second storage
		//
	} else {
		e = e1
		if nil == e2 {
			p = p2
		}
	}

	return p, e
} // Read()

// `Rename()` renames a posting from its old ID to a new ID.
//
// Parameters:
//   - aOldID: The unique identifier of the posting to be renamed.
//   - aNewID: The new unique identifier for the posting.
//
// Returns:
//   - `error`: An error if the operation fails, or `nil` on success.
func (tp TeePersistence) Rename(aOldID, aNewID uint64) error {
	var (
		e, e1, e2 error
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		e1 = tp.first.Rename(aOldID, aNewID)
	}()

	go func() {
		defer wg.Done()
		e2 = tp.second.Rename(aOldID, aNewID)
	}()

	wg.Wait()
	if nil != e1 {
		e = e1
	} else if nil != e2 {
		e = e2
	}

	return e
} // Rename()

// `Search()` retrieves a list of postings based on a search term.
//
// A zero value of `aLimit` means: no practical limit at all.
//
// The returned `TPostList` type is a slice of `TPosting` instances, where
// `TPosting` is a struct representing a single posting. If the returned
// slice is an empty list then no matching postings were found; if it is
// `nil` it means there was an error retrieving the matches.
//
// Parameters:
//   - `aText`: The search query string.
//   - `aOffset`: An offset in the database result set of the search results.
//   - `aLimit`: The maximum number of search results to return.
//
// Returns:
//   - `*TPostList`: The list of search results, or `nil` in case of errors.
//   - `error`: If the search operation fails, or `nil` on success.
func (tp TeePersistence) Search(aText string, aOffset, aLimit uint) (*TPostList, error) {
	var (
		e, e1, e2    error
		pl, pl1, pl2 *TPostList
		wg           sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		pl1, e1 = tp.first.Search(aText, aOffset, aLimit)
	}()

	go func() {
		defer wg.Done()
		pl2, e2 = tp.second.Search(aText, aOffset, aLimit)
	}()

	wg.Wait()
	if nil == e1 {
		pl = pl1
	} else {
		e = e1
		if nil == e2 {
			pl = pl2
		}
	}
	//
	//TODO: COMPARE the lists and add possibly missing elements
	//

	return pl, e
} // Search()

// `Update()` updates the article's Markdown on disk.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
// - `aPost`: A `TPosting` instance with the article's updated data.
//
// Returns:
// - `int`: The number of bytes written to the file.
// - 'error`:` A possible I/O error, or `nil` on success.
func (tp TeePersistence) Update(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, se.Wrap(ErrEmptyPosting, 1)
	}
	var (
		e, e1, e2 error
		i, i1, i2 int
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		i1, e1 = tp.first.Update(aPost)
	}()

	go func() {
		defer wg.Done()
		i2, e2 = tp.second.Update(aPost)
	}()

	wg.Wait()
	if nil != e1 {
		e = e1
		i = i1
	} else if nil != e2 {
		e = e2
		i = i2
	}
	//
	//TODO: try to CREATE the one failing from the successful one
	//

	return i, e
} // Update()

// `Walk()` visits all existing postings, calling `aWalkFunc`
// for each posting.
//
// Parameters:
//   - `aWalkFunc`: The function to call for each posting.
//
// Returns:
//   - `error`: a possible error occurring the traversal process.
func (tp TeePersistence) Walk(aWalkFunc TWalkFunc) error {
	var (
		e, e1, e2 error
		wg        sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		e1 = tp.first.Walk(aWalkFunc)
	}()

	go func() {
		defer wg.Done()
		e2 = tp.second.Walk(aWalkFunc)
	}()

	wg.Wait()
	if nil != e1 {
		e = e1
	} else if nil != e2 {
		e = e2
	}

	return e
} // Walk()

/* _EoF_ */
