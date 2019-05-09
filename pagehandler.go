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
	"time"

	passlist "github.com/mwat56/go-passlist"
)

type (
	// TPageHandler provides the handling of HTTP request/response.
	TPageHandler struct {
		addr     string              // listen address ("1.2.3.4:5678")
		fileH    http.Handler        // static file handler
		lang     string              // default language
		realm    string              // host/domain to secure by BasicAuth
		theme    string              // `dark` or `light` display theme
		ul       *passlist.TPassList // user/password list
		viewList *TViewList          // list of template/views
	}
)

// Address returns the configured `IP:Port` address to use for listening.
func (ph *TPageHandler) Address() string {
	return ph.addr
} // Address()

// `basicPageData()` returns a list of common Head entries.
func (ph *TPageHandler) basicPageData() *TDataList {
	css := fmt.Sprintf(`<link rel="stylesheet" type="text/css" title="mwat's styles" href="/css/stylesheet.css" /><link rel="stylesheet" type="text/css" href="/css/%s.css" />`, ph.theme)
	pageData := NewDataList().
		Set("CSS", template.HTML(css)).
		Set("Lang", ph.lang).
		Set("Robots", "index,follow")

	return pageData
} // basicPageData()

// `check4lang()` looks for a CGI value of `lang` and adds it to `aPD` if found.
func check4lang(aPD *TDataList, aRequest *http.Request) *TDataList {
	if l := aRequest.FormValue("lang"); 0 < len(l) {
		return aPD.Set("Lang", l)
	}
	return aPD
} // check4lang()

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
	switch path {
	case "a", "ap": // add a new post
		pageData = check4lang(pageData, aRequest).
			Set("BackURL", aRequest.URL.Path). // for POSTing back
			Set("Robots", "noindex,nofollow")
		ph.viewList.Render("ap", aWriter, pageData)

	case "css":
		ph.fileH.ServeHTTP(aWriter, aRequest)

	case "d", "dp": // change date
		if 0 < len(tail) {
			y, mo, d := time.Now().Date()
			now := fmt.Sprintf("%d-%02d-%02d", y, mo, d)
			t := timeID(tail)
			y, mo, d = t.Date()
			pageData = check4lang(pageData, aRequest).
				Set("BackURL", aRequest.URL.Path). // for POSTing back
				Set("NOW", now).
				Set("HMS", fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())).
				Set("YMD", fmt.Sprintf("%d-%02d-%02d", y, mo, d)).
				Set("Robots", "noindex,nofollow")
			ph.viewList.Render("dc", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "e", "ep": // edit a single posting
		if 0 < len(tail) {
			p := newPosting(tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Set("Manuscript", template.HTML(txt)).
					Set("BackURL", aRequest.URL.Path). // for POSTing back
					Set("Robots", "noindex,nofollow")
				ph.viewList.Render("ed", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "faq", "faq.html":
		ph.viewList.Render("faq", aWriter, check4lang(pageData, aRequest))

	case "img", "favicon.ico":
		ph.fileH.ServeHTTP(aWriter, aRequest)

	case "imprint", "imprint.html":
		ph.viewList.Render("imprint", aWriter, check4lang(pageData, aRequest))

	case "index", "index.html":
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "js":
		ph.fileH.ServeHTTP(aWriter, aRequest)

	case "licence", "license":
		ph.viewList.Render("licence", aWriter, pageData)

	case "m": // handle a given month
		if 0 < len(tail) {
			y, m, _ := getYMD(tail)
			pl := NewPostList().Month(y, m)
			pageData = check4lang(pageData, aRequest).
				Set("Robots", "noindex,follow").
				Set("Matches", pl.Len()).
				Set("Postings", pl.Sort())
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
		pl := NewPostList()
		pl.Newest(num) // ignore fs errors here
		pageData = check4lang(pageData, aRequest).
			Set("Robots", "noindex,follow").
			Set("Postings", pl.Sort())
		ph.viewList.Render("index", aWriter, pageData)

	case "p": // handle a single posting
		if 0 < len(tail) {
			p := newPosting(tail)
			if err := p.Load(); nil == err {
				pageData = check4lang(pageData, aRequest).
					Set("Postings", p)
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
			p := newPosting(tail)
			txt := p.Markdown()
			if 0 < len(txt) {
				pageData = check4lang(pageData, aRequest).
					Set("Manuscript", template.HTML(txt)).
					Set("BackURL", aRequest.URL.Path). // for POSTing back
					Set("Robots", "noindex,nofollow")
				ph.viewList.Render("rp", aWriter, pageData)
				return
			}
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "s": // handle a query/search
		if 0 < len(tail) {
			pl := SearchPostings(postingBaseDirectory, regexp.QuoteMeta(tail))
			pageData = check4lang(pageData, aRequest).
				Set("Robots", "noindex,follow").
				Set("Matches", pl.Len()).
				Set("Postings", pl.Sort())
			ph.viewList.Render("searchresult", aWriter, pageData)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "static":
		ph.fileH.ServeHTTP(aWriter, aRequest)

	case "w": // handle a given week
		if 0 < len(tail) {
			y, m, d := getYMD(tail)
			pl := NewPostList().Week(y, m, d)
			pageData = check4lang(pageData, aRequest).
				Set("Robots", "noindex,follow").
				Set("Matches", pl.Len()).
				Set("Postings", pl.Sort())
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

// `handlePOST()` process the HTTP POST requests.
func (ph *TPageHandler) handlePOST(aWriter http.ResponseWriter, aRequest *http.Request) {
	path, tail := URLparts(aRequest.URL.Path)
	switch path {
	case "a", "ap": // add a new post
		m := replCRLF([]byte(aRequest.FormValue("manuscript")))
		if 0 < len(m) {
			p := NewPosting()
			p.Set(m)
			if _, err := p.Store(); nil != err {
				log.Printf("handlePOST(a): %v\n", err)
				//TODO better error handling
			}
			tail = p.ID() + "?z=" + p.Date()
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "d", "dp": // change date
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
			if err := os.Rename(opn, npn); nil != err {
				log.Printf("handlePOST(d): %v\n", err)
				//TODO better error handling
			}
			tail = np.ID() + fmt.Sprintf("?z=%d%02d%02d%02d%02d%02d%04d", y, mo, d, h, mi, s, n)
			// dummy CGI argument to confuse the browser chache
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "e", "ep": // edit posting
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
			tail += "?z=" + p.ID()
			http.Redirect(aWriter, aRequest, "/p/"+tail, http.StatusSeeOther)
			return
		}
		http.Redirect(aWriter, aRequest, "/n/", http.StatusSeeOther)

	case "r", "rp": // remove posting
		if 0 < len(tail) {
			p := newPosting(tail)
			if err := p.Delete(); nil != err {
				log.Printf("handlePOST(r): %v\n", err)
				//TODO better error handling
			}
			tail = p.Date() + "?z=" + p.ID()
			http.Redirect(aWriter, aRequest, "/m/"+tail, http.StatusSeeOther)
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

	if s, err = AppArguments.Get("datadir"); nil != err {
		return nil, err
	}
	result.fileH = http.FileServer(http.Dir(s + "/"))
	if result.viewList, err = newViewList(s + "/views"); nil != err {
		return nil, err
	}

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

	if s, err = AppArguments.Get("realm"); nil == err {
		result.realm = s
	}

	if s, err = AppArguments.Get("theme"); nil != err {
		return nil, err
	}
	result.theme = s

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
