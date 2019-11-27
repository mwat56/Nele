/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides functions and methods to handle HTTP requests.
 */

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/mwat56/apachelogger"
	"github.com/mwat56/hashtags"
	"github.com/mwat56/jffs"
	"github.com/mwat56/passlist"
	"github.com/mwat56/uploadhandler"
)

type (
	// TPageHandler provides the handling of HTTP request/response.
	TPageHandler struct {
		addr     string                        // listen address ("1.2.3.4:5678")
		bn       string                        // the blog's name
		dataDir  string                        // datadir: base dir for data
		hashList *hashtags.THashList           // #hashtags/@mentions list
		iup      *uploadhandler.TUploadHandler // `img` upload handler
		lang     string                        // default language
		logStack bool                          // log stack trace
		mfs      int64                         // max. size of uploaded files
		realm    string                        // host/domain to secure by BasicAuth
		sfh      http.Handler                  // `static` file handler
		sup      *uploadhandler.TUploadHandler // `static` upload handler
		theme    string                        // `dark` or `light` display theme
		userList *passlist.TPassList           // user/password list
		viewList *TViewList                    // list of template/views
	}
)

// `check4lang()` looks for a CGI value of `lang` and adds it to `aData` if found.
func check4lang(aData *TDataList, aRequest *http.Request) *TDataList {
	if l := aRequest.FormValue("lang"); 0 < len(l) {
		return aData.Set("Lang", l)
	}

	return aData
} // check4lang()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// NewPageHandler returns a new `TPageHandler` instance.
func NewPageHandler() (*TPageHandler, error) {
	var (
		err error
		s   string
	)
	result := new(TPageHandler)

	if s, err = AppArguments.Get("blogname"); nil == err {
		result.bn = s
	}

	if s, err = AppArguments.Get("datadir"); nil != err {
		return nil, err
	}
	result.dataDir = s
	result.sfh = jffs.FileServer(http.Dir(result.dataDir + `/`))
	if result.viewList, err = newViewList(filepath.Join(result.dataDir, "views")); nil != err {
		return nil, err
	}

	if s, err = AppArguments.Get("hashfile"); nil == err {
		result.hashList, _ = hashtags.New("")
		result.hashList.SetFilename(s)
		InitHashlist(result.hashList)
	}

	if s, err = AppArguments.Get("lang"); nil == err {
		result.lang = s
	}

	// an empty value means: listen on all interfaces:
	result.addr, _ = AppArguments.Get("listen")

	if s, err = AppArguments.Get("logStack"); nil == err {
		result.logStack = ("true" == s)
	}

	s, _ = AppArguments.Get("mfs")
	if mfs, _ := strconv.Atoi(s); 0 < mfs {
		result.mfs = int64(mfs)
	} else {
		result.mfs = 10485760 // 10 MB
	}

	if s, err = AppArguments.Get("pageview"); nil == err {
		if pageview := ("true" == s); pageview {
			InitPageImages(PostingBaseDirectory(), "./img/")
		}
	}

	if s, err = AppArguments.Get("port"); nil != err {
		return nil, err
	}
	result.addr += ":" + s

	if s, err = AppArguments.Get("uf"); nil != err {
		log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
	} else if result.userList, err = passlist.LoadPasswords(s); nil != err {
		log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
		result.userList = nil
	}

	if s, err = AppArguments.Get("realm"); nil == err {
		result.realm = s
	}

	if s, err = AppArguments.Get("theme"); (nil == err) && (0 < len(s)) {
		result.theme = s
	} else {
		result.theme = "dark"
	}

	return result, nil
} // NewPageHandler()

// `newViewList()` returns a list of views found in `aDirectory`
// and a possible I/O error.
func newViewList(aDirectory string) (*TViewList, error) {
	var v *TView
	result := NewViewList()

	files, err := filepath.Glob(aDirectory + "/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, fName := range files {
		fName := filepath.Base(fName[:len(fName)-7]) // remove extension
		if v, err = NewView(aDirectory, fName); nil != err {
			return nil, err
		}
		result = result.Add(v)
	}

	return result, nil
} // newViewList()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// Address returns the configured `IP:Port` address to use for listening.
func (ph *TPageHandler) Address() string {
	return ph.addr
} // Address()

// `basicPageData()` returns a list of common Head entries.
func (ph *TPageHandler) basicPageData() *TDataList {
	y, m, d := time.Now().Date()
	date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
	pageData := NewDataList().
		Set("Blogname", ph.bn).
		Set("CSS", template.HTML(`<link rel="stylesheet" type="text/css" title="mwat's styles" href="/css/stylesheet.css"><link rel="stylesheet" type="text/css" href="/css/`+ph.theme+`.css"><link rel="stylesheet" type="text/css" href="/css/fonts.css">`)).
		Set("Lang", ph.lang).
		Set("monthURL", "/m/"+date).
		Set("Robots", "index,follow").
		Set("Taglist", MarkupCloud(ph.hashList)).
		Set("Title", ph.realm+": "+date).
		Set("weekURL", "/w/"+date) // #nosec G203

	return pageData
} // basicPageData()

// GetErrorPage returns an error page for `aStatus`,
// implementing the `TErrorPager` interface.
func (ph *TPageHandler) GetErrorPage(aData []byte, aStatus int) []byte {
	var empty []byte

	pageData := ph.basicPageData().
		Set("Robots", "noindex,follow")

	switch aStatus {
	case 404:
		if page, err := ph.viewList.RenderedPage("404", pageData); nil == err {
			return page
		}

	//TODO implement other status codes

	default:
		pageData = pageData.Set("Error", template.HTML(aData)) // #nosec G203
		if page, err := ph.viewList.RenderedPage("error", pageData); nil == err {
			return page
		}
	}

	return empty
} // GetErrorPage()

// `handleGET()` processes the HTTP GET requests.
func (ph *TPageHandler) handleGET(aWriter http.ResponseWriter, aRequest *http.Request) {

	pageData := ph.basicPageData()
	path, tail := URLparts(aRequest.URL.Path)
	switch strings.ToLower(path) { // handle URLs case-insensitive

	case "a", "ap": // add a new post
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.handleReply("ap", aWriter, pageData)

	case "certs": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "css":
		ph.sfh.ServeHTTP(aWriter, aRequest)

	case "d", "dp": // change date
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
		y, mo, d := time.Now().Date()
		now := fmt.Sprintf("%d-%02d-%02d", y, mo, d)
		p := NewPosting(tail)
		txt := p.Markdown()
		t := p.Time()
		y, mo, d = t.Date()
		pageData = check4lang(pageData, aRequest).
			Set("HMS", fmt.Sprintf("%02d:%02d:%02d",
				t.Hour(), t.Minute(), t.Second())).
			Set("ID", p.ID()).
			Set("Manuscript", template.HTML(txt)).
			Set("NOW", now).
			Set("Robots", "noindex,nofollow").
			Set("YMD", fmt.Sprintf("%d-%02d-%02d", y, mo, d)) // #nosec G203
		ph.handleReply("dp", aWriter, pageData)

	case "e", "ep": // edit a single posting
		if 0 < len(tail) {
			p := NewPosting(tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				t := p.Time()
				date := p.Date()
				pageData = check4lang(pageData, aRequest).
					Set("HMS",
						fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())).
					Set("ID", p.ID()).
					Set("Manuscript", template.HTML(txt)).
					Set("monthURL", "/m/"+date).
					Set("Robots", "noindex,nofollow").
					Set("weekURL", "/w/"+date).
					Set("YMD", date) // #nosec G203
				ph.handleReply("ep", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "faq", "faq.html":
		ph.handleReply("faq", aWriter, check4lang(pageData, aRequest))

	case "favicon.ico":
		http.Redirect(aWriter, aRequest, "/img/"+path, http.StatusMovedPermanently)

	case "fonts":
		ph.sfh.ServeHTTP(aWriter, aRequest)

	case "hl": // #hashtag list
		if 0 < len(tail) {
			ph.handleHashtag(tail, pageData, aWriter, aRequest)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "img":
		ph.sfh.ServeHTTP(aWriter, aRequest)

	case "imprint", "impressum":
		ph.handleReply("imprint", aWriter, check4lang(pageData, aRequest))

	case "index", "index.html":
		ph.handleRoot("20", pageData, aWriter, aRequest)

		/*
			case "js":
				ph.sfh.ServeHTTP(aWriter, aRequest)
		*/

	case "licence", "license", "lizenz":
		ph.handleReply("licence", aWriter, check4lang(pageData, aRequest))

	case "m", "mm": // handle a given month
		var y, d int
		var m time.Month
		if 0 == len(tail) {
			y, m, d = time.Now().Date()
		} else {
			y, m, d = getYMD(tail)
			if 0 == d {
				d = 1
			}
		}
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		pl := NewPostList().Month(y, m)
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,follow").
			Set("Matches", pl.Len()).
			Set("Postings", pl.Sort()).
			Set("monthURL", "/m/"+date).
			Set("weekURL", "/w/"+date)
		ph.handleReply("searchresult", aWriter, pageData)

	case "ml": // @mention list
		if 0 < len(tail) {
			ph.handleMention(tail, pageData, aWriter, aRequest)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "n", "np": // handle newest postings
		ph.handleRoot(tail, pageData, aWriter, aRequest)

	case "p", "pp": // handle a single posting
		if 0 < len(tail) {
			p := NewPosting(tail)
			if err := p.Load(); nil != err {
				apachelogger.Err("TPageHandler.handleGET()",
					fmt.Sprintf("TPosting.Load(%s): %v", p.ID(), err))
			} else {
				date := p.Date()
				pageData = check4lang(pageData, aRequest).
					Set("Posting", p).
					Set("monthURL", "/m/"+date).
					Set("weekURL", "/w/"+date)
				ph.handleReply("article", aWriter, pageData)
				return
			}
		}
		http.NotFound(aWriter, aRequest)

	case "postings": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "privacy", "datenschutz":
		ph.handleReply("privacy", aWriter, check4lang(pageData, aRequest))

	case "q":
		http.Redirect(aWriter, aRequest, "/s/"+tail, http.StatusMovedPermanently)

	case "r", "rp": // posting's removal
		if 0 < len(tail) {
			p := NewPosting(tail)
			txt := p.Markdown()
			date := p.Date()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Set("Manuscript", template.HTML(txt)).
					Set("ID", p.ID()).
					Set("monthURL", "/m/"+date).
					Set("weekURL", "/w/"+date).
					Set("Robots", "noindex,nofollow") // #nosec G203
				ph.handleReply("rp", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "s": // handle a query/search
		if 0 < len(tail) {
			ph.handleSearch(tail, pageData, aWriter, aRequest)
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
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.handleReply("si", aWriter, pageData)

	case "ss": // store static
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.handleReply("ss", aWriter, pageData)

	case "static":
		ph.sfh.ServeHTTP(aWriter, aRequest)

	case "views": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "w", "ww": // handle a given week
		var y, d int
		var m time.Month
		if 0 == len(tail) {
			y, m, d = time.Now().Date()
		} else {
			y, m, d = getYMD(tail)
			if 0 == d {
				d = 1
			}
		}
		date := fmt.Sprintf("%d-%02d-%02d", y, m, d)
		pl := NewPostList().Week(y, m, d)
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,follow").
			Set("Matches", pl.Len()).
			Set("Postings", pl.Sort()).
			Set("monthURL", "/m/"+date).
			Set("weekURL", "/w/"+date)
		ph.handleReply("searchresult", aWriter, pageData)

	case "":
		if ht := aRequest.FormValue("ht"); 0 < len(ht) {
			ph.handleHashtag(ht, pageData, aWriter, aRequest)
		} else if m := aRequest.FormValue("m"); 0 < len(m) {
			ph.reDir(aWriter, aRequest, "/m/"+m)
		} else if mt := aRequest.FormValue("mt"); 0 < len(mt) {
			ph.handleMention(mt, pageData, aWriter, aRequest)
		} else if n := aRequest.FormValue("n"); 0 < len(n) {
			ph.reDir(aWriter, aRequest, "/n/"+n)
		} else if p := aRequest.FormValue("p"); 0 < len(p) {
			ph.reDir(aWriter, aRequest, "/p/"+p)
		} else if q := aRequest.FormValue("q"); 0 < len(q) {
			ph.handleSearch(q, pageData, aWriter, aRequest)
		} else if s := aRequest.FormValue("s"); 0 < len(s) {
			ph.handleSearch(s, pageData, aWriter, aRequest)
		} else if s := aRequest.FormValue("share"); 0 < len(s) {
			if 0 < len(aRequest.URL.RawQuery) {
				// we need this for e.g. YouTube URLs
				s += "?" + aRequest.URL.RawQuery
			}
			ph.handleShare(s, aWriter, aRequest)
		} else if w := aRequest.FormValue("w"); 0 < len(w) {
			ph.reDir(aWriter, aRequest, "/w/"+w)
		} else {
			ph.handleRoot("20", pageData, aWriter, aRequest)
		}

	default:
		// if nothing matched (above) reply to the request
		// with an HTTP 404 not found error.
		http.NotFound(aWriter, aRequest)
	} // switch
} // handleGET()

func (ph *TPageHandler) handleHashtag(aTag string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	tagList := ph.hashList.HashList("#" + aTag)

	ph.handleTagMentions(tagList, aData, aWriter, aRequest)
} // handleHashtag()

func (ph *TPageHandler) handleMention(aMention string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	mentionList := ph.hashList.MentionList("@" + aMention)

	ph.handleTagMentions(mentionList, aData, aWriter, aRequest)
} // handleMention()

// `handleReply()` sends `aPage` with `aData` to `aWriter`.
func (ph *TPageHandler) handleReply(aPage string, aWriter http.ResponseWriter, aData *TDataList) {
	if err := ph.viewList.Render(aPage, aWriter, aData); nil != err {
		apachelogger.Err("TPageHandler.handleReply()",
			fmt.Sprintf("viewList.Render('%s'): %v", aPage, err))
	}
} // handleReply()

// `handleShare()` serves the edit page for a shared URL.
func (ph *TPageHandler) handleShare(aShare string, aWriter http.ResponseWriter, aRequest *http.Request) {
	p := NewPosting("")
	p.Set([]byte("\n\n> [ ](" + aShare + ")\n"))
	if _, err := p.Store(); nil != err {
		apachelogger.Err("TPageHandler.handleShare()",
			fmt.Sprintf("TPosting.Store('%s'): %v", aShare, err))
	}

	ph.reDir(aWriter, aRequest, "/e/"+p.ID())
} // handleShare()

func (ph *TPageHandler) handleTagMentions(aList []string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	pl := NewPostList()
	if 0 < len(aList) {
		for _, id := range aList {
			p := NewPosting(id)
			if err := p.Load(); nil != err {
				apachelogger.Err("TPageHandler.handleTagMentions()",
					fmt.Sprintf("TPosting.Load('%s'): %v", id, err))
				continue
			}
			pl.Add(p)
		}
	}

	aData = check4lang(aData, aRequest).
		Set("Robots", "index,follow").
		Set("Matches", pl.Len()).
		Set("Postings", pl.Sort())
	ph.handleReply("searchresult", aWriter, aData)
} // handleTagMentions()

// `handleUpload()` processes a file upload.
func (ph *TPageHandler) handleUpload(aWriter http.ResponseWriter, aRequest *http.Request, isImage bool) {
	var (
		status          int
		fName, img, txt string
	)
	if isImage {
		img = "!"
		txt, status = ph.iup.ServeUpload(aWriter, aRequest)
	} else {
		txt, status = ph.sup.ServeUpload(aWriter, aRequest)
	}

	if 200 == status {
		fName = strings.TrimPrefix(txt, ph.dataDir)
		p := NewPosting("")
		p.Set([]byte("\n\n\n> " + img + "[" + fName + "](" + fName + ")\n\n"))
		if _, err := p.Store(); nil != err {
			apachelogger.Err("TPageHandler.handleUpload()",
				fmt.Sprintf("TPosting.Store(%s): %v", p.ID(), err))
		}
		http.Redirect(aWriter, aRequest, "/e/"+p.ID(), http.StatusSeeOther)
	} else {
		http.Error(aWriter, txt, status)
	}
} // handleUpload()

// `handlePOST()` process the HTTP POST requests.
func (ph *TPageHandler) handlePOST(aWriter http.ResponseWriter, aRequest *http.Request) {
	// Here we can't use
	//	ph.reDir(aWriter, aRequest, "/somethingelse/")
	// because we change the POST '/something/` URL to
	// GET `/somethingelse/` which would confuse the browser.
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case "a": // add a new post
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if m := replCRLF([]byte(aRequest.FormValue("manuscript"))); 0 < len(m) {
			p := NewPosting("")
			p.Set(m)
			if _, err := p.Store(); nil != err {
				apachelogger.Err("TPageHandler.handlePOST()",
					fmt.Sprintf("TPosting.Store(%s): %v", p.ID(), err))
			}
			go goAddID(ph.hashList, p.ID(), p.Markdown())

			http.Redirect(aWriter, aRequest, "/p/"+p.ID(), http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "d": // change date of posting
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
		op := NewPosting(tail)
		t := op.Time()
		y, mo, d := t.Date()
		if ymd := aRequest.FormValue("ymd"); 0 < len(ymd) {
			y, mo, d = getYMD(ymd)
		}
		h, mi, s, n := t.Hour(), t.Minute(), t.Second(), t.Nanosecond()
		if hms := aRequest.FormValue("hms"); 0 < len(hms) {
			h, mi, s = getHMS(hms)
		}
		opn := op.PathFileName()
		t = time.Date(y, mo, d, h, mi, s, n, time.Local)
		np := NewPosting(newID(t))
		npn := np.PathFileName()
		// ensure existence of directory:
		if _, err := np.makeDir(); nil != err {
			apachelogger.Err("TPageHandler.handlePOST()",
				fmt.Sprintf("np.makeDir(%s): %v", np.ID(), err))
		}
		if err := os.Rename(opn, npn); nil != err {
			apachelogger.Err("TPageHandler.handlePOST()",
				fmt.Sprintf("os.Rename(%s, %s): %v", opn, npn, err))
		}
		go goRenameID(ph.hashList, tail, np.ID())

		http.Redirect(aWriter, aRequest, "/p/"+np.ID(), http.StatusSeeOther)

	case "e": // edit posting
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
		var old []byte
		m := replCRLF([]byte(aRequest.FormValue("manuscript")))
		p := NewPosting(tail)
		if err := p.Load(); nil != err {
			apachelogger.Err("TPageHandler.handlePOST()",
				fmt.Sprintf("TPosting.Load(%s): %v", p.ID(), err))
		} else {
			old = p.Markdown()
		}
		p.Set(m)
		if bw, err := p.Store(); nil != err {
			if bw < int64(len(m)) {
				// let's hope for the best …
				_, _ = p.Set(old).Store()
			}
		}
		go goUpdateID(ph.hashList, tail, m)

		tail += "?z=" + p.ID() // kick the browser cache
		http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)

	case "r": // posting removal
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 == len(tail) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
		p := NewPosting(tail)
		if err := p.Delete(); nil != err {
			apachelogger.Err("TPageHandler.handlePOST()",
				fmt.Sprintf("TPosting.Delete(%s): %v", p.ID(), err))
		}
		go goRemoveID(ph.hashList, tail)

		http.Redirect(aWriter, aRequest, "/m/"+p.Date(), http.StatusSeeOther)

	case "si": // store image
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if nil == ph.iup { // lazy initialisation
			ph.iup = uploadhandler.NewHandler(filepath.Join(ph.dataDir, "/img/"),
				"imgFile", ph.mfs)
		}
		ph.handleUpload(aWriter, aRequest, true)

	case "ss": // store static
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if nil == ph.sup { // lazy initialisation
			ph.sup = uploadhandler.NewHandler(filepath.Join(ph.dataDir, "/static/"),
				"statFile", ph.mfs)
		}
		ph.handleUpload(aWriter, aRequest, false)

	default:
		// if nothing matched (above) reply to the request
		// with an HTTP 404 "not found" error.
		http.NotFound(aWriter, aRequest)
	}
} // handlePOST()

// `handleRoot()` serves the logical web-root directory.
func (ph *TPageHandler) handleRoot(aNumStr string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	num, start := numStart(aNumStr)
	if 0 == num {
		num = 30
	}
	pl := NewPostList()
	_ = pl.Newest(num, start) // ignore fs errors here
	aData = check4lang(aData, aRequest).
		Set("Postings", pl.Sort()).
		Set("Robots", "noindex,follow")
	if pl.Len() >= num {
		aData.Set("nextLink", fmt.Sprintf("/n/%d,%d", num, num+start+1))
	}
	ph.handleReply("index", aWriter, aData)
} // handleRoot()

// `handleSearch()` serves the search results.
func (ph *TPageHandler) handleSearch(aTerm string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	pl := SearchPostings(regexp.QuoteMeta(aTerm))
	aData = check4lang(aData, aRequest).
		Set("Robots", "noindex,follow").
		Set("Matches", pl.Len()).
		Set("Postings", pl.Sort())
	ph.handleReply("searchresult", aWriter, aData)
} // handleSearch()

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
	case "a", "ap", // add new post
		"d", "dp", // change post's date
		"e", "ep", // edit post
		"r", "rp", // posting's removal
		"share",    // share another URL
		"si", "ss": // store images, store static data
		return true
	}

	if s := aRequest.FormValue("share"); 0 < len(s) {
		return true
	}
	if s := aRequest.FormValue("si"); 0 < len(s) {
		return true
	}
	if s := aRequest.FormValue("ss"); 0 < len(s) {
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
	defer func() {
		// make sure a `panic` won't kill the program
		if err := recover(); err != nil {
			var msg string
			if ph.logStack {
				msg = fmt.Sprintf("caught panic: %v – %s", err, debug.Stack())
			} else {
				msg = fmt.Sprintf("caught panic: %v", err)
			}
			apachelogger.Err("TPageHandler.ServeHTTP()", msg)
		}
	}()

	aWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST")
	if ph.NeedAuthentication(aRequest) {
		if nil == ph.userList {
			passlist.Deny(ph.realm, aWriter)
			return
		}
		if !ph.userList.IsAuthenticated(aRequest) {
			passlist.Deny(ph.realm, aWriter)
			return
		}
	}

	switch aRequest.Method {
	case "GET":
		ph.handleGET(aWriter, aRequest)

	case "POST":
		ph.handlePOST(aWriter, aRequest)

	default:
		http.Error(aWriter, "HTTP Method Not Allowed", http.StatusMethodNotAllowed)
	}
} // ServeHTTP()

/* _EoF_ */
