/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package blog

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

	passlist "github.com/mwat56/go-passlist"
)

type (
	// TPageHandler provides the handling of HTTP request/response.
	TPageHandler struct {
		addr     string              // listen address ("1.2.3.4:5678")
		basedir  string              // base of directories holding the posts
		cssD     string              // configured CSS directory
		cssH     http.Handler        // static CSS file handler
		imgD     string              // configured image directopry
		imgH     http.Handler        // static image file handler
		jsD      string              // configured JavaScript directory
		jsH      http.Handler        // static JS files
		lang     string              // default language
		realm    string              // host/domain to secure by BasicAuth
		staticD  string              // configured static directory
		staticH  http.Handler        // other static files
		ul       *passlist.TPassList // user/password list
		viewList *TViewList
	}

	tBoolMap map[bool]string
)

var (
	// simple lookup table for daily changing style sheets
	styles = tBoolMap{
		true:  "light",
		false: "dark",
	}
)

// Address returns the configured `IP:Port` address to use for listening.
func (ph *TPageHandler) Address() string {
	return ph.addr
} // Address()

// `basicPageData()` returns a list of common Head entries.
func (ph *TPageHandler) basicPageData() *TDataList {
	day := time.Now().Day()
	css := fmt.Sprintf(`<link rel="stylesheet" type="text/css" title="mwat's styles" href="/css/stylesheet.css" /><link rel="stylesheet" type="text/css" href="/css/%s.css" />`, styles[1 == day&1])
	pageData := NewDataList().
		Add("CSS", template.HTML(css)).
		Add("Lang", ph.lang).
		Add("Robots", "index,follow")

	return pageData
} // basicPageData()

// check4lang() looks for a CGI value of `lang` and adds it to `aPD` if found.
func check4lang(aPD *TDataList, aRequest *http.Request) *TDataList {
	if l := aRequest.FormValue("lang"); 0 < len(l) {
		return aPD.Add("Lang", l)
	}
	return aPD
} // check4lang()

// GetErrorPage returns an error page for `aStatus`,
// implementing the `TErrorPager` interface.
func (ph *TPageHandler) GetErrorPage(aData []byte, aStatus int) []byte {
	var empty []byte

	pageData := ph.basicPageData().
		Add("Robots", "noindex,follow")

	switch aStatus {
	case 404:
		if page, err := ph.viewList.RenderedPage("404", pageData); nil == err {
			return page
		}

	//TODO implement other status codes

	default:
		pageData = pageData.Add("Error", template.HTML(aData))
		if page, err := ph.viewList.RenderedPage("error", pageData); nil == err {
			return page
		}
	}

	return empty
} // GetErrorPage()

// handleGET() processes the HTTP GET requests.
func (ph *TPageHandler) handleGET(aWriter http.ResponseWriter, aRequest *http.Request) {

	pageData := ph.basicPageData()
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case "a", "ap": // add a new post
		pageData = check4lang(pageData, aRequest).
			Add("BackURL", aRequest.URL.Path). // for POSTing back
			Add("Robots", "noindex,nofollow")
		ph.viewList.Render("ap", aWriter, pageData)

	case ph.cssD:
		if 0 < len(tail) {
			aRequest.URL.Path = tail
		}
		ph.cssH.ServeHTTP(aWriter, aRequest)

	case "d", "dp": // change date
		if 0 < len(tail) {
			y, mo, d := time.Now().Date()
			now := fmt.Sprintf("%d-%02d-%02d", y, mo, d)
			t := timeID(tail)
			y, mo, d = t.Date()
			pageData = check4lang(pageData, aRequest).
				Add("BackURL", aRequest.URL.Path). // for POSTing back
				Add("NOW", now).
				Add("HMS", fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())).
				Add("YMD", fmt.Sprintf("%d-%02d-%02d", y, mo, d)).
				Add("Robots", "noindex,nofollow")
			ph.viewList.Render("dc", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "e", "ep": // edit a single posting
		if 0 < len(tail) {
			p := newPosting(ph.basedir, tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Add("Manuscript", template.HTML(txt)).
					Add("BackURL", aRequest.URL.Path). // for POSTing back
					Add("Robots", "noindex,nofollow")
				ph.viewList.Render("ed", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "faq", "faq.html":
		ph.viewList.Render("faq", aWriter, check4lang(pageData, aRequest))

	case ph.imgD, "favicon.ico":
		if 0 < len(tail) {
			aRequest.URL.Path = tail
		}
		ph.imgH.ServeHTTP(aWriter, aRequest)

	case "imprint", "imprint.html":
		ph.viewList.Render("imprint", aWriter, check4lang(pageData, aRequest))

	case "index", "index.html":
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case ph.jsD:
		if 0 < len(tail) {
			aRequest.URL.Path = tail
		}
		ph.jsH.ServeHTTP(aWriter, aRequest)

	case "licence", "license":
		ph.viewList.Render("licence", aWriter, pageData)

	case "m": // handle a given month
		if 0 < len(tail) {
			y, m, _ := getYMD(tail)
			pl := NewPostList(ph.basedir).Month(y, m)
			pageData = check4lang(pageData, aRequest).
				Add("Robots", "noindex,follow").
				Add("Matches", pl.Len()).
				Add("Postings", pl.Sort())
			ph.viewList.Render("searchresult", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "n": // handle newest postings
		num := 20
		if 0 < len(tail) {
			var err error
			if num, err = strconv.Atoi(tail); nil != err {
				num = 20
			}
		}
		pl := NewPostList(ph.basedir)
		pl.Newest(num) // ignore fs errors here
		pageData = check4lang(pageData, aRequest).
			Add("Robots", "noindex,follow").
			Add("Postings", pl.Sort())
		ph.viewList.Render("index", aWriter, pageData)

	case "p": // handle a single posting
		if 0 < len(tail) {
			p := newPosting(ph.basedir, tail)
			if err := p.Load(); nil == err {
				pageData = check4lang(pageData, aRequest).
					Add("Postings", p)
				ph.viewList.Render("article", aWriter, pageData)
				return
			}
		}
		http.NotFound(aWriter, aRequest)

	case "privacy", "privacy.html":
		ph.viewList.Render("privacy", aWriter, check4lang(pageData, aRequest))

	case "q":
		http.Redirect(aWriter, aRequest, "/s/"+tail, http.StatusSeeOther)

	case "r", "rp": // remove posting
		if 0 < len(tail) {
			p := newPosting(ph.basedir, tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Add("Manuscript", template.HTML(txt)).
					Add("BackURL", aRequest.URL.Path). // for POSTing back
					Add("Robots", "noindex,nofollow")
				ph.viewList.Render("rp", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "s": // handle a query/search
		if 0 < len(tail) {
			pl := SearchPostings(ph.basedir, regexp.QuoteMeta(tail))
			pageData = check4lang(pageData, aRequest).
				Add("Robots", "noindex,follow").
				Add("Matches", pl.Len()).
				Add("Postings", pl.Sort())
			ph.viewList.Render("searchresult", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case ph.staticD:
		if 0 < len(tail) {
			aRequest.URL.Path = tail
		}
		ph.staticH.ServeHTTP(aWriter, aRequest)

	case "w": // handle a given week
		if 0 < len(tail) {
			y, m, d := getYMD(tail)
			pl := NewPostList(ph.basedir).Week(y, m, d)
			pageData = check4lang(pageData, aRequest).
				Add("Robots", "noindex,follow").
				Add("Matches", pl.Len()).
				Add("Postings", pl.Sort())
			ph.viewList.Render("searchresult", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "":
		if m := aRequest.FormValue("m"); 0 < len(m) {
			http.Redirect(aWriter, aRequest, "/m/"+m, http.StatusSeeOther)
		} else if n := aRequest.FormValue("n"); 0 < len(n) {
			http.Redirect(aWriter, aRequest, "/n/"+n, http.StatusSeeOther)
		} else if p := aRequest.FormValue("p"); 0 < len(p) {
			http.Redirect(aWriter, aRequest, "/p/"+p, http.StatusSeeOther)
		} else if q := aRequest.FormValue("q"); 0 < len(q) {
			http.Redirect(aWriter, aRequest, "/s/"+q, http.StatusSeeOther)
		} else if s := aRequest.FormValue("s"); 0 < len(s) {
			http.Redirect(aWriter, aRequest, "/s/"+s, http.StatusSeeOther)
		} else if w := aRequest.FormValue("w"); 0 < len(w) {
			http.Redirect(aWriter, aRequest, "/w/"+w, http.StatusSeeOther)
		} else {
			http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)
		}
	default:
		// if nothing matched (above) reply to the request
		// with an HTTP 404 not found error.
		http.NotFound(aWriter, aRequest)
	} // switch
} // handleGET()

// handlePOST() process the HTTP POST requests.
func (ph *TPageHandler) handlePOST(aWriter http.ResponseWriter, aRequest *http.Request) {
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case "a", "ap": // add a new post
		m := replCRLF([]byte(aRequest.FormValue("manuscript")))
		if 0 < len(m) {
			p := NewPosting(ph.basedir)
			p.Set(m)
			if _, err := p.Store(); nil != err {
				err = nil //TODO better error handling
			}
			http.Redirect(aWriter, aRequest, "/p/"+p.ID(), http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "d", "dp": // change date
		if 0 < len(tail) {
			op := newPosting(ph.basedir, tail)
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
			np := newPosting(ph.basedir, newID(t))
			npn := np.PathFileName()
			if err := os.Rename(opn, npn); nil != err {
				log.Printf("handlePost(d): %v\n", err)
				//TODO better error handling
			}
			tail = strings.TrimPrefix(npn, ph.basedir+"/")
			// remove leading directory and trailing extension:
			tail = tail[4:len(tail)-3] +
				fmt.Sprintf("?z=%d%02d%02d%02d%02d%02d%04d", y, mo, d, h, mi, s, n)
			// dummy CGI argument to confuse the browser chache
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "e", "ep": // edit posting
		if 0 < len(tail) {
			var old []byte
			m := replCRLF([]byte(aRequest.FormValue("manuscript")))
			p := newPosting(ph.basedir, tail)
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
			http.Redirect(aWriter, aRequest, "/p/"+tail+"?z="+p.Date(), http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "r", "rp": // remove posting
		if 0 < len(tail) {
			p := newPosting(ph.basedir, tail)
			tail = p.Date()
			if err := p.Delete(); nil != err {
				log.Printf("handlePost(r): %v\n", err)
				//TODO better error handling
			}
			http.Redirect(aWriter, aRequest, "/m/"+tail+"?z="+tail, http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	default:
		// if nothing matched (above) reply to the request
		// with an HTTP 404 not found error.
		http.NotFound(aWriter, aRequest)
	}
} // handlePOST()

// Len returns the lenght of the internal view list.
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
		"r", "rp": // remove post
		return true
	}

	return false
} // NeedAuthentication()

// ServeHTTP handles the incoming HTTP requests.
func (ph TPageHandler) ServeHTTP(aWriter http.ResponseWriter, aRequest *http.Request) {
	if (nil != ph.ul) && ph.NeedAuthentication(aRequest) {
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

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// NewPageHandler returns a new `TPageHandler` instance.
func NewPageHandler() (*TPageHandler, error) {
	var (
		err error
		s   string
	)
	result := new(TPageHandler)

	if s, err = AppArguments.Get("css"); nil != err {
		return nil, err
	}
	result.cssD = filepath.Base(s)
	result.cssH = http.FileServer(http.Dir(s))

	if s, err = AppArguments.Get("img"); nil != err {
		return nil, err
	}
	result.imgD = filepath.Base(s)
	result.imgH = http.FileServer(http.Dir(s))

	if s, err = AppArguments.Get("js"); nil != err {
		return nil, err
	}
	result.jsD = filepath.Base(s)
	result.jsH = http.FileServer(http.Dir(s))

	if s, err = AppArguments.Get("lang"); nil == err {
		result.lang = s
	}

	if s, err = AppArguments.Get("listen"); nil != err {
		return nil, err
	}
	result.addr = s

	if s, err = AppArguments.Get("port"); nil != err {
		return nil, err
	}
	result.addr += ":" + s

	if s, err = AppArguments.Get("uf"); nil == err {
		if result.ul, err = passlist.LoadPasswords(s); nil != err {
			log.Printf("NewPageHandler(): %v\nAUTHENTICATION DISABLED!", err)
			result.ul = nil
		}
	}

	if s, err = AppArguments.Get("postdir"); nil != err {
		return nil, err
	}
	result.basedir = s

	if s, err = AppArguments.Get("realm"); nil == err {
		result.realm = s
	}

	if s, err = AppArguments.Get("static"); nil != err {
		return nil, err
	}
	result.staticD = filepath.Base(s)
	result.staticH = http.FileServer(http.Dir(s))

	if s, err = AppArguments.Get("tpldir"); nil != err {
		return nil, err
	}
	if result.viewList, err = newViewList(s); nil != err {
		return nil, err
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

/* _EoF_ */
