/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides some template/view related functions and methods.
 */

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/mwat56/whitespace"
)

const (
	// replacement text for `reHrefRE`
	reHrefReplace = ` target="_extern" $1`
)

var (
	// RegEx to HREF= tag attributes
	reHrefRE = regexp.MustCompile(` (href="http)`)
)

// `addExternURLtargets()` adds a TARGET attribute to HREFs.
func addExternURLtargets(aPage []byte) []byte {
	return reHrefRE.ReplaceAll(aPage, []byte(reHrefReplace))
} // addExternURLtargets()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// Internal type to track changes in certain template vars.
	tDataChange struct {
		current string
	}
)

// Changed returns whether `aNext` is the same as the last value.
func (c *tDataChange) Changed(aNext string) bool {
	if c.current == aNext {
		return false
	}
	c.current = aNext

	return true
} // Changed()

// `newChange()` returns a new change structure.
func newChange() *tDataChange {
	return &tDataChange{
		`{{$}}`, // ensure that first change is recognised
	}
} // newChange()

// `htmlSafe()` returns `aText` as template.HTML.
//
// _Note_ that this is just a typecast without any tests.
func htmlSafe(aText string) template.HTML {
	return template.HTML(aText) // #nosec G203
} // htmlSafe()

var (
	viewFunctionMap = template.FuncMap{
		"change":   newChange, // a new change structure
		"htmlSafe": htmlSafe,  // returns `aText` as template.HTML
	}
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// TView combines a template and its logical name.
type TView struct {
	// The view's symbolic name.
	tvName string

	// The template as returned by a `NewView()` function call.
	tvTpl *template.Template
}

// NewView returns a new `TView` with `aName`.
//
//	`aBaseDir` is the path to the directory storing the template files.
//
//	`aName` is the name of the template file providing the page's main
// body without the filename extension (i.e. w/o `.gohtml`). `aName`
// serves as both the main template's name as well as the view's name.
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
		tvName: aName,
		tvTpl:  tpl,
	}, nil
} // NewView()

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

	if err := v.tvTpl.ExecuteTemplate(buf, v.tvName, aData); nil != err {
		return nil, err
	}

	return buf.Bytes(), nil
} // RenderedPage()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// Map indexed by a name pointing to a view instance.
	tViewList map[string]*TView

	// TViewList is a list of `TView` instances (to be used as a template pool).
	TViewList tViewList
)

// NewViewList returns a new (empty) `TViewList` instance.
func NewViewList() *TViewList {
	result := make(TViewList, 16)

	return &result
} // NewViewlist()

// Add appends `aView` to the list.
//
// `aView` is the view to add to this list.
//
// The view's name (as specified in the `NewView()` function call)
// is used as the view's key in this list.
func (vl *TViewList) Add(aView *TView) *TViewList {
	(*vl)[aView.tvName] = aView

	return vl
} // Add()

// Get returns the view with `aName`.
//
// `aName` is the name (key) of the `TView` object to retrieve.
//
// If `aName` doesn't exist, the return value is `nil`.
// The second value (ok) is a `bool` that is `true` if `aName`
// exists in the list, and `false` if not.
func (vl *TViewList) Get(aName string) (*TView, bool) {
	if result, ok := (*vl)[aName]; ok {
		return result, true
	}

	return nil, false
} // Get()

// `render()` is the core of `Render()` with a slightly different API
// (`io.Writer` instead of `http.ResponseWriter`) for easier testing.
func (vl *TViewList) render(aName string, aWriter io.Writer, aData *TemplateData) error {
	if view, ok := (*vl)[aName]; ok {
		return view.render(aWriter, aData)
	}

	return fmt.Errorf("template/view '%s' not found", aName)
} // render()

// Render executes the template with the key `aName`.
//
// `aName` is the name of the template/view to use.
//
// `aWriter` is a `http.ResponseWriter` to handle the executed template.
//
// `aData` is a list of data to be injected into the template.
//
// If an error occurs executing the template or writing its output,
// execution stops, and the method returns without writing anything
// to the output `aWriter`.
func (vl *TViewList) Render(aName string, aWriter http.ResponseWriter, aData *TemplateData) error {
	return vl.render(aName, aWriter, aData)
} // Render()

// RenderedPage returns the rendered template/page with the key `aName`.
//
// `aName` is the name of the template/view to use.
//
// `aData` is a list of data to be injected into the template.
func (vl *TViewList) RenderedPage(aName string, aData *TemplateData) ([]byte, error) {

	if view, ok := (*vl)[aName]; ok {
		return view.RenderedPage(aData)
	}

	return nil, fmt.Errorf("template/view '%s' not found", aName)
} // RenderedPage()

/* _EoF_ */
