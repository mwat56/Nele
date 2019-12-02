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
	// Checking for the not-existence od e leading `!` should exclude
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

// `checkForImgageURL()` tests whether `aTxt` contains an external link with
// an embedded page preview image, returning a list of the two link URLs.
//
//	`aTxt` is the text to search.
func checkForImgageURL(aTxt []byte) (rList tImgURLlist) {
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

// `checkPageImages()` checks whether the image file referenced
// in the text of `aPosting` exists in `aImageDir`.
//
//	`aPosting` is the posting the text of which is searched for
// a local image link.
//	`aImageURLdir` is the local URL directory for page preview images.
//	`aImageDir` is the directory used to store the generated images.
func checkPageImages(aPosting *TPosting, aImageURLdir, aImageDir string) {
	if nil == aPosting {
		return
	}

	list := checkForImgageURL(aPosting.Markdown())
	for _, pair := range list {
		// `pageview.CreateImage()` checks whether the file exists
		// and whether it's too old.
		imgName, err := pageview.CreateImage(pair.pageURL)
		if (nil != err) || (0 == len(imgName)) {
			apachelogger.Err("checkPageImages()",
				fmt.Sprintf("pageview.CreateImage(%s): '%v'", pair.pageURL, err))
			continue
		}
	}
} // checkPageImages()

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

		if 0 < len(filenames) {
			for _, postName := range filenames {
				fName := filepath.Base(postName)
				p := NewPosting(fName[:len(fName)-3]) // strip name extension
				if err := p.Load(); nil == err {
					setPostingLinkViews(p, aImageURLdir, aImageDir)
				}
			}
		}
	}
} // goUpdateAllLinkPreviews()

// `preparePost()` creates an preview image and updates `aLinkMatch`
// in the text of `aPosting` to embedd a link to the image.
//
//	`aPosting` is the posting the text of which is going to be updated.
//	`aLink` contains the link parts to use.
//	`aImageURLdir` is the URL directory for page preview images.
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

// RemoveImages deletes the images used in `aPosting`.
//
//	`aPosting` The posting the image(s) of which are going to be deleted.
func RemoveImages(aPosting *TPosting) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	list := checkForImgageURL(aPosting.Markdown())
	for _, pair := range list {
		fName := pageview.PathFile(pair.pageURL)
		if fi, err := os.Stat(fName); (nil != err) || fi.IsDir() {
			continue
		}
		_ = os.Remove(fName)
	}
} // RemoveImages()

// `setPostingLinkViews()` changes the external links in the text
// of `aPosting` to include a page preview image (if available).
//
//	`aPosting` is the posting the text of which is going to be processed.
//	`aImageURLdir` is the URL directory for page preview images.
//	`aImageDir` is the directory used to store the generated images.
func setPostingLinkViews(aPosting *TPosting, aImageURLdir, aImageDir string) {
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
	checkPageImages(aPosting, aImageURLdir, aImageDir)
} // setPostingLinkViews()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// UpdateLinkPreviews starts the process to update the preview images
// in all postings.
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the URL directory for page preview images.
func UpdateLinkPreviews(aPostingBaseDir, aImageURLdir string) {
	go goUpdateAllLinkPreviews(aPostingBaseDir, aImageURLdir, pageview.ImageDirectory())
	runtime.Gosched()
} // UpdateLinkPreviews()

// PrepareLinkPreviews updates the external link(s) in `aPosting`
// to include page preview image(s) (if available).
//
//	`aPosting` is the posting the text of which is going to be processed.
//	`aImageURLdir` is the URL directory for page preview images.
func PrepareLinkPreviews(aPosting *TPosting, aImageURLdir string) {
	setPostingLinkViews(aPosting, aImageURLdir, pageview.ImageDirectory())
} // PrepareLinkPreviews()

/* _EoF_ */
