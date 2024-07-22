/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	se "github.com/mwat56/sourceerror"
	"github.com/mwat56/whitespace"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

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

//go:embed views/*
var viewsFS embed.FS

// --------------------------------------------------------------------------
// constructor functions:

// `NewView()` returns a new `TView` with `aName`.
//
// `aName` serves as both the main template's name as well as the
// view's name; it's given here without the filename extension (i.e.
// w/o `.gohtml`).
//
// Parameters:
//   - `aName`: The name of the template file providing the page's main body.
//
// Returns:
//   - `*TView`: A new `TView` instance.
//   - `error`: A possible error during processing.
func NewView(aName string) (*TView, error) {
	var (
		err   error
		fc    []byte // file contents
		fn    string // file name
		files []string
		tpl   *template.Template
	)

	// Get the files defining the overall page layout
	if files, err = filepath.Glob("views/layout/*.gohtml"); nil != err {
		return nil, se.Wrap(err, 1)
	}
	files = append(files, `views/`+aName+`.gohtml`)

	tpl = template.New(aName)
	for _, fn = range files {
		if fc, err = viewsFS.ReadFile(fn); nil != err {
			return nil, se.Wrap(err, 1)
		}
		if tpl, err = tpl.New(fn).
			Funcs(viewFunctionMap).
			Parse(string(fc)); nil != err {
			return nil, se.Wrap(err, 3)
		}
	}

	return &TView{
		vName: aName,
		vTpl:  tpl,
	}, nil
} // NewView()

// `newChange()` returns a new change structure.
func newChange() *tDataChange {
	return &tDataChange{
		`{{$}}`, // ensure that first change is recognised
	}
} // newChange()

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

// `equals()` compares the current `TView` with another `TView` for
// equality. It checks if the symbolic names of both views are identical.
//
// Parameters:
//
//   - `aView: The `TView` instance to compare with the current one.
//
// Returns:
//
//   - `bool`: `true` if the symbolic names of both views are identical.
func (v *TView) equals(aView *TView) bool {
	if nil == v {
		return (nil == aView)
	}
	if v.vName == aView.vName {
		return true
	}

	//TODO: check the embedded template

	return false
} // equals()

// `render()` is the core of `Render()` with a slightly different API
// (`io.Writer` instead of `http.ResponseWriter`) for easier testing.
//
// Parameters:
//   - `aName`: The view's name to render
//   - `aWriter`: A `http.ResponseWriter` to handle the executed template.
//   - `aData`: The data to be put into the view.
//
// Returns:
//   - `error`: A possible error during processing.
func (v *TView) render(aWriter io.Writer, aData *TemplateData) error {
	var (
		err  error
		page []byte
	)

	if page, err = v.RenderedPage(aData); nil != err {
		return err
	}

	if _, err = aWriter.Write(addExternURLtargets(whitespace.Remove(page))); nil != err {
		return se.Wrap(err, 1)
	}

	return nil
} // render()

// `Render()` executes the template using the TView's properties.
//
// `aWriter` is a http.ResponseWriter, or e.g. `os.Stdout` in console apps.
//
// If an error occurs executing the template or writing its output,
// execution stops, and the method returns without writing anything
// to the output `aWriter`.
//
// Parameters:
//   - `aWriter`: A `http.ResponseWriter` to handle the executed template.
//   - `aData`: A list of data to be injected into the template.
//
// Returns:
//   - `error`: A possible error during processing.
func (v *TView) Render(aWriter http.ResponseWriter, aData *TemplateData) error {
	return v.render(aWriter, aData)
} // Render()

// `RenderedPage()` returns the rendered template/page and a possible Error
// executing the template.
//
// Parameters:
// `aData` is a list of data to be injected into the template.
//
// Returns:
//   - `error`: A possible error during processing.
func (v *TView) RenderedPage(aData *TemplateData) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := v.vTpl.ExecuteTemplate(buf, v.vName, aData); nil != err {
		return nil, se.Wrap(err, 1)
	}

	return buf.Bytes(), nil
} // RenderedPage()

/* _EoF_ */
