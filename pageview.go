/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/pageview"
)

var (
	// R/O RegEx to find an URL for the preview image's.
	pvImageRE = regexp.MustCompile(
		`(?s)\[\s*\!\[[^\]]*?\]\s*\(\s*([^\)]+?)\s*\)\s*\]\s*\(\s*([^\)]+?)\s*\)`)
	//                                 111111111                  222222222
	// `[![alt-text](image-URL)](link-url)`
	// 1 : local image URL
	// 2 : remote page URL

	// R/O RegEx to extract link-text and link-URL from markup.
	// Checking for the not-existence of the leading `!` should exclude
	// embedded image links.
	pvLinkRE = regexp.MustCompile(`(?s)([^\[]]\s*[^\!]\s*)?\[([^\[\)]+?)\]\s*\(([^\]]+?)\)`)
	//                                 1111111111111111111   2222222222        33333333333
	// `[link-text](link-url)`
	// 1 : lead in (ignored)
	// 2 : link text
	// 3 : remote page URL

	// R/O simple RegEx to check whether an URL starts with a scheme.
	pvSchemeRE = regexp.MustCompile(`^\w+://`)
)

type (
	// `tImgURL` represents a pair of image name and page URL.
	tImgURL struct {
		imgURL  string
		pageURL string
	}
	tImgURLlist []tImgURL

	// `tLink` represents a group of link, link text and link URL.
	tLink struct {
		link     string // the whole markdown link
		linkText string
		linkURL  string
	}
)

// `checkForImageURL()` tests whether `aTxt` contains an external link with
// an embedded page preview image, returning a list of the two link URLs.
//
//	`aTxt` is the text to search.
func checkForImageURL(aTxt []byte) (rList tImgURLlist) {
	matches := pvImageRE.FindAllSubmatch(aTxt, -1)
	for idx := 0; idx < len(matches); idx++ {
		pair := tImgURL{
			filepath.Base(string(matches[idx][1])),
			string(matches[idx][2]),
		}
		rList = append(rList, pair)
	}
	return
} // checkForImageURL()

// `checkPageImages()` checks whether the image file referenced
// in the text of `aPosting` exists in `aImageDir`.
//
//	`aPosting` The posting the text of which is searched for
// a local image link.
func checkPageImages(aPosting *TPosting) {
	if nil == aPosting {
		return
	}

	list := checkForImageURL(aPosting.Markdown())
	for _, pair := range list {
		goCreatePreview(pair.pageURL)
	}
} // checkPageImages()

// `goCreatePreview()` generates a preview of `aURL` in background.
//
//	`aURL` The URL for which to create a preview image.
func goCreatePreview(aURL string) {
	// `pageview.CreateImage()` checks whether the file exists
	// and whether it's too old.
	imgName, err := pageview.CreateImage(aURL)
	if (nil != err) || (0 == len(imgName)) {
		apachelogger.Err("goCreatePreview()",
			fmt.Sprintf("pageview.CreateImage(%s): '%v'", aURL, err))
	}
} // goCreatePreview()

// `goUpdateAllLinkPreviews()` prepares the external links in all postings
// to use a page preview image (if available).
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the local URL directory for page preview images.
func goUpdateAllLinkPreviews(aPostingBaseDir, aImageURLdir string) {
	dirnames, err := filepath.Glob(aPostingBaseDir + "/*")
	if nil != err {
		return // we can't recover from this :-(
	}

	for _, mdName := range dirnames {
		filenames, err := filepath.Glob(mdName + "/*.md")
		if nil != err {
			continue // it might be a file (not a directory) …
		}

		if 0 < len(filenames) {
			for _, postName := range filenames {
				fName := filepath.Base(postName)
				p := NewPosting(fName[:len(fName)-3]) // strip name extension
				if err := p.Load(); nil == err {
					setPostingLinkViews(p, aImageURLdir)
				}
			}
		}
	}
} // goUpdateAllLinkPreviews()

// `preparePost()` creates a preview image and updates `aLink`
// in the text of `aPosting` to embed a link into the image.
//
//	`aPosting` The posting the text of which is going to be updated.
//	`aLink` The link parts to use.
//	`aImageURLdir` The URL directory for page preview images.
func preparePost(aPosting *TPosting, aLink *tLink, aImageURLdir string) {
	imgName, err := pageview.CreateImage(aLink.linkURL)
	if (nil != err) || (0 == len(imgName)) {
		apachelogger.Err("preparePost()",
			fmt.Sprintf("pageview.CreateImage(%s): '%v'", aLink.linkURL, err))
		return
	}
	if 0 == len(aImageURLdir) {
		return
	}
	if '/' != aImageURLdir[0] {
		aImageURLdir = `/` + aImageURLdir
	}

	// replace `[link-text](link-url)` by
	// `[![alt-text](image-URL)](link-URL)`
	search := regexp.QuoteMeta(aLink.link)
	if re, err := regexp.Compile(search); nil == err {
		replace := "[![" + aLink.linkText +
			"](" + filepath.Join(aImageURLdir, imgName) +
			")](" + aLink.linkURL + ")"
		txt := re.ReplaceAllLiteral(aPosting.Markdown(), []byte(replace))
		_, _ = aPosting.Set(txt).Store()
	}
} // preparePost()

// RemovePagePreviews deletes the images used in `aPosting`.
//
//	`aPosting` The posting the image(s) of which are going to be deleted.
func RemovePagePreviews(aPosting *TPosting) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	list := checkForImageURL(aPosting.Markdown())
	for _, pair := range list {
		fName := pageview.PathFile(pair.pageURL)
		if fi, err := os.Stat(fName); (nil != err) || fi.IsDir() {
			continue
		}
		_ = os.Remove(fName)
	}
} // RemovePagePreviews()

// `setPostingLinkViews()` changes the external links in the text
// of `aPosting` to include a page preview image (if available).
//
//	`aPosting` The posting the text of which is going to be processed.
//	`aImageURLdir` The URL directory for page preview images.
func setPostingLinkViews(aPosting *TPosting, aImageURLdir string) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	linkMatches := pvLinkRE.FindAllSubmatch(aPosting.Markdown(), -1)
	if (nil != linkMatches) && (0 < len(linkMatches)) {
		for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
			url := string(linkMatches[cnt][3])
			if !pvSchemeRE.MatchString(url) {
				continue // skip local links
			}
			preparePost(aPosting,
				&tLink{
					link:     string(linkMatches[cnt][0]),
					linkText: string(linkMatches[cnt][2]),
					linkURL:  url,
				},
				aImageURLdir)
		}
	}
	// In case we didn't find any normal links we check
	// the links with an embedded page preview image:
	checkPageImages(aPosting)
} // setPostingLinkViews()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// PrepareLinkPreviews updates the external link(s) in `aPosting`
// to include page preview image(s) (if available).
//
//	`aPosting` is the posting the text of which is going to be processed.
//	`aImageURLdir` is the URL directory for page preview images.
func PrepareLinkPreviews(aPosting *TPosting, aImageURLdir string) {
	setPostingLinkViews(aPosting, aImageURLdir)
} // PrepareLinkPreviews()

// UpdatePreviews starts the process to update the preview images
// in all postings.
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the URL directory for page preview images.
func UpdatePreviews(aPostingBaseDir, aImgURLdir string) {
	go goUpdateAllLinkPreviews(aPostingBaseDir, aImgURLdir)
	runtime.Gosched() // get the background operation started
} // UpdateLinkPreviews()

/* _EoF_ */
