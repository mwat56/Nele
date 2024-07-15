/*
Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"fmt"
	"io"
	"net/http"
)

type (
	// `TViewList` is a list of `TView` instances (to be used as a
	// template pool).
	// It's a map indexed by a name pointing to a view instance.
	TViewList map[string]*TView
)

// --------------------------------------------------------------------------
// constructor function:

// `NewViewList()` returns a new (empty) `TViewList` instance.
func NewViewList() *TViewList {
	result := make(TViewList, 16)

	return &result
} // NewViewlist()

// --------------------------------------------------------------------------
// TViewList methods

// `Add()` appends `aView` to the list.
//
// The view's name (as specified in the `NewView()` function call)
// is used as the view's key in this list.
//
// Parameters:
//
//   - `aView` is the view to add to this list.
//
// Returns:
func (vl *TViewList) Add(aView *TView) *TViewList {
	(*vl)[aView.vName] = aView

	return vl
} // Add()

// `equals()` compares the current `TViewList` with another `TViewList` for
// equality. It checks if the symbolic names of both views are identical.
//
// Parameters:
//
//   - `aViewList`: The `TView` instance to compare with the current one.
//
// Returns:
//
//   - `bool`: `true` if the symbolic names of both viewlists are identical.
func (vl *TViewList) equals(aViewList *TViewList) bool {
	if nil == vl {
		return (nil == aViewList)
	}
	if len(*vl) != len(*aViewList) {
		return false
	}

	// Check if the values are equal for each key
	for key, myView := range *vl {
		otherView, ok := (*aViewList)[key]
		if !ok {
			return false
		}
		if !myView.equals(otherView) {
			return false
		}
	}

	return true
} // equals()

// `Get()` returns the view with `aName`.
//
// If `aName` doesn't exist, the return value is `nil`.
// The second value (ok) is a `bool` that is `true` if `aName`
// exists in the list, and `false` if not.
//
// Parameters:
//
//   - `aName` is the name (key) of the `TView` object to retrieve.
//
// Returns:
func (vl *TViewList) Get(aName string) (*TView, bool) {
	if result, ok := (*vl)[aName]; ok {
		return result, true
	}

	return nil, false
} // Get()

// `render()` is the core of `Render()` with a slightly different API
// (`io.Writer` instead of `http.ResponseWriter`) for easier testing.
//
// Parameters:
//
// Returns:
func (vl *TViewList) render(aName string, aWriter io.Writer, aData *TemplateData) error {
	if view, ok := (*vl)[aName]; ok {
		return view.render(aWriter, aData)
	}

	return fmt.Errorf("template/view '%s' not found", aName)
} // render()

// `Render()` executes the template with the key `aName`.
//
// If an error occurs executing the template or writing its output,
// execution stops, and the method returns without writing anything
// to the output `aWriter`.
//
// Parameters:
//
//   - `aName`: The name of the template/view to use.
//   - `aWriter`: A `http.ResponseWriter` to handle the executed template.
//   - `aData`: A list of data to be injected into the template.
//
// Returns:
func (vl *TViewList) Render(aName string, aWriter http.ResponseWriter, aData *TemplateData) error {
	return vl.render(aName, aWriter, aData)
} // Render()

// `RenderedPage()` returns the rendered template/page with the key `aName`.
//
// Parameters:
//
//   - `aName` is the name of the template/view to use.
//   - `aData` is a list of data to be injected into the template.
//
// Returns:
func (vl *TViewList) RenderedPage(aName string, aData *TemplateData) ([]byte, error) {

	if view, ok := (*vl)[aName]; ok {
		return view.RenderedPage(aData)
	}

	return nil, fmt.Errorf("template/view '%s' not found", aName)
} // RenderedPage()

/* _EoF_ */
