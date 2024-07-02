/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

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
	//    "github.com/mwat56/apachelogger"
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
func PostingBaseDirectory() string {
	return poPostingBaseDirectory
} // PostingBaseDirectory()

// `SetPostingBaseDirectory()` sets the base directory used for
// storing the articles/postings.
//
//	`aBaseDir` The base directory to use for storing articles/postings.
func SetPostingBaseDirectory(aBaseDir string) error {
	dir, err := filepath.Abs(aBaseDir)
	if nil != err {
		return err
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
//	`aFileName` The name of the file to delete.
func delFile(aFileName string) error {
	err := os.Remove(aFileName)
	if nil != err {
		if e, ok := err.(*os.PathError); ok && e.Err == syscall.ENOENT {
			return nil
		}
	}
	// atomic.StoreUint32(&µCountCache, 0) // invalidate count cache

	return err
} // delFile()

func id2dir(aID uint64) string {
	fname := id2str(aID)
	// Using the id's first three characters leads to
	// directories worth about 52 days of data.
	// We use the year to guard against ID overflows in a directory.
	dir := fmt.Sprintf(`%04d%s`, id2time(aID).Year(), fname[:3])

	return path.Join(poPostingBaseDirectory, dir)
} // id2dir()

func id2filename(aID uint64) string {
	dir := id2dir(aID)
	fname := id2str(aID)

	return path.Join(dir, fname) + `.md`
} // id2filename()

func id2str(aID uint64) (id string) {
	id = fmt.Sprintf("%x", aID)
	if 16 > len(id) {
		id = strings.Repeat("0", 16-len(id)) + id
	}

	return
} // id2str()

// `id2time()` returns a date/time represented by `aID`.
//
//	`aID` is a posting's ID.
func id2time(aID uint64) time.Time {
	return time.Unix(0, int64(aID))
} // id2time()

// `mkDir()` creates the directory for storing an article
// returning the created directory.
//
// The directory is created with filemode `0775` (`drwxrwxr-x`).
func mkDir(aID uint64) (string, error) {
	fMode := os.ModeDir | 0775
	dirname := id2dir(aID)
	if err := os.MkdirAll(filepath.FromSlash(dirname), fMode); nil != err {
		return "", err
	}

	//    return path.Join(dirname, p.id), nil
	return dirname, nil
} // mkDir()

/* _EoF_ */
