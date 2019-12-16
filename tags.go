/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
* This file provides functions related to #hashtags/@mentions
 */

import (
	"bytes"
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

// `walkAllPosts()` visits all existing postings and calling `aWalkFunc`
// for each article.
//
//	`aList` The hashlist to use/update.
//	`aWalkFunc` The function to call for each posting.
func walkAllPosts(aList *hashtags.THashList, aWalkFunc tWalkPostFunc) {
	dirnames, err := filepath.Glob(PostingBaseDirectory() + "/*")
	if nil != err {
		return // we can't recover from this :-(
	}
	for _, mdName := range dirnames {
		filesnames, err := filepath.Glob(mdName + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 == len(filesnames) {
			continue // skip empty directory
		}
		for _, postName := range filesnames {
			fName := strings.TrimPrefix(postName, mdName+"/")
			aWalkFunc(aList, NewPosting(fName[:len(fName)-3])) // strip name extension
		}
	}
	_, _ = aList.Store()
} // walkAllPosts()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

var (
	// RegEx to find PREformatted parts in an HTML page.
	htAHrefRE = regexp.MustCompile(`(?si)(<a[^>]*>.*?</a>)`)

	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`(#[0-9]+;)`)

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(`(?i)([@#][§\wÄÖÜß-]+)(.?|$)`)
)

// AddTagID checks a newly added `aPosting` for #hashtags and @mentions.
//
//	`aList` The hashlist to use (update).
//	`aPosting` The new posting to handle.
func AddTagID(aList *hashtags.THashList, aPosting *TPosting) {
	go func() {
		aList.IDparse(aPosting.ID(), aPosting.Markdown())
	}()
	runtime.Gosched() // get the background operation started
} // AddTagID()

// `initHashlist()` initialises the hash list.
//
//	`aList` The list of #hashtags/@mentions to update.
func initHashlist(aList *hashtags.THashList) {
	doInitHashlist := func(aHL *hashtags.THashList, aPosting *TPosting) {
		if 0 < aPosting.Len() {
			aHL.IDparse(aPosting.ID(), aPosting.Markdown())
		}
	} // doInitHashlist()

	go walkAllPosts(aList, doInitHashlist)
	runtime.Gosched() // get the background operation started
} // initHashlist()

// InitHashlist initialises the hash list.
//
//	`aList` The list of #hashtags/@mentions to update.
func InitHashlist(aList *hashtags.THashList) {
	if _, err := aList.Load(); (nil == err) && (0 < aList.Len()) {
		// `doCheckPost()` returns whether there is a file identified
		// by `aID` containing `aHash`.
		//
		// The function's result is `false` (1) if the file associated
		// with `aID` doesnt't exist, or (2) if the file can't be
		// read, or (3) the given `aHash` can't be found in the
		// posting's text.
		//
		//	`aHash` The hashtag to check for.
		//	`aID` The ID of the posting to handle.
		doCheckPost := func(aHash, aID string) bool {
			p := NewPosting(aID)
			if !p.Exists() {
				return false
			}
			if err := p.Load(); nil != err {
				return false
			}
			txt := bytes.ToLower(p.Markdown())

			return (0 <= bytes.Index(txt, []byte(aHash)))
		} // doCheckPost()

		go aList.Walk(doCheckPost)
		return // assume everything is up-to-date
	}

	initHashlist(aList)
} // InitHashlist()

// MarkupCloud returns a list with the markup of all existing
// #hashtags/@mentions.
//
//	`aList` The list of #hashtags/@mentions to use.
func MarkupCloud(aList *hashtags.THashList) []template.HTML {
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
			url = "/hl/" + item.Tag[1:]
		} else {
			url = "/ml/" + item.Tag[1:]
		}
		tl[idx] = template.HTML(` <a href="` + url + `" class="` + class + `" title=" ` + fmt.Sprintf("%d * %s", item.Count, item.Tag[1:]) + ` ">` + item.Tag + `</a> `) // #nosec G203
	}

	return tl
} // MarkupCloud()

// MarkupTags returns `aPage` with all #hashtags/@mentions marked up
// as a HREF links.
//
//	`aPage` The HTML page to process.
func MarkupTags(aPage []byte) []byte {
	var repl, search string
	// (0) Check whether there are any links present:
	linkMatches := htAHrefRE.FindAll(aPage, -1)
	if (nil != linkMatches) || (0 < len(linkMatches)) {
		// (1) replace the links with a dummy text:
		for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
			search = regexp.QuoteMeta(string(linkMatches[cnt]))
			if re, err := regexp.Compile(search); nil == err {
				repl = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
				aPage = re.ReplaceAllLiteral(aPage, []byte(repl))
			}
		}
	}

	// (2) markup the #hashtags/@mentions:
	result := htHashMentionRE.ReplaceAllStringFunc(string(aPage),
		func(aString string) string {
			sub := htHashMentionRE.FindSubmatch([]byte(aString))
			if (nil == sub) || (0 >= len(sub)) || (0 >= len(sub[1])) {
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
					if '"' == sub[2][0] {
						// double quote following a possible hashtag: most
						// probably an URL#fragment, hence leave it as is
						return aString
					}
					if (';' == sub[2][0]) && htEntityRE.MatchString(aString) {
						// leave HTML entities as is
						return aString
					}
				}
				url = "/hl/" + hash[1:]
			} else {
				url = "/ml/" + hash[1:]
			}
			if 0 < len(sub[2]) {
				suffix += string(sub[2])
			}

			return `<a href="` + url + `" class="smaller">` + hash + `</a>` + suffix
		})

	// (3) replace the link dummies with the real markup:
	for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
		search = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
		if re, err := regexp.Compile(search); nil == err {
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
	initHashlist(aList.Clear())
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

	go walkAllPosts(aList, doReplaceTag)
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
