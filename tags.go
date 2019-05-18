/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package blog

/*
* This files provides functions related to #hashtags/@mentions
 */

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	hashtags "github.com/mwat56/go-hashtags"
)

/*
* All functions starting with `go` are supposed to run in the background.
 */

// `goAddID()` checks a newly added posting for #hashtags and @mentions.
func goAddID(aList *hashtags.THashList, aFilename string, aID string, aText []byte) {
	oldlen := aList.LenTotal()

	aList.IDparse(aID, aText)

	if aList.LenTotal() != oldlen {
		aList.Store(aFilename)
	}
} // goAddID()

// `goInitHashlist()` initialises the hash list.
func goInitHashlist(aList *hashtags.THashList, aFilename string) {
	if _, err := aList.Load(aFilename); nil == err {
		return // assume everything is uptodate
	}
	dirnames, err := filepath.Glob(postingBaseDirectory + "/*")
	if nil != err {
		return // we can't recover from this :-(
	}
	for _, dName := range dirnames {
		filesnames, err := filepath.Glob(dName + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filesnames) {
			continue // skip empty directory
		}
		for _, postName := range filesnames {
			id := strings.TrimPrefix(postName, dName+"/")
			if txt, err := ioutil.ReadFile(postName); nil == err {
				aList.IDparse(id[:len(id)-3], txt) // strip name extension
			}
		}
	}
	if 0 < aList.Len() {
		aList.Store(aFilename)
	}
} // goInitHashlist()

// `goRemoveID()` removes `aID` from `aList`s items.
func goRemoveID(aList *hashtags.THashList, aFilename, aID string) {
	oldLen := aList.LenTotal()

	aList.IDremove(aID)

	if aList.LenTotal() != oldLen {
		aList.Store(aFilename)
	}
} // goRemoveID()

// `goRenameID()` renames all references of `aOldID` to `aNewID`.
func goRenameID(aList *hashtags.THashList, aFilename, aOldID, aNewID string) {
	oldLen := aList.LenTotal()

	aList.IDrename(aOldID, aNewID)

	if aList.LenTotal() != oldLen {
		aList.Store(aFilename)
	}
} // goRenameID()

// `goUpdateID()` updates the #hashtag/@mention references of `aID`.
func goUpdateID(aList *hashtags.THashList, aFilename, aID string, aText []byte) {
	oldLen := aList.LenTotal() //FIXME this doesn't catch cases when …
	// … the number of removals equals the number of additions.

	aList.IDupdate(aID, aText)

	if aList.LenTotal() != oldLen {
		aList.Store(aFilename)
	}
} // goUpdateID()

// `markupCloud()` returns a string mith the markup of all existing
// #hashtags/@mentions
func markupCloud(aHashList *hashtags.THashList) (rTags template.HTML) {
	var (
		class, url string
	)
	list := aHashList.CountedList()
	for _, item := range list {
		if 5 > item.Count {
			class = "tc5"
		} else if 25 > item.Count {
			class = "tc25"
		} else if 50 > item.Count {
			class = "tc50"
		} else {
			class = "tc99"
		}
		if '#' == item.Tag[0] {
			url = "/ht/" + item.Tag[1:]
		} else {
			url = "/mt/" + item.Tag[1:]
		}
		rTags += template.HTML(` <a href="` + url + `" class="` + class + `" title=" ` + item.Tag + ` ">` + item.Tag + `</a> `)
	}
	return
} // markupCloud()

/* _EoF_ */
