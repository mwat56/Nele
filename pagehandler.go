/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

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
	"strconv"
	"strings"
	"time"

	"github.com/mwat56/hashtags"
	"github.com/mwat56/passlist"
	"github.com/mwat56/uploadhandler"
)

type (
	// TPageHandler provides the handling of HTTP request/response.
	TPageHandler struct {
		addr     string                        // listen address ("1.2.3.4:5678")
		bn       string                        // the blog's name
		dd       string                        // datadir: base dir for data
		fh       http.Handler                  // static file handler
		hl       *hashtags.THashList           // #hashtags/@mentions list
		iup      *uploadhandler.TUploadHandler // `img` upload handler
		lang     string                        // default language
		mfs      int64                         // max. size of uploaded files
		realm    string                        // host/domain to secure by BasicAuth
		sup      *uploadhandler.TUploadHandler // `static` upload handler
		theme    string                        // `dark` or `light` display theme
		ul       *passlist.TPassList           // user/password list
		viewList *TViewList                    // list of template/views
	}
)

// `check4lang()` looks for a CGI value of `lang` and adds it to `aPD` if found.
func check4lang(aData *TDataList, aRequest *http.Request) *TDataList {
	if l := aRequest.FormValue("lang"); 0 < len(l) {
		return aData.Set("Lang", l)
	}

	return aData
} // check4lang()

// `handleShare()` serves the edit page for a shared URL.
func handleShare(aShare string, aWriter http.ResponseWriter, aRequest *http.Request) {
	p := NewPosting()
	p.Set([]byte("\n\n> [ ](" + aShare + ")\n"))
	if _, err := p.Store(); nil != err {
		log.Printf("handleShare('%s'): %v", aShare, err)
		//TODO better error handling
	}
	http.Redirect(aWriter, aRequest, "/e/"+p.ID(), http.StatusSeeOther)
} // handleShare()

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
	result.dd = s
	result.fh = http.FileServer(http.Dir(result.dd + "/"))
	if result.viewList, err = newViewList(filepath.Join(result.dd, "views")); nil != err {
		return nil, err
	}

	if s, err = AppArguments.Get("hashfile"); nil == err {
		result.hl, _ = hashtags.New("")
		result.hl.SetFilename(s)
		go goInitHashlist(result.hl)
	}

	if s, err = AppArguments.Get("lang"); nil == err {
		result.lang = s
	}

	if s, err = AppArguments.Get("listen"); nil != err {
		return nil, err
	}
	result.addr = s

	if s, err = AppArguments.Get("mfs"); nil != err {
		return nil, err
	}
	if mfs, _ := strconv.Atoi(s); 0 < mfs {
		result.mfs = int64(mfs)
	} else {
		result.mfs = 10485760 // 10 MB
	}

	if s, err = AppArguments.Get("port"); nil != err {
		return nil, err
	}
	result.addr += ":" + s

	if s, err = AppArguments.Get("uf"); nil != err {
		log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
	} else if result.ul, err = passlist.LoadPasswords(s); nil != err {
		log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
		result.ul = nil
	}

	if s, err = AppArguments.Get("realm"); nil == err {
		result.realm = s
	}

	if s, err = AppArguments.Get("theme"); nil != err {
		result.theme = "dark"
	} else {
		result.theme = s
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
		Set("Taglist", markupCloud(ph.hl)).
		Set("Title", ph.realm+": "+date).
		Set("weekURL", "/w/"+date)

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
		pageData = pageData.Set("Error", template.HTML(aData))
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
	// log.Printf("head: `%s`: tail: `%s`", path, tail) //FIXME REMOVE
	switch path {
	case "a", "ap": // add a new post
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.viewList.Render("ap", aWriter, pageData)

	case "certs": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "css":
		ph.fh.ServeHTTP(aWriter, aRequest)

	case "d", "dp": // change date
		if 0 < len(tail) {
			y, mo, d := time.Now().Date()
			now := fmt.Sprintf("%d-%02d-%02d", y, mo, d)
			p := newPosting(tail)
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
				Set("YMD", fmt.Sprintf("%d-%02d-%02d", y, mo, d))
			ph.viewList.Render("dc", aWriter, pageData)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "e", "ep": // edit a single posting
		if 0 < len(tail) {
			p := newPosting(tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				date := p.Date()
				pageData = check4lang(pageData, aRequest).
					Set("ID", p.ID()).
					Set("Manuscript", template.HTML(txt)).
					Set("monthURL", "/m/"+date).
					Set("Robots", "noindex,nofollow").
					Set("weekURL", "/w/"+date)
				ph.viewList.Render("ed", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "faq", "faq.html":
		ph.viewList.Render("faq", aWriter, check4lang(pageData, aRequest))

	case "favicon.ico":
		http.Redirect(aWriter, aRequest, "/img/"+path, http.StatusMovedPermanently)

	case "fonts":
		ph.fh.ServeHTTP(aWriter, aRequest)

	case "ht": // #hashtag search
		if 0 < len(tail) {
			ph.handleHashtag(tail, pageData, aWriter, aRequest)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "img":
		ph.fh.ServeHTTP(aWriter, aRequest)

	case "imprint", "impressum":
		ph.viewList.Render("imprint", aWriter, check4lang(pageData, aRequest))

	case "index", "index.html":
		ph.handleRoot("20", pageData, aWriter, aRequest)

		/*
			case "js":
				ph.fh.ServeHTTP(aWriter, aRequest)
		*/

	case "licence", "license", "lizenz":
		ph.viewList.Render("licence", aWriter, pageData)

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
		ph.viewList.Render("searchresult", aWriter, pageData)

	case "mt": // @mention search
		if 0 < len(tail) {
			ph.handleMention(tail, pageData, aWriter, aRequest)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "n", "np": // handle newest postings
		ph.handleRoot(tail, pageData, aWriter, aRequest)

	case "p", "pp": // handle a single posting
		if 0 < len(tail) {
			p := newPosting(tail)
			if err := p.Load(); nil == err {
				date := p.Date()
				pageData = check4lang(pageData, aRequest).
					Set("Posting", p).
					Set("monthURL", "/m/"+date).
					Set("weekURL", "/w/"+date)
				ph.viewList.Render("article", aWriter, pageData)
				return
			}
		}
		http.NotFound(aWriter, aRequest)

	case "postings": // this files are handled internally
		http.Redirect(aWriter, aRequest, "/n/", http.StatusMovedPermanently)

	case "privacy", "datenschutz":
		ph.viewList.Render("privacy", aWriter, check4lang(pageData, aRequest))

	case "q":
		http.Redirect(aWriter, aRequest, "/s/"+tail, http.StatusMovedPermanently)

	case "r", "rp": // remove posting
		if 0 < len(tail) {
			p := newPosting(tail)
			txt := p.Markdown()
			date := p.Date()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Set("Manuscript", template.HTML(txt)).
					Set("ID", p.ID()).
					Set("monthURL", "/m/"+date).
					Set("weekURL", "/w/"+date).
					Set("Robots", "noindex,nofollow")
				ph.viewList.Render("rp", aWriter, pageData)
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
		if 0 < len(tail) {
			if 0 < len(aRequest.URL.RawQuery) {
				// we need this for e.g. YouTube
				tail += "?" + aRequest.URL.RawQuery
			}
			handleShare(tail, aWriter, aRequest)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "si": // store images
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.viewList.Render("si", aWriter, pageData)

	case "ss": // store static
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,nofollow")
		ph.viewList.Render("ss", aWriter, pageData)

	case "static":
		ph.fh.ServeHTTP(aWriter, aRequest)

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
		ph.viewList.Render("searchresult", aWriter, pageData)

	case "":
		if ht := aRequest.FormValue("ht"); 0 < len(ht) {
			ph.handleHashtag(ht, pageData, aWriter, aRequest)
		} else if m := aRequest.FormValue("m"); 0 < len(m) {
			http.Redirect(aWriter, aRequest, "/m/"+m, http.StatusSeeOther)
		} else if mt := aRequest.FormValue("mt"); 0 < len(mt) {
			ph.handleMention(mt, pageData, aWriter, aRequest)
		} else if n := aRequest.FormValue("n"); 0 < len(n) {
			http.Redirect(aWriter, aRequest, "/n/"+n, http.StatusSeeOther)
		} else if p := aRequest.FormValue("p"); 0 < len(p) {
			http.Redirect(aWriter, aRequest, "/p/"+p, http.StatusSeeOther)
		} else if q := aRequest.FormValue("q"); 0 < len(q) {
			ph.handleSearch(q, pageData, aWriter, aRequest)
		} else if s := aRequest.FormValue("s"); 0 < len(s) {
			ph.handleSearch(s, pageData, aWriter, aRequest)
		} else if s := aRequest.FormValue("share"); 0 < len(s) {
			if 0 < len(aRequest.URL.RawQuery) {
				// we need this for e.g. YouTube
				s += "?" + aRequest.URL.RawQuery
			}
			handleShare(s, aWriter, aRequest)
		} else if w := aRequest.FormValue("w"); 0 < len(w) {
			http.Redirect(aWriter, aRequest, "/w/"+w, http.StatusSeeOther)
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
	tagList := ph.hl.HashList("#" + aTag)
	ph.handleTagMentions(tagList, aData, aWriter, aRequest)
} // handleHashtag()

func (ph *TPageHandler) handleMention(aMention string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	mentionList := ph.hl.MentionList("@" + aMention)

	ph.handleTagMentions(mentionList, aData, aWriter, aRequest)
} // handleMention()

func (ph *TPageHandler) handleTagMentions(aList []string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	pl := NewPostList()
	if 0 < len(aList) {
		for _, id := range aList {
			p := newPosting(id)
			err := p.Load()
			if nil == err {
				pl.Add(p)
			}
		}
	}

	aData = check4lang(aData, aRequest).
		Set("Robots", "index,follow").
		Set("Matches", pl.Len()).
		Set("Postings", pl.Sort())
	ph.viewList.Render("searchresult", aWriter, aData)

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
		fName = strings.TrimPrefix(txt, ph.dd)
		p := NewPosting()
		p.Set([]byte("\n\n\n> " + img + "[" + fName + "](" + fName + ")\n\n"))
		if _, err := p.Store(); nil != err {
			log.Printf("handlePOST(): %v", err)
			//TODO better error handling
		}
		http.Redirect(aWriter, aRequest, "/e/"+p.ID(), http.StatusSeeOther)
	} else {
		// aWriter.WriteHeader(status)
		// aWriter.Write([]byte(txt))
		http.Error(aWriter, txt, status)
	}
} // handleUpload()

// `handlePOST()` process the HTTP POST requests.
func (ph *TPageHandler) handlePOST(aWriter http.ResponseWriter, aRequest *http.Request) {
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case "a": // add a new post
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if m := replCRLF([]byte(aRequest.FormValue("manuscript"))); 0 < len(m) {
			p := NewPosting()
			p.Set(m)
			if _, err := p.Store(); nil != err {
				log.Printf("handlePOST(a): %v\n", err)
				//TODO better error handling
			}
			go goAddID(ph.hl, p.ID(), p.Markdown())

			// tail = p.ID() + "?z=" + p.Date()
			http.Redirect(aWriter, aRequest, "/p/"+p.ID(), http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "d": // change date
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 < len(tail) {
			op := newPosting(tail)
			t := op.Time()
			y, mo, d := t.Date()
			if ymd := aRequest.FormValue("ymd"); 0 < len(ymd) {
				y, mo, d = getYMD(ymd)
			}
			h := t.Hour()
			mi := t.Minute()
			s := t.Second()
			n := t.Nanosecond()
			if hms := aRequest.FormValue("hms"); 0 < len(hms) {
				h, mi, s = getHMS(hms)
			}
			opn := op.PathFileName()
			t = time.Date(y, mo, d, h, mi, s, n, time.Local)
			np := newPosting(newID(t))
			npn := np.PathFileName()
			// ensure existence of directory:
			if _, err := np.makeDir(); nil != err {
				log.Printf("handlePOST(d1): %v\n", err)
				//TODO better error handling
			}
			if err := os.Rename(opn, npn); nil != err {
				log.Printf("handlePOST(d2): %v\n", err)
				//TODO better error handling
			}
			go goRenameID(ph.hl, tail, np.ID())

			http.Redirect(aWriter, aRequest, "/p/"+np.ID(), http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "e": // edit posting
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 < len(tail) {
			var old []byte
			m := replCRLF([]byte(aRequest.FormValue("manuscript")))
			p := newPosting(tail)
			if err := p.Load(); nil == err {
				old = p.Markdown()
			}
			p.Set(m)
			if bw, err := p.Store(); nil != err {
				if bw < int64(len(m)) {
					// let's hope for the best …
					p.Set(old).Store()
				}
			}
			go goUpdateID(ph.hl, tail, m)

			tail += "?z=" + p.ID() // kick the browser cache
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "r": // remove posting
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		if 0 < len(tail) {
			p := newPosting(tail)
			if err := p.Delete(); nil != err {
				log.Printf("handlePOST(r): %v\n", err)
				//TODO better error handling
			}
			go goRemoveID(ph.hl, tail)

			http.Redirect(aWriter, aRequest, "/m/"+p.Date(), http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}

	case "si": // store image
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if nil == ph.iup { // lazy initialisation
			ph.iup = uploadhandler.NewHandler(filepath.Join(ph.dd, "/img/"),
				"imgFile", ph.mfs)
		}
		ph.handleUpload(aWriter, aRequest, true)

	case "ss": // store static
		if a := aRequest.FormValue("abort"); 0 < len(a) {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
			return
		}
		if nil == ph.sup { // lazy initialisation
			ph.sup = uploadhandler.NewHandler(filepath.Join(ph.dd, "/static/"),
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
		num = 20
	}
	pl := NewPostList()
	pl.Newest(num, start) // ignore fs errors here
	aData = check4lang(aData, aRequest).
		Set("Postings", pl.Sort()).
		Set("Robots", "noindex,follow")
	if pl.Len() >= num {
		aData.Set("nextLink", fmt.Sprintf("/n/%d,%d", num, num+start+1))
	}
	ph.viewList.Render("index", aWriter, aData)
} // handleRoot()

// `handleSearch()` serves the search results.
func (ph *TPageHandler) handleSearch(aTerm string, aData *TDataList, aWriter http.ResponseWriter, aRequest *http.Request) {
	pl := SearchPostings(regexp.QuoteMeta(aTerm))
	aData = check4lang(aData, aRequest).
		Set("Robots", "noindex,follow").
		Set("Matches", pl.Len()).
		Set("Postings", pl.Sort())
	ph.viewList.Render("searchresult", aWriter, aData)
} // handleSearch()

// Len returns the length of the internal view list.
func (ph *TPageHandler) Len() int {
	return len(*((*ph).viewList))
} // Len()

// NeedAuthentication returns `true` if authentication is needed,
// or `false` otherwise.
//
// `aURL` is the URL to check.
func (ph *TPageHandler) NeedAuthentication(aRequest *http.Request) bool {
	path, _ := URLparts(aRequest.URL.Path)
	switch path {
	case "a", "ap", // add new post
		"d", "dp", // change post's date
		"e", "ep", // edit post
		"r", "rp", // remove post
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

// ServeHTTP handles the incoming HTTP requests.
func (ph TPageHandler) ServeHTTP(aWriter http.ResponseWriter, aRequest *http.Request) {
	if ph.NeedAuthentication(aRequest) {
		if nil == ph.ul {
			passlist.Deny(ph.realm, aWriter)
			return
		}
		if !ph.ul.IsAuthenticated(aRequest) {
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
