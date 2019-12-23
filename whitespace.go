/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This files provides a function to remove redundant whitespace
 * and comments from a HTML page.
 */

import (
	"bytes"
	"fmt"
	"regexp"
)

// `trimPREmatches()` removes leading/trailing whitespace from list entries.
func trimPREmatches(aList [][]byte) [][]byte {
	for idx, hit := range aList {
		if aList[idx] = bytes.TrimSpace(hit); nil == aList[idx] {
			aList[idx] = []byte(``)
		}
	}

	return aList
} // trimPREmatches()

// Internal list of regular expressions used by
// the `RemoveWhiteSpace()` function.
type (
	tReItem struct {
		search  string
		replace string
		regEx   *regexp.Regexp
	}
)

var (
	// RegEx to find PREformatted parts in an HTML page.
	wsPreRE = regexp.MustCompile(`(?si)\s*<pre[^>]*>.*?</pre>\s*`)

	// List of regular expressions matching different sets of HTML whitespace.
	wsREs = []tReItem{
		// comments
		{`(?s)<!--.*?-->`, ``, nil},
		// HTML and HEAD elements:
		{`(?i)\s*(</?(body|\!DOCTYPE|head|html|link|meta|script|style|title)[^>]*>)\s*`, `$1`, nil},
		// block elements:
		{`(?i)\s*(</?(article|blockquote|div|footer|h[1-6]|header|nav|p|section)[^>]*>)\s*`, `$1`, nil},
		// lists:
		{`(?i)\s*(</?([dou]l|li|d[dt])[^>]*>)\s*`, `$1`, nil},
		// table:
		{`(?i)\s*(</?(col|t(able|body|foot|head|[dhr]))[^>]*>)\s*`, `$1`, nil},
		// forms:
		{`(?i)\s*(</?(form|fieldset|legend|opt(group|ion))[^>]*>)\s*`, `$1`, nil},
		// BR / HR:
		{`(?i)\s*(<[bh]r[^>]*>)\s*`, `$1`, nil},
		// whitespace after opened anchor:
		{`(?i)(<a\s+[^>]*>)\s+`, `$1`, nil},
		// preserve empty table cells:
		{`(?i)(<td(\s+[^>]*)?>)\s+(</td>)`, `$1&#160;$3`, nil},
		// remove empty paragraphs:
		{`(?i)<(p)(\s+[^>]*)?>\s*</$1>`, ``, nil},
		// whitespace before closing GT:
		{`\s+>`, `>`, nil},
	}
)

var (
	// Initialise the `whitespaceREs` list.
	_ = func() int {
		result := 0
		for idx, re := range wsREs {
			wsREs[idx].regEx = regexp.MustCompile(re.search)
			result++
		}
		result++

		return result
	}()
)

// RemoveWhiteSpace returns `aPage` with al.l HTML comments and
// unnecessary whitespace removed.
//
// This function removes all unneeded/redundant whitespace
// and HTML comments from the given <tt>aPage</tt>.
// This can reduce significantly the amount of data to send to
// the remote user agent thus saving bandwidth and transfer time.
//
//	`aPage` The web page's HTML markup to process.
func RemoveWhiteSpace(aPage []byte) []byte {
	var repl, search string

	// fmt.Println("Page0:", string(aPage))
	// (0) Check whether there are PREformatted parts:
	preMatches := wsPreRE.FindAll(aPage, -1)
	if (nil == preMatches) || (0 >= len(preMatches)) {
		// no PRE hence only the other REs to perform
		for _, reEntry := range wsREs {
			aPage = reEntry.regEx.ReplaceAll(aPage, []byte(reEntry.replace))
		}
		return aPage
	}
	preMatches = trimPREmatches(preMatches)

	// Make sure PREformatted parts remain as-is.
	// (1) replace the PRE parts with a dummy text:
	for lLen, cnt := len(preMatches), 0; cnt < lLen; cnt++ {
		search = fmt.Sprintf(`\s*%s\s*`, regexp.QuoteMeta(string(preMatches[cnt])))
		if re, err := regexp.Compile(search); nil == err {
			repl = fmt.Sprintf(`</-%d-%d-%d-%d-/>`, cnt, cnt, cnt, cnt)
			aPage = re.ReplaceAllLiteral(aPage, []byte(repl))
		}
	}
	// fmt.Println("Page1:", string(aPage))

	// (2) traverse through all the whitespace REs:
	for _, re := range wsREs {
		aPage = re.regEx.ReplaceAll(aPage, []byte(re.replace))
	}
	// fmt.Println("Page2:", string(aPage))

	// (3) replace the PRE dummies with the real markup:
	for lLen, cnt := len(preMatches), 0; cnt < lLen; cnt++ {
		search = fmt.Sprintf(`\s*</-%d-%d-%d-%d-/>\s*`, cnt, cnt, cnt, cnt)
		if re, err := regexp.Compile(search); nil == err {
			aPage = re.ReplaceAllLiteral(aPage, preMatches[cnt])
		}
	}
	// fmt.Println("Page3:", string(aPage))

	return aPage
} // RemoveWhiteSpace()

/* _EoF_ */
