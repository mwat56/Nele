package nele

import (
	"io/ioutil"
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

func TestMDtoHTML(t *testing.T) {
	m1 := []byte(`# head

a single paragraph[^1].

another "test" paragraph.

    some preformatted
	text

[^1]: the footnote text
`)
	w1 := []byte(`<h1>head</h1>

<p>a single paragraph<sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup>.</p>

<p>another &ldquo;test&rdquo; paragraph.</p><pre>
some preformatted
text
</pre><div class="footnotes">

<hr>

<ol>
<li id="fn:1">the footnote text <a class="footnote-return" href="#fnref:1"><sup>[return]</sup></a></li>
</ol>

</div>
`)
	type args struct {
		aMarkdown []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 1", args{m1}, w1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MDtoHTML(tt.args.aMarkdown); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MDtoHTML() = %s, want %s", got, tt.want)
			}
		})
	}
} // TestMDtoHTML()

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
	SetPostingBaseDirectory("/tmp/postings/")
	bd := PostingBaseDirectory()
	prepareTestFiles()
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
			if got := SearchPostings(tt.args.aText); got.Len() != tt.want {
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
		{" 8", args{"/static/https://github.com/"}, "static", "https://github.com/"},
		{" 9", args{"/ht/kurzerklärt"}, "ht", "kurzerklärt"},
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
		{" 7", args{"23:2:73"}, 23, 2, 7},
		{" 8", args{"1:82:3"}, 1, 0, 0},
		{" 9", args{"1:0:1"}, 1, 0, 1},
		{"10", args{"01:02:03/#anchor"}, 1, 2, 3},
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
		{" 5", args{"1914Sep01"}, 1914, 0, 0},
		{" 6", args{"2019-10-18"}, 2019, 10, 18},
		{" 7", args{"1914 09 28"}, 1914, 9, 28},
		{" 8", args{"20191308"}, 2019, 0, 0},
		{" 9", args{"20190332"}, 2019, 3, 0},
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

func Test_handlePreCode(t *testing.T) {
	m1 := []byte(`<p>test</p>

	<pre><code>pre
part</code></pre>

<p>line 2</p>
`)
	w1 := []byte(`<p>test</p><pre>
pre
part
</pre><p>line 2</p>
`)
	m2 := []byte(`<p>test</p>

	<pre><code class="language-go">pre
part</code></pre>

<p>line 2</p>
`)
	w2 := []byte(`<p>test</p><pre class="language-go">
pre
part
</pre><p>line 2</p>
`)
	type args struct {
		aMarkdown []byte
	}
	tests := []struct {
		name      string
		args      args
		wantRHTML []byte
	}{
		// TODO: Add test cases.
		{" 1", args{m1}, w1},
		{" 2", args{m2}, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRHTML := handlePreCode(tt.args.aMarkdown); !reflect.DeepEqual(gotRHTML, tt.wantRHTML) {
				t.Errorf("handlePreCode() = %s, want %s", gotRHTML, tt.wantRHTML)
			}
		})
	}
} // Test_handlePreCode()
