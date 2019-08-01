/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"html/template"
	"reflect"
	"testing"
)

func Test_cachedHTML(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	md1 := []byte(`
---

Here is an example of AppleScript:

	tell application "Foo"
		beep
	end tell

---

Whatever the reason, the compression doesn't seem to take much off.
I wonder whether it's really worth the effort.

But I guess we will find out soon: By adding more text to this test data the compression ratio should improve.

At least, that's what I think …
`)
	p1 := &TPosting{
		id:       "15a200c9e8b51bd1",
		markdown: md1,
	}
	p1.Store()
	w1 := template.HTML(`<hr>

<p>Here is an example of AppleScript:</p><pre>
tell application &quot;Foo&quot;
	beep
end tell
</pre><hr>

<p>Whatever the reason, the compression doesn&rsquo;t seem to take much off.
I wonder whether it&rsquo;s really worth the effort.</p>

<p>But I guess we will find out soon: By adding more text to this test data the compression ratio should improve.</p>

<p>At least, that&rsquo;s what I think …</p>
`)
	p2 := &TPosting{
		id:       "15a200c9e8b51bd2",
		markdown: []byte("Just a test sentence."),
	}
	p2.Store()
	w2 := template.HTML(`<p>Just a test sentence.</p>
`)
	p3 := &TPosting{
		id: "15a200c9e8b51bd3",
		markdown: []byte(`Just a test sentence.

And another one of those.
`),
	}
	p3.Store()
	w3 := template.HTML(`<p>Just a test sentence.</p>

<p>And another one of those.</p>
`)
	p4 := &TPosting{
		id: "15a200c9e8b51bd4",
		markdown: []byte(`Just a test sentence.

And another one of those.

And yet another sentence for testing.
`),
	}
	p4.Store()
	w4 := template.HTML(`<p>Just a test sentence.</p>

<p>And another one of those.</p>

<p>And yet another sentence for testing.</p>
`)
	type args struct {
		aPost *TPosting
	}
	tests := []struct {
		name string
		args args
		want template.HTML
	}{
		// TODO: Add test cases.
		{" 1", args{p1}, w1},
		{" 2", args{p2}, w2},
		{" 3", args{p3}, w3},
		{" 4", args{p4}, w4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := template.HTML(cachedHTML(tt.args.aPost)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cachedHTML() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // Test_cachedHTML()
