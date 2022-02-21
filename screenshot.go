/*
   Copyright © 2022 M.Watermann, 10247 Berlin, Germany
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
	"github.com/mwat56/screenshot"
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

// `checkScreenshots()` checks whether the image file referenced
// in the text of `aPosting` actually exists locally.
//
//	`aPosting` The posting the text of which is searched for
// a local image link.
func checkScreenshots(aPosting *TPosting) {
	if nil == aPosting {
		return
	}

	var pair tImgURL // re-use variable
	list := checkScreenshotURLs(aPosting.Markdown())
	for _, pair = range list {
		goCreateScreenshot(pair.pageURL)
	}
} // checkScreenshots()

var (
	// R/O RegEx to find an URL for the screenshot image's.
	ssImageRE = regexp.MustCompile(
		`(?s)\[\s*\!\[[^\]]*?\]\s*\(\s*([^\)]+?)\s*\)\s*\]\s*\(\s*([^\)]+?)\s*\)`)
	//                                 111111111                  222222222
	// `[![alt-text](image-URL)](link-url)`
	// 1 : local image URL
	// 2 : remote page URL
)

// `checkScreenshotURLs()` tests whether `aTxt` contains an external
// link with an embedded page screenshot image,
// returning a list of the two link URLs.
//
//	`aTxt` is the text to search.
func checkScreenshotURLs(aTxt []byte) (rList tImgURLlist) {
	matches := ssImageRE.FindAllSubmatch(aTxt, -1)
	for idx, l := 0, len(matches); idx < l; idx++ {
		rList = append(rList, tImgURL{
			filepath.Base(string(matches[idx][1])),
			string(matches[idx][2]),
		})
	}
	return
} // checkScreenshotURLs()

// `goCreateScreenshot()` generates a screenshot of `aURL` in background.
//
//	`aURL` The URL for which to create a screenshot image.
func goCreateScreenshot(aURL string) {
	// `screenshot.CreateImage()` checks whether the file exists
	// and whether it's too old.
	imgName, err := screenshot.CreateImage(aURL)
	if (nil != err) || (0 == len(imgName)) {
		apachelogger.Err("goCreateScreenshot()",
			fmt.Sprintf("screenshot.CreateImage(%s): '%v'", aURL, err))
	}
} // goCreateScreenshot()

// `goSetLinkScreenshots()` changes the external links in the text
// of `aPosting` to include a page screenshot image (if available).
//
//	`aPosting` The posting the text of which is going to be processed.
//	`aImageURLdir` The URL directory for page screenshot images.
func goSetLinkScreenshots(aPosting *TPosting, aImageURLdir string) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	url := `` // re-use in loop below
	linkMatches := ssLinkRE.FindAllSubmatch(aPosting.Markdown(), -1)
	if (nil != linkMatches) && (0 < len(linkMatches)) {
		for l, cnt := len(linkMatches), 0; cnt < l; cnt++ {
			url = string(linkMatches[cnt][3])
			if !ssSchemeRE.MatchString(url) {
				continue // skip local links
			}
			preparePost(aPosting,
				&tLink{
					link:     string(linkMatches[cnt][1]),
					linkText: string(linkMatches[cnt][2]),
					linkURL:  url,
				},
				aImageURLdir)
		}
	}

	// In case we didn't find any normal links we check
	// the links with an embedded page screenshot image:
	checkScreenshots(aPosting)
} // goSetLinkScreenshots()

// `goUpdateAllLinkScreenshots()` prepares the external links in
// all postings to use a page screenshot image (if available).
//
//	`aPostingBaseDir` is the base directory used for storing
// articles/postings.
//	`aImageURLdir` is the local URL directory for page screenshot images.
func goUpdateAllLinkScreenshots(aPostingBaseDir, aImageURLdir string) {
	var ( // re-use variables in loops below
		err                 error
		dName, fName, pName string
		dNames, fNames      []string
		p                   *TPosting
	)
	if dNames, err = filepath.Glob(aPostingBaseDir + "/*"); nil != err {
		return // we can't recover from this :-(
	}

	for _, dName = range dNames {
		if fNames, err = filepath.Glob(dName + "/*.md"); (nil != err) || (0 == len(fNames)) {
			continue // no files found
		}

		for _, pName = range fNames {
			fName = filepath.Base(pName)
			p = NewPosting(fName[:len(fName)-3]) // strip name extension
			if err = p.Load(); nil == err {
				goSetLinkScreenshots(p, aImageURLdir)
			} // ELSE ignore the error here …
		}
	}
} // goUpdateAllLinkScreenshots()

// `preparePost()` creates a screenshot image and updates `aLink`
// in the text of `aPosting` to embed a link into the image.
//
//	`aPosting` The posting the text of which is going to be updated.
//	`aLink` The link parts to use.
//	`aImageURLdir` The URL directory for page screenshot images.
func preparePost(aPosting *TPosting, aLink *tLink, aImageURLdir string) {
	imgName, err := screenshot.CreateImage(aLink.linkURL)
	if (nil != err) || (0 == len(imgName)) {
		apachelogger.Err("preparePost()",
			fmt.Sprintf("screenshot.CreateImage(%s): '%v'", aLink.linkURL, err))
		return
	}
	if 0 == len(aImageURLdir) {
		return
	}

	if '/' != aImageURLdir[0] {
		aImageURLdir = `/` + aImageURLdir
	}
	if txt := prepPostText(aPosting.Markdown(), aLink, imgName, aImageURLdir); 0 < len(txt) {
		_, _ = aPosting.Set(txt).Store()
	}
} // preparePost()

// `prepPostText()` gets called by `preparePost()` to replace
// `[link-text](link-url)` by `[![alt-text](image-URL)](link-URL)`
//
//	`aPosting` The posting the text of which is going to be updated.
//	`aLink` The link parts to use.
//	`aImageURLdir` The URL directory for page screenshot images.
func prepPostText(aPosting []byte, aLink *tLink, aImageName, aImageURLdir string) (rText []byte) {
	search := regexp.QuoteMeta(aLink.link)
	if re, err := regexp.Compile(search); nil == err {
		replace := "[![" + aLink.linkText +
			"](" + filepath.Join(aImageURLdir, aImageName) +
			")](" + aLink.linkURL + ")"
		rText = re.ReplaceAllLiteral(aPosting, []byte(replace))
	}

	return
} // prepPostText()

// RemovePageScreenshots deletes the images used in `aPosting`.
//
//	`aPosting` The posting the image(s) of which are going to be deleted.
func RemovePageScreenshots(aPosting *TPosting) {
	if (nil == aPosting) || (0 == aPosting.Len()) {
		return
	}

	var ( // re-use variables
		err   error
		fi    os.FileInfo
		fName string
		list  tImgURLlist
		pair  tImgURL
	)

	if list = checkScreenshotURLs(aPosting.Markdown()); 0 == len(list) {
		return
	}

	for _, pair = range list {
		fName = screenshot.PathFile(pair.pageURL)
		if fi, err = os.Stat(fName); (nil != err) || fi.IsDir() {
			continue
		}
		_ = os.Remove(fName)
	}
} // RemovePageScreenshots()

var (
	// R/O RegEx to extract link-text and link-URL from markup.
	// Checking for the not-existence of the leading `!` should exclude
	// embedded image links.
	ssLinkRE = regexp.MustCompile(
		`(?m)(?:^\s*\>[\t ]*)((?:[^\!\n\>][\t ]*)?\[([^\[]+?)\]\s*\(([^\]]+?)\))`)
	//                                            11222222222111111133333333311
	// `[link-text](link-url)`
	// 0 : complete RegEx match
	// 1 : markdown link markup
	// 2 : link text
	// 3 : remote page URL

	// R/O simple RegEx to check whether an URL starts with a scheme.
	ssSchemeRE = regexp.MustCompile(`^\w+://`)
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// CreateScreenshot generates a screenshot of `aURL` in background.
//
//	`aURL` The URL for which to create a screenshot image.
func CreateScreenshot(aURL string) {
	go goCreateScreenshot(aURL)
	runtime.Gosched() // get the background operation started
} // CreateScreenshot

// PrepareLinkScreenshots updates the external link(s) in `aPosting`
// to include page screenshot image(s) (if available).
//
//	`aPosting` The posting the text of which is going to be processed.
//	`aImageURLdir`The URL directory for page screenshot images.
func PrepareLinkScreenshots(aPosting *TPosting, aImageURLdir string) {
	go goSetLinkScreenshots(aPosting, aImageURLdir)
	runtime.Gosched() // get the background operation started
} // PrepareLinkScreenshots()

// UpdateScreenshots starts the process to update the screenshot images
// in all postings.
//
//	`aPostingBaseDir` The base directory used for storing
// articles/postings.
//	`aImageURLdir` The URL directory for page screenshot images.
func UpdateScreenshots(aPostingBaseDir, aImgURLdir string) {
	go goUpdateAllLinkScreenshots(aPostingBaseDir, aImgURLdir)
	runtime.Gosched() // get the background operation started
} // UpdateLinkScreenshots()

/* _EoF_ */
