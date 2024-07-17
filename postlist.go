/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/mwat56/apachelogger"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

/* * /
const (
	//TODO make this a `TPostList` property:
	plAscending  int = 1
	plDescending int = -1
	plNone       int = 0
)
/* */

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

// NewPostList returns a new (empty) TPostList instance.
func NewPostList() *TPostList {
	result := make(TPostList, 0, 32)

	return &result
} // NewPostList()

// --------------------------------------------------------------------------
// TPostList methods:

// `Add()` appends `aPosting` to the list.
//
// Parameters:
//   - `aPosting` contains the actual posting's text.
//
// Returns:
//   - `*TPostList`: The updated list.
func (pl *TPostList) Add(aPosting *TPosting) *TPostList {
	*pl = append(*pl, *aPosting)

	// We need an explicit return value (despite the in-place
	// modification of `pl`) to allow for command chaining like
	// `list := NewPostList().Add(p3).Add(p1).Add(p2)`.
	return pl.Sort()
} // Add()

// `Day()` adds all postings of the current day to the list.
//
// Returns:
//   - `*TPostList`: A list with the postings of the current day.
func (pl *TPostList) Day() *TPostList {
	t := time.Now()
	y, m, d := t.Year(), t.Month(), t.Day()

	tLo := time.Date(y, m, d, 0, 0, 0, -1, time.Local)
	tHi := time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)

	return pl.doTimeWalk(tLo, tHi)
} // Day()

// `Delete()` removes `aPosting` from the list, returning the (possibly
// modified) list and whether the operation war successful.
//
// Parameters:
//   - `aPosting`: The posting to remove from this list.
//
// Returns:
//   - `*TPostList`: The updated list.
//   - `bool`: Whether `aPosting` was successfully removed.
func (pl *TPostList) Delete(aPosting *TPosting) (*TPostList, bool) {
	idx := pl.Index(aPosting)
	if 0 > idx {
		return pl, false // `aPosting` not found in list
	}
	if 0 == idx {
		*pl = (*pl)[1:] // remove first list entry
	} else if (len(*pl) - 1) == idx { // len - 1: because list is zero-based
		*pl = (*pl)[:idx] // omit last list entry
	} else {
		*pl = append((*pl)[:idx], (*pl)[idx+1:]...)
	}

	return pl, true // `aPosting` found and removed
} // Delete()

// `doTimeWalk()` computes the first and last posting to process.
//
// Parameters:
//   - `aLo` is the earliest ID time to use.
//   - `aHi` is the latest ID time to use.
//
// Returns:
//   - `*TPostList`: A list with postings between `aLo` and `aHi`.
func (pl *TPostList) doTimeWalk(aLo, aHi time.Time) *TPostList {
	if tn := time.Now(); tn.Before(aHi) {
		aHi = tn // exclude postings from the future ;-)
	}

	wf := func(aID uint64) error {
		tID := id2time(aID)
		if tID.After(aLo) && tID.Before(aHi) {
			bgAddPosting(pl, aID)
		}

		return nil
	} // wf()
	poPersistence.Walk(wf)

	return pl
} // doTimeWalk()

// `Index()` returns the 0-based list index of `aPosting`.
// In case `aPosting` was not found in list the return value
// will be `-1`.
//
// Parameters:
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

// `IsSorted()` returns whether the list is sorted (in descending order).
//
// Returns:
//   - `bool`: `true` if the list is sorted in descending order.
func (pl *TPostList) IsSorted() bool {
	return sort.SliceIsSorted(*pl, func(i, j int) bool {
		// return ((*pl)[i].id < (*pl)[j].id) // ascending
		return ((*pl)[i].id > (*pl)[j].id) // descending
	})
} // IsSorted()

// `Len()` returns the number of postings stored in this list.
//
// Returns:
//
//	`int`: The number of postings in the current list.
func (pl *TPostList) Len() int {
	return len(*pl)
} // Len()

// `Month()` adds all postings of `aMonth` to the list.
//
// Parameters:
//   - `aYear`: The year to lookup; if zero the current year is used.
//   - `aMonth`: The year's month to lookup; if zero the current month is used.
//
// Returns:
//   - `*TPostList`: A list with the postings of the given year and month.
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

	return pl.doTimeWalk(tLo, tHi)
} // Month()

// `Newest()` adds the last `aNumber` of postings to the list.
//
// The resulting list is sorted in descending order (newest first)
// with at most `aNumber` posts.
//
// Parameters:
//   - `aNumber` The number of articles to show.
//   - `aStart` The start number to use.
//
// Returns:
//   - `error`: A possible error during processing of the request.
func (pl *TPostList) Newest(aNumber, aStart int) error {
	wf := func(aID uint64) error {
		if len(*pl) < aNumber {
			bgAddPosting(pl, aID)
		} else {
			return ErrSkipAll
		}

		return nil
	} // wf()

	return poPersistence.Walk(wf)
} // Newest()

// `Sort()` returns the list sorted by posting IDs (i.e. date/time)
// in descending order.
//
// Returns:
//   - `*TPostList`: The current list of postings in descending order.
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
//   - `aYear` The year to lookup; if zero the current year is used.
//   - `aMonth` The year's month to lookup; if zero the current month is used.
//   - `aDay` The month's day to lookup; if zero the current day is used.
//
// Returns:
//   - `*TPostList`: A list with the postings of the given year, month, and day.
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

	return pl.doTimeWalk(tLo, tHi)
} // Week()

// --------------------------------------------------------------------------
// utility functions:

// `bgAddPosting()` adds a new posting (with `aID`) to `aPostList`.
//
// The data associated with `aID` is loaded from storage.
//
// Parameters:
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
	// errors are ignored since we can't do anything about it here.
} // bgAddPosting()

// `SearchPostings()` traverses the directories holding the postings
// looking for `aText` in all article files.
//
// The returned `TPostList` can be empty because (a) `aText` could
// not be compiled into a regular expression, (b) no files to
// search were found, or (c) no files matched `aText`.
//
// Parameters:
//   - `aText`: The text to look for in the postings.
//
// Returns:
//   - `*TPostList`: The found list.
func SearchPostings(aText string) *TPostList {
	result := NewPostList()

	pattern, err := regexp.Compile(fmt.Sprintf("(?s)%s", aText))
	if nil != err {
		return result // empty list
	}

	wf := func(aID uint64) error {
		post := NewPosting(aID, "")
		if err := post.Load(); nil != err {
			return nil
		}
		if !pattern.Match(post.markdown) {
			return nil
		}
		result.Add(post)

		return nil
	} // wf()

	poPersistence.Walk(wf)

	return result
} // SearchPostings()

/* _EoF_ */
