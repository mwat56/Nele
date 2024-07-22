/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func Test_addExternURLtargets(t *testing.T) {
	prepareTestFiles()

	t1 := ` bla <a href="https://site/page">bla</a> `
	p1 := []byte(t1)
	w1 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> `)

	t2 := t1 + `bla <a href="/page">bla</a>`
	p2 := []byte(t2)
	w2 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> bla <a href="/page">bla</a>`)

	t3 := t1 + `bla <a href="http://site.com/page">bla</a>`
	p3 := []byte(t3)
	w3 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> bla <a target="_extern" href="http://site.com/page">bla</a>`)

	tests := []struct {
		name string
		page []byte
		want []byte
	}{
		{" 1", p1, w1},
		{" 2", p2, w2},
		{" 3", p3, w3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addExternURLtargets(tt.page); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: addExternURLtargets() = \n%s\n>>> want >>>\n%s",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_addExternURLtargets()

func Test_NewView(t *testing.T) {
	prepareTestFiles()

	tests := []struct {
		name     string
		tplName  string
		wantView bool
		wantErr  bool
	}{
		{"1", "test1", false, true},
		{"2", "index", true, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewView(tt.tplName)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: NewView() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if (got != nil) != tt.wantView {
				t.Errorf("NewView() = %v, want %v", got, tt.wantView)
			}
		})
	}
} // Test_NewView()

func Test_TView_equals(t *testing.T) {
	prepareTestFiles()

	tv1, _ := NewView("index")
	tv2, _ := NewView("404")

	tests := []struct {
		name string
		tv   *TView
		view *TView
		want bool
	}{
		{"1", tv1, tv2, false},
		{"2", tv2, tv1, false},
		{"3", tv1, tv1, true},
		{"4", tv2, tv2, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tv.equals(tt.view); got != tt.want {
				t.Errorf("%q: TView.equals() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TView_equals()

func Test_TView_render(t *testing.T) {
	prepareTestFiles()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1, "View_render: Oh dear! This is a first posting.")

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "View_render: Hi there! This is another posting.")

	id3 := time2id(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p3 := NewPosting(id3, "View_render: Oh dear! This is a single posting.")

	v1, _ := NewView("index")
	pl1 := NewPostList().Add(p1).Add(p2).Sort()
	dl1 := NewTemplateData().
		Set("Title", "this is the title").
		Set("Headline", "This is an interesting issue").
		Set("Postings", pl1)

	v2, _ := NewView("article")
	dl2 := NewTemplateData().
		Set("Title", "this is the article title").
		Set("Headline", "Tis is an important topic").
		Set("Lang", "en").
		Set("Postings", p3).
		Set("ToBeIgnored", "! Ignore Me !")

	type tArgs struct {
		aWriter io.Writer
		aData   *TemplateData
	}
	tests := []struct {
		name    string
		view    TView
		args    tArgs
		wantErr bool
	}{
		{"1", *v1, tArgs{os.Stdout, dl1}, false},
		{"2", *v2, tArgs{os.Stdout, dl2}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.view.render(tt.args.aWriter, tt.args.aData); (err != nil) != tt.wantErr {
				t.Errorf("%q: TView.render() error = %v,\nwantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
		})
	}
} // Test_TView_render()

/* _EoF_ */
