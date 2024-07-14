/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/mwat56/whitespace"
)

type (
	// `TView` combines a template and its logical name.
	TView struct {
		// The view's symbolic name.
		vName string

		// The template as returned by a `NewView()` function call.
		vTpl *template.Template
	}

	// Internal type to track changes in certain template vars.
	tDataChange struct {
		current string
	}
)

// --------------------------------------------------------------------------
// constructor functions:

// `newChange()` returns a new change structure.
func newChange() *tDataChange {
	return &tDataChange{
		`{{$}}`, // ensure that first change is recognised
	}
} // newChange()

// `NewView()` returns a new `TView` with `aName`.
//
// `aName` serves as both the main template's name as well as the
// view's name; it's given here without the filename extension (i.e.
// w/o `.gohtml`).
//
// Parameters:
//
//	`aBaseDir` is the path to the directory storing the template files.
//	- `aName` is the name of the template file providing the page's main body.
func NewView(aBaseDir, aName string) (*TView, error) {
	var (
		bd    string
		err   error
		files []string
		tpl   *template.Template
	)

	if bd, err = filepath.Abs(aBaseDir); nil != err {
		return nil, err
	}

	if files, err = filepath.Glob(bd + "/layout/*.gohtml"); nil != err {
		return nil, err
	}

	files = append(files, bd+`/`+aName+`.gohtml`)

	if tpl, err = template.New(aName).
		Funcs(viewFunctionMap).
		ParseFiles(files...); nil != err {
		return nil, err
	}

	return &TView{
		vName: aName,
		vTpl:  tpl,
	}, nil
} // NewView()

// --------------------------------------------------------------------------
// helper data and functions

const (
	// replacement text for `reHrefRE`
	reHrefReplace = ` target="_extern" $1`
)

var (
	// RegEx to HREF= tag attributes
	reHrefRE = regexp.MustCompile(` (href="http)`)

	viewFunctionMap = template.FuncMap{
		"change":   newChange, // a new change structure
		"htmlSafe": htmlSafe,  // returns `aText` as template.HTML
	}
)

// `addExternURLtargets()` adds a TARGET attribute to HREFs.
func addExternURLtargets(aPage []byte) []byte {
	return reHrefRE.ReplaceAll(aPage, []byte(reHrefReplace))
} // addExternURLtargets()

// `htmlSafe()` returns `aText` as template.HTML.
//
// _Note_ that this is just a typecast without any tests.
func htmlSafe(aText string) template.HTML {
	return template.HTML(aText) // #nosec G203
} // htmlSafe()

// --------------------------------------------------------------------------
// tDataChange methods

// Changed returns whether `aNext` is the same as the last value.
func (c *tDataChange) Changed(aNext string) bool {
	if c.current == aNext {
		return false
	}
	c.current = aNext

	return true
} // Changed()

// --------------------------------------------------------------------------
// TView methods

// `render()` is the core of `Render()` with a slightly different API
// (`io.Writer` instead of `http.ResponseWriter`) for easier testing.
func (v *TView) render(aWriter io.Writer, aData *TemplateData) (rErr error) {
	var page []byte

	if page, rErr = v.RenderedPage(aData); nil != rErr {
		return
	}
	_, rErr = aWriter.Write(addExternURLtargets(whitespace.Remove(page)))

	return
} // render()

// Render executes the template using the TView's properties.
//
// `aWriter` is a http.ResponseWriter, or e.g. `os.Stdout` in console apps.
//
// `aData` is a list of data to be injected into the template.
//
// If an error occurs executing the template or writing its output,
// execution stops, and the method returns without writing anything
// to the output `aWriter`.
func (v *TView) Render(aWriter http.ResponseWriter, aData *TemplateData) error {
	return v.render(aWriter, aData)
} // Render()

// RenderedPage returns the rendered template/page and a possible Error
// executing the template.
//
// `aData` is a list of data to be injected into the template.
func (v *TView) RenderedPage(aData *TemplateData) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := v.vTpl.ExecuteTemplate(buf, v.vName, aData); nil != err {
		return nil, err
	}

	return buf.Bytes(), nil
} // RenderedPage()

/* _EoF_ */
