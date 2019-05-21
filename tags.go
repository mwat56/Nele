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
	oldCRC := aList.Checksum()

	aList.IDparse(aID, aText)

	if aList.Checksum() != oldCRC {
		aList.Store(aFilename)
	}
} // goAddID()

// `doCheckPost` returns whether there is a file identified
// by `aID` containing `aHash`.
func doCheckPost(aHash, aID string) bool {
	p := newPosting(aID)
	if !p.Exists() {
		return false
	}
	if err := p.Load(); nil != err {
		return false
	}
	if txt := strings.ToLower(string(p.Markdown())); 0 > strings.Index(txt, aHash) {
		return false
	}

	return true
} // doCheckPost()

// `goCheckHashes()` walks all postings referenced by `aList`.
func goCheckHashes(aList *hashtags.THashList) {
	aList.Walk(doCheckPost)
} // goCheckHashes()

// `goInitHashlist()` initialises the hash list.
func goInitHashlist(aList *hashtags.THashList, aFilename string) {
	if _, err := aList.Load(aFilename); nil == err {
		go goCheckHashes(aList)
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
	oldCRC := aList.Checksum()

	aList.IDremove(aID)

	if aList.Checksum() != oldCRC {
		aList.Store(aFilename)
	}
} // goRemoveID()

// `goRenameID()` renames all references of `aOldID` to `aNewID`.
func goRenameID(aList *hashtags.THashList, aFilename, aOldID, aNewID string) {
	oldCRC := aList.Checksum()

	aList.IDrename(aOldID, aNewID)

	if aList.Checksum() != oldCRC {
		aList.Store(aFilename)
	}
} // goRenameID()

// `goUpdateID()` updates the #hashtag/@mention references of `aID`.
func goUpdateID(aList *hashtags.THashList, aFilename, aID string, aText []byte) {
	oldCRC := aList.Checksum()

	aList.IDupdate(aID, aText)

	if aList.Checksum() != oldCRC {
		aList.Store(aFilename)
	}
} // goUpdateID()

// `markupCloud()` returns a list with the markup of all existing
// #hashtags/@mentions.
func markupCloud(aHashList *hashtags.THashList) []template.HTML {
	var (
		class, url string
	)
	list := aHashList.CountedList()
	tl := make([]template.HTML, len(list))
	for idx, item := range list {
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
		tl[idx] = template.HTML(` <a href="` + url + `" class="` + class + `" title=" ` + item.Tag + ` ">` + item.Tag + `</a> `)
	}

	return tl
} // markupCloud()

/* _EoF_ */
