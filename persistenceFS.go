/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	//     "strconv"

	//    "sync/atomic"

	"time"
	//    "github.com/mwat56/apachelogger"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

type (

	/* Defined in `persistence.go`:

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
		Create(aPost *TPosting) (int, error)
		Read(aID uint64) (*TPosting, error)
		Update(aPost *TPosting) (int, error)
		Delete(aID uint64) error

		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		PostingCount() uint32
	}

	*/

	// `TFSpersistence` is a file-based `IPersistence` implementation.
	TFSpersistence struct {
		_ struct{}
	}
)

func init() {
	var (
		_ IPersistence = TFSpersistence{}
		_ IPersistence = (*TFSpersistence)(nil)
	)
} // init()

// --------------------------------------------------------------------------
// TFSpersistence methods

func (fsp TFSpersistence) Create(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, errors.New("empty post not created")
	}

	return fsp.store(aPost, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
} // Create()

// `Delete()` removes the posting/article from the filesystem
// returning a possible I/O error.
func (fs TFSpersistence) Delete(aID uint64) error {
	return delFile(id2filename(aID))
} // Delete()

// `Exists()` returns whether there is a file with more than zero bytes.
func (fsp TFSpersistence) Exists(aID uint64) bool {
	fName := id2filename(aID)
	fi, err := os.Stat(fName)
	if (nil != err) || (!fi.Mode().IsRegular()) {
		return false
	}

	return (0 < fi.Size())
} // Exists()

// `PathFileName()` returns the posting's complete path-/filename.
func (fsp TFSpersistence) PathFileName(aID uint64) string {
	return id2filename(aID)
} // PathFileName()

// `PostingCount()` returns the number of postings currently available.
//
// In case of I/O errors the return value will be 0.
func (fsp TFSpersistence) PostingCount() uint32 {
	/*
	   if rCount = atomic.LoadUint32(&µCountCache); 0 < rCount {
	           return
	   }
	*/
	var ( // re-use variable
		err            error
		dName          string
		dNames, fNames []string
		rCount         uint32
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
	//        atomic.StoreUint32(&µCountCache, rCount)

	return rCount
} // PostingCount()

// `Read()` reads the posting from disk, returning a possible I/O error.
func (fsp TFSpersistence) Read(aID uint64) (*TPosting, error) {
	var ( // re-use variables
		bs  []byte
		err error
		fi  os.FileInfo
	)
	fName := id2filename(aID)

	if fi, err = os.Stat(fName); nil != err {
		return nil, err // probably ENOENT
	}

	if bs, err = ioutil.ReadFile(fName); /* #nosec G304 */ nil != err {
		return nil, err
	}

	p := &TPosting{
		id:           aID,
		lastModified: fi.ModTime(),
		markdown:     bytes.TrimSpace(bs),
		persistence:  fsp,
	}
	if nil == p.markdown {
		// `bytes.TrimSpace()` returns `nil` instead of an empty slice
		p.markdown = []byte(``)
	}

	return p, nil
} // Read()

// Store writes the article's Markdown to disk returning
// the number of bytes written and a possible I/O error.
//
// The file is created on disk with mode `0640` (`-rw-r-----`).
func (fsp TFSpersistence) store(aPost *TPosting, aMode int) (int, error) {
	var ( // re-use variables
		err    error
		mdFile *os.File
	)
	if _, err = mkDir(aPost.id); nil != err {
		// without an appropriate directory we can't save anything …
		return 0, err
	}

	if 0 == len(aPost.markdown) {
		return 0, fsp.Delete(aPost.id)
	}

	fName := id2filename(aPost.id)
	mdFile, err = os.OpenFile(fName, aMode, 0640) /* #nos
	ec G302 */
	if nil != err {
		return 0, err
	}
	defer mdFile.Close()

	// atomic.StoreUint32(&µCountCache, 0) // invalidate count cache
	aPost.lastModified = time.Now()

	return mdFile.Write(aPost.markdown)
} // store()

func (fsp TFSpersistence) Update(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, errors.New("empty post not updated")
	}

	return fsp.store(aPost, os.O_WRONLY|os.O_TRUNC)
} // Update()

// --------------------------------------------------------------------------
// constructor functions

func NewFSpersistence() *TFSpersistence {
	return &TFSpersistence{
		//        basepath: aBasepath,
	}
} // NewFSpersistence()

// --------------------------------------------------------------------------

/* _EoF_ */
