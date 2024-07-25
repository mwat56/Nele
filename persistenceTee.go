/*
Copyright © 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

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
		first  *IPersistence
		second *IPersistence
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
func NewTeePersistence(aFirst, aSecond *IPersistence) *TeePersistence {
	return &TeePersistence{
		first:  aFirst,
		second: aSecond,
	}
} // NewTeePersistence()

// --------------------------------------------------------------------------
// TeePersistence methods

// `Count()` returns the number of postings currently available.
//
// NOTE: This method is very resource intensive as it has to count all the
// posts stored in the filesystem.
//
// Returns:
//   - `int32`: The number of available postings, or `0` in case of I/O errors.
func (tp TeePersistence) Count() int {
	var result int

	return result
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
//   - `int`: The number of bytes written to the file.
//   - 'error`:` A possible error, or `nil` on success.
func (tp TeePersistence) Create(aPost *TPosting) (int, error) {

	return 0, nil
} // Create()

// `Delete()` removes the posting/article from the filesystem
// and returns a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to delete.
//
// Returns:
//   - 'error`: A possible I/O error, or `nil` on success.
//
// Side Effects:
//   - Invalidates the internal count cache.
func (tp TeePersistence) Delete(aID uint64) error {

	return nil
} // Delete()

// `Exists()` checks if a file with the given ID exists in the filesystem.
//
// It returns a boolean value indicating whether the file exists.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to check.
//
// Returns:
//   - `bool`: `true` if the file exists, `false` otherwise.
func (tp TeePersistence) Exists(aID uint64) bool {

	return false
} // Exists()

// `PathFileName()` returns the posting's complete path-/filename.
//
// The returned path-/filename is in the format:
//
//	<base_directory>/<posting_id>.md
//
// Parameters:
//   - `aID`: The unique identifier of the posting to handle.
//
// Returns:
//   - `*string`: The path-/filename associated with `aID`.
func (tp TeePersistence) PathFileName(aID uint64) string {

	return ""
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

	return nil, nil
} // Read()

// `Rename()` renames a posting from its old ID to a new ID.
//
// Parameters:
//   - aOldID: The unique identifier of the posting to be renamed.
//   - aNewID: The new unique identifier for the new posting.
//
// Returns:
//   - `error`: An error if the operation fails, or `nil` on success.
func (tp TeePersistence) Rename(aOldID, aNewID uint64) error {

	return nil
} // Rename()

// `Search()` retrieves a list of postings based on a search term.
//
// A zero value of `aLimit` means: no limit alt all.
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

	return nil, nil
} // Search()

// `Update()` updates the article's Markdown on disk.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
// - `aPost`: A `TPosting` instance containing the article's data.
//
// Returns:
// - `int`: The number of bytes written to the file.
// - 'error`:` A possible I/O error, or `nil` on success.
func (tp TeePersistence) Update(aPost *TPosting) (int, error) {

	return 0, nil
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

	return nil
} // Walk()

/* _EoF_ */
