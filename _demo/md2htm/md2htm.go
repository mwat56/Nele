/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
               All rights reserved
           EMail : <support@mwat.de>
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"strings"

	"github.com/mwat56/go-blog"
)

type (
	// Container holding the data to insert into the pageTemplate.
	tData struct {
		Title string        // page title (from cmdline)
		Page  template.HTML // converted HTML data
	}
)

const (
	// The overall structure of the resulting HTML page.
	pageTemplate = `<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta http-equiv="Window-target" content="_top">
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<title>{{.Title}}</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes">
	<style type="text/css">
	html,body{font-size:101%;font-family:serif,monospace;}
	code,kbd,tt,pre{background-color:#f0f0f0;font-family:monospace;font-size:0.987em;letter-spacing:normal;line-height:1.123;white-space:normal;}
	pre,xmp{border-left:thin solid #999;display:block;font-size:88%;font-weight:inherit;line-height:1.36;margin:1ex 1ex 1ex 0.3ex;padding:0.5ex 1ex;overflow:visible;text-align:left;text-indent:0;white-space:pre-wrap;width:98%;}
	small{color:inherit;font-size:79%;font-weight:inherit;}
	</style>
</head><body>
	{{.Page}}
</body>
</html>
`
)

/*
// dtString() returns a string representation of aTime.
func dtString(aTime time.Time) string {
	y, mo, d := aTime.Date()
	h, mi, s := aTime.Clock()
	zn, ofs := aTime.Zone()
	sgn := "+"
	if 0 > ofs {
		ofs = -ofs
		sgn = "-"
	}
	ofs /= 60             // seconds to minutes
	zh := ofs / 60        // minutes to hours
	zm := ofs - (60 * zh) // remaining minutes

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d %s%02d%02d %s",
		y, mo, d, h, mi, s, sgn, zh, zm, zn)
} // dtString()
*/

// ShowHelp lists the usage information to `Stderr`.
func showHelp() {
	fmt.Fprintf(os.Stderr,
		"Usage: %s 'a page title text' < infile.md > outfile.html\n\n",
		os.Args[0])
} // ShowHelp()

// Run reads the contents from StdIn, insert the data into the
// pageTemplate and writes the resulting HTML page to StdOut.
func Run() error {
	flag.Usage = showHelp
	data, err := bufio.NewReader(os.Stdin).ReadBytes(0x03)
	if (nil != err) && (io.EOF != err) {
		return err
	}
	t := template.New("HTMLpage")
	if t, err = t.Parse(pageTemplate); nil != err {
		return err
	}

	err = t.Execute(os.Stdout,
		tData{
			strings.Join(os.Args[1:], " "),
			template.HTML(blog.RemoveWhiteSpace(blog.MDtoHTML(data))),
		},
	)

	return err
} // Run()

func main() {
	err := Run()
	if nil != err {
		showHelp()
		panic(err)
	}
} // main()
