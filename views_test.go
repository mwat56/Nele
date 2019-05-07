package blog

import (
	"html/template"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewView(t *testing.T) {
	type args struct {
		aBaseDir string
		aName    string
	}
	tests := []struct {
		name     string
		args     args
		wantView bool
		wantErr  bool
	}{
		// TODO: Add test cases.
		{" 1", args{"./views/", "test1"}, false, true},
		{" 2", args{"./views/", "index"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewView(tt.args.aBaseDir, tt.args.aName)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewView() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantView {
				t.Errorf("NewView() = %v, want %v", got, tt.wantView)
			}
		})
	}
} // TestNewView()

func TestNewDataList(t *testing.T) {
	dl := TDataList{}
	tests := []struct {
		name string
		want *TDataList
	}{
		// TODO: Add test cases.
		{" 1", &dl},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDataList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataList() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestNewDataList()

func TestTDataList_Set(t *testing.T) {
	dl := NewDataList()
	rl := NewDataList()
	(*rl)["Title"] = "Testing"
	type args struct {
		aKey   string
		aValue interface{}
	}
	tests := []struct {
		name string
		d    *TDataList
		args args
		want *TDataList
	}{
		// TODO: Add test cases.
		{" 1", dl, args{"Title", "Testing"}, rl},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Set(tt.args.aKey, tt.args.aValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TDataList.Add() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTDataList_Set()

func TestTView_render(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"

	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := newPosting( /* bd, */ id1).
		Set([]byte("View_render: Oh dear! This is a first posting."))

	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := newPosting( /* bd, */ id2).
		Set([]byte("View_render: Hi there! This is another posting."))

	id3 := newID(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p3 := newPosting( /* bd,  */ id3).
		Set([]byte("View_render: Oh dear! This is a single posting."))

	v1, _ := NewView("./views/", "index")
	pl1 := NewPostList( /* bd */ ).Add(p1).Add(p2).Sort()
	dl1 := NewDataList().
		Set("Title", "this is the title").
		Set("Headline", "This is an interesting issue").
		Set("Postings", pl1)

	v2, _ := NewView("./views/", "article")
	dl2 := NewDataList().
		Set("Title", "this is the article title").
		Set("Headline", "Tis is an important topic").
		Set("Lang", "en").
		Set("Postings", p3).
		Set("ToBeIgnored", "! Ignore Me !")

	type fields struct {
		name string
		tpl  *template.Template
	}
	type args struct {
		aWriter io.Writer
		aData   *TDataList
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields(*v1), args{os.Stdout, dl1}, false},
		{" 2", fields(*v2), args{os.Stdout, dl2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &TView{
				name: tt.fields.name,
				tpl:  tt.fields.tpl,
			}
			if err := v.render(tt.args.aWriter, tt.args.aData); (err != nil) != tt.wantErr {
				t.Errorf("TView.render() error = %v,\nwantErr %v", err, tt.wantErr)
				return
			}
		})
	}
} // TestTView_render()

func TestNewViewList(t *testing.T) {
	vl := TViewList{}
	tests := []struct {
		name string
		want *TViewList
	}{
		// TODO: Add test cases.
		{" 1", &vl},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewViewList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewViewList() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestNewViewList()

func TestTViewList_Add(t *testing.T) {
	vname1 := "index"
	vw1, _ := NewView("./views/", vname1)
	vl1 := NewViewList()
	rl1 := NewViewList().Add(vw1)
	type args struct {
		aName string
		aView *TView
	}
	tests := []struct {
		name string
		vl   *TViewList
		args args
		want *TViewList
	}{
		// TODO: Add test cases.
		{" 1", vl1, args{vname1, vw1}, rl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.vl.Add(tt.args.aView); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TViewList.Add() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTViewList_Add()

func TestTViewList_render(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	vname1, vname2 := "index", "article"
	vw1, _ := NewView("./views/", vname1)
	vw2, _ := NewView("./views/", vname2)
	vl1 := NewViewList().
		Add(vw1).
		Add(vw2)

	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := newPosting( /* bd,  */ id1).
		Set([]byte("ViewList_render: Oh dear! This is a first posting."))

	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := newPosting( /* bd,  */ id2).
		Set([]byte("ViewList_render: Hi there! This is another posting."))

	id3 := newID(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p3 := newPosting( /* bd, */ id3).
		Set([]byte("ViewList_render: Oh dear! This is a single posting."))

	pl1 := NewPostList( /* bd */ ).
		Add(p1).
		Add(p2)
	dl1 := NewDataList().
		Set("Lang", "en").
		Set("Title", "this is the index title").
		Set("Postings", *pl1)
	pl2 := NewPostList( /* bd */ ).
		Add(p3)
	dl2 := NewDataList().
		Set("Lang", "en").
		Set("Title", "this is the article title").
		Set("Postings", *pl2)
	type args struct {
		aName string
		aData *TDataList
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
