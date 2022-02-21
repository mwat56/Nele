/*
   Copyright © 2019, 2022 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"testing"
	"time"
)

func TestNewPageHandler(t *testing.T) {
	// `InitConfig()` calls `flag.parse()` which in turn will cause
	// errors when run with `go test …`.
	readAppArgsDebug()
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
		{"14", args{"///asterisk/admin/config.php"}, "asterisk", "admin/config.php"},
		{"15", args{"/p/15ee22f54a6f700e"}, "p", "15ee22f54a6f700e"},
		{"16", args{"/ml/edward_snowden's"}, "ml", "edward_snowden's"},
		{"17", args{"/ml/paul_o'hare"}, "ml", "paul_o'hare"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRHead, gotRTail := URLparts(tt.args.aURL)
			if gotRHead != tt.wantRHead {
				t.Errorf("URLpath1() gotRHead = {%v},\nwant {%v}", gotRHead, tt.wantRHead)
			}
			if gotRTail != tt.wantRTail {
				t.Errorf("URLpath1() gotRTail = {%v},\nwant {%v}", gotRTail, tt.wantRTail)
			}
		})
	}
} // TestURLparts()

func Test_getHMS(t *testing.T) {
	type args struct {
		aTime string
	}
	tests := []struct {
		name        string
		args        args
		wantRHour   int
		wantRMinute int
		wantRSecond int
	}{
		// TODO: Add test cases.
		{" 1", args{"1:2:3"}, 1, 2, 3},
		{" 2", args{"01:2:3"}, 1, 2, 3},
		{" 3", args{"01:02:3"}, 1, 2, 3},
		{" 4", args{"01:02:03"}, 1, 2, 3},
		{" 5", args{"23:02:03"}, 23, 2, 3},
		{" 6", args{"24:02:03"}, 0, 0, 0},
		{" 7", args{"23:2:73"}, 23, 2, 7},
		{" 8", args{"1:82:3"}, 1, 8, 0},
		{" 9", args{"1:0:1"}, 1, 0, 1},
		{"10", args{"01:02"}, 1, 2, 0},
		{"11", args{"01:02:"}, 1, 2, 0},
		{"12", args{"01:02:03/#anchor"}, 1, 2, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRHour, gotRMinute, gotRSecond := getHMS(tt.args.aTime)
			if gotRHour != tt.wantRHour {
				t.Errorf("getHMS() gotRHour = %v, want %v", gotRHour, tt.wantRHour)
			}
			if gotRMinute != tt.wantRMinute {
				t.Errorf("getHMS() gotRMinute = %v, want %v", gotRMinute, tt.wantRMinute)
			}
			if gotRSecond != tt.wantRSecond {
				t.Errorf("getHMS() gotRSecond = %v, want %v", gotRSecond, tt.wantRSecond)
			}
		})
	}
} // Test_getHMS()

func Test_getYMD(t *testing.T) {
	type args struct {
		aDate string
	}
	tests := []struct {
		name       string
		args       args
		wantRYear  int
		wantRMonth time.Month
		wantRDay   int
	}{
		// TODO: Add test cases.
		{" 0", args{""}, 0, 0, 0},
		{" 1", args{"2019"}, 2019, 0, 1},
		{" 2", args{"201909"}, 2019, 9, 1},
		{" 3", args{"20191009"}, 2019, 10, 9},
		{" 4", args{"WTF"}, 0, 0, 0},
		{" 5", args{"1914Sep01"}, 1914, 0, 1},
		{" 6", args{"2019-10-18"}, 2019, 10, 18},
		{" 7", args{"1914 09 28"}, 1914, 9, 28},
		{" 8", args{"20191308"}, 2019, 0, 1},
		{" 9", args{"20190332"}, 2019, 3, 1},
		{"10", args{"20191021#p156b52d99af4b401"}, 2019, 10, 21},
		{"11", args{"20181128/#p156b52d99af4b401"}, 2018, 11, 28},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRYear, gotRMonth, gotRDay := getYMD(tt.args.aDate)
			if gotRYear != tt.wantRYear {
				t.Errorf("getYMD() gotRYear = %v, want %v", gotRYear, tt.wantRYear)
			}
			if gotRMonth != tt.wantRMonth {
				t.Errorf("getYMD() gotRMonth = %v, want %v", gotRMonth, tt.wantRMonth)
			}
			if gotRDay != tt.wantRDay {
				t.Errorf("getYMD() gotRDay = %v, want %v", gotRDay, tt.wantRDay)
			}
		})
	}
} // Test_getYMD()
