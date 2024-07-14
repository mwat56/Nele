/*
   Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
)

func Test_NewTemplateData(t *testing.T) {
	td1 := &TemplateData{}

	tests := []struct {
		name string
		want *TemplateData
	}{
		{"1", td1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTemplateData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemplateData() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // Test_NewTemplateData()

func Test_TemplateData_Get(t *testing.T) {
	td1 := NewTemplateData()
	(*td1)[`key1`] = `val1`
	(*td1)[`key3`] = false
	(*td1)[`key4`] = 123.456

	tests := []struct {
		name       string
		dl         TemplateData
		key        string
		wantRValue any
		wantROK    bool
	}{
		{"1", *td1, `key1`, `val1`, true},
		{"2", *td1, `key2`, `val2`, false},
		{"3", *td1, `key3`, false, true},
		{"4", *td1, `key4`, 123.456, true},
		{"5", *td1, `key5`, false, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRValue, gotROK := tt.dl.Get(tt.key)
			if gotROK != tt.wantROK {
				t.Errorf("TemplateData.Get() gotROK = %v, want %v", gotROK, tt.wantROK)
				return
			}
			if gotROK && !reflect.DeepEqual(gotRValue, tt.wantRValue) {
				t.Errorf("TemplateData.Get() gotRValue = %v,\nwant %v", gotRValue, tt.wantRValue)
			}
		})
	}
} // Test_TemplateData_Get()

func Test_TemplateData_Set(t *testing.T) {
	td1 := NewTemplateData()
	wd1 := NewTemplateData()
	(*wd1)["Title"] = "Testing"

	type tArgs struct {
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name string
		d    *TemplateData
		args tArgs
		want *TemplateData
	}{
		{" 1", td1, tArgs{"Title", "Testing"}, wd1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Set(tt.args.aKey, tt.args.aValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateData.Add() = %v, want\n%v", got, tt.want)
			}
		})
	}
} // Test_TemplateData_Set()

/* _EoF_ */
