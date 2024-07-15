/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	ht "github.com/mwat56/hashtags"
	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

type (
	// `TPosting` is a single article/posting..
	TPosting struct {
		id           uint64        // integer representation of date/time
		lastModified time.Time     // file modification time
		markdown     []byte        // article contents in Markdown markup
		mtx          *sync.RWMutex // pointer to avoid copying warnings
	}

	// `TPostList` is a list of postings to be injected
	// into a template/view.
	TPostList []TPosting

	// This function type is used by `WalkPostings()`.
	//
	//	`aList` The hashlist to use (update).
	//	`aPosting` The posting to handle.
	TWalkPostFunc func(aList *ht.THashTags, aPosting *TPosting)

	// This function type is used by `Walk()`.
	//
	// Parameters:
	//	- `aID`: The ID of the posting to handle.
	TWalkFunc func(aID uint64) error

	// `IPersistence` defines a persistence layer for storing `TPosting`
	// objects.
	// It uses a CRUD interface with some additional methods as documented
	// below.
	IPersistence interface {
		//
		// `Create()` creates a new persistent posting.
		//
		// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
		// is returned.
		//
		// Parameters:
		//
		//	- `aPost`: The `TPosting` instance containing the article's data.
		//
		// Returns:
		//
		//	- `int`: The number of bytes written.
		//	- 'error`: A possible error, or `nil` on success.
		Create(aPost *TPosting) (int, error)

		//
		// `Delete()` removes the posting/article from the persistence layer
		// and returns a possible I/O error.
		//
		// Parameters:
		//
		//	- `aID`: The unique identifier of the posting to delete.
		//
		// Returns:
		//
		//	- 'error`: A possible error, or `nil` on success.
		Read(aID uint64) (*TPosting, error)

		//
		// `Update()` updates the article's data in the persistence layer.
		//
		// It returns the number of bytes written and a possible I/O error.
		//
		// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
		// is returned.
		//
		// Parameters:
		//
		//	- `aPost`: A `TPosting` instance containing the article's data.
		//
		// Returns:
		//
		//	- `int`: The number of bytes written.
		//	- 'error`:` A possible error, or `nil` on success.
		Update(aPost *TPosting) (int, error)

		//
		// `Delete()` removes the posting/article from the persistence layer.
		//
		// Parameters:
		//
		//	- `aID`: The unique identifier of the posting to delete.
		//
		// Returns:
		//
		//	 - 'error`: A possible error, or `nil` on success.
		Delete(aID uint64) error

		//
		// `Exists()` checks if a post with the given ID exists in the
		// persistence layer.
		//
		// Parameters:
		//
		//	- `aID`: The unique identifier of the posting to check.
		//
		// Returns:
		//
		//	- `bool`: `true` if the post exists, `false` otherwise.
		Exists(aID uint64) bool

		//
		// `PathFileName()` returns the posting's complete path-/filename.
		//
		// NOTE: The actual definition of the path-/filename depends on
		// the implementation of this interface. In a file-based persistence
		// layer it would be a `/path/directory/filename` string.
		// However, in a database-based persistence layer it would be the
		// `/path/file` of the database file.
		//
		// Parameters:
		//
		//	- `aID`: The unique identifier of the posting to handle.
		//
		// Returns:
		//
		//	- `string`: The path-/filename associated with `aID`.
		PathFileName(aID uint64) string

		//
		// `PostingCount()` returns the number of postings available.
		//
		// Returns:
		//
		//	 - `uint32`: The number of available postings, or `0`
		// in case of errors.
		PostingCount() uint32

		//
		// `Rename()` renames a posting from its old ID to a new ID.
		//
		// Parameters:
		//
		//	- aOldID: The unique identifier of the posting to be renamed.
		//	- aNewID: The new unique identifier for the new posting.
		//
		// Returns:
		//
		//	- `error`: An error if the operation fails, or `nil` on success.
		Rename(aOldID, aNewID uint64) error

		//
		// `SearchPostings()` is looking for `aText` in all article files.
		//
		// The returned `TPostList` can be empty because (a) `aText` could
		// not be compiled into a regular expression, (b) no files to
		// search were found, or (c) no files matched `aText`.
		//
		// Parameters:
		//
		//	`aText`: The text to look for in the postings.
		//
		// Returns:
		//
		//	`*TPostList`: A list of found postings.
		SearchPostings(aText string) *TPostList

		//
		// `Walk()` visits all existing postings, calling `aWalkFunc`
		// for each posting.
		//
		// Parameters:
		//	- `aWalkFunc`: The function to call for each posting.
		//
		// Returns:
		//	- `error`: a possible error occurring the traversal process.
		Walk(aWalkFunc TWalkFunc) error
	}
)

var (
	// `ErrEmptyPosting` is returned when a `nil` posting is passed to
	// a method.
	ErrEmptyPosting = errors.New("empty post")

	// `poPostingBaseDirectory` is the base directory for storing articles.
	//
	// This variable's value must be set initially before creating any
	// `TPosting` or `TPostList` instances.
	// After that it should be considered `read/only`.
	// Its default value is `./postings`.
	//
	// - see `PostingBaseDirectory()`, `SetPostingBaseDirectory()`
	poPostingBaseDirectory = func() string {
		dir, _ := filepath.Abs(`./postings`)
		return dir
	}()

	// The persistence layer to actually use:
	poPersistence IPersistence
)

// --------------------------------------------------------------------------
// public utility functions:

// `Persistence()` returns the persistence layer to actually use creating,
// updating, deleting, searching, and walking through postings.
//
// Returns:
//   - `IPersistence`: The persistence layer to use for storing/retrieving postings.
func Persistence() IPersistence {
	return poPersistence
} // Persistence()

// `PostingBaseDirectory()` returns the base directory used for
// storing the postings.
//
// Returns:
// - `string`: The base directory tu use.
func PostingBaseDirectory() string {
	return poPostingBaseDirectory
} // PostingBaseDirectory()

// `SetPersistence()` sets the persistence layer to actually use.
//
// Parameters:
//   - `aPersistence`: The persistence layer to use for storing/retrieving postings.
func SetPersistence(aPersistence IPersistence) {
	poPersistence = aPersistence
} // SetPersistence()

// `SetPostingBaseDirectory()` sets the base directory used for
// storing the postings.
//
// Parameters:
// - `aBaseDir` The base directory to use for storing articles/postings.
//
// Returns:
// - `error`: Any error that occurred during the setting of the base directory.
//
// Example:
//
//	// Set the base directory to "/path/to/new/base/directory"
//	err := nele.SetPostingBaseDirectory("/path/to/new/base/directory")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Get the current base directory
//	fmt.Println(nele.PostingBaseDirectory())
//
//	// Set the base directory back to the default value
//	err = nele.SetPostingBaseDirectory("./postings")
//	if err != nil {
//		log.Fatal(err)
//	}
func SetPostingBaseDirectory(aBaseDir string) error {
	dir, err := filepath.Abs(aBaseDir)
	if nil != err {
		return se.Wrap(err, 2)
	}

	fMode := os.ModeDir | 0775
	if err := os.MkdirAll(dir, fMode); nil != err {
		return se.Wrap(err, 1)
	}

	poPostingBaseDirectory = dir

	return nil
} // SetPostingBaseDirectory()

// `id2time()` returns a date/time represented by `aID`.
//
// Parameters:
// - `aID`: A posting's ID to be converted to a `time.Time`.
//
// Returns:
// - `time.Time`: The UnixNano value of the provided time.Time.
func id2time(aID uint64) time.Time {
	return time.Unix(0, int64(aID))
} // id2time()

// `time2id()` converts a given `aTime` to an integer representation
//
// The function returns the UnixNano value of the provided time.Time.
//
// Parameters:
// - `aTime` (time.Time) The time to be converted to a uint64 integer.
//
// Return Value:
// - `uint64`: The UnixNano value of the provided time.Time.
func time2id(aTime time.Time) uint64 {
	return uint64(aTime.UnixNano())
} // time2id()

/* _EoF_ */
