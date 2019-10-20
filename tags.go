/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
* This files provides functions related to #hashtags/@mentions
 */

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"regexp"
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

// `doCheckPost()` returns whether there is a file identified
// by `aID` containing `aHash`.
//
// The function's result is `false` (1) if the file associated with
// `aID` doesnt't exist, or (2) if the file can't be read, or (3)
// the given `aHash` can't be found in the posting's text.
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
	// go goCacheCleanup()
} // goCheckHashes()

// `goInitHashlist()` initialises the hash list.
func goInitHashlist(aList *hashtags.THashList) {
	if _, err := aList.Load(); (nil == err) && (0 < aList.Len()) {
		go goCheckHashes(aList)
		return // assume everything is up-to-date
	}

	dirnames, err := filepath.Glob(PostingBaseDirectory() + "/*")
	if nil != err {
		return // we can't recover from this :-(
	}
	for _, mdName := range dirnames {
		filesnames, err := filepath.Glob(mdName + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filesnames) {
			continue // skip empty directory
		}
		for _, postName := range filesnames {
			id := strings.TrimPrefix(postName, mdName+"/")
			if txt, err := ioutil.ReadFile(postName); /* #nosec G304 */ nil == err {
				aList.IDparse(id[:len(id)-3], txt) // strip name extension
			}
		}
	}

	_, _ = aList.Store()
	// go goCacheCleanup()
} // goInitHashlist()

// `goRemoveID()` removes `aID` from `aList's` items.
func goRemoveID(aList *hashtags.THashList, aID string) {
	aList.IDremove(aID)
	// go goCacheCleanup()
} // goRemoveID()

// `goRenameID()` renames all references of `aOldID` to `aNewID`.
func goRenameID(aList *hashtags.THashList, aOldID, aNewID string) {
	aList.IDrename(aOldID, aNewID)
	// go goCacheCleanup()
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
			url = "/hl/" + item.Tag[1:]
		} else {
			url = "/ml/" + item.Tag[1:]
		}
		tl[idx] = template.HTML(` <a href="` + url + `" class="` + class + `" title=" ` + fmt.Sprintf("%d * %s", item.Count, item.Tag[1:]) + ` ">` + item.Tag + `</a> `) // #nosec G203
	}

	return tl
} // markupCloud()

var (
	// RegEx to find PREformatted parts in an HTML page.
	htAHrefRE = regexp.MustCompile(`(?si)(<a[^>]*>.*?</a>)`)

	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`(#[0-9]+;)`)

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(`(?i)([@#][§\wÄÖÜß-]+)(.?|$)`)
)

// `markupTags()` returns `aPage` with all #hashtags/@mentions marked
// up as a HREF links.
func markupTags(aPage []byte) []byte {
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
} // markupTags()

/* _EoF_ */
