/*
   Copyright © 2020 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides a function to convert NarkDown to HTML.
 */

import (
	"bytes"
	"regexp"
	"sync"

	bf "github.com/russross/blackfriday/v2"
	// bf "gopkg.in/russross/blackfriday.v2"
)

var (
	// Instead of creating this objects with every call to `MDtoHTML()`
	// we use some prepared global instances.
	bfExtensions = bf.WithExtensions(
		bf.Autolink |
			bf.BackslashLineBreak |
			bf.DefinitionLists |
			bf.FencedCode |
			bf.Footnotes |
			bf.HeadingIDs |
			bf.NoIntraEmphasis |
			bf.SpaceHeadings |
			bf.Strikethrough |
			bf.Tables)

	// WithRenderer allows overriding the default `blackfriday` renderer.
	bfRenderer = bf.WithRenderer(
		bf.NewHTMLRenderer(bf.HTMLRendererParameters{
			Flags: bf.FootnoteReturnLinks |
				bf.Smartypants |
				bf.SmartypantsFractions |
				bf.SmartypantsDashes |
				bf.SmartypantsLatexDashes,
		}),
	)

	// Guard against race conditions in `MDtoHTML()` when calling
	// `blackfriday.Run()` concurrently.
	bfMtx = new(sync.Mutex)

	// Text to recognise a PREformatted section.
	bfPre = []byte("</pre>")

	// Text to recognise a PREformatted section with redundant CODE block.
	bfPreCode = []byte("<pre><code ")

	// RegEx to correct redundant markup created by 'bf';
	// see `mdToHTML()`
	bfPreCodeRE1 = regexp.MustCompile(`(?s)\s*(<pre>)<code>(.*?)\s*</code>(</pre>)\s*`)

	bfPreCodeRE2 = regexp.MustCompile(`(?s)\s*(<pre)><code (class="language-\w+")>(.*?)\s*</code>(</pre>)\s*`)

	// RegEx to correct back markup since Blackfriday v2.1.0';
	// see `mdToHTML()`
	bfSupRE = regexp.MustCompile(`<span aria-label='Return'>.*</span>`)
)

// MDtoHTML converts the `aMarkdown` data and returns HTML data.
//
//	`aMarkdown` The raw Markdown text to convert.
func MDtoHTML(aMarkdown []byte) (rHTML []byte) {
	var i int // re-use variable
	bfMtx.Lock()
	defer bfMtx.Unlock()

	rHTML = bytes.TrimSpace(bf.Run(aMarkdown, bfRenderer, bfExtensions))
	rHTML = bfSupRE.ReplaceAll(rHTML, []byte("<sup>[return]</sup>"))

	// Testing for PRE first makes this implementation twice as fast
	// if there's no PRE in the generated HTML and about the same
	// speed if there actually is a PRE part.
	if i = bytes.Index(rHTML, bfPre); 0 > i {
		return // no need for further RegEx execution
	}

	rHTML = bfPreCodeRE1.ReplaceAll(rHTML, []byte("$1\n$2\n$3"))

	if i = bytes.Index(rHTML, bfPreCode); 0 > i {
		return // no need for the second RegEx execution
	}

	return bfPreCodeRE2.ReplaceAll(rHTML, []byte("$1 $2>\n$3\n$4"))
} // MDtoHTML()

/* _EoF_ */
