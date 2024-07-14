/*
Copyright © 2019, 2023 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions
//lint:file-ignore ST1005 - I prefer capitalisation

/*
 * This file provides functions and methods to handle HTTP requests.
 */

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/cssfs"
	ht "github.com/mwat56/hashtags"
	"github.com/mwat56/jffs"
	"github.com/mwat56/passlist"
	"github.com/mwat56/uploadhandler"
)

type (
	// TPageHandler provides the handling of HTTP request/response.
	TPageHandler struct {
		cssFS    http.Handler                  // CSS file server
		hashList *ht.THashTags                 // #hashtags/@mentions list
		imgUp    *uploadhandler.TUploadHandler // `img` upload handler
		staticFS http.Handler                  // `static` file server
		staticUp *uploadhandler.TUploadHandler // `static` upload handler
		userList *passlist.TPassList           // user/password list
		viewList *TViewList                    // list of template/views
	}
)

// --------------------------------------------------------------------------
// constructor function:

// `NewPageHandler()` returns a new `TPageHandler` instance.
//
// The returned object implements the `errorhandler.TErrorPager`,
// `http.Handler`, `and `passlist.TAuthDecider` interfaces.
func NewPageHandler() (*TPageHandler, error) {
	var (
		err error
		msg string
	)
	result := new(TPageHandler)

	if result.viewList, err = newViewList(
		filepath.Join(AppArgs.DataDir, "views")); nil != err {
		msg = fmt.Sprintf("Error: views problem: %v", err)
		log.Println(`NewPageHandler()`, msg)

		// Without our views we can't generate any web-page.
		return nil, err
	}

	result.cssFS = cssfs.FileServer(AppArgs.DataDir + `/`)

	if 0 < len(AppArgs.HashFile) {
		// hashtags.UseBinaryStorage = false //TODO REMOVE
		if result.hashList, err = ht.New(AppArgs.HashFile, true); nil != err {
			result.hashList = nil
		} else {
			InitHashlist(result.hashList) // background operation
		}
	}
	if nil == result.hashList {
		if nil == err {
			err = errors.New(`Error: missing hashFile`)
		}
		msg = fmt.Sprintf("%v", err)
		log.Println(`NewPageHandler()`, msg)
		return nil, err
	}

	result.staticFS = jffs.FileServer(AppArgs.DataDir + `/`)

	if AppArgs.Screenshot {
		UpdateScreenshots(PostingBaseDirectory()) // background operation
	}

	if 0 == len(AppArgs.UserFile) {
		log.Println("NewPageHandler(): missing password file\nAUTHENTICATION DISABLED!")
	} else if result.userList, err = passlist.LoadPasswords(AppArgs.UserFile); nil != err {
		log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
		result.userList = nil
	}

	return result, nil
} // NewPageHandler()

// --------------------------------------------------------------------------
// helper functions:

var (
	// RegEx to match hh:mm:ss
	phHmsRE = regexp.MustCompile(`^(([01]?[0-9])|(2[0-3]))[^0-9](([0-5]?[0-9])([^0-9]([0-5]?[0-9]))?)?[^0-9]?|$`)
)

// `getHMS()` splits up `aTime` into `rHour`, `rMinute`, and `rSecond`.
func getHMS(aTime string) (rHour, rMinute, rSecond int) {
	matches := phHmsRE.FindStringSubmatch(aTime)
	if 1 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all Atoi() errors.
		rHour, _ = strconv.Atoi(matches[1])
		if 0 < len(matches[5]) {
			rMinute, _ = strconv.Atoi(matches[5])
			if 0 < len(matches[7]) {
				rSecond, _ = strconv.Atoi(matches[7])
			}
		}
	}

	return
} // getHMS()

var (
	// RegEx to match YYYY(MM)(DD).
	// Invalid values for month or day result in a `0` result.
	// This is just a pattern test, it doesn't check whether
	// the date is valid.
	phYmdRE = regexp.MustCompile(`^([0-9]{4})([^0-9]?(0[1-9]|1[012])([^0-9]?(0[1-9]|[12][0-9]|3[01])?)?)?[^0-9]?`)
)

// `getYMD()` splits up `aDate` into `rYear`, `rMonth`, and `rDay`.
//
// This is just a pattern test: the function doesn't check whether
// the date as such is a valid date.
func getYMD(aDate string) (rYear int, rMonth time.Month, rDay int) {
	matches := phYmdRE.FindStringSubmatch(aDate)
	if 1 < len(matches) {
		// The RegEx only matches digits so we can
		// safely ignore all `Atoi()` errors.
		rYear, _ = strconv.Atoi(matches[1])
		if 0 < len(matches[3]) {
			m, _ := strconv.Atoi(matches[3])
			rMonth = time.Month(m)
			if 0 < len(matches[5]) {
				rDay, _ = strconv.Atoi(matches[5])
			} else {
				rDay = 1
			}
		} else {
			rDay = 1
		}
	}

	return
} // getYMD()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// `newViewList()` returns a list of views found in `aDirectory`
// and a possible I/O error.
func newViewList(aDirectory string) (*TViewList, error) {
	var ( // re-use variables
		err   error
		files []string
		fName string
		v     *TView
	)
	result := NewViewList()

	if files, err = filepath.Glob(aDirectory + "/*.gohtml"); err != nil {
		return nil, err
	}

	for _, fName = range files {
		fName = filepath.Base(fName[:len(fName)-7]) // remove extension
		if v, err = NewView(aDirectory, fName); nil != err {
			return nil, err
		}
		result = result.Add(v)
	}

	return result, nil
} // newViewList()

var (
	// RegEx to extract number and start of articles shown
	phNumStartRE = regexp.MustCompile(`^(\d*)(\D*(\d*)?)?`)
)

// `numStart()` extracts two numbers from `aString`.
func numStart(aString string) (rNum, rStart int) {
	matches := phNumStartRE.FindStringSubmatch(aString)
	if 3 < len(matches) {
		if 0 < len(matches[1]) {
			rNum, _ = strconv.Atoi(matches[1])
		}
		if 0 < len(matches[3]) {
			rStart, _ = strconv.Atoi(matches[3])
		}
	}

	return
} // numStart()

var (
	// RegEx to replace CR/LF by LF
	phCrLfRE = regexp.MustCompile("\r\n")
)

// `replCRLF()` replaces all CR/LF pairs by a single LF.
func replCRLF(aText []byte) []byte {
	return phCrLfRE.ReplaceAllLiteral(aText, []byte("\n"))
} // replCRLF()

var (
	// RegEx to find path and possible added path components
	phURLpartsRE = regexp.MustCompile(
		`(?i)^/*([\p{L}\d_.-]+)?/*([\p{L}\d_§.?!=:;/,@# ’'-]*)?`)
	//           1111111111111     22222222222222222222222222
)

// `URLparts()` returns two parts: `rDir` holds the base-directory of `aURL`,
// `rPath` holds the remaining part of `aURL`.
//
// Depending on the actual value of `aURL` both return values may be
// empty or both may be filled; none of both will hold a leading slash.
func URLparts(aURL string) (rDir, rPath string) {
	if result, err := url.QueryUnescape(aURL); nil == err {
		aURL = result
	}

	matches := phURLpartsRE.FindStringSubmatch(aURL)
	if 2 < len(matches) {
		return matches[1], strings.TrimSpace(matches[2])
	}

	return aURL, ""
} // URLparts()

// --------------------------------------------------------------------------
// TPageHandler methods

// `basicPageData()` returns a list of data to be inserted into the
// `view`/templates.
func (ph *TPageHandler) basicPageData(aRequest *http.Request) *TemplateData {
	lang, theme := AppArgs.Lang, AppArgs.Theme
	if nil != aRequest {
		var val string // re-use variable
		if val = strings.ToLower(aRequest.FormValue(`lang`)); 0 < len(val) {
			switch val {
			case `de`, `en`:
				lang = val
			}
		}
		if val = strings.ToLower(aRequest.FormValue(`theme`)); 0 < len(val) {
			switch val {
			case `dark`, `light`:
				theme = val
			}
		}
	}

	y, m, d := time.Now().Date()
	now := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	pageData := NewTemplateData().
		Set("Blogname", AppArgs.BlogName).
		Set(`CSS`, template.HTML(`<link rel="stylesheet" type="text/css" title="mwat's styles" href="/css/stylesheet.css"><link rel="stylesheet" type="text/css" href="/css/`+theme+`.css"><link rel="stylesheet" type="text/css" href="/css/fonts.css">`)).
		Set("HashCount", ph.hashList.HashCount()).
		Set(`Lang`, lang).
		Set("MentionCount", ph.hashList.MentionCount()).
		Set("monthURL", "/m/"+now).
		Set("NOW", now).
		Set("PostingCount", poPersistence.PostingCount()).
		Set("Robots", "index,follow").
		Set("Taglist", MarkupCloud(ph.hashList)).
		Set("Title", AppArgs.Realm+": "+now).
		Set("weekURL", "/w/"+now) // #nosec G203

	return pageData
} // basicPageData()

// GetErrorPage returns an error page for `aStatus`,
// implementing the `TErrorPager` interface.
func (ph *TPageHandler) GetErrorPage(aData []byte, aStatus int) []byte {
	var (
		empty  []byte
		err    error
		result []byte
	)

	pageData := ph.basicPageData(nil).Set(`Robots`, `noindex,follow`)

	switch aStatus {
	case 404:
		if result, err = ph.viewList.RenderedPage("404", pageData); nil == err {
			return result
		}

	//TODO implement other status codes

	default:
		pageData = pageData.Set("Error", template.HTML(aData)) // #nosec G203
		if result, err = ph.viewList.RenderedPage("error", pageData); nil == err {
			return result
		}
	}

	return empty
} // GetErrorPage()

// `handleGET()` processes the HTTP GET requests.
func (ph *TPageHandler) handleGET(aWriter http.ResponseWriter, aRequest *http.Request) {
	pageData := ph.basicPageData(aRequest)
	path, tail := URLparts(aRequest.URL.Path)
	switch strings.ToLower(path) { // handle URLs case-insensitive

	case "a", "ap": // add a new post
		ph.handleReply(`ap`, aWriter,
			pageData.Set(`Robots`, `noindex,nofollow`))

	case "certs": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "css":
		ph.cssFS.ServeHTTP(aWriter, aRequest)

	case "d", "dp": // change date
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		p := NewPosting(0, tail)
		if !p.Exists() {
			http.NotFound(aWriter, aRequest)
			return
		}
		t := p.Time()
		date := p.Date()
		pageData = pageData.Set(`HMS`, fmt.Sprintf("%02d:%02d:%02d",
			t.Hour(), t.Minute(), t.Second())).
			Set("ID", p.IDstr()).
			Set("Manuscript", template.HTML(p.Markdown())).
			Set("monthURL", "/m/"+date).
			Set("Robots", "noindex,nofollow").
			Set("weekURL", "/w/"+date).
			Set("YMD", date) // #nosec G203
		ph.handleReply("dp", aWriter, pageData)

	case "e", "ep": // edit a single posting
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		p := NewPosting(0, tail)
		if !p.Exists() {
			http.NotFound(aWriter, aRequest)
			return
		}
		t := p.Time()
		date := p.Date()
		pageData = pageData.Set(`HMS`, fmt.Sprintf("%02d:%02d:%02d",
			t.Hour(), t.Minute(), t.Second())).
			Set("ID", p.IDstr()).
			Set("Manuscript", template.HTML(p.Markdown())).
			Set("monthURL", "/m/"+date).
			Set("Robots", "noindex,nofollow").
			Set("weekURL", "/w/"+date).
			Set("YMD", date) // #nosec G203
		ph.handleReply("ep", aWriter, pageData)

	case "faq", "faq.html":
		ph.handleReply(`faq`, aWriter, pageData)

	case "favicon.ico":
		http.Redirect(aWriter, aRequest,
			`/img/`+path,
			http.StatusMovedPermanently)

	case "fonts":
		ph.staticFS.ServeHTTP(aWriter, aRequest)

	case "hl": // #hashtag list
		if 0 < len(tail) {
			ph.handleTagMentions(ph.hashList.HashList(string(ht.MarkHash)+tail),
				pageData, aWriter)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case `i`, `il`: // (re-)init the hashList
		if nil != ph.hashList {
			ph.handleReply(`il`, aWriter,
				pageData.Set(`Robots`, `noindex,nofollow`))
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "img":
		ph.staticFS.ServeHTTP(aWriter, aRequest)

	case "imprint", "impressum":
		ph.handleReply(`imprint`, aWriter, pageData)

	case `index`, `index.html`, `index.php`, `index.shtml`:
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

		/*
			case "js":
				ph.sfh.ServeHTTP(aWriter, aRequest)
		*/

	case "licence", "license", "lizenz":
		ph.handleReply(`licence`, aWriter, pageData)

	case `m`, `mw`: // handle a given month
		var (
			d, y   int
			m      time.Month
			robots string = `noindex,follow`
		)
		if 0 == len(tail) {
			y, m, d = time.Now().Date()
		} else {
			y, m, d = getYMD(tail)
			// Allow indexing for lists older 30 days
			t := time.Date(y, m, d+30, 0, 0, 0, 0, time.Local)
			if t.Before(time.Now()) {
				robots = `index,follow`
			}
		}
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		pl := NewPostList().Month(y, m)
		ph.handleReply(`searchresult`, aWriter,
			pageData.Set(`Matches`, pl.Len()).
				Set(`monthURL`, "/m/"+date).
				Set(`Postings`, pl.Sort()).
				Set(`Robots`, robots).
				Set(`weekURL`, "/w/"+date))

	case "ml": // @mention list
		if 0 < len(tail) {
			ph.handleTagMentions(ph.hashList.MentionList("@"+tail),
				pageData, aWriter)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "n", "np": // handle newest postings
		ph.handleRoot(tail, pageData, aWriter /*, aRequest*/)

	case "p", "pp": // handle a single posting
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		p := NewPosting(0, tail)
		if err := p.Load(); nil != err {
			apachelogger.Err("TPageHandler.handleGET()",
				fmt.Sprintf("TPosting.Load('%s'): %v", p.IDstr(), err))
			http.NotFound(aWriter, aRequest)
			return
		}
		date := p.Date()
		err := ph.userList.IsAuthenticated(aRequest)
		aWriter.Header().Set(`Cache-Control`, `private, max-age=864000`) // 10 days
		aWriter.Header().Set(`Last-Modified`, p.LastModified())

		pageData = pageData.Set(`isAuth`, nil == err).
			Set(`monthURL`, `/m/`+date).
			Set("Posting", p).
			Set("weekURL", "/w/"+date)
		ph.handleReply("article", aWriter, pageData)

	case "postings": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "privacy", "datenschutz":
		ph.handleReply(`privacy`, aWriter, pageData)

	case `pv`, `v`: // update the Screenshot images
		if AppArgs.Screenshot {
			ph.handleReply(`pv`, aWriter,
				pageData.Set(`Robots`, `noindex,nofollow`))
		} else {
			http.Redirect(aWriter, aRequest, `/n/`,
				http.StatusMovedPermanently)
		}

	case "q":
		http.Redirect(aWriter, aRequest, "/s/"+tail, http.StatusMovedPermanently)

	case "r", "rp": // posting's removal
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, `/n/`,
				http.StatusSeeOther)
			return
		}
		p := NewPosting(0, tail)
		if !p.Exists() {
			http.NotFound(aWriter, aRequest)
			return
		}
		t := p.Time()
		date := p.Date()
		pageData = pageData.Set(`HMS`,
			fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())).
			Set(`ID`, p.IDstr()).
			Set("Manuscript", template.HTML(p.Markdown())).
			Set("monthURL", "/m/"+date).
			Set(`Robots`, `noindex,nofollow`).
			Set("weekURL", "/w/"+date).
			Set("YMD", date) // #nosec G203
		ph.handleReply("rp", aWriter, pageData)

	case "robots.txt":
		ph.staticFS.ServeHTTP(aWriter, aRequest)

	case "s": // handle a query/search
		if 0 < len(tail) {
			ph.handleSearch(tail, pageData, aWriter /*, aRequest*/)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "share":
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
		if 0 < len(aRequest.URL.RawQuery) {
			// we need this for e.g. YouTube URLs
			tail += "?" + aRequest.URL.RawQuery
		}
		if 0 < len(aRequest.URL.Fragment) {
			tail += "#" + aRequest.URL.Fragment
		}
		ph.handleShare(tail, aWriter, aRequest)

	case "si": // store images
		ph.handleReply(`si`, aWriter,
			pageData.Set(`Robots`, `noindex,nofollow`))

	case "ss": // store static
		ph.handleReply(`ss`, aWriter,
			pageData.Set(`Robots`, `noindex,nofollow`))

	case "static": // deliver a static resource
		ph.staticFS.ServeHTTP(aWriter, aRequest)

	case "views": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case `w`, `ww`: // handle a given week
		var (
			d, y   int
			m      time.Month
			robots string = `noindex,follow`
		)
		if 0 == len(tail) {
			y, m, d = time.Now().Date()
		} else {
			y, m, d = getYMD(tail)
			// Allow indexing for lists older 30 days
			t := time.Date(y, m, d+30, 0, 0, 0, 0, time.Local)
			if t.Before(time.Now()) {
				robots = `index,follow`
			}
		}
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		pl := NewPostList().Week(y, m, d)
		ph.handleReply(`searchresult`, aWriter,
			pageData.Set(`Matches`, pl.Len()).
				Set(`monthURL`, `/m/`+date).
				Set(`Postings`, pl.Sort()).
				Set(`Robots`, robots).
				Set(`weekURL`, "/w/"+date))

	case `x`, `xp`, `xt`: // eXchange #tags/@mentions
		ph.handleReply(`xt`, aWriter,
			pageData.Set(`Robots`, `noindex,nofollow`))

	case ``:
		var val string // re-use variable
		if val = aRequest.FormValue("ht"); 0 < len(val) {
			ph.handleTagMentions(ph.hashList.HashList("#"+val),
				pageData, aWriter)
		} else if val = aRequest.FormValue("m"); 0 < len(val) {
			ph.reDir(aWriter, aRequest, "/m/"+val)
		} else if val = aRequest.FormValue("mt"); 0 < len(val) {
			ph.handleTagMentions(ph.hashList.MentionList("@"+val),
				pageData, aWriter)
		} else if val = aRequest.FormValue("n"); 0 < len(val) {
			ph.reDir(aWriter, aRequest, "/n/"+val)
		} else if val = aRequest.FormValue("p"); 0 < len(val) {
			ph.reDir(aWriter, aRequest, "/p/"+val)
		} else if val = aRequest.FormValue("q"); 0 < len(val) {
			ph.handleSearch(val, pageData, aWriter /*, aRequest*/)
		} else if val = aRequest.FormValue("s"); 0 < len(val) {
			ph.handleSearch(val, pageData, aWriter /*, aRequest*/)
		} else if val = aRequest.FormValue("share"); 0 < len(val) {
			if 0 < len(aRequest.URL.RawQuery) {
				// we need this for e.g. YouTube URLs
				val += "?" + aRequest.URL.RawQuery
			}
			ph.handleShare(val, aWriter, aRequest)
		} else if val = aRequest.FormValue("w"); 0 < len(val) {
			ph.reDir(aWriter, aRequest, "/w/"+val)
		} else {
			ph.handleRoot("30", pageData, aWriter /*, aRequest*/)
		}

	case `admin`, `echo.php`, `cgi-bin`, `config`, `console`, `.env`, `vendor`, `wp-content`:
		// Redirect spyware to the NSA:
		http.Redirect(aWriter, aRequest, `https://www.nsa.gov/`,
			http.StatusMovedPermanently)

	default:
		// if nothing matched (above) reply to the request
		// with an HTTP 404 not found error.
		http.NotFound(aWriter, aRequest)

	} // switch
} // handleGET()

// `handlePOST()` processes the HTTP POST requests.
func (ph *TPageHandler) handlePOST(aWriter http.ResponseWriter, aRequest *http.Request) {
	// Here we can't use
	//	ph.reDir(aWriter, aRequest, "/somethingelse/")
	// because we change the POST '/something/` URL to
	// GET `/somethingelse/` which would confuse the browser.
	var (
		bb  []byte
		err error
		i   int
		val string
	)
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case `ap`: // add a new post
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}

		if bb = replCRLF([]byte(aRequest.FormValue("manuscript"))); 0 < len(bb) {
			p := NewPosting(0, "").Set(bb)
			if _, err = p.Store(); nil != err {
				apachelogger.Err("TPageHandler.handlePOST('a')",
					fmt.Sprintf("TPosting.Store(%s): %v", p.IDstr(), err))
			}
			if AppArgs.Screenshot {
				PrepareLinkScreenshots(p)
			}
			AddTagID(ph.hashList, p)

			http.Redirect(aWriter, aRequest, "/p/"+p.IDstr(), http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case `dp`: // change date of posting
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}

		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

		op := NewPosting(str2id(tail), "")
		oid := op.ID()
		t := op.Time()

		y, mo, d := t.Date()
		if ymd := aRequest.FormValue("ymd"); 0 < len(ymd) {
			y, mo, d = getYMD(ymd)
		}

		h, mi, s, n := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
		if hms := aRequest.FormValue("hms"); 0 < len(hms) {
			h, mi, s = getHMS(hms)
		}

		t = time.Date(y, mo, d, h, mi, s, n, time.Local)
		nid := time2id(t)
		if err = poPersistence.Rename(oid, nid); nil != err {
			apachelogger.Err("TPageHandler.handlePOST('dp')",
				fmt.Sprintf("Persistence.Rename(%d, %d): %v", oid, nid, err))
		}

		RenameIDTags(ph.hashList, oid, nid)

		np := NewPosting(nid, "")
		np.Load()
		if AppArgs.Screenshot {
			PrepareLinkScreenshots(np)
		}

		http.Redirect(aWriter, aRequest, "/p/"+np.IDstr(), http.StatusSeeOther)

	case `ep`: // edit posting
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}

		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

		var old []byte
		nTxt := replCRLF([]byte(aRequest.FormValue("manuscript")))

		p := NewPosting(str2id(tail), "")
		if err = p.Load(); nil != err {
			apachelogger.Err("TPageHandler.handlePOST('e')",
				fmt.Sprintf("TPosting.Load(%s): %v", p.IDstr(), err))
		} else {
			old = p.Markdown()
		}

		if i, err = p.Set(nTxt).Store(); nil != err {
			if i < len(nTxt) {
				// let's hope for the best …
				_, _ = p.Set(old).Store()
			}
		}
		if AppArgs.Screenshot {
			PrepareLinkScreenshots(p)
		}
		UpdateTags(ph.hashList, p)

		tail += "?z=" + p.IDstr() // kick the browser cache
		http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)

	case `il`: // init hash list
		if nil != ph.hashList {
			if val = aRequest.FormValue("abort"); 0 < len(val) {
				http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
				return
			}

			ReadHashlist(ph.hashList)
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)
		}

	case `pv`: // update page previews/screenshots
		if AppArgs.Screenshot {
			if val = aRequest.FormValue("abort"); 0 < len(val) {
				http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
				return
			}

			UpdateScreenshots(PostingBaseDirectory())
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)
		}

	case `rp`: // posting removal
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}

		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

		post := NewPosting(str2id(tail), "")
		RemovePageScreenshots(post)
		if err = post.Delete(); nil != err {
			apachelogger.Err("TPageHandler.handlePOST('r')",
				fmt.Sprintf("TPosting.Delete(%s): %v", post.IDstr(), err))
		}
		RemoveIDTags(ph.hashList, str2id(tail))

		http.Redirect(aWriter, aRequest, "/m/"+post.Date(), http.StatusSeeOther)

	case `si`: // store image
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}

		if nil == ph.imgUp { // lazy initialisation
			ph.imgUp = uploadhandler.NewHandler(filepath.Join(AppArgs.DataDir, "/img/"),
				"imgFile", AppArgs.MaxFileSize)
		}
		ph.handleUpload(aWriter, aRequest, true)

	case `ss`: // store static
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}

		if nil == ph.staticUp { // lazy initialisation
			ph.staticUp = uploadhandler.NewHandler(filepath.Join(AppArgs.DataDir, "/static/"),
				"statFile", AppArgs.MaxFileSize)
		}
		ph.handleUpload(aWriter, aRequest, false)

	case `xt`: // eXchange #tags/@mentions
		if val = aRequest.FormValue("abort"); 0 < len(val) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}

		if val = aRequest.FormValue("search"); 0 < len(val) {
			if r := aRequest.FormValue("replace"); 0 < len(r) {
				ReplaceTag(ph.hashList,
					strings.TrimSpace(val),
					strings.TrimSpace(r)) // background operation
			}
		}
		http.Redirect(aWriter, aRequest, "/x/", http.StatusSeeOther)

	default:
		// // If nothing matched (above) reply to the request
		// // with an HTTP 404 "not found" error.
		// http.NotFound(aWriter, aRequest)

		// Redirect all invalid URLs to the NSA:
		http.Redirect(aWriter, aRequest, "https://www.nsa.gov/", http.StatusMovedPermanently)
	}
} // handlePOST()

// `handleReply()` sends `aPage` with `aData` to `aWriter`.
func (ph *TPageHandler) handleReply(aPage string, aWriter http.ResponseWriter, aData *TemplateData) {
	if err := ph.viewList.Render(aPage, aWriter, aData); nil != err {
		apachelogger.Err("TPageHandler.handleReply()",
			fmt.Sprintf("viewList.Render('%s'): %v", aPage, err))
	}
} // handleReply()

// `handleRoot()` serves the logical web-root directory.
func (ph *TPageHandler) handleRoot(aNumStr string,
	aData *TemplateData,
	aWriter http.ResponseWriter /* ,aRequest *http.Request*/) {
	num, start := numStart(aNumStr)
	if 0 == num {
		num = 30
	}

	pl := NewPostList()
	_ = pl.Newest(num, start) // ignore fs errors here

	aData = aData.Set(`Postings`, pl.Sort()).
		Set("Robots", "noindex,follow")
	if pl.Len() >= num {
		aData.Set("nextLink", fmt.Sprintf("/n/%d,%d", num, num+start+1))
	}
	ph.handleReply("index", aWriter, aData)
} // handleRoot()

// `handleSearch()` serves the search results.
func (ph *TPageHandler) handleSearch(aTerm string,
	aData *TemplateData,
	aWriter http.ResponseWriter /*, aRequest *http.Request*/) {
	pl := SearchPostings(regexp.QuoteMeta(strings.Trim(aTerm, `"`)))

	ph.handleReply(`searchresult`, aWriter,
		aData.Set(`Robots`, `noindex,follow`).
			Set(`Matches`, pl.Len()).
			Set(`Postings`, pl.Sort()))
} // handleSearch()

// `handleShare()` serves the edit page for a shared URL.
//
//	`aShare` The URL to share with the new posting.
//	`aWriter` The writer to respond to the remote user.
func (ph *TPageHandler) handleShare(aShare string, aWriter http.ResponseWriter, aRequest *http.Request) {
	p := NewPosting(0, "\n\n> [ ]("+aShare+")\n")
	if _, err := p.Store(); nil != err {
		apachelogger.Err("TPageHandler.handleShare()",
			fmt.Sprintf("TPosting.Store('%s'): %v", aShare, err))
	}

	CreateScreenshot(aShare) // background operation
	ph.reDir(aWriter, aRequest, "/e/"+p.IDstr())
} // handleShare()

// `handleTagMentions()` add the hashtag/mention list to `aData`
func (ph *TPageHandler) handleTagMentions(aList []uint64, aData *TemplateData, aWriter http.ResponseWriter) {
	var ( // re-use variables
		err  error
		id   uint64
		post *TPosting
	)
	pl := NewPostList()
	if 0 < len(aList) {
		for _, id = range aList {
			post = NewPosting(id, "")
			if err = post.Load(); nil != err {
				apachelogger.Err("TPageHandler.handleTagMentions()",
					fmt.Sprintf("TPosting.Load('%s'): %v", id2str(id), err))
				continue
			}
			pl.Add(post)
		}
	}

	ph.handleReply(`searchresult`, aWriter,
		aData.Set(`Robots`, `index,follow`).
			Set(`Matches`, pl.Len()).
			Set(`Postings`, pl.Sort()))
} // handleTagMentions()

// `handleUpload()` processes a file upload.
func (ph *TPageHandler) handleUpload(aWriter http.ResponseWriter, aRequest *http.Request, isImage bool) {
	var (
		status   int
		img, txt string
	)
	if isImage {
		img = "!"
		txt, status = ph.imgUp.ServeUpload(aWriter, aRequest)
	} else {
		txt, status = ph.staticUp.ServeUpload(aWriter, aRequest)
	}

	if 200 == status {
		fName := strings.TrimPrefix(txt, AppArgs.DataDir)
		post := NewPosting(0, "\n\n\n> "+img+"["+fName+"]("+fName+")\n\n")
		if _, err := post.Store(); nil != err {
			apachelogger.Err("TPageHandler.handleUpload()",
				fmt.Sprintf("TPosting.Store(%s): %v", post.IDstr(), err))
		}
		http.Redirect(aWriter, aRequest, "/e/"+post.IDstr(), http.StatusSeeOther)
	} else {
		http.Error(aWriter, txt, status)
	}
} // handleUpload()

// Len returns the length of the internal views list.
func (ph *TPageHandler) Len() int {
	return len(*(ph.viewList))
} // Len()

// NeedAuthentication returns `true` if authentication is needed,
// or `false` otherwise.
//
//	`aRequest` is the request to check.
func (ph *TPageHandler) NeedAuthentication(aRequest *http.Request) bool {
	path, _ := URLparts(aRequest.URL.Path)
	switch path {
	case `a`, `ap`, // add new post
		`d`, `dp`, // change post's date
		`e`, `ep`, // edit post
		`i`, `il`, // init hash list
		`r`, `rp`, // posting's removal
		`share`,    // share another URL
		`si`, `ss`, // store images, store static data
		`v`, `pv`, // update Screenshot
		`x`, `xp`, `xt`: // eXchange #tags/@mentions
		return true
	}

	var s string // re-use variable
	if s = aRequest.FormValue("share"); 0 < len(s) {
		return true
	}
	if s = aRequest.FormValue("si"); 0 < len(s) {
		return true
	}
	if s = aRequest.FormValue("ss"); 0 < len(s) {
		return true
	}

	return false
} // NeedAuthentication()

// `reDir()` continues handling the current `aRequest` by changing
// the requested URL to `aURL`.
func (ph *TPageHandler) reDir(aWriter http.ResponseWriter, aRequest *http.Request, aURL string) {
	aRequest.URL.Path = aURL
	ph.handleGET(aWriter, aRequest)
} // reDir()

// ServeHTTP handles the incoming HTTP requests.
func (ph *TPageHandler) ServeHTTP(aWriter http.ResponseWriter, aRequest *http.Request) {
	defer func() { // make sure a `panic` won't kill the program
		if err := recover(); err != nil {
			var msg string
			if AppArgs.LogStack {
				msg = fmt.Sprintf("caught panic: %v – %s", err, debug.Stack())
			} else {
				msg = fmt.Sprintf("caught panic: %v", err)
			}
			apachelogger.Err("TPageHandler.ServeHTTP()", msg)
		}
	}()

	aWriter.Header().Set(`Access-Control-Allow-Methods`, `GET, HEAD, POST`)
	if ph.NeedAuthentication(aRequest) {
		if nil == ph.userList {
			passlist.Deny(AppArgs.Realm, aWriter)
			return
		}
		if err := ph.userList.IsAuthenticated(aRequest); nil != err {
			passlist.Deny(AppArgs.Realm, aWriter)
			return
		}
	}

	switch aRequest.Method {
	case `GET`:
		ph.handleGET(aWriter, aRequest)

	case `HEAD`:
		ph.handleGET(aWriter, aRequest)

	case `OPTIONS`:
		aWriter.WriteHeader(http.StatusOK)

	case `POST`:
		ph.handlePOST(aWriter, aRequest)

	default:
		http.Error(aWriter, `HTTP Method Not Allowed`,
			http.StatusMethodNotAllowed)
	}
} // ServeHTTP()

/* _EoF_ */
