/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"compress/zlib"
	"os"
	"path/filepath"
	"time"
)

/*
This file provides functions to maintain a simple cache for articles.

The original posting is stored in Markdown (i.e. text) which has on
each request converted to HTML before showing it to the remote user.
To minimise the need for Markdown/HTML conversion each generated HTML
is stored in compressed form to a separate file.
The next time the posting's markup is requested the prepared HTML is
returned instead of converting the Markdown again.

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
func cachedHTML(aPost *TPosting) []byte {
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
		return txt
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
					return txt
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

	return txt
} // cachedHTML()

// `goCacheCleanup()` is intended to be run in background and
// searches for stale and outdated cache files and deletes them.
func goCacheCleanup() {
	var (
		dirnames, filenames []string
		err                 error
		fi                  os.FileInfo
		mdTime              time.Time
	)
	if dirnames, err = filepath.Glob(postingBaseDirectory + "/*"); nil != err {
		return
	}
	for _, dirname := range dirnames {
		if filenames, err = filepath.Glob(dirname + "/*.ht"); nil != err {
			continue // it might be a file (not a directory) …
		}
		if 0 >= len(filenames) {
			continue // skip empty directory
		}
		for _, htName := range filenames {
			if fi, err = os.Stat(htName); nil != err {
				//TODO better error handling
				continue
			}
			htTime := fi.ModTime()

			mdName := htName[:len(htName)-2] + "md"
			if fi, err = os.Stat(mdName); nil != err {
				mdTime = htTime.AddDate(1, 1, 1)
			} else {
				mdTime = fi.ModTime()
			}
			if htTime.Before(mdTime) {
				os.Remove(htName)
			}
		}
	}
} // goCacheCleanup()

// `goCacheHTML()` is intended to be run in background and writes
// `aText` in compressed form to `aFilename`.
func goCacheHTML(aFilename string, aText []byte) {
	file, err := os.OpenFile(aFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if nil != err {
		os.Remove(aFilename) // what else could we do here?
	}
	defer file.Close()

	w, _ := zlib.NewWriterLevel(file, zlib.BestCompression)
	defer w.Close()

	w.Write(aText)
} // goCacheHTML()

/* _EoF_ */
