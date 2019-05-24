/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

/*
* This files provides functions related to #hashtags/@mentions
 */

import (
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/mwat56/hashtags"
)

/*
* All functions starting with `go` are supposed to run in the background.
 */

// `goAddID()` checks a newly added posting for #hashtags and @mentions.
func goAddID(aList *hashtags.THashList, aID string, aText []byte) {
	aList.IDparse(aID, aText)
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
func goInitHashlist(aList *hashtags.THashList) {
	if _, err := aList.Load(); (nil == err) && (0 < aList.Len()) {
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

	aList.Store()
} // goInitHashlist()

// `goRemoveID()` removes `aID` from `aList`s items.
func goRemoveID(aList *hashtags.THashList, aID string) {
	aList.IDremove(aID)
} // goRemoveID()

// `goRenameID()` renames all references of `aOldID` to `aNewID`.
func goRenameID(aList *hashtags.THashList, aOldID, aNewID string) {
	aList.IDrename(aOldID, aNewID)
} // goRenameID()

// `goUpdateID()` updates the #hashtag/@mention references of `aID`.
func goUpdateID(aList *hashtags.THashList, aID string, aText []byte) {
	aList.IDupdate(aID, aText)
} // goUpdateID()

// `markupCloud()` returns a list with the markup of all existing
// #hashtags/@mentions.
func markupCloud(aList *hashtags.THashList) []template.HTML {
	var (
		class, url string
	)
	list := aList.CountedList()
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
