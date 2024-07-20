/*
Copyright © 2020, 2024  M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"reflect"
	"testing"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestMDtoHTML(t *testing.T) {
	md1 := []byte(`
---

Here is an example of AppleScript:

	tell application "Foo"
		beep
	end tell

---
`)
	ht1 := []byte(`<hr>

<p>Here is an example of AppleScript:</p><pre>
tell application &quot;Foo&quot;
	beep
end tell
</pre><hr>`)

	md2 := []byte(`
---

    tell application "Foo"
      beep
	end tell

That's an example of AppleScript

---
`)
	ht2 := []byte(`<hr><pre>
tell application &quot;Foo&quot;
  beep
end tell
</pre><p>That&rsquo;s an example of AppleScript</p>

<hr>`)

	md3 := []byte("Hello `world`!")
	ht3 := []byte(`<p>Hello <code>world</code>!</p>`)

	md4 := []byte(`# head

a single paragraph[^1].

another "test" paragraph.

    some preformatted
	text

[^1]: the footnote text
`)
	ht4 := []byte(`<h1>head</h1>

<p>a single paragraph<sup class="footnote-ref" id="fnref:1"><a href="#fn:1">1</a></sup>.</p>

<p>another &ldquo;test&rdquo; paragraph.</p><pre>
some preformatted
text
</pre><div class="footnotes">

<hr>

<ol>
<li id="fn:1">the footnote text <a class="footnote-return" href="#fnref:1"><sup>[return]</sup></a></li>
</ol>

</div>`)

	type args struct {
		aMarkdown []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"1", args{md1}, ht1},
		{"2", args{md2}, ht2},
		{"3", args{md3}, ht3},
		{"4", args{md4}, ht4},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MDtoHTML(tt.args.aMarkdown); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: MDtoHTML() = \n%s\n>>> want >>>\n%s",
					tt.name, got, tt.want)
			}
		})
	}
} // TestMDtoHTML()

/* _EoF_ */
