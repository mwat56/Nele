/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides list_of_post related functions and methods.
 */

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mwat56/apachelogger"
)

/* * /
const (
	//TODO make this a `TPostList` property:
	plAscending  int = 1
	plDescending int = -1
	plNone       int = 0
)
/* */

// --------------------------------------------------------------------------
// constructor function:

// NewPostList returns a new (empty) TPostList instance.
func NewPostList() *TPostList {
	result := make(TPostList, 0, 32)

	return &result
} // NewPostList()

// --------------------------------------------------------------------------
// TPostList methods:

// Add appends `aPosting` to the list.
//
//	`aPosting` contains the actual posting's text.
func (pl *TPostList) Add(aPosting *TPosting) *TPostList {
	*pl = append(*pl, *aPosting)

	// We need an explicit return value (despite the in-place
	// modification of `pl`) ro allow for command chaining like
	// `list := NewPostList().Add(p3).Add(p1).Add(p2)`.
	return pl
} // Add()

/*
//TODO remove unused method
// Article adds the posting identified by `aID` to the list.
//
//	`aID` is the ID of the posting to add to this list.
func (pl *TPostList) Article(aID uint64) *TPostList {
	bgAddPosting(pl, aID)

	return pl
} // Article()
*/

// Day adds all postings of the current day to the list.
func (pl *TPostList) Day() *TPostList {
	t := time.Now()
	y, m, d := t.Year(), t.Month(), t.Day()
	tLo := time.Date(y, m, d, 0, 0, 0, -1, time.Local)
	tHi := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)

	return pl.prepareWalk(tLo, tHi)
} // Day()

// Delete removes `aPosting` from the list, returning the (possibly
// modified) list and whether the operation war successful.
//
//	`aPosting` is the posting o remove from this list.
func (pl *TPostList) Delete(aPosting *TPosting) (*TPostList, bool) {
	if idx := pl.Index(aPosting); 0 <= idx {
		if 0 == idx {
			*pl = (*pl)[1:] // remove first list entry
		} else if (len(*pl) - 1) == idx { // len - 1: because list is zero-based
			*pl = (*pl)[:idx] // omit last list entry
		} else {
			*pl = append((*pl)[:idx], (*pl)[idx+1:]...)
		}

		return pl, true // `aPosting` found and removed
	}

	return pl, false // `aPosting` not found in list
} // Delete()

var (
	// RegEx to check a posting's filename
	plFilenameRE = regexp.MustCompile(`[0-9a-fA-F]{16}\.md`)
)

// `doWalk()` traverses `aActDir` adding every valid posting
// to the list.
//
// Parameters:
//
//   - `aActDir`: The root directory for the traversal.
//   - `aLo`: The earliest ID time to use.
//   - `aHi`: The latest ID time to use.
func (pl *TPostList) doWalk(aActDir string, aLo, aHi time.Time) {

	//
	//TODO: delegate this to persistence layer
	//

	// We ignore all possible errors since we can't do anything about
	// them anyway and simply exclude those files from our list.
	walk := func(aPath string, aInfo os.FileInfo, aErr error) error {
		if (nil != aErr) || (0 == aInfo.Size()) || (aInfo.IsDir()) {
			return aErr
		}
		if plFilenameRE.Match([]byte(aInfo.Name())) {
			fName := aInfo.Name()
			fName = fName[:len(fName)-3] // w/o dot/file extension
			fID := filename2id(fName)
			tID := id2time(fID)
			if tID.After(aLo) && tID.Before(aHi) {
				bgAddPosting(pl, fID)
			}
		}
		return nil
	} // walk

	_ = filepath.Walk(aActDir, walk)
} // doWalk()

// `Index()` returns the 0-based list index of `aPosting`.
// In case `aPosting` was not found in list the return value
// will be `-1`.
//
// Parameters:
//
//   - `aPosting` is the posting to lookup in this list.
//
// Returns:
//   - `int`: The index of `aPosting` in this list, or `-1` if not found.
func (pl *TPostList) Index(aPosting *TPosting) int {
	for idx, post := range *pl {
		if post.id == aPosting.id {
			return idx
		}
	}

	return -1
} // Index()

// `IsSorted()` returns `true` if the list is sorted (in descending order),
// or `false` otherwise.
//
// Returns:
func (pl *TPostList) IsSorted() bool {
	return sort.SliceIsSorted(*pl, func(i, j int) bool {
		// return ((*pl)[i].id < (*pl)[j].id) // ascending
		return ((*pl)[i].id > (*pl)[j].id) // descending
	})
} // IsSorted()

// `Len()` returns the number of postings stored in this list.
//
// Returns:
func (pl *TPostList) Len() int {
	return len(*pl)
} // Len()

// `Month()` adds all postings of `aMonth` to the list.
//
// Parameters:
//
//   - `aYear`: The year to lookup; if zero the current year is used.
//   - `aMonth`: The year's month to lookup; if zero the current month is used.
//
// Returns:
func (pl *TPostList) Month(aYear int, aMonth time.Month) *TPostList {
	var (
		y int
		m time.Month
	)
	tLo := time.Now()

	if 0 < aYear {
		y = aYear
	} else {
		y = tLo.Year()
	}
	if (0 < aMonth) && (13 > aMonth) {
		m = aMonth
	} else {
		m = tLo.Month()
	}

	tLo = time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
	tHi := time.Date(y, m+1, 1, 0, 0, 0, -1, time.Local)

	return pl.prepareWalk(tLo, tHi)
} // Month()

// Newest adds the last `aNumber` of postings to the list.
//
// The resulting list is sorted in descending order (newest first)
// with at most `aNumber` posts.
//
// Parameters:
//
//   - `aNumber` The number of articles to show.
//   - `aStart` The start number to use.
//
// Returns:
func (pl *TPostList) Newest(aNumber, aStart int) error {
	var ( // re-use variables
		counter        int
		dName, pName   string
		dNames, fNames []string
		err            error
	)

	//
	//TODO: delegate this to persistence layer
	//

	if dNames, err = filepath.Glob(PostingBaseDirectory() + "/*"); nil != err {
		return err
	}

	// Sort the directory names to have the youngest entry first:
	sort.Slice(dNames,
		func(i, j int) bool {
			return (dNames[i] > dNames[j]) // descending
		})

	for _, dName = range dNames {
		if fNames, err = filepath.Glob(dName + "/*.md"); (nil != err) || (0 == len(fNames)) {
			continue // skip empty directory
		}

		// Sort the file names to have the youngest post first:
		sort.Slice(fNames,
			func(i, j int) bool {
				return (fNames[i] > fNames[j]) // descending
			})
		for _, pName = range fNames {
			counter++
			if counter <= aStart {
				continue
			}
			pName = strings.TrimPrefix(pName, dName+"/")
			id := filename2id(pName[:len(pName)-3]) // strip name extension
			bgAddPosting(pl, id)
			if len(*pl) >= aNumber {
				return nil
			}
		}
	}

	// Reaching this point of execution means:
	// there are less than `aNumber` posts available.
	return nil
} // Newest()

// `prepareWalk()` computes the first and last directory to process.
//
// Parameters:
//
//   - `aLo` is the earliest ID time to use.
//   - `aHi` is the latest ID time to use.
//
// Returns:
func (pl *TPostList) prepareWalk(aLo, aHi time.Time) *TPostList {
	if tn := time.Now(); tn.Before(aHi) {
		aHi = tn // exclude postings from the future ;-)
	}

	//
	// TODO: refer to the persistence layer
	//

	lID, hID := time2id(aLo), time2id(aHi)
	dirLo, dirHi := id2dir(lID), id2dir(hID)

	if dirLo == dirHi {
		// both, the first and last postings, are in the same directory
		pl.doWalk(dirLo, aLo, aHi)
	} else {
		pl.doWalk(dirLo, aLo, aHi)
		// A single directory holds about 52 days worth of files;
		// so even a whole month can't use more than two directories.
		pl.doWalk(dirHi, aLo, aHi)
	}

	return pl
} // prepareWalk()

// `Sort()` returns the list sorted by posting IDs (i.e. date/time)
// in descending order.
func (pl *TPostList) Sort() *TPostList {
	sort.Slice(*pl, func(i, j int) bool {
		// return ((*pl)[i].id < (*pl)[j].id) // ascending
		return ((*pl)[i].id > (*pl)[j].id) // descending
	})

	return pl
} // Sort()

// `Week()` adds all postings of the current week to the list.
//
// Parameters:
//
//   - `aYear` The year to lookup; if zero the current year is used.
//   - `aMonth` The year's month to lookup; if zero the current month is used.
//   - `aDay` The month's day to lookup; if zero the current day is used.
//
// Returns:
func (pl *TPostList) Week(aYear int, aMonth time.Month, aDay int) *TPostList {
	var y, d int
	var m time.Month
	tLo := time.Now()

	if 0 < aYear {
		y = aYear
	} else {
		y = tLo.Year()
	}
	if (0 < aMonth) && (13 > aMonth) {
		m = aMonth
	} else {
		m = tLo.Month()
	}
	if (0 < aDay) && (32 > aDay) {
		d = aDay
	} else {
		d = tLo.Day()
	}

	tLo = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	wd := tLo.Weekday() // day of the week (Sunday = 0, ...).
	if 0 == wd {
		d -= 6
	} else {
		d -= (int(wd) - 1)
	}

	tLo = time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	tHi := tLo.Add((time.Hour * 24 * 7) + 1)

	return pl.prepareWalk(tLo, tHi)
} // Week()

// --------------------------------------------------------------------------
// utility functions:

// `bgAddPosting()` adds a new posting (with `aID`) to `aPostList`.
//
// The data associated with `aID` is loaded from storage.
//
// Parameters:
//
//   - `aPostList`: The `TPostList` instance to add to.
//   - `aID` is the identifier of the new posting to add.
func bgAddPosting(aPostList *TPostList, aID uint64) {
	post := NewPosting(aID, "")

	if err := post.Load(); nil != err {
		apachelogger.Err("TPostList.bgAddPosting()",
			fmt.Sprintf("TPosting.Load(%s): %v", id2str(aID), err))
	} else {
		aPostList.Add(post)
	}
	// `Load()` errors are ignored since we can't do anything about it here.
} // bgAddPosting()

// `SearchPostings()` traverses the directories holding the postings
// looking for `aText` in all article files.
//
// The returned `TPostList` can be empty because (a) `aText` could
// not be compiled into a regular expression, (b) no files to
// search were found, or (c) no files matched `aText`.
//
// Parameters:
//
//   - `aText`: The text to look for in the postings.
//
// Returns:
func SearchPostings(aText string) *TPostList {
	return poPersistence.SearchPostings(aText)
} // SearchPostings()

/* _EoF_ */
