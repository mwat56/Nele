/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"compress/zlib"
	"html/template"
	"os"
	"time"
)

/*

intercept a `posting`'s HTML generation

*/

// `cachedHTML()` returns the HTML generated from `aPost` Markdown.
//
// This functions maintains a cache containing the pre-prepared HTML
// markup in compressed files to avoid converting the same Markdown
// over and over again.
// However, it should be noted that the compression improves the longer
// the text is.
// HTML generated from Markdown shorter than approximately 100 bytes
// might even grow due to the compression.
//
// `aPost` is the posting whose Markdown is used.
func cachedHTML(aPost *TPosting) template.HTML {
	var (
		err    error
		fi     os.FileInfo
		file   *os.File
		htTime time.Time
		txt    []byte
	)
	mdName := aPost.PathFileName()
	if fi, err = os.Stat(mdName); nil != err {
		// return empty result
		return template.HTML(txt)
	}
	mdTime := fi.ModTime()

	htName := mdName[:len(mdName)-2] + "ht"
	if fi, err = os.Stat(htName); nil != err {
		htTime = mdTime.AddDate(-1, -1, -1)
	} else {
		htTime = fi.ModTime()
	}
	if htTime.After(mdTime) {
		// read cached/compressed HTML file
		htSize := fi.Size()

		if file, err = os.OpenFile(htName, os.O_RDONLY, 0644); nil == err {
			defer file.Close()

			if reader, err := zlib.NewReader(file); nil == err {
				defer reader.Close()

				txt = make([]byte, htSize*4)
				if _, err = reader.Read(txt); nil == err {
					return template.HTML(txt)
				}
			}
		}
	}
	// Reaching this point of execution means that (1) there is no HTML
	// file available, or (2) it's older than the Markdown, or (3) the
	// HTML file couldn't be read without errors.
	// So, get the current Markdown and convert it to HTML:
	txt = MDtoHTML(aPost.Markdown())

	// Write the compressed HTML to a file so it's available the
	// next time this article is requested.
	go goCacheHTML(htName, txt)

	return template.HTML(txt)
} // cachedHTML()

func goCacheHTML(aFilename string, aText []byte) {
	// Now, write the compressed HTML to a file so it's available the
	// next time this article is requested.
	file, err := os.OpenFile(aFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if nil != err {
		os.Remove(aFilename) // what else could we do here?
	}
	defer file.Close()

	w, _ := zlib.NewWriterLevel(file, zlib.BestCompression)
	defer w.Close()

	w.Write(aText)
} // goCacheHTML()

/*
func cachedHTML2(aPost *TPosting) template.HTML {
	var (
		err    error
		fi     os.FileInfo
		htTime time.Time
		txt    []byte
	)
	mdName := aPost.PathFileName()
	htName := mdName[:len(mdName)-2] + "ht"
	if fi, err = os.Stat(mdName); nil != err {
		// return empty result
		return template.HTML(txt)
	}
	mdTime := fi.ModTime()

	if fi, err = os.Stat(htName); nil != err {
		htTime = mdTime.AddDate(-1, -1, -1)
	} else {
		htTime = fi.ModTime()
	}
	if htTime.After(mdTime) {
		// read cached/compressed HTML file
		var buf []byte
		if buf, err = ioutil.ReadFile(htName); nil == err {
			return template.HTML(buf)
		}
	}

	// get the current Markdown and convert it into HTML
	txt = MDtoHTML(aPost.Markdown())

	if err = ioutil.WriteFile(htName, txt, 0644); nil != err {
		os.Remove(htName) // what else could we do here?
	}

	return template.HTML(txt)
} // cachedHTML2()
*/

/* _EoF_ */
