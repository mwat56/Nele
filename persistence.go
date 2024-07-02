/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	//     "strconv"
	"sync"
	//    "sync/atomic"
	"syscall"
	"time"

	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

type (
	// TPosting is a single article/posting..
	TPosting struct {
		id           uint64    // integer representation of date/time
		lastModified time.Time // file modification time
		markdown     []byte    // article contents in Markdown markup
		persistence  IPersistence
		mtx          *sync.RWMutex
	}

	//
	IPersistence interface {
		// `Create()` writes `aPost` to the persistence layer.
		Create(aPost *TPosting) (int, error)

		//
		Read(aID uint64) (*TPosting, error)

		//
		Update(aPost *TPosting) (int, error)

		//
		Delete(aID uint64) error

		// `Exists()` returns whether there is a file with more than zero bytes.
		Exists(aID uint64) bool

		// `PathFileName()` returns the posting's complete path-/filename.
		PathFileName(aID uint64) string

		// `PostingCount()` returns the number of postings currently available.
		//
		// In case of I/O errors the return value will be 0.
		PostingCount() uint32
	}
)

var (
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

	/*
	   // Cache of last/current posting count.
	   // see `delFile()`, `PostingCount()`, `TPosting.Store()`
	   µCountCache uint32
	*/

	// The persistence layer to actually use:
	poPersistence IPersistence
)

// --------------------------------------------------------------------------
// public utility functions:

// `PostingBaseDirectory()` returns the base directory used for
// storing the articles/postings.
//
// Returns:
// - `string`: The base directory tu use.
func PostingBaseDirectory() string {
	return poPostingBaseDirectory
} // PostingBaseDirectory()

// `SetPostingBaseDirectory()` sets the base directory used for
// storing the articles/postings.
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

	poPostingBaseDirectory = dir
	return nil
} // SetPostingBaseDirectory()

// --------------------------------------------------------------------------
// private helper functions:

// `delFile()` removes `aFileName` from the filesystem
// returning a possible I/O error.
//
// A non-existing file is not considered an error here.
//
// Parameters:
// - `aFileName` The name of the file to delete.
//
// Returns:
// - `error`: Any error that occurred during the deletion process.
func delFile(aFileName string) error {
	err := os.Remove(aFileName)
	if nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
		return se.Wrap(err, 5)
	}
	// atomic.StoreUint32(&µCountCache, 0) // invalidate count cache

	return nil
} // delFile()

// `id2dir()` converts a given uint64 to a directory name.
//
// The function returns a directory name based on the provided uint64.
// The directory name is constructed by joining the year and the hexadecimal
// string representation of the provided uint64, with a "/" separator.
//
// Parameters:
// - `aID`: A posting's ID to be converted to a directory name.
//
// Return Value:
// - `string`: The directory name based on the `aID`.
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
// - `aID`: A posting's ID to be converted to a file name.
//
// Return Value:
// - `string`: The file name based on the provided uint64.
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
// - `aID`: The uint64 value to be converted to a hexadecimal string.
//
// Return Value:
// - `string`: The hexadecimal string representation of `aID`.
func id2str(aID uint64) (rStr string) {
	rStr = fmt.Sprintf("%x", aID)
	if 16 > len(rStr) {
		rStr = strings.Repeat("0", 16-len(rStr)) + rStr
	}

	return
} // id2str()

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

// `mkDir()` creates the directory for storing an article
// returning the created directory.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
//
// Parameters:
// - `aID`: The posting's ID.
//
// Returns:
// - `string`: The created directory,
// - `error`: TAny error that occurred during the creation process.
func mkDir(aID uint64) (string, error) {
	fMode := os.ModeDir | 0775
	dirname := id2dir(aID)
	if err := os.MkdirAll(filepath.FromSlash(dirname), fMode); nil != err {
		return "", se.Wrap(err, 1)
	}

	return dirname, nil
} // mkDir()

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
