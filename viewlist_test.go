/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"io"
	"os"
	"testing"
	"time"
)

func Test_NewViewList(t *testing.T) {
	prep4Tests()
	vl1 := make(TViewList, 16)

	tests := []struct {
		name string
		want *TViewList
	}{
		{"1", &vl1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewViewList(); !got.equals(tt.want) {
				t.Errorf("NewViewList() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_NewViewList()

func Test_TViewList_Add(t *testing.T) {
	prep4Tests()

	vw1, _ := NewView("./views/", "index")

	vl1 := NewViewList()
	rl1 := NewViewList().Add(vw1)

	tests := []struct {
		name string
		vl   *TViewList
		view *TView
		want *TViewList
	}{
		// TODO: Add test cases.
		{" 1", vl1, vw1, rl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vl.Add(tt.view); !got.equals(tt.want) {
				t.Errorf("TViewList.Add() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_TViewList_Add()

func Test_TViewList_equals(t *testing.T) {
	prep4Tests()

	vl0 := NewViewList()
	tv1, _ := NewView("./views/", "index")
	vl1 := NewViewList().Add(tv1)

	tv2, _ := NewView("./views/", "404")
	vl2 := NewViewList().Add(tv2)

	tests := []struct {
		name  string
		mList *TViewList
		oList *TViewList
		want  bool
	}{
		{"0", vl0, vl0, true},
		{"1", vl0, vl1, false},
		{"2", vl1, vl2, false},
		{"3", vl2, vl2, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mList.equals(tt.oList); got != tt.want {
				t.Errorf("TViewList.equals() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_TViewList_equals()

func TestTViewList_render(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	vname1, vname2 := "index", "article"
	vw1, _ := NewView("./views/", vname1)
	vw2, _ := NewView("./views/", vname2)
	vl1 := NewViewList().
		Add(vw1).
		Add(vw2)

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1, "ViewList_render: Oh dear! This is a first posting.")

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "ViewList_render: Hi there! This is another posting.")

	id3 := time2id(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p3 := NewPosting(id3, "ViewList_render: Oh dear! This is a single posting.")

	pl1 := NewPostList().
		Add(p1).
		Add(p2)
	dl1 := NewTemplateData().
		Set("Lang", "en").
		Set("Title", "this is the index title").
		Set("Postings", *pl1)
	pl2 := NewPostList().
		Add(p3)
	dl2 := NewTemplateData().
		Set("Lang", "en").
		Set("Title", "this is the article title").
		Set("Postings", *pl2)

	type args struct {
		aName string
		aData *TemplateData
	}
	tests := []struct {
		name    string
		vl      *TViewList
		args    args
		aWriter io.Writer
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", vl1, args{vname1, dl1}, os.Stdout, false},
		{" 2", vl1, args{vname2, dl2}, os.Stdout, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.vl.render(tt.args.aName, tt.aWriter, tt.args.aData); (err != nil) != tt.wantErr {
				t.Errorf("TViewList.render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
} // TestTViewList_render()

/* _EoF_ */
