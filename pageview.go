/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/pageview"
)

var (
	// R/O RegEx to extract link-text and link-URL from markup.
	pvLinksRE = regexp.MustCompile(`(?s)\[([^\)]+)\]\s*\(([^\]]+)\)`)
	//                                    1              2
	// 1 : link text
	// 2 : link URL

	// RegEx to check whether an URL starts with a scheme.
	pvSchemeRE = regexp.MustCompile(`^\w+:`)
)

// goSetPostingLinkViews sets the external links in `aPosting`
// to include a page preview image (if available).
//
//	`aPosting` is the path/filename of the posting to process.
//	`aImageDirectory` is the URL directory for page preview images.
func goSetPostingLinkViews(aPostName, aImageDirectory string) {
	if 0 == len(aPostName) {
		return
	}
	fName := filepath.Base(aPostName)
	id := fName[:len(fName)-3] // strip name extension
	p := NewPosting(id)
	txt := p.Markdown()
	if 0 == len(txt) {
		return
	}

	linkMatches := pvLinksRE.FindAll(txt, -1)
	if (nil == linkMatches) || (0 == len(linkMatches)) {

		//TODO check for image URL constructs and make sure
		// the referenced file does in fact exist.
		return
	}
	for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
		linkParts := pvLinksRE.FindSubmatch(linkMatches[cnt])
		linkTxt, linkURL := string(linkParts[1]), string(linkParts[2])
		if !pvSchemeRE.MatchString(linkURL) {
			continue // skip local links
		}

		imgName, err := pageview.CreateImage(linkURL)
		if (nil != err) || (0 == len(imgName)) {
			apachelogger.Err("goSetPostingLinkViews()",
				fmt.Sprintf("pageview.CreateImage(%s): %v", linkURL, err))
			// log.Printf("pageview.CreateImage(%s): %v", linkURL, err) //TODO REMOVE
			continue
		}
		imgName = filepath.Join(aImageDirectory, imgName)

		// replace `[link-text](link-url)` by
		// `[![link-text](image-URL)](link-url)`
		search := regexp.QuoteMeta(string(linkParts[0]))
		if re, err := regexp.Compile(search); nil == err {
			replace := "[![" + linkTxt + "](/" + imgName + ")](" + linkURL + ")"
			txt = re.ReplaceAllLiteral(txt, []byte(replace))
			_, _ = p.Set(txt).Store()
		}
	}
} // goSetPostingLinkViews()

// `goUpdateAllLinkViews()` prepares the external links in all postings
// to use a page preview image (if available).
//
//	`aPostingBaseDir` is the base directory used for
// storing articles/postings.
//	`aImageDirectory` is the URL directory for page preview images.
func goUpdateAllLinkViews(aPostingBaseDir, aImageDirectory string) {
	dirnames, err := filepath.Glob(aPostingBaseDir + "/*")
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
			go goSetPostingLinkViews(postName, aImageDirectory)
		}
	}
} // goUpdateAllLinkViews()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// InitPageImages starts the process to update the preview images
// in all postings.
//
//	`aPostingBaseDir` is the base directory used for
// storing articles/postings.
//	`aImageDirectory` is the URL directory for page preview images.
func InitPageImages(aPostingBaseDir, aImageDirectory string) {
	go goUpdateAllLinkViews(aPostingBaseDir, aImageDirectory)
} // InitPageImages()

// SetPostingLinkViews updates the external links in `aPosting`
// to include a page preview image (if available).
//
//	`aPosting` is the path/filename of the posting to process.
//	`aImageDirectory` is the URL directory for page preview images.
func SetPostingLinkViews(aPostName, aImageDirectory string) {
	goSetPostingLinkViews(aPostName, aImageDirectory)
} // SetPostingLinkViews()

/* _EoF_ */
