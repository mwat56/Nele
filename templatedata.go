/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

type (
	// `TemplateData` is a list of values to be injected into a template.
	TemplateData map[string]any
)

// --------------------------------------------------------------------------
// constructor function:

// `NewTemplateData()` returns a new (empty) `TemplateData` instance.
//
// Returns:
//   - `*TemplateData`: The new data list.
func NewTemplateData() *TemplateData {
	result := make(TemplateData, 32)

	return &result
} // NewTemplateData()

// --------------------------------------------------------------------------
// TemplateData methods

// `Get()` returns the value associated with `aKey` and `true`.
//
// If `aKey` is not present in the list then the `bool` return
// value will be `false`.
//
// Parameters:
//   - `aKey`: The value's identifier (as used as placeholder in the template).
//
// Returns:
//   - `any`: The value associated with `aKey`.
//   - `bool`: Indicator whether `aKey` was found.
func (dl TemplateData) Get(aKey string) (rValue any, rOK bool) {
	rValue, rOK = dl[aKey]

	return
} // Get()

// `Set()` inserts `aValue` identified by `aKey` to the list.
//
// If there's already a list entry with `aKey` its current value
// gets replaced by `aValue`.
//
// Parameters:
//   - `aKey`: The value's identifier (as used as placeholder in the template).
//   - `aValue`: The data entry's value.
//
// Returns:
//   - `*TemplateData`: The updated data list.
func (dl *TemplateData) Set(aKey string, aValue any) *TemplateData {
	(*dl)[aKey] = aValue

	return dl
} // Set()

/* _EoF_ */
