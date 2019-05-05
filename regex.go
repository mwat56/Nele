/*
    Copyright © 2019  M.Watermann, 10247 Berlin, Germany
                All rights reserved
            EMail : <support@mwat.de>
*/

package blog

/*
 * This files provides a few RegEx based functions.
 */

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/russross/blackfriday.v2"
)

var (
	// RegEx to match hh:mm:ss
	hmsRE = regexp.MustCompile("^(([01]?[0-9])|(2[0-3]))[^0-9](([0-5]?[0-9])[^0-9]([0-5]?[0-9])?)?$")
)

// getHMS() splits up `aTime` into `rHour`, `rMinute`, and `rSecond`.
func getHMS(aTime string) (rHour, rMinute, rSecond int) {
	matches := hmsRE.FindStringSubmatch(aTime)
	if 0 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all Atoi() errors.
		rHour, _ = strconv.Atoi(matches[1])
		if 0 < len(matches[5]) {
			rMinute, _ = strconv.Atoi(matches[5])
			if 0 < len(matches[6]) {
				rSecond, _ = strconv.Atoi(matches[6])
			}
		}
	}

	return
} // getHMS()

var (
	// RegEx to match YYYY(MM)(DD)
	// Invalid values for month or day result in a `0,0,0` result.
	ymdRE = regexp.MustCompile("^([0-9]{4})[^0-9]?(((0?[0-9])|(1[0-2]))[^0-9]?((0?[0-9])?|([12][0-9])?|(3[01])?)?)?$")
)

// getYMD() splits up `aDate` into `rYear`, `rMonth`, and `rDay`.
func getYMD(aDate string) (rYear int, rMonth time.Month, rDay int) {
	matches := ymdRE.FindStringSubmatch(aDate)
	if 0 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all Atoi() errors.
		rYear, _ = strconv.Atoi(matches[1])
		if 0 < len(matches[3]) {
			m, _ := strconv.Atoi(matches[3])
			rMonth = time.Month(m)
			if 0 < len(matches[6]) {
				rDay, _ = strconv.Atoi(matches[6])
			}
		}
	}

	return
} // getYMD()

func init() {
	initWSre()
} // init()

// Initialise the `whitespaceREs` list.
func initWSre() int {
	result := 0
	for idx, re := range whitespaceREs {
		whitespaceREs[idx].regEx = regexp.MustCompile(re.search)
		result++
	}
	result++

	return result
} // initWSre()

var (
	// RegEx to correct wrong markup created by 'blackfriday';
	// see MDtoHTML()
	bfPreCodeRE = regexp.MustCompile(`(?s)\s*(<pre>)<code>(.*?)\s*</code>(</pre>)\s*`)
)

// MDtoHTML converts the `aMarkdown` data returning HTML data.
func MDtoHTML(aMarkdown []byte) []byte {
	result := blackfriday.Run(aMarkdown)
	if i := bytes.Index(result, []byte("</pre>")); 0 > i {
		// no need for RegEx execution
		return result
	}
	// Testing for PRE makes this implementation twice as fast
	// if there's no PRE in the generated HTML and about the
	// same speed if there actually is a PRE part.

	return bfPreCodeRE.ReplaceAll(result, []byte("$1\n$2\n$3"))
} // MDtoHTML()

/*
// MDtoHTv1 converts the `aMarkdown` data returning HTML data.
// Note: This is just an implementation for benchmark purposes.
func MDtoHTv1(aMarkdown []byte) []byte {
	return bfPreCodeRE.ReplaceAll(
		blackfriday.Run(aMarkdown),
		[]byte("$1\n$2\n$3"))
} // MDtoHTv1()

// MDtoHTv0 converts the `aMarkdown` data returning HTML data.
// Note: This is just an implementation for benchmark purposes.
func MDtoHTv0(aMarkdown []byte) []byte {
	// This actually the fastest implementation but it
	// leaves all the wrong PRE/CODE markup in the result.
	return blackfriday.Run(aMarkdown)
} // MDtoHTv0()
*/

// `trimPREmatches()` removes leading/trailing whitespace from list entries.
func trimPREmatches(aList [][]byte) [][]byte {
	for idx, hit := range aList {
		aList[idx] = bytes.TrimSpace(hit)
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
	tReList []tReItem
)

var (
	// RegEx to find PREformatted parts in an HTML page.
	preRE = regexp.MustCompile(`(?si)\s*<pre[^>]*>.*?</pre>\s*`)

	// List of regular expressions matching different sets of HTML whitespace.
	whitespaceREs = tReList{
		// comments
		{`(?s)<!--.*?-->`, ``, nil},
		// HTML and HEAD elements
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

// RemoveWhiteSpace removes HTML comments and unneccessary whitespace.
//
// This function removes all unneeded/redundant whitespace
// and HTML comments from the given <tt>aPage</tt>.
// This can reduce significantly the amount of data to send to
// the remote user agent thus saving bandwidth.
func RemoveWhiteSpace(aPage []byte) []byte {
	var repl, search string

	// fmt.Println("Page0:", string(aPage))
	// (0) Check whether there are PREformatted parts:
	preMatches := preRE.FindAll(aPage, -1)
	if (nil == preMatches) || (0 >= len(preMatches)) {
		// no PRE hence only the other REs to perform
		for _, reEntry := range whitespaceREs {
			aPage = reEntry.regEx.ReplaceAll(aPage, []byte(reEntry.replace))
		}
		return aPage
	}
	preMatches = trimPREmatches(preMatches)

	// Make sure PREformatted parts remain as-is.
	// (1) replace the PRE parts with a dummy text:
	for l, cnt := len(preMatches), 0; cnt < l; cnt++ {
		search := fmt.Sprintf("\\s*%s\\s*", regexp.QuoteMeta(string(preMatches[cnt])))
		if re, err := regexp.Compile(search); nil == err {
			repl = fmt.Sprintf("</@%d@%d@%d@%d@/>", cnt, cnt, cnt, cnt)
			aPage = re.ReplaceAllLiteral(aPage, []byte(repl))
		}
	}
	// fmt.Println("Page1:", string(aPage))

	// (2) traverse through all the whitespace REs:
	for _, re := range whitespaceREs {
		aPage = re.regEx.ReplaceAll(aPage, []byte(re.replace))
	}
	// fmt.Println("Page2:", string(aPage))

	// (3) replace the PRE dummies with the real markup:
	for l, cnt := len(preMatches), 0; cnt < l; cnt++ {
		search = fmt.Sprintf("\\s*</@%d@%d@%d@%d@/>\\s*", cnt, cnt, cnt, cnt)
		if re, err := regexp.Compile(search); nil == err {
			aPage = re.ReplaceAllLiteral(aPage, preMatches[cnt])
		}
	}
	// fmt.Println("Page3:", string(aPage))

	return aPage
} // RemoveWhiteSpace()

var (
	// RegEx to replace CR/LF by LF
	crlfRE = regexp.MustCompile("\r\n")
)

// `replCRLF()` replaces all CR/LF pairs by a single LF.
func replCRLF(aText []byte) []byte {
	return crlfRE.ReplaceAllLiteral(aText, []byte("\n"))
} // replCRLF()

// SearchPostings traverses the sub-directories of `aBaseDir` looking
// for `aText` in all posting files.
//
// `aBaseDir` is the directory of which all subdirectories are scanned.
//
// `aText` is the text to look for in the postings.
//
// The returned `TPostList` can be empty because (a) `aText` could not be
// compiled into a regular expression, (b) no files to search were found,
// or (c) no files matched `aText`.
func SearchPostings(aBaseDir, aText string) *TPostList {
	bd, _ := filepath.Abs(aBaseDir)
	pl := NewPostList(bd)

	// search := fmt.Sprintf("(?si)%s", regexp.QuoteMeta(string(aText)))
	// pattern, err := regexp.Compile(fmt.Sprintf("(?si)%s", search))
	pattern, err := regexp.Compile(fmt.Sprintf("(?s)%s", aText))
	if err != nil {
		return pl // empty list
	}

	files, err := filepath.Glob(bd + "/*/*.md")
	if nil != err {
		return pl // empty list
	}

	for _, fName := range files {
		fTxt, err := ioutil.ReadFile(fName)
		if (nil != err) || (!pattern.Match(fTxt)) {
			// We 'eat' possible errors here, indirectly assuming
			// them to be a no-match.
			continue
		}
		id := path.Base(fName)
		p := newPosting(bd, id[:len(id)-3]) // exclude file extension
		pl.Add(p.Set(fTxt))
	}

	return pl
} // SearchPostings()

/*
// SearchRubric traverses the sub-directories of `aBaseDir` looking
// for `aRubric` in all posting files.
//
// `aBaseDir` is the directory of which all subdirectories are scanned.
//
// `aRubric` is the rubric to look for in the postings.
//
// The returned `TPostList` can be empty because (a) `aText` could not be
// compiled into a regular expression, (b) no files to search were found,
// or (c) no files matched `aText`.
func SearchRubric(aBaseDir, aRubric string) *TPostList {
	// the required markup would be "* _`text`_" on a single line
	rubric := fmt.Sprintf("^\\*\\s+_`%s`_\\s*$", aRubric)

	return SearchPostings(aBaseDir, rubric)
} // SearchRubric()
*/

// var (
// RegEx to find the leading part of an URL;
// see `URLpath0()`
// 	routeRE0 = regexp.MustCompile("^/?[a-z0-9]*/?")
// )
/*
// urlPath0 returns the leading part of an `aURL`.
func urlPath0(aURL string) string {
	match := routeRE0.FindString(aURL)
	l := len(match)
	if ('/' == match[0]) && (1 < l) {
		match = match[1:]
	}

	return match
} // urlPath0()
*/

var (
	// RegEx to find path and possible added path components
	routeRE = regexp.MustCompile("^/?([\\w\\._-]+)?/?([\\w\\.\\?\\=_-]*)?")
)

// URLparts returns two parts: `rDir` holds the base-directory of `aURL`,
// `rPath` holds the remaining part of `aURL`.
//
// Depending on the actual value of `aURL` both return values might be empty
// or both may be filled. None of both will hold a leading slash.
func URLparts(aURL string) (rDir, rPath string) {
	matches := routeRE.FindStringSubmatch(aURL)
	if 2 < len(matches) {
		return matches[1], matches[2]
	}

	return aURL, ""
} // URLparts()

// var (
// 	routeRE2 = regexp.MustCompile("^/([\\w\\._-]+)?(/[\\w\\._-]*)?(\\?[=\\w\\._-]+)?")
// )
//
// // URLpath2 is a testing proc
// func URLpath2(aURL string) (rHead, rTail, rQuery string) {
// 	matches := routeRE2.FindStringSubmatch(aURL)
// 	rHead = matches[1]
// 	rTail = matches[2]
// 	rQuery = matches[3]
//
// 	return
// } // URLpath()

// ShiftPath splits off the first component of p, which will be
// cleaned of relative components before processing.
//
// `rHead` will never contain a slash and `rTail` will always
// be a rooted path without trailing slash.
//
// see https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func ShiftPath(aPath string) (rHead, rTail string) {
	aPath = path.Clean("/" + aPath)
	i := strings.Index(aPath[1:], "/") + 1
	if 0 > i {
		return aPath[1:], "/"
	}
	if 0 == i {
		return "", aPath[i:]

	}

	return aPath[1:i], aPath[i:]
} // ShiftPath()

/* _EoF_ */
