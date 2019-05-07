package blog

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func Test_initWSre(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{" 1", 12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initWSre(); got != tt.want {
				t.Errorf("initWSre() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_initWSre()

func TestRemoveWhiteSpace(t *testing.T) {
	txtIn1 := []byte(`<hr />
	<p>Here is an example of AppleScript:</p>
	<pre>
	tell application &quot;Foo&quot;
	  beep
	end tell
	</pre>
	<hr />`)
	txtOut1 := []byte(`<hr /><p>Here is an example of AppleScript:</p><pre>
	tell application &quot;Foo&quot;
	  beep
	end tell
	</pre><hr />`)
	// pageIn2, _ := ioutil.ReadFile("./TODO.html")
	// pageOut2 := []byte(`<!DOCTYPE html><html><head><meta charset="UTF-8"><meta http-equiv="Window-target" content="_top"><meta http-equiv="Content-Type" content="text/html; charset=utf-8"><title>ToDo</title><meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=yes"><style type="text/css">html,body{font-size:101%;font-family:serif,monospace;}
	// code,kbd,tt,pre{background-color:#f0f0f0;font-family:monospace;font-size:0.987em;letter-spacing:normal;line-height:1.123;white-space:normal;}
	// pre,xmp{border-left:thin solid #999;display:block;font-size:88%;font-weight:inherit;line-height:1.36;margin:1ex 1ex 1ex 0.3ex;padding:0.5ex 1ex;overflow:visible;text-align:left;text-indent:0;white-space:pre-wrap;width:98%;}
	// small{color:inherit;font-size:79%;font-weight:inherit;}</style></head><body><h1>ToDo</h1><h2>Some unsorted things that came to mind:</h2><ul><li>Internationalisation of text(-fragments)</li><li>cmdline arguments for:<ul><li><del>port (to listen to),</del></li><li><del>ip (to listen at),</del></li><li>hostname (to accept requests for) ??,</li><li><del>path/name for INI file,</del></li><li><del>path/dir for template files,</del></li></ul></li><li><del>static routes for CSS. IMG, JS</del></li><li><del>logic for reading *.tpl files</del></li><li>config data (INI file) for:<ul><li>CERT files (for TLS)</li><li><del>ip (to listen at),</del></li><li><del>port (to listen to),</del></li><li>hostname (to accept requests for),</li><li><del>name of template directory (for *.tpl files)</del></li><li>(opt.) DB username/pass??</li><li><del>CSS dir</del></li><li><del>IMG dir</del></li></ul></li></ul><hr /></body></html>`)
	// page3, _ := ioutil.ReadFile("./Markdown_Syntax.html")
	type args struct {
		aPage []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 1", args{txtIn1}, txtOut1},
		// {" 2", args{pageIn2}, pageOut2},
		// {" 3", args{page3}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveWhiteSpace(tt.args.aPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveWhiteSpace() = [%s],\nwant [%s]", got, tt.want)
			}
		})
	}
} // TestRemoveWhiteSpace()

func TestSearchPostings(t *testing.T) {
	bd, _ := filepath.Abs("/tmp/postings/")
	prepareTestFiles(bd)
	type args struct {
		aBaseDir string
		aText    string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", args{bd, "16"}, 113},
		{" 2", args{bd, "8"}, 151},
		{" 3", args{bd, "1\\d+"}, 153},
		{" 4", args{bd, "10\\d+"}, 30},
		{" 5", args{bd, "08\\s+08"}, 2},
		{" 6", args{bd, bd}, 333},
		{" 7", args{bd, "postings"}, 333},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SearchPostings(tt.args.aBaseDir, tt.args.aText); got.Len() != tt.want {
				t.Errorf("Search() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestSearchPostings()

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
		{" 1b", args{"index/"}, "index", ""},
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

func BenchmarkMDtoHTML(b *testing.B) {
	page, _ := ioutil.ReadFile("./Markdown_Syntax.md")
	// page, _ := ioutil.ReadFile("./TODO.md")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if html := MDtoHTML(page); nil == html {
			continue
		}
	}
} // BenchmarkMDtoHTML()

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
		{" 7", args{"23:2:73"}, 0, 0, 0},
		{" 8", args{"0:82:3"}, 0, 0, 0},
		{" 9", args{"0:0:1"}, 0, 0, 1},
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
		{" 1", args{"2019"}, 2019, 0, 0},
		{" 2", args{"201909"}, 2019, 9, 0},
		{" 3", args{"20191009"}, 2019, 10, 9},
		{" 4", args{"WTF"}, 0, 0, 0},
		{" 5", args{"1914Sep01"}, 0, 0, 0},
		{" 6", args{"2019-10-18"}, 2019, 10, 18},
		{" 7", args{"1914 09 28"}, 1914, 9, 28},
		{" 8", args{"20191308"}, 0, 0, 0},
		{" 9", args{"20190332"}, 0, 0, 0},
		{"10", args{"20191021"}, 2019, 10, 21},
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
}
