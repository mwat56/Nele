/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
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
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mwat56/apachelogger"
)

type (
	// TPostList is a list of postings to be injected
	// into a template/view.
	TPostList []TPosting
)

// Add appends `aPosting` to the list.
//
//	`aPosting` contains the actual posting's text.
func (pl *TPostList) Add(aPosting *TPosting) *TPostList {
	*pl = append(*pl, *aPosting)

	return pl
} // Add()

// Article adds the posting identified by `aID` to the list.
//
//	`aID` is the ID of the posting to add to this list.
func (pl *TPostList) Article(aID string) *TPostList {
	bgAddPosting(pl, aID)

	return pl
} // Article()

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
	// we ignore the very first list entry (hold the base directory)
	if idx := pl.Index(aPosting); 0 < idx {
		if (len(*pl) - 1) == idx { // len - 1: because list is zero-based
			*pl = (*pl)[:idx] // last list entry
		} else {
			*pl = append((*pl)[:idx], (*pl)[idx+1:]...)
		}
		return pl, true
	}

	return pl, false
} // Delete()

var (
	// RegEx to check a posting's filename
	plFilenameRE = regexp.MustCompile(`[0-9a-fA-F]{16}\.md`)
)

// `doWalk()` traverses `aActDir` adding every valid posting
// to the list.
//
//	`aActDir` the root directory for the traversal.
//
//	`aLo` is the earliest ID time to use.
//
//	`aHi` is the latest ID time to use.
func (pl *TPostList) doWalk(aActDir string, aLo, aHi time.Time) {
	// We ignore all possible errors since we can't do anything about
	// them anyway and simply exclude those files from our list.
	_ = filepath.Walk(aActDir,
		func(aPath string, aInfo os.FileInfo, aErr error) error {
			if (nil != aErr) || (0 == aInfo.Size()) || (aInfo.IsDir()) {
				return aErr
			}
			if plFilenameRE.Match([]byte(aInfo.Name())) {
				fName := aInfo.Name()
				fName = fName[:len(fName)-3] // w/o dot/file extension
				fID := timeID(fName)
				if fID.After(aLo) && fID.Before(aHi) {
					bgAddPosting(pl, fName)
				}
			}
			return nil
		})
} // doWalk()

// internal method for unit testing
func (pl *TPostList) in() *TPostList {
	for i, p := range *pl {
		fmt.Fprintf(os.Stdout, "[%d] %v\n", i, p)
	}
	return pl
} // in()

// Index returns the 0-based list index of `aPosting`.
// In case `aPosting` was not found in list the return value
// will be `-1`.
//
//	`aPosting` is the posting to lookup in this list.
func (pl *TPostList) Index(aPosting *TPosting) int {
	for idx, p := range *pl {
		if p.id == aPosting.id {
			return idx
		}
	}

	return -1
} // Index()

// IsSorted returns `true` if the list is sorted (in descending order),
// or `false` otherwise.
func (pl *TPostList) IsSorted() bool {
	return sort.SliceIsSorted(*pl, func(i, j int) bool {
		// return ((*pl)[i].id < (*pl)[j].id) // ascending
		return ((*pl)[i].id > (*pl)[j].id) // descending
	})
} // IsSorted()

// Len returns the number of postings stored in this list.
func (pl *TPostList) Len() int {
	return len(*pl)
} // Len()

// Month adds all postings of `aMonth` to the list.
//
//	`aYear` the year to lookup; if `0` (zero) the current year
// is used.
//
//	`aMonth` the year's month to lookup; if `0` (zero) the
// current month is used.
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
//	`aNumber` The number of articles to show.
//
//	`aStart` The start number to use.
func (pl *TPostList) Newest(aNumber, aStart int) error {
	dirnames, err := filepath.Glob(PostingBaseDirectory() + "/*")
	if nil != err {
		return err
	}
	// Sort the directory names to have the youngest entry first:
	sort.Slice(dirnames,
		func(i, j int) bool {
			return (dirnames[i] > dirnames[j]) // descending
		})
	counter := 0
	for _, dirname := range dirnames {
		filesnames, err := filepath.Glob(dirname + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filesnames) {
			continue // skip empty directory
		}
		// Sort the file names to have the youngest post first:
		sort.Slice(filesnames,
			func(i, j int) bool {
				return (filesnames[i] > filesnames[j]) // descending
			})
		for _, postname := range filesnames {
			counter++
			if counter <= aStart {
				continue
			}
			postname = strings.TrimPrefix(postname, dirname+"/")
			bgAddPosting(pl, postname[:len(postname)-3]) // strip name extension
			if len(*pl) >= aNumber {
				return nil
			}
		}
	}

	// Reaching this point of execution means:
	// there are less than `aNumber` posts available.
	return nil
} // Newest()

// prepareWalk() computes the first and last directory to process.
//
//	`aLo` is the earliest ID time to use.
//
//	`aHi` is the latest ID time to use.
func (pl *TPostList) prepareWalk(aLo, aHi time.Time) *TPostList {
	tn := time.Now()
	if tn.Before(aHi) {
		aHi = tn // exclude postings from the future ;-)
	}
	dirLo := path.Dir(pathname(newID(aLo)))
	dirHi := path.Dir(pathname(newID(aHi)))
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

// Sort returns the list sorted by posting IDs (i.e. date/time)
// in descending order.
func (pl *TPostList) Sort() *TPostList {
	sort.Slice(*pl, func(i, j int) bool {
		// return ((*pl)[i].id < (*pl)[j].id) // ascending
		return ((*pl)[i].id > (*pl)[j].id) // descending
	})

	return pl
} // Sort()

// Week adds all postings of the current week to the list.
//
//	`aYear` the year to lookup; if `0` (zero) the current year
// is used.
//
//	`aMonth` the year's month to lookup; if `0` (zero) the current
// month is used.
//
//	`aDay` the month's day to lookup; if `0` (zero) the current day is used.
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

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `bgAddPosting()` is a function adding a new posting to `aPostList`.
//
//	`aPostList` is the `TPostList` instance to add to.
//
//	`aID` is the identifier of the new posting to add;
// the associated file's contents are loaded from storage.
func bgAddPosting(aPostList *TPostList, aID string) {
	p := newPosting(aID)
	if err := p.Load(); nil != er
		apachelogger.Err("TPostList.bgAddPosting()",
			fmt.Sprintf("TPosting.Load(%s): %v", aID, err))r {
	} else {
		aPostList.Add(p)
	}
	// `Load()` errors are ignored since we can't do anything about it here.
} // bgAddPosting()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// NewPostList returns a new (empty) TPostList instance.
func NewPostList() *TPostList {
	result := make(TPostList, 0, 32)

	return &result
} // NewPostList()

/* _EoF_ */
