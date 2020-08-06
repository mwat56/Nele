/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
* This file provides functions related to #hashtags/@mentions.
 */

import (
	"fmt"
	"html/template"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mwat56/hashtags"
)

type (
	// This function type is used by `walkAllPosts()`.
	//
	//	`aList` The hashlist to use (update).
	//	`aPosting` The posting to handle.
	tWalkPostFunc func(aList *hashtags.THashList, aPosting *TPosting)
)

// `goWalkAllPosts()` visits all existing postings and calling `aWalkFunc`
// for each article.
//
//	`aList` The hashlist to use/update.
//	`aWalkFunc` The function to call for each posting.
func goWalkAllPosts(aList *hashtags.THashList, aWalkFunc tWalkPostFunc) {
	var ( // re-use variables
		dName, id, pName string
		dNames, fNames   []string
		err              error
	)
	if dNames, err = filepath.Glob(PostingBaseDirectory() + "/*"); nil != err {
		return // we can't recover from this :-(
	}

	for _, dName = range dNames {
		if fNames, err = filepath.Glob(dName + "/*.md"); nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 == len(fNames) {
			continue // skip empty directory
		}
		for _, pName = range fNames {
			id = strings.TrimPrefix(pName, dName+"/")
			aWalkFunc(aList, NewPosting(id[:len(id)-3])) // strip name extension
		}
	}
	_, _ = aList.Store()
} // goWalkAllPosts()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// RegEx to find PREformatted parts in an HTML page.
	htAHrefRE = regexp.MustCompile(`(?si)(<a[^>]*>.*?</a>)`)

	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`(#[0-9]+;)`)

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(`(?i)([@#][\p{L}\d_§-]+)(.?|$)`)
	//                                         11111111111111111  2222
)

// AddTagID checks a newly added `aPosting` for #hashtags and @mentions.
//
//	`aList` The hashlist to use (update).
//	`aPosting` The new posting to handle.
func AddTagID(aList *hashtags.THashList, aPosting *TPosting) {
	go aList.IDparse(aPosting.ID(), aPosting.Markdown())

	runtime.Gosched() // get the background operation started
} // AddTagID()

// InitHashlist initialises the hash list.
//
//	`aList` The list of #hashtags/@mentions to update.
func InitHashlist(aList *hashtags.THashList) {
	doInitHashlist := func(aHL *hashtags.THashList, aPosting *TPosting) {
		if 0 < aPosting.Len() {
			aHL.IDparse(aPosting.ID(), aPosting.Markdown())
		}
	} // doInitHashlist()

	go goWalkAllPosts(aList, doInitHashlist)
	runtime.Gosched() // get the background operation started
} // InitHashlist()

var (
	// RegEx to match texts like `#----`.
	htHyphenRE = regexp.MustCompile(`#[^-]*--`)

	// Lookup table for URL to use in `MarkupCloud()`.
	htListLookup = map[bool]string{
		true:  `/hl/`,
		false: `/ml/`,
	}
)

// MarkupCloud returns a list with the markup of all existing
// #hashtags/@mentions.
//
//	`aList` The list of #hashtags/@mentions to use.
func MarkupCloud(aList *hashtags.THashList) []template.HTML {
	var (
		class string // re-use variable
		idx   int
		item  hashtags.TCountItem
	)
	list := aList.CountedList()
	tl := make([]template.HTML, len(list))
	for idx, item = range list {
		if 7 > item.Count { // b000111
			class = "tc5"
		} else if 31 > item.Count { // b011111
			class = "tc25"
		} else if 63 > item.Count { // b111111
			class = "tc50"
		} else {
			class = "tc99"
		}

		tl[idx] = template.HTML(` <a href="` +
			htListLookup['#' == item.Tag[0]] + item.Tag[1:] +
			`" class="` + class + `" title=" ` +
			fmt.Sprintf("%d * %s", item.Count, item.Tag[1:]) +
			` ">` + item.Tag + `</a>`) // #nosec G203
	}

	return tl
} // MarkupCloud()

// MarkupTags returns `aPage` with all #hashtags/@mentions marked up
// as a HREF links.
//
//	`aPage` The HTML page to process.
func MarkupTags(aPage []byte) []byte {
	var ( // re-use variables
		cnt, l       int
		err          error
		re           *regexp.Regexp
		repl, search string
	)
	// (0) Check whether there are any links present:
	linkMatches := htAHrefRE.FindAll(aPage, -1)
	if (nil != linkMatches) || (0 < len(linkMatches)) {
		// (1) replace the links with a dummy text:
		for l, cnt = len(linkMatches), 0; cnt < l; cnt++ {
			search = regexp.QuoteMeta(string(linkMatches[cnt]))
			if re, err = regexp.Compile(search); nil == err {
				repl = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
				aPage = re.ReplaceAllLiteral(aPage, []byte(repl))
			}
		}
	}

	// (2) markup the #hashtags/@mentions:
	result := htHashMentionRE.ReplaceAllStringFunc(string(aPage),
		func(aString string) string {
			sub := htHashMentionRE.FindSubmatch([]byte(aString))
			if (nil == sub) || (0 == len(sub)) || (0 == len(sub[1])) {
				return aString
			}

			var suffix, url string

			hash := string(sub[1])
			// '_' can be both, part of the hashtag and italic
			// markup so we must remove it if it's at the end:
			if '_' == hash[len(hash)-1] {
				hash = hash[:len(hash)-1]
				suffix = `_`
			}
			if '#' == hash[0] {
				if 0 < len(sub[2]) {
					switch sub[2][0] {
					case '"':
						// double quote following a possible hashtag: most
						// probably an URL#fragment, hence leave it as is
						return aString
					case ')':
						// This is a tricky one: it can either be a
						// normal right round bracket or the end of
						// a Markdown link. Here we assume that it's
						// the latter one and ignore this match:
						return aString
					case ';':
						if htEntityRE.MatchString(aString) {
							// leave HTML entities as is
							return aString
						}
					}
				}
				if htHyphenRE.MatchString(hash) {
					return aString
				}
				url = "/hl/" + strings.ToLower(hash[1:])
			} else {
				url = "/ml/" + strings.ToLower(hash[1:])
			}
			if 0 < len(sub[2]) {
				suffix += string(sub[2])
			}

			return `<a href="` + url + `" class="smaller">` + hash + `</a>` + suffix
		})

	// (3) replace the link dummies with the real markup:
	for l, cnt = len(linkMatches), 0; cnt < l; cnt++ {
		search = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
		if re, err = regexp.Compile(search); nil == err {
			result = re.ReplaceAllLiteralString(result, string(linkMatches[cnt]))
		}
	}

	return []byte(result)
} // MarkupTags()

// ReadHashlist reads all postings to (re-)build the list of
// #hashtags/@mentions disregarding any pre-existing list.
//
//	`aList` The list of #hashtags/@mentions to build.
func ReadHashlist(aList *hashtags.THashList) {
	InitHashlist(aList.Clear())
} // ReadHashlist()

// RemoveIDTags removes `aID` from `aList's` items.
//
//	`aList` The hashlist to update.
//	`aID` The ID of the posting to remove.
func RemoveIDTags(aList *hashtags.THashList, aID string) {
	go aList.IDremove(aID)
	runtime.Gosched() // get the background operation started
} // RemoveIDTags()

// RenameIDTags renames all references of `aOldID` to `aNewID`.
//
//	`aList` The hashlist to update.
//	`aOldID` The posting's old ID.
//	`aNewID` The posting's new ID.
func RenameIDTags(aList *hashtags.THashList, aOldID, aNewID string) {
	go aList.IDrename(aOldID, aNewID)
	runtime.Gosched() // get the background operation started
} // RenameIDTags()

// ReplaceTag replaces the #tags/@mentions in `aList`.
//
//	`aList` The hashlist to update.
//	`aSearchTag` The old #tag/@mention to find.
//	`aReplaceTag` The new #tag/@mention to use.
func ReplaceTag(aList *hashtags.THashList, aSearchTag, aReplaceTag string) {
	if (nil == aList) || (0 == len(aSearchTag)) || (0 == len(aReplaceTag)) {
		return
	}

	switch aSearchTag[0] {
	case '#', '@':
		switch aReplaceTag[0] {
		case '#', '@':
			// nothing to do
		default:
			return
		}
	default:
		return
	}

	doReplaceTag := func(aHL *hashtags.THashList, aPosting *TPosting) {
		if 0 == aPosting.Len() {
			return
		}
		searchRE, err := regexp.Compile(`(?i)\` + aSearchTag)
		if nil != err {
			return
		}
		if !searchRE.Match(aPosting.Markdown()) {
			return
		}

		txt := searchRE.ReplaceAllLiteral(aPosting.Markdown(), []byte(aReplaceTag))
		_, _ = aPosting.Set(txt).Store()
		aHL.IDremove(aPosting.ID()).IDparse(aPosting.ID(), txt)
	} // doReplaceTag()

	go goWalkAllPosts(aList, doReplaceTag)
	runtime.Gosched() // get the background operation started
} // ReplaceTag()

// UpdateTags updates the #hashtag/@mention references of `aPosting`.
//
//	`aList` The hashlist to update.
//	`aPosting` The new posting to process.
func UpdateTags(aList *hashtags.THashList, aPosting *TPosting) {
	go aList.IDupdate(aPosting.ID(), aPosting.Markdown())
	runtime.Gosched() // get the background operation started
} // UpdateTags()

/* _EoF_ */
