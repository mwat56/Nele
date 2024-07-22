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
			if got, _ := NewViewList(); got.equals(tt.want) {
				t.Errorf("%q: NewViewList() = \n%v\n>>> want >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_NewViewList()

func Test_TViewList_add(t *testing.T) {
	prep4Tests()

	vw1, _ := NewView("index")

	vl1, _ := NewViewList()
	rl1, _ := NewViewList()
	vl1.add(vw1)

	tests := []struct {
		name string
		vl   *TViewList
		view *TView
		want *TViewList
	}{
		{" 1", vl1, vw1, rl1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vl.add(tt.view); !got.equals(tt.want) {
				t.Errorf("TViewList.add() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_TViewList_add()

func Test_TViewList_equals(t *testing.T) {
	prep4Tests()

	vl0, _ := NewViewList()

	tv1, _ := NewView("index")
	vl1, _ := NewViewList()
	vl1.add(tv1)

	tv2, _ := NewView("503")
	vl2, _ := NewViewList()
	vl2.add(tv2)

	tests := []struct {
		name  string
		mList *TViewList
		oList *TViewList
		want  bool
	}{
		// since all viewlists contain all views,
		// all tests will return true ...
		{"0", vl0, vl0, true},
		{"1", vl0, vl1, true},
		{"2", vl0, vl2, true},
		{"3", vl2, vl2, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mList.equals(tt.oList); got != tt.want {
				t.Errorf("%q: TViewList.equals() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TViewList_equals()

func TestTViewList_render(t *testing.T) {
	prepareTestFiles()

	vname1, vname2 := "index", "article"
	vw1, _ := NewView(vname1)
	vw2, _ := NewView(vname2)
	vl1, _ := NewViewList()
	vl1.add(vw1).
		add(vw2)

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

	type tArgs struct {
		aName string
		aData *TemplateData
	}
	tests := []struct {
		name    string
		vl      *TViewList
		args    tArgs
		aWriter io.Writer
		wantErr bool
	}{
		{" 1", vl1, tArgs{vname1, dl1}, os.Stdout, false},
		{" 2", vl1, tArgs{vname2, dl2}, os.Stdout, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vl.render(tt.args.aName, tt.aWriter, tt.args.aData)
			if (err != nil) != tt.wantErr {
				t.Errorf("TViewList.render() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
		})
	}
} // TestTViewList_render()

/* _EoF_ */
