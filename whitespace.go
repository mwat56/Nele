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
		aList[idx] = bytes.TrimSpace(hit)
		// Since a list entry here can never be completely empty we
		// don't have to check for `nil` after `bytes.TrimSpace()`.
	}

	return aList
} // trimPREmatches()

// Internal list of regular expressions used by
// the `RemoveWhiteSpace()` function.
type (
	tReItem struct {
		replace string
		regEx   *regexp.Regexp
	}
)

var (
	// RegEx to find PREformatted parts in an HTML page.
	wsPreRE = regexp.MustCompile(`(?si)\s*<pre[^>]*>.*?</pre>\s*`)

	wsREs = []tReItem{
		// comments:
		{``, regexp.MustCompile(`(?s)<!--.*?-->`)},
		// HTML and HEAD elements:
		{`$1`, regexp.MustCompile(`(?si)\s*(</?(body|\!DOCTYPE|head|html|link|meta|script|style|title)[^>]*>)\s*`)},
		// block elements:
		{`$1`, regexp.MustCompile(`(?si)\s+(</?(article|blockquote|div|footer|h[1-6]|header|nav|p|section)[^>]*>)`)},
		{`$1`, regexp.MustCompile(`(?si)(</?(article|blockquote|div|footer|h[1-6]|header|nav|p|section)[^>]*>)\s+`)},
		// lists:
		{`$1`, regexp.MustCompile(`(?si)\s+(</?([dou]l|li|d[dt])[^>]*>)`)},
		{`$1`, regexp.MustCompile(`(?si)(</?([dou]l|li|d[dt])[^>]*>)\s+`)},
		// table elements:
		{`$1`, regexp.MustCompile(`(?si)\s+(</?(col|t(able|body|foot|head|[dhr]))[^>]*>)`)},
		{`$1`, regexp.MustCompile(`(?si)(</?(col|t(able|body|foot|head|[dhr]))[^>]*>)\s+`)},
		// form elements:
		{`$1`, regexp.MustCompile(`(?si)\s+(</?(form|fieldset|legend|opt(group|ion))[^>]*>)`)},
		{`$1`, regexp.MustCompile(`(?si)(</?(form|fieldset|legend|opt(group|ion))[^>]*>)\s+`)},
		// BR / HR:
		{`$1`, regexp.MustCompile(`(?i)\s*(<[bh]r[^>]*>)\s*`)},
		// whitespace after opened anchor:
		{`$1`, regexp.MustCompile(`(?si)(<a\s+[^>]*>)\s+`)},
		// preserve empty table cells:
		{`$1&#160;$3`, regexp.MustCompile(`(?i)(<td(\s+[^>]*)?>)\s+(</td>)`)},
		// remove empty paragraphs:
		{``, regexp.MustCompile(`(?i)<(p)(\s+[^>]*)?>\s*</$1>`)},
		// whitespace before closing GT:
		{`>`, regexp.MustCompile(`\s+>`)},
	}
)

// RemoveWhiteSpace returns `aPage` with all HTML comments and
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

	// (2) traverse through all the whitespace REs:
	for _, re := range wsREs {
		aPage = re.regEx.ReplaceAll(aPage, []byte(re.replace))
	}

	// (3) replace the PRE dummies with the real markup:
	for lLen, cnt := len(preMatches), 0; cnt < lLen; cnt++ {
		search = fmt.Sprintf(`\s*</-%d-%d-%d-%d-/>\s*`, cnt, cnt, cnt, cnt)
		if re, err := regexp.Compile(search); nil == err {
			aPage = re.ReplaceAllLiteral(aPage, preMatches[cnt])
		}
	}

	return aPage
} // RemoveWhiteSpace()

/* _EoF_ */
