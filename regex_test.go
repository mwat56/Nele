/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"io/ioutil"
	"reflect"
	"testing"
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

func BenchmarkMDtoHTML(b *testing.B) {
	page, _ := ioutil.ReadFile("./Markdown_Syntax.md")
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if html := MDtoHTML(page); nil == html {
			continue
		}
	}
} // BenchmarkMDtoHTML()

func Test_addExternURLtagets(t *testing.T) {
	t1 := ` bla <a href="https://site/page">bla</a> `
	p1 := []byte(t1)
	w1 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> `)
	t2 := t1 + `bla <a href="/page">bla</a>`
	p2 := []byte(t2)
	w2 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> bla <a href="/page">bla</a>`)
	t3 := t1 + `bla <a href="http://site.com/page">bla</a>`
	p3 := []byte(t3)
	w3 := []byte(` bla <a target="_extern" href="https://site/page">bla</a> bla <a target="_extern" href="http://site.com/page">bla</a>`)
	type args struct {
		aPage []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 1", args{p1}, w1},
		{" 2", args{p2}, w2},
		{" 3", args{p3}, w3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addExternURLtagets(tt.args.aPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addExternURLtagets() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_addExternURLtagets()

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
