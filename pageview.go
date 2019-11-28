/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/pageview"
)

var (
	// R/O RegEx to find an URL for the preview image's.
	pvImageRE = regexp.MustCompile(
		`(?s)\[\!\[[^\]]*\]\(([^\)]+)\)\]\(([^\)]+)\)`)
	//                      11111         22222
	//       `[![alt-text](image-URL)](link-url)`
	// 1 : local image URL
	// 2 : remote page URL

	// R/O RegEx to extract link-text and link-URL from markup.
	pvLinksRE = regexp.MustCompile(`(?s)\[([^\)]+)\]\s*\(([^\]]+)\)`)
	//                                    11111          22222222
	// 1 : link text
	// 2 : remote link URL

	// R/O RegEx to check whether an URL starts with a scheme.
	pvSchemeRE = regexp.MustCompile(`^\w+:`)
)

type (
	// tImgURL represents a pair of image name and page URL.
	tImgURL struct {
		imgURL  string
		pageURL string
	}
	tImgURLlist []tImgURL
)

// checkForImgURL tests whether `aTxt` contains an external link with an
// embedded page preview image, returning a list of the two link URLs.
//
//	`aTxt` is the text to search.
func checkForImgURL(aTxt []byte) (rList tImgURLlist) {
	matches := pvImageRE.FindAllSubmatch(aTxt, -1)
	for idx := 0; idx < len(matches); idx++ {
		pair := tImgURL{
			filepath.Base(string(matches[idx][1])),
			string(matches[idx][2]),
		}
		rList = append(rList, pair)
	}
	return
} // checkForImgURL()

// `goCheckPageImages()` checks whether the image file referenced
// in the text of `aPosting` exists in `aImageDir`.
//
//	`aPosting` is the posting the text of which is searched for
// a local image link.
//	`aImageURLdir` is the local URL directory for page preview images.
//	`aImageDir` is the directory used to store the generated images.
func goCheckPageImages(aPosting *TPosting, aImageURLdir, aImageDir string) {
	if nil == aPosting {
		return
	}
	list := checkForImgURL(aPosting.Markdown())
	if 0 == len(list) {
		return
	}

	for _, pair := range list {
		// `pageview.CreateImage()` checks whether the file exists
		// and whether it's too old.
		imgName, err := pageview.CreateImage(pair.pageURL)
		if (nil != err) || (0 == len(imgName)) {
			apachelogger.Err("goCheckPageImage()",
				fmt.Sprintf("pageview.CreateImage(%s): %v", pair.pageURL, err))
			log.Printf("pageview.CreateImage(%s): %v",
				pair.pageURL, err) //TODO REMOVE
			continue
		}
	}
} // goCheckPageImages()

// `goPreparePost()` creates an preview image and updates `aLinkMatch`
// in the text of `aPosting` to embedd a link to the image.
//
//	`aPosting` is the posting the text of which is going to be updated.
//	`aLinkMatch` is the matchd link markdown text.
//	`aImageURLdir` is the URL directory for page preview images.
func goPreparePost(aPosting *TPosting, aLinkMatch []byte, aImageURLdir string) {
	linkParts := pvLinksRE.FindSubmatch(aLinkMatch)
	linkTxt, linkURL := string(linkParts[1]), string(linkParts[2])
	if !pvSchemeRE.MatchString(linkURL) {
		return // skip local links
	}

	imgName, err := pageview.CreateImage(linkURL)
	if (nil != err) || (0 == len(imgName)) {
		apachelogger.Err("goPreparePost()",
			fmt.Sprintf("pageview.CreateImage(%s): %v", linkURL, err))
		log.Printf("pageview.CreateImage(%s): %v", linkURL, err) //TODO REMOVE
		return
	}

	// replace `[link-text](link-url)` by
	// `[![alt-text](image-URL)](link-url)`
	search := regexp.QuoteMeta(string(linkParts[0]))
	if re, err := regexp.Compile(search); nil == err {
		replace := "[![" +
			linkTxt +
			"](/" +
			filepath.Join(aImageURLdir, imgName) +
			")](" +
			linkURL +
			")"
		txt := aPosting.Markdown()
		txt = re.ReplaceAllLiteral(txt, []byte(replace))
		_, _ = aPosting.Set(txt).Store()
	}
} // goPreparePost()

// `goSetPostingLinkViews()` sets the external links in the text of
// `aPosting` to include a page preview image (if available).
//
//	`aPosting` is the posting the text of which is going to be processed.
//	`aImageURLdir` is the URL directory for page preview images.
//	`aImageDir` is the directory used to store the generated images.
func goSetPostingLinkViews(aPosting *TPosting, aImageURLdir, aImageDir string) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	linkMatches := pvLinksRE.FindAll(aPosting.Markdown(), -1)
	if (nil == linkMatches) || (0 == len(linkMatches)) {
		// Here we didn't find any normal links so check
		// the links with an embedded page preview image.
		go goCheckPageImages(aPosting, aImageURLdir, aImageDir)
		return
	}

	for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
		go goPreparePost(aPosting, linkMatches[cnt], aImageURLdir)
	}
} // goSetPostingLinkViews()

// `goUpdateAllLinkPreviews()` prepares the external links in all postings
// to use a page preview image (if available).
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the local URL directory for page preview images.
//	`aImageDir` is the directory used to store the generated images.
func goUpdateAllLinkPreviews(aPostingBaseDir, aImageURLdir, aImageDir string) {
	dirnames, err := filepath.Glob(aPostingBaseDir + "/*")
	if nil != err {
		return // we can't recover from this :-(
	}

	for _, mdName := range dirnames {
		filenames, err := filepath.Glob(mdName + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filenames) {
			continue // skip empty directory
		}
		for _, postName := range filenames {
			fName := filepath.Base(postName)
			go goSetPostingLinkViews(
				NewPosting(fName[:len(fName)-3]), // strip name extension
				aImageURLdir, aImageDir)
		}
	}
} // goUpdateAllLinkPreviews()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// InitPageImages starts the process to update the preview images
// in all postings.
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the URL directory for page preview images.
func InitPageImages(aPostingBaseDir, aImageURLdir string) {
	go goUpdateAllLinkPreviews(aPostingBaseDir, aImageURLdir, pageview.ImageDirectory())
} // InitPageImages()

// SetPostingLinkViews updates the external links in `aPosting`
// to include a page preview image (if available).
//
//	`aPosting` is the posting the text of which is going to be processed.
//	`aImageURLdir` is the URL directory for page preview images.
func SetPostingLinkViews(aPosting *TPosting, aImageURLdir string) {
	goSetPostingLinkViews(aPosting, aImageURLdir, pageview.ImageDirectory())
} // SetPostingLinkViews()

/* _EoF_ */
