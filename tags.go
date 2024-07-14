/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

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
	"regexp"
	"runtime"
	"strings"

	ht "github.com/mwat56/hashtags"
)

var (
	// RegEx to find PREformatted parts in an HTML page.
	htAHrefRE = regexp.MustCompile(`(?si)(<a[^>]*>.*?</a>)`)

	// RegEx to identify a numeric HTML entity.
	htEntityRE = regexp.MustCompile(`(#[0-9]+;)`)

	// match: #hashtag|@mention
	htHashMentionRE = regexp.MustCompile(`(?i)([@#][\p{L}’'\d_§-]+)(.?|$)`)
	//                                         1111111111111111111  2222
	// NOTE: compare with `github.com/mwat56/hashtags/`
)

// --------------------------------------------------------------------------

// AddTagID checks a newly added `aPosting` for #hashtags and @mentions.
//
// Parameters:
//
//	`aList`: The hashlist to use (update).
//	`aPosting`: The new posting to handle.
func AddTagID(aList *ht.THashTags, aPosting *TPosting) {
	go aList.IDparse(aPosting.ID(), aPosting.Markdown())

	runtime.Gosched() // get the background operation started
} // AddTagID()

// InitHashlist initialises the hash list.
//
// Parameters:
//
//	`aList`: The list of #hashtags/@mentions to update.
func InitHashlist(aList *ht.THashTags) {
	wf := func(aID uint64) error {
		post := NewPosting(aID, "")
		if err := post.Load(); nil != err {
			// we ignore the error here ...
			return nil
		}

		if 0 < post.Len() {
			aList.IDparse(aID, post.Markdown())
		}
		return nil
	} // wf()

	go poPersistence.Walk(wf)
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
func MarkupCloud(aList *ht.THashTags) []template.HTML {
	var (
		class string // re-use variable
		idx   int
		item  ht.TCountItem
	)
	list := aList.List()
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
			htListLookup[ht.MarkHash == item.Tag[0]] + item.Tag[1:] +
			`" class="` + class + `" title=" ` +
			fmt.Sprintf("%d * %s", item.Count, item.Tag[1:]) +
			` ">` + item.Tag + `</a>`) // #nosec G203
	}

	return tl
} // MarkupCloud()

// MarkupTags returns `aPage` with all #hashtags/@mentions marked up
// as a HREF links.
//
// Parameters:
//
//	`aPage`: The HTML page to process.
func MarkupTags(aPage []byte) []byte {
	var ( // re-use variables
		cnt, hits    int
		err          error
		re           *regexp.Regexp
		repl, search string
	)
	// (0) Check whether there are any links present:
	linkMatches := htAHrefRE.FindAll(aPage, -1)
	if (nil != linkMatches) || (0 < len(linkMatches)) {
		// (1) replace the links with a dummy text:
		for hits, cnt = len(linkMatches), 0; cnt < hits; cnt++ {
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

			if ht.MarkHash == hash[0] {
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
	for hits, cnt = len(linkMatches), 0; cnt < hits; cnt++ {
		search = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
		if re, err = regexp.Compile(search); nil == err {
			result = re.ReplaceAllLiteralString(result,
				string(linkMatches[cnt]))
		}
	}

	return []byte(result)
} // MarkupTags()

// `ReadHashlist()` reads all postings to (re-)build the list of
// #hashtags/@mentions disregarding any pre-existing list.
//
// Parameters:
//
//	`aList`: The list of #hashtags/@mentions to build.
func ReadHashlist(aList *ht.THashTags) {
	InitHashlist(aList.Clear())
} // ReadHashlist()

// RemoveIDTags removes `aID` from `aList's` items.
//
// Parameters:
//
//	`aList`: The hashlist to update.
//	`aID`: The ID of the posting to remove.
func RemoveIDTags(aList *ht.THashTags, aID uint64) {
	go aList.IDremove(aID)

	runtime.Gosched() // get the background operation started
} // RemoveIDTags()

// RenameIDTags renames all references of `aOldID` to `aNewID`.
//
// Parameters:
//
//	`aList`: The hashlist to update.
//	`aOldID`: The posting's old ID.
//	`aNewID`: The posting's new ID.
func RenameIDTags(aList *ht.THashTags, aOldID, aNewID uint64) {
	go aList.IDrename(aOldID, aNewID)

	runtime.Gosched() // get the background operation started
} // RenameIDTags()

// ReplaceTag replaces the #tags/@mentions in `aList`.
//
// Parameters:
//
//	`aList`: The hashlist to update.
//	`aSearchTag`: The old #tag/@mention to find.
//	`aReplaceTag`: The new #tag/@mention to use.
func ReplaceTag(aList *ht.THashTags, aSearchTag, aReplaceTag string) {
	if (nil == aList) || (0 == len(aSearchTag)) || (0 == len(aReplaceTag)) {
		return
	}

	switch aSearchTag[0] {
	case ht.MarkHash, ht.MarkMention:
		switch aReplaceTag[0] {
		case ht.MarkHash, ht.MarkMention:
			// nothing to do
		default:
			return
		}
	default:
		return
	}

	searchRE, err := regexp.Compile(`(?i)\` + aSearchTag)
	if nil != err {
		return //se.Wrap(err, 2)
	}

	wf := func(aID uint64) error {
		post := NewPosting(aID, "")
		if err := post.Load(); nil != err {
			// no contents, no joy ...
			return nil
		}

		if 0 == post.Len() {
			// no contents, no joy ...
			return nil
		}

		if !searchRE.Match(post.Markdown()) {
			return nil
		}

		nMarkdown := searchRE.ReplaceAllLiteral(
			post.Markdown(),
			[]byte(aReplaceTag))

		post.Set(nMarkdown).Store()

		aList.IDupdate(aID, nMarkdown)

		return nil
	} // wf()

	poPersistence.Walk(wf)
	// runtime.Gosched() // get the background operation started
} // ReplaceTag()

// UpdateTags updates the #hashtag/@mention references of `aPosting`.
//
// Parameters:
//
//	`aList`: The hashlist to update.
//	`aPosting`: The new posting to process.
func UpdateTags(aList *ht.THashTags, aPosting *TPosting) {
	go aList.IDupdate(aPosting.ID(), aPosting.Markdown())

	runtime.Gosched() // get the background operation started
} // UpdateTags()

/* _EoF_ */
