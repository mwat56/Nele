/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

/*
 * This file provides an object whose properties are to be inserted
 * into templates.
 */

type (
	// TemplateData is a list of values to be injected into a template.
	TemplateData map[string]interface{}
)

// Get returns the value associated with `aKey` and `true`.
// If `aKey` is not present in the list then the `bool` return
// value will be `false`.
//
//	`aKey` The value's identifier (as used as placeholder in the template).
func (dl TemplateData) Get(aKey string) (rValue interface{}, rOK bool) {
	rValue, rOK = dl[aKey]

	return
} // Get()

// Set inserts `aValue` identified by `aKey` to the list.
//
// If there's already a list entry with `aKey` its current value
// gets replaced by `aValue`.
//
//	`aKey` The value's identifier (as used as placeholder in the template).
//	`aValue` The data entry's value.
func (dl *TemplateData) Set(aKey string, aValue interface{}) *TemplateData {
	(*dl)[aKey] = aValue

	return dl
} // Set()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

// NewTemplateData returns a new (empty) `TDataList` instance.
func NewTemplateData() *TemplateData {
	result := make(TemplateData, 32)

	return &result
} // NewTemplateData()

/* _EoF_ */
