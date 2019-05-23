/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

/*
 * This file provides some template/views related functions and methods.
 */

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

type (
	// TDataList is a list of values to be injected into a template.
	TDataList map[string]interface{}
)

// NewDataList returns a new (empty) TDataList instance.
func NewDataList() *TDataList {
	result := make(TDataList, 32)

	return &result
} // NewDatalist()

// Set inserts `aValue` identified by `aKey` to the list.
//
// If there's already a list entry with `aKey` its current value
// gets replaced by `aValue`.
//
// `aKey` is the values's identifier (as used as placeholder in the template).
//
// `aValue` contains the data entry's value.
func (dl *TDataList) Set(aKey string, aValue interface{}) *TDataList {
	(*dl)[aKey] = aValue

	return dl
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	// Internal type to track changes in certain template vars.
	tChange struct {
		current string
	}
)

// Changed returns whether `aNext` is the same as the last value.
func (c *tChange) Changed(aNext string) bool {
	if c.current == aNext {
		return false
	}
	c.current = aNext

	return true
} // Changed()

// `newChange()` returns a new change structure.
func newChange() *tChange {
	return &tChange{
		"{{$}}", // ensure that first change is recognised
	}
} // newChange()

// `htmlSafe()` returns `aText` as template.HTML.
func htmlSafe(aText string) template.HTML {
	return template.HTML(aText)
} // htmlSafe()

var (
	fMap = template.FuncMap{
		"change":   newChange, // a new change structure
		"htmlSafe": htmlSafe,  // returns `aText` as template.HTML
	}
)

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// TView combines a template and its logical name.
type TView struct {
	// The view's symbolic name.
	name string
	// The template as returned by a `NewView()` function call.
	tpl *template.Template
}

// NewView returns a new `TView` with `aName`.
//
// `aBaseDir` is the path to the directory storing the template files.
//
// `aName` is the name of the template file providing the page's main
// body without the filename extension (i.e. w/o ".gohtml"). `aName`
// serves as both the main template's name as well as the view's name.
func NewView(aBaseDir, aName string) (*TView, error) {
	bd, err := filepath.Abs(aBaseDir)
	if nil != err {
		return nil, err
	}
	files, err := filepath.Glob(fmt.Sprintf("%s/layout/*.gohtml", bd))
	if nil != err {
		return nil, err
	}
	files = append(files, fmt.Sprintf("%s/%s.gohtml", bd, aName))

	templ, err := template.New(aName).
		Funcs(fMap).
		ParseFiles(files...)
	if nil != err {
		return nil, err
	}

	return &TView{
		name: aName,
		tpl:  templ,
	}, nil
} // NewView()

// `render()` is the core of `Render()` with a slightly different API
// (`io.Writer` instead of `http.ResponseWriter`) for easier testing.
func (v *TView) render(aWriter io.Writer, aData *TDataList) (rErr error) {
	var page []byte

	if page, rErr = v.RenderedPage(aData); nil != rErr {
		return
	}

	// if _, rErr := aWriter.Write(page); nil != rErr {
	if _, rErr := aWriter.Write(RemoveWhiteSpace(page)); nil != rErr {
		return rErr
	}

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
func (v *TView) Render(aWriter http.ResponseWriter, aData *TDataList) error {
	return v.render(aWriter, aData)
} // Render()

// RenderedPage returns the rendered template/page and a possible Error
// executing the template.
//
// `aData` is a list of data to be injected into the template.
func (v *TView) RenderedPage(aData *TDataList) (rBytes []byte, rErr error) {
	buf := &bytes.Buffer{}

	if rErr = v.tpl.ExecuteTemplate(buf, v.name, aData); nil != rErr {
		return
	}

	return buf.Bytes(), nil
} // RenderedPage()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

type (
	tViewList map[string]*TView

	// TViewList is a list of `TView` instances (to be used as a template pool).
	TViewList tViewList
)

// NewViewList returns a new (empty) `TViewList` instance.
func NewViewList() *TViewList {
	result := make(TViewList, 8)

	return &result
} // NewViewlist()

// Add appends `aView` to the list.
//
// `aView` is the view to add to this list.
//
// The view's name (as specified in the `NewView()` function call)
// is used as the view's key in this list.
func (vl *TViewList) Add(aView *TView) *TViewList {
	(*vl)[aView.name] = aView

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
func (vl *TViewList) render(aName string, aWriter io.Writer, aData *TDataList) error {
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
func (vl *TViewList) Render(aName string, aWriter http.ResponseWriter, aData *TDataList) error {
	return vl.render(aName, aWriter, aData)
} // Render()

// RenderedPage returns the rendered template/page with the key `aName`.
//
// `aName` is the name of the template/view to use.
//
// `aData` is a list of data to be injected into the template.
func (vl *TViewList) RenderedPage(aName string, aData *TDataList) (rBytes []byte, rErr error) {

	if view, ok := (*vl)[aName]; ok {
		return view.RenderedPage(aData)
	}

	return rBytes, fmt.Errorf("template/view '%s' not found", aName)
} // RenderedPage()

/* _EoF_ */
