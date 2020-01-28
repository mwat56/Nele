/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func Test_mdToHTML(t *testing.T) {
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
		// TODO: Add test cases.
		{" 1", args{md1}, ht1},
		{" 2", args{md2}, ht2},
		{" 3", args{md3}, ht3},
		{" 4", args{md4}, ht4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mdToHTML(tt.args.aMarkdown); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mdToHTML() = [%s],\nwant [%s]", got, tt.want)
			}
		})
	}
} // Test_mdToHTML()

func Test_newID(t *testing.T) {
	ct000 := time.Date(2019, 10, 22, 0, 0, 0, 0, time.Local)
	ct001 := time.Date(2019, 10, 23, 0, 0, 0, 0, time.Local)
	ct052 := time.Date(2019, 12, 13, 0, 0, 0, 0, time.Local)
	ct053 := time.Date(2019, 12, 14, 0, 0, 0, 0, time.Local)
	ct104 := time.Date(2020, 2, 3, 0, 0, 0, 0, time.Local)
	ct105 := time.Date(2020, 2, 4, 0, 0, 0, 0, time.Local)
	ct158 := time.Date(2020, 3, 27, 0, 0, 0, 0, time.Local)
	ct159 := time.Date(2020, 3, 28, 0, 0, 0, 0, time.Local)
	ct209 := time.Date(2020, 5, 18, 0, 0, 0, 0, time.Local)
	ct210 := time.Date(2020, 5, 19, 0, 0, 0, 0, time.Local)
	ct261 := time.Date(2020, 7, 9, 0, 0, 0, 0, time.Local)
	ct262 := time.Date(2020, 7, 10, 0, 0, 0, 0, time.Local)
	ct313 := time.Date(2020, 8, 30, 0, 0, 0, 0, time.Local)
	ct314 := time.Date(2020, 8, 31, 0, 0, 0, 0, time.Local)
	ct365 := time.Date(2020, 10, 21, 0, 0, 0, 0, time.Local)
	ct366 := time.Date(2020, 10, 22, 0, 0, 0, 0, time.Local)
	type args struct {
		aTime time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"000", args{ct000}, "15cfc8750b2fc000"},
		{"001", args{ct001}, "15d017099c7ec000"},
		{"052", args{ct052}, "15dfc1e8bff46000"},
		{"053", args{ct053}, "15e0107d51436000"},
		{"104", args{ct104}, "15efb81644006000"},
		{"105", args{ct105}, "15f006aad54f6000"},
		{"158", args{ct158}, "15fffcd8595b6000"},
		{"159", args{ct159}, "16004b6ceaaa6000"},
		{"209", args{ct209}, "160fefbfacaec000"},
		{"210", args{ct210}, "16103e543dfdc000"},
		{"261", args{ct261}, "161fe5ed30bac000"},
		{"262", args{ct262}, "16203481c209c000"},
		{"313", args{ct313}, "162fdc1ab4c6c000"},
		{"314", args{ct314}, "16302aaf4615c000"},
		{"365", args{ct365}, "163fd24838d2c000"},
		{"366", args{ct366}, "164020dcca21c000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newID(tt.args.aTime); got != tt.want {
				t.Errorf("newID() = [%v], want [%v]", got, tt.want)
			}
		})
	}
} // Test_newID()

func Test_newPost(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	var bs []byte
	id1 := "12345678"
	rp1 := &TPosting{
		id1,
		bs,
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		args args
		want *TPosting
	}{
		// TODO: Add test cases.
		{" 1", args{id1}, rp1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPosting(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPost() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_newPost()

func TestPostingCount(t *testing.T) {
	SetPostingBaseDirectory("./postings/")
	atomic.StoreUint32(&µCountCache, 0) // invalidate count cache
	tests := []struct {
		name       string
		wantRCount int
	}{
		// TODO: Add test cases.
		{" 1", 1026},
		{" 2", 1026},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRCount := PostingCount(); gotRCount != tt.wantRCount {
				t.Errorf("PostingCount() = %v, want %v", gotRCount, tt.wantRCount)
			}
		})
	}
} // TestPostingCount()

func TestTPosting_After(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2)
	type fields struct {
		p *TPosting
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, args{id2}, false},
		{" 2", fields{p2}, args{id1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.After(tt.args.aID); got != tt.want {
				t.Errorf("TPosting.After() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTPosting_After()

func TestTPosting_Before(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2)
	type fields struct {
		p *TPosting
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, args{id2}, true},
		{" 2", fields{p2}, args{id1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.Before(tt.args.aID); got != tt.want {
				t.Errorf("TPosting.(Before) = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTPosting_Before()

func TestTPosting_Clear(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id := newID(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id)
	rp := p1.clone()
	md2 := []byte("Oh dear! This is a posting.")
	p2 := NewPosting(id).Set(md2)
	p3 := NewPosting(id)
	p3.Set(md2).Len()

	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name   string
		fields fields
		want   *TPosting
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, rp},
		{" 2", fields{p2}, rp},
		{" 3", fields{p3}, rp},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.Clear(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.Clear() = '%v',\n\t\t\twant '%v'", got, tt.want)
			}
		})
	}
} // TestTPosting_Clear()

func TestTPosting_clone(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id := newID(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id).
		Set([]byte("Oh dear! This is a posting."))
	wp1 := NewPosting(id).
		Set([]byte("Oh dear! This is a posting."))
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name   string
		fields fields
		want   *TPosting
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, wp1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.clone() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // TestTPosting_clone()

func TestTPosting_Delete(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2).
		Set([]byte("just a dummy"))
	p2.Store() // create a file
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, false},
		{" 2", fields{p2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if err := p.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("TPosting.Delete() error = %v, wantErr '%v'", err, tt.wantErr)
			}
		})
	}
} // TestTPosting_Delete()

func TestTPosting_Equal(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2)
	type fields struct {
		p *TPosting
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, args{id2}, false},
		{" 2", fields{p2}, args{id2}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.Equal(tt.args.aID); got != tt.want {
				t.Errorf("TPosting.Equal() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTPosting_Equal()

func TestTPosting_Exists(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 1, 1, 0, 0, 0, 1, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 1, 1, 0, 0, 0, 2, time.Local))
	p2 := NewPosting(id2)
	id3 := newID(time.Date(2019, 1, 1, 0, 0, 0, 3, time.Local))
	p3 := NewPosting(id3).Set([]byte("Hello World"))
	p3.Store()
	type fields struct {
		id       string
		markdown []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1.id, p1.markdown}, false},
		{" 2", fields{p2.id, p2.markdown}, false},
		{" 3", fields{p3.id, p3.markdown}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &TPosting{
				id:       tt.fields.id,
				markdown: tt.fields.markdown,
			}
			if got := p.Exists(); got != tt.want {
				t.Errorf("TPosting.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_Exists()

func TestTPosting_Load(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	md2 := []byte("Load: this is more nonsense")
	p2 := NewPosting(id2).Set(md2)
	p2.Store()
	p2.Clear()
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, true},
		{" 2", fields{p2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if err := p.Load(); (err != nil) != tt.wantErr {
				t.Errorf("TPosting.LoadMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.fields.p.Delete() // clean up
	}
} // TestTPosting_Load()

func TestTPosting_makeDir(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1)
	rp1 := "/tmp/postings/2019158/158d2fcc0ff16000"
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2)
	rp2 := "/tmp/postings/2019159/159b4b37fb6ac000"
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, rp1, false},
		{" 2", fields{p2}, rp2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			got, err := p.makeDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("TPosting.MakeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TPosting.MakeDir() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_makeDir()

func TestTPosting_Markdown(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	md1 := []byte("Markdown: this is a nonsensical posting")
	p1 := NewPosting(id1).Set(md1)
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	md2 := []byte("Markdown: this is more nonsense")
	p2 := NewPosting(id2).Set(md2)
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, md1, false},
		{" 2", fields{p2}, md2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			got := p.Markdown()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.Markdown() = [%s], want [%s]", got, tt.want)
			}
		})
	}
} // TestTPosting_Markdown()

func TestTPosting_pathFileName(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1)
	rp1 := "/tmp/postings/2019158/158d2fcc0ff16000.md"
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2)
	rp2 := "/tmp/postings/2019159/159b4b37fb6ac000.md"
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, rp1},
		{" 2", fields{p2}, rp2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.PathFileName(); got != tt.want {
				t.Errorf("TPosting.pathName() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_pathFileName()

func TestTPosting_Set(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	md1 := []byte("Set: this is obviously nonsense")
	p1 := NewPosting(id1)
	rp1 := p1.clone()
	rp1.markdown = md1
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	md2 := []byte("Set: this is more nonsense")
	p2 := NewPosting(id2)
	rp2 := p2.clone()
	rp2.markdown = md2
	type fields struct {
		p *TPosting
	}
	type args struct {
		aMarkdown []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *TPosting
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, args{md1}, rp1},
		{" 2", fields{p2}, args{md2}, rp2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			got := p.Set(tt.args.aMarkdown)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.Set() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_Set

func TestTPosting_Store(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	var len1 int64
	id1 := newID(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1)
	id2 := newID(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	md2 := []byte("Store: this is more nonsense")
	p2 := NewPosting(id2).
		Set(md2)
	len2 := int64(len(md2))
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name    string
		fields  fields
		want    int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, len1, false},
		{" 2", fields{p2}, len2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			got, err := p.Store()
			if (err != nil) != tt.wantErr {
				t.Errorf("TPosting.Store() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TPosting.Store() = '%v',\nwant '%v'", got, tt.want)
			}
		})
		tt.fields.p.Delete() // clean up
	}
} // TestTPosting_Store()

func TestTPosting_Time(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	tm1 := time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local)
	id1 := newID(tm1)
	p1 := NewPosting(id1)
	tm2 := time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local)
	id2 := newID(tm2)
	p2 := NewPosting(id2)
	tm3 := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	id3 := newID(tm3)
	p3 := NewPosting(id3)
	p3.id = ""
	type fields struct {
		p *TPosting
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		// TODO: Add test cases.
		{" 1", fields{p1}, tm1},
		{" 2", fields{p2}, tm2},
		{" 3", fields{p3}, tm3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.p
			if got := p.Time(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.Time() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_Time()
