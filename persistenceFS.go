/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
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

type (
	// `TFSpersistence` is a file-based `IPersistence` implementation.
	TFSpersistence struct {
		_   struct{}
		mtx *sync.RWMutex // pointer to avoid copying warnings
	}
)

var (
	// Cache of last/current posting count.
	// see `[delFile]`, `[Count]`, `[TPosting.Store]`
	µCountCache uint32
)

// --------------------------------------------------------------------------
// private helper functions:

// `delFile()` removes `aFileName` from the filesystem
// returning a possible I/O error.
//
// A non-existing file is not considered an error here.
//
// Parameters:
//   - `aFileName` The name of the file to delete.
//
// Returns:
//   - `error`: Any error that occurred during the deletion process.
func delFile(aFileName string) error {
	if err := os.Remove(aFileName); nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
		return se.Wrap(err, 4)
	}

	return nil
} // delFile()

/*
// `filename2id()` converts a given file name to a `uint64` identifier.
//
// The file name is expected to be in hexadecimal format.
//
// The function first extracts the base name of the file using `path.Base()`.
// Then, it attempts to parse the base name as a `uint64` using
// `strconv.ParseUint()`.
//
// Parameters:
//   - `aFilename`: The file name to be converted.
//
// Returns:
//   - (uint64) The uint64 identifier corresponding to the input file name.
//   - (0) If an error occurs during parsing.
func filename2id(aFilename string) uint64 {
	if aFilename = strings.TrimSpace(aFilename); 0 == len(aFilename) {
		return 0 // empty filename
	}

	if aFilename = path.Base(aFilename); 16 > len(aFilename) {
		return 0 // invalid filename
	}

	if aFilename = strings.TrimSuffix(aFilename, ".md"); 16 > len(aFilename) {
		return 0 // invalid filename
	}

	return str2id(aFilename)
} // filename2id()
*/

/*
// `filename2time()` converts a given file name to a `Time` value.
//
// The file name is expected to be in hexadecimal format.
//
// The function first extracts the base name of the file using path.Base().
// Then, it attempts to parse the base name as a uint64 using
// strconv.ParseUint().
//
// The parsed uint64 is then used to create a Time value using time.Unix.
//
// Parameters:
//   - `aFilename`: The file name to be converted.
//
// Return Value:
//   - (time.Time) The Time value corresponding to the input file name.
func filename2time(aFilename string) time.Time {
	id := filename2id(aFilename)

	return time.Unix(0, int64(id))
} // filename2time()
*/

// `id2dir()` converts a given uint64 to a directory name.
//
// The function returns a directory name based on the provided uint64.
// The directory name is constructed by joining the year and the hexadecimal
// string representation of the provided uint64, with a "/" separator.
//
// Parameters:
//   - `aID`: A posting's ID to be converted to a directory name.
//
// Return Value:
//   - `string`: The directory name based on the `aID`.
func id2dir(aID uint64) string {
	fname := id2str(aID)
	// Using the id's first four characters leads to
	// directories worth about 52 days of data.
	// We use the year to guard against ID overflows in a directory.
	dir := fmt.Sprintf(`%04d%s`, id2time(aID).Year(), fname[:3])

	return path.Join(poPostingBaseDirectory, dir)
} // id2dir()

// `id2filename()` converts a given uint64 to a file name.
//
// The function returns a file name based on `aID`.
// The file name is constructed by joining the directory and the hexadecimal
// string representation of the provided uint64, with a ".md" extension.
//
// Parameters:
//   - `aID`: A posting's ID to be converted to a file name.
//
// Return Value:
//   - `string`: The file name based on the provided uint64.
func id2filename(aID uint64) string {
	dir := id2dir(aID)
	fname := id2str(aID)

	return path.Join(dir, fname) + `.md`
} // id2filename()

// `id2str()` converts a given uint64 to a hexadecimal string.
//
// The function returns a hexadecimal string representation of the
// provided uint64.
//
// Parameters:
//   - `aID`: The uint64 value to be converted to a hexadecimal string.
//
// Returns:
//   - `string`: The hexadecimal string representation of `aID`.
func id2str(aID uint64) (rStr string) {
	rStr = fmt.Sprintf("%x", aID)
	if 16 > len(rStr) {
		rStr = strings.Repeat("0", 16-len(rStr)) + rStr
	}

	return
} // id2str

// `mkDir()` creates the directory for storing an article
// returning the created directory.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
//
// Parameters:
//   - `aID`: The posting's ID.
//
// Returns:
//   - `string`: The created directory,
//   - `error`: TAny error that occurred during the creation process.
func mkDir(aID uint64) (string, error) {
	fMode := os.ModeDir | 0775
	dirname := id2dir(aID)
	if err := os.MkdirAll(dirname, fMode); nil != err {
		return "", se.Wrap(err, 1)
	}

	return dirname, nil
} // mkDir()

// `str2id()` converts a given hexadecimal string to a `uint64` integer.
//
// The function takes a hexadecimal string representation of a `uint64`
// value and attempts to parse that string into a `uint64` value.
//
// Parameters:
//   - `aHexString`: The string to be converted.
//
// Returns:
//   - (uint64) The `uint64` identifier corresponding to the input string.
//   - (0) If an error occurs during parsing.
func str2id(aHexString string) uint64 {
	if aHexString = strings.TrimSpace(aHexString); 16 > len(aHexString) {
		return 0 // invalid string
	}

	if ui64, err := strconv.ParseUint(aHexString, 16, 64); nil == err {
		return ui64
	}

	return 0
} // str2id()

// --------------------------------------------------------------------------

func init() {
	// ensure proper interface implementation
	var (
		_ IPersistence = TFSpersistence{}
		_ IPersistence = (*TFSpersistence)(nil)
	)
} // init()

// --------------------------------------------------------------------------
// constructor function

// `NewFSpersistence()` creates a new instance of `TFSpersistence`.
//
// It does not take any parameters.
//
// Returns:
//   - `*TFSpersistence`: A persistence instance instance.
func NewFSpersistence() *TFSpersistence {
	return &TFSpersistence{
		mtx: new(sync.RWMutex),
	}
} // NewFSpersistence()

// --------------------------------------------------------------------------
// TFSpersistence methods

// `Count()` returns the number of postings currently available.
//
// NOTE: This method is very resource intensive as it has to count all the
// posts stored in the filesystem.
//
// Returns:
//   - `uint32`: The number of available postings, or `0` in case of I/O errors.
//
// Side Effects:
//   - Updates the count cache.
func (fsp TFSpersistence) Count() (rCount uint32) {
	fsp.mtx.RLock()
	defer fsp.mtx.RUnlock()

	if rCount = atomic.LoadUint32(&µCountCache); 0 < rCount {
		return
	}

	var ( // re-use variable
		err            error
		dName          string
		dNames, fNames []string
	)
	// Apparently there's no current value ready so we compute a new one.
	// Instead of doing this in _one_ glob we actually do it in two
	// thus trading memory consumption with processing time and so
	// we're being better prepared for huge amounts of postings.
	if dNames, err = filepath.Glob(poPostingBaseDirectory + `/*`); nil != err {
		return 0 // we can't recover from this :-(
	}
	for _, dName = range dNames {
		if fNames, err = filepath.Glob(dName + `/*.md`); nil == err {
			rCount += uint32(len(fNames))
		}
	}
	atomic.StoreUint32(&µCountCache, rCount)

	return
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
//
// Side Effects:
//   - Invalidates the internal count cache.
func (fsp TFSpersistence) Create(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, se.Wrap(ErrEmptyPosting, 1)
	}
	fsp.mtx.Lock()
	defer fsp.mtx.Unlock()

	return fsp.store(aPost, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
} // Create()

// called by [Delete] and [store], both of which are already locked
func (fsp TFSpersistence) delete(aID uint64) error {
	err := delFile(id2filename(aID))
	if nil == err {
		atomic.StoreUint32(&µCountCache, 0) // invalidate count cache
	}

	return err
} // delete()

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
func (fsp TFSpersistence) Delete(aID uint64) error {
	fsp.mtx.Lock()
	defer fsp.mtx.Unlock()

	return fsp.delete(aID)
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
func (fsp TFSpersistence) Exists(aID uint64) bool {
	fsp.mtx.RLock()
	defer fsp.mtx.RUnlock()

	fName := id2filename(aID)
	fi, err := os.Stat(fName)
	if (nil != err) || (!fi.Mode().IsRegular()) {
		return false
	}

	return (0 < fi.Size())
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
func (fsp TFSpersistence) PathFileName(aID uint64) string {
	return id2filename(aID)
} // PathFileName()

// `Read()` reads the posting from disk, returning a possible I/O error.
//
// Parameters:
//
//   - `aID`: The unique identifier of the posting to be read.
//
// Returns:
//
//   - `*TPosting`: The `TPosting` instance containing the article's data, or `nil` if the file does not exist.
//   - 'error`: A possible I/O error, or `nil` on success.
func (fsp TFSpersistence) Read(aID uint64) (*TPosting, error) {
	fsp.mtx.Lock()
	defer fsp.mtx.Unlock()

	var ( // re-use variables
		bs  []byte
		err error
		fi  os.FileInfo
	)

	fName := id2filename(aID)
	if fi, err = os.Stat(fName); nil != err {
		return nil, se.Wrap(err, 1) // probably ENOENT
	}

	if bs, err = os.ReadFile(fName); /* #nosec G304 */ nil != err {
		return nil, se.Wrap(err, 1)
	}

	post := &TPosting{
		id:           aID,
		lastModified: fi.ModTime(),
		markdown:     bytes.TrimSpace(bs),
	}
	if nil == post.markdown {
		// `bytes.TrimSpace()` returns `nil` instead of an empty slice
		post.markdown = []byte(``)
	}

	return post, nil
} // Read()

// `Rename()` renames a posting from its old ID to a new ID.
//
// Parameters:
//   - aOldID: The unique identifier of the posting to be renamed.
//   - aNewID: The new unique identifier for the new posting.
//
// Returns:
//   - `error`: An error if the operation fails, or `nil` on success.
func (fsp TFSpersistence) Rename(aOldID, aNewID uint64) error {
	fsp.mtx.Lock()
	defer fsp.mtx.Unlock()

	oName := id2filename(aOldID)
	nName := id2filename(aNewID)
	nDir := id2dir(aNewID)

	fMode := os.ModeDir | 0775
	if err := os.MkdirAll(filepath.FromSlash(nDir), fMode); nil != err {
		apachelogger.Err("TFSpersistence.Rename()",
			fmt.Sprintf("os.Rename(%s, %s): %v", oName, nName, err))

		return se.Wrap(err, 4)
	}

	if err := os.Rename(oName, nName); nil != err {
		apachelogger.Err("TFSpersistence.Rename()",
			fmt.Sprintf("os.Rename(%s, %s): %v", oName, nName, err))

		return se.Wrap(err, 4)
	}

	return nil
} // Rename()

// `store()` writes the article's Markdown to disk returning
// the number of bytes written and a possible I/O error.
//
// Parameters:
// - `aPost`: A `TPosting` instance containing the article's data.
// - `aMode`: An integer representing the file mode for the OpenFile function.
//
// Returns:
// - `int`: The number of bytes written to the file.
// - 'error`:` A possible I/O error.
//
// Side Effects:
// - Invalidates the internal count cache.
func (fsp TFSpersistence) store(aPost *TPosting, aMode int) (int, error) {
	// Locking is done by `Create()` and `Update()`.
	var ( // re-use variables
		err    error
		mdFile *os.File
	)
	if _, err = mkDir(aPost.id); nil != err {
		// without an appropriate directory we can't save anything …
		return 0, err // err is already wrapped
	}

	if 0 == len(aPost.markdown) {
		return 0, fsp.delete(aPost.id)
	}

	fName := id2filename(aPost.id)
	mdFile, err = os.OpenFile(fName, aMode, 0640) /* #nosec G302 */
	if nil != err {
		return 0, se.Wrap(err, 2)
	}
	defer mdFile.Close()

	atomic.StoreUint32(&µCountCache, 0) // invalidate count cache
	aPost.lastModified = time.Now()

	return mdFile.Write(aPost.markdown)
} // store()

// `Update()` updates the article's Markdown on disk.
//
// It returns the number of bytes written to the file and a possible I/O error.
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
//
// Side Effects:
// - Invalidates the internal count cache.
func (fsp TFSpersistence) Update(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, se.Wrap(ErrEmptyPosting, 1)
	}
	fsp.mtx.Lock()
	defer fsp.mtx.Unlock()

	return fsp.store(aPost, os.O_WRONLY|os.O_TRUNC)
} // Update()

// `Walk()` visits all existing postings, calling `aWalkFunc`
// for each posting.
//
// Parameters:
//   - `aWalkFunc`: The function to call for each posting.
//
// Returns:
//   - `error`: a possible error occurring the traversal process.
func (fsp TFSpersistence) Walk(aWalkFunc TWalkFunc) error {
	// fsp.mtx.Lock()
	// defer fsp.mtx.Unlock()

	var (
		// RegEx to check a posting's filename
		filenameRE = regexp.MustCompile(`[0-9a-fA-F]{16}\.md`)
	)

	dNames, err := filepath.Glob(poPostingBaseDirectory + "/*")
	if nil != err {
		return se.Wrap(err, 1)
	}

	// Sort the directory names to have the youngest entry first:
	sortStr := func(a, b string) int {
		if a < b {
			return 1
		}
		if a > b {
			return -1
		}
		return 0
	}
	slices.SortFunc(dNames, sortStr)

dirLoop:
	for _, dName := range dNames {
		fNames, err := filepath.Glob(dName + "/*.md")
		if (nil != err) || (0 == len(fNames)) {
			continue // no files found
		}
		slices.SortFunc(fNames, sortStr)

	fileLoop:
		for _, fName := range fNames {
			fn := path.Base(fName)
			if !filenameRE.Match([]byte(fn)) {
				continue // no proper filename
			}
			fn = fn[:len(fn)-3] // exclude extension `.md`

			if err := aWalkFunc(str2id(fn)); nil != err {
				if errors.Is(err, ErrSkipAll) {
					break dirLoop
				}
				if errors.Is(err, ErrSkipFiles) {
					break fileLoop
				}
				return se.Wrap(err, 2)
			}
		}
	}

	return nil
} // Walk()

/* _EoF_ */
