/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
              All rights reserved
          EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
)

func TestNewTemplateData(t *testing.T) {
	d1 := &TemplateData{}
	tests := []struct {
		name string
		want *TemplateData
	}{
		// TODO: Add test cases.
		{" 1", d1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTemplateData(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTemplateData() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestNewTemplateData()

func TestTemplateData_Get(t *testing.T) {
	d1 := NewTemplateData()
	(*d1)[`key1`] = `val1`
	(*d1)[`key3`] = false
	(*d1)[`key4`] = true

	type args struct {
		aKey string
	}
	tests := []struct {
		name       string
		dl         TemplateData
		args       args
		wantRValue interface{}
		wantROK    bool
	}{
		// TODO: Add test cases.
		{" 1", *d1, args{`key1`}, `val1`, true},
		{" 2", *d1, args{`key2`}, `val2`, false},
		{" 3", *d1, args{`key3`}, false, true},
		{" 4", *d1, args{`key4`}, true, true},
		{" 5", *d1, args{`key5`}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRValue, gotROK := tt.dl.Get(tt.args.aKey)
			if gotROK != tt.wantROK {
				t.Errorf("TemplateData.Get() gotROK = %v, want %v", gotROK, tt.wantROK)
				return
			}
			if gotROK && !reflect.DeepEqual(gotRValue, tt.wantRValue) {
				t.Errorf("TemplateData.Get() gotRValue = %v,\nwant %v", gotRValue, tt.wantRValue)
			}
		})
	}
} // TestTemplateData_Get()

func TestTemplateData_Set(t *testing.T) {
	d1 := NewTemplateData()
	w1 := NewTemplateData()
	(*w1)["Title"] = "Testing"
	type args struct {
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name string
		d    *TemplateData
		args args
		want *TemplateData
	}{
		// TODO: Add test cases.
		{" 1", d1, args{"Title", "Testing"}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Set(tt.args.aKey, tt.args.aValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TemplateData.Add() = %v, want\n%v", got, tt.want)
			}
		})
	}
} // TestTemplateData_Set()
