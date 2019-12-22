/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"testing"
)

func TestNewPageHandler(t *testing.T) {
	InitConfig()
	SetPostingBaseDirectory("/tmp/postings/")
	prepareTestFiles()
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", 18, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageHandler()
			if (nil != err) != tt.wantErr {
				t.Errorf("NewPageHandler() error = %v,\nwantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != got.Len() {
				t.Errorf("NewPageHandler() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestNewPageHandler()

func Test_numStart(t *testing.T) {
	type args struct {
		aString string
	}
	tests := []struct {
		name       string
		args       args
		wantRNum   int
		wantRStart int
	}{
		// TODO: Add test cases.
		{" 0", args{","}, 0, 0},
		{" 1", args{"10"}, 10, 0},
		{" 2", args{"10-10"}, 10, 10},
		{" 3", args{"10,"}, 10, 0},
		{" 4", args{",10"}, 0, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRNum, gotRStart := numStart(tt.args.aString)
			if gotRNum != tt.wantRNum {
				t.Errorf("numStart() gotRNum = %v, want %v", gotRNum, tt.wantRNum)
			}
			if gotRStart != tt.wantRStart {
				t.Errorf("numStart() gotRStart = %v, want %v", gotRStart, tt.wantRStart)
			}
		})
	}
} // Test_numStart()

func TestURLparts(t *testing.T) {
	type args struct {
		aURL string
	}
	tests := []struct {
		name      string
		args      args
		wantRHead string
		wantRTail string
	}{
		// TODO: Add test cases.
		{" 1", args{"/"}, "", ""},
		{" 1a", args{""}, "", ""},
		{" 1b", args{"index/ "}, "index", ""},
		{" 2", args{"/css"}, "css", ""},
		{" 2a", args{"css"}, "css", ""},
		{" 3", args{"/css/styles.css"}, "css", "styles.css"},
		{" 3a", args{"css/styles.css"}, "css", "styles.css"},
		{" 4", args{"/?q=searchterm"}, "", "?q=searchterm"},
		{" 4a", args{"?q=searchterm"}, "", "?q=searchterm"},
		{" 5", args{"/article/abcdef1122334455"},
			"article", "abcdef1122334455"},
		{" 6", args{"/q/searchterm"}, "q", "searchterm"},
		{" 6a", args{"/q/?s=earchterm"}, "q", "?s=earchterm"},
		{" 7", args{"/q/search?s=term"}, "q", "search?s=term"},
		{" 8", args{"/static/https://github.com/"}, "static", "https://github.com/"},
		{" 9", args{"/ht/kurzerklärt"}, "ht", "kurzerklärt"},
		{"10", args{`share/https://utopia.de/ratgeber/pink-lady-das-ist-faul-an-dieser-apfelsorte/#main_content`}, `share`, `https://utopia.de/ratgeber/pink-lady-das-ist-faul-an-dieser-apfelsorte/#main_content`},
		{"11", args{"/s/search term"}, "s", "search term"},
		{"12", args{"/ml/antoni_comín"}, "ml", "antoni_comín"},
		{"13", args{"/s/Änderungen erklären"}, "s", "Änderungen erklären"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRHead, gotRTail := URLparts(tt.args.aURL)
			if gotRHead != tt.wantRHead {
				t.Errorf("URLpath1() gotRHead = {%v}, want {%v}", gotRHead, tt.wantRHead)
			}
			if gotRTail != tt.wantRTail {
				t.Errorf("URLpath1() gotRTail = {%v}, want {%v}", gotRTail, tt.wantRTail)
			}
		})
	}
} // TestURLparts()
