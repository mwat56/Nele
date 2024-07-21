/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
	"time"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func Test_NewPosting(t *testing.T) {
	prep4Tests()

	var md []byte
	id1 := uint64(time.Now().UnixNano())
	wp1 := &TPosting{
		id:           id1,
		lastModified: time.Now(),
		markdown:     md,
	}

	tests := []struct {
		name string
		id   uint64
		want *TPosting
	}{
		{"1", id1, wp1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPosting(tt.id, "")
			if (got.id != tt.want.id) ||
				(got.String() != tt.want.String()) {
				t.Errorf("%q: NewPosting() = %v,\nwant %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_NewPosting()

func Test_TPosting_After(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1, "")
	id2 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "")

	tests := []struct {
		name string
		post *TPosting
		id   uint64
		want bool
	}{
		// TODO: Add test cases.
		{" 1", p1, id2, false},
		{" 2", p2, id1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if got := p.After(tt.id); got != tt.want {
				t.Errorf("TPosting.After() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // Test_TPosting_After()

func Test_TPosting_Before(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1, "")
	id2 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "")

	tests := []struct {
		name string
		post *TPosting
		id   uint64
		want bool
	}{
		{"1", p1, id2, true},
		{"2", p2, id1, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if got := p.Before(tt.id); got != tt.want {
				t.Errorf("TPosting.Before() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // Test_TPosting_Before()

func TestTPosting_ChangeID(t *testing.T) {
	prep4Tests()

	p1 := NewPosting(0, "> one")
	p1.Store()
	id2 := time2id(time.Now()) + 1

	tests := []struct {
		name    string
		post    *TPosting
		id      uint64
		wantErr bool
	}{
		{"1", p1, id2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &TPosting{
				id:           tt.post.id,
				lastModified: tt.post.lastModified,
				markdown:     tt.post.markdown,
			}
			if err := p.ChangeID(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("TPosting.ChangeID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTPosting_ChangeID()

func Test_TPosting_Clear(t *testing.T) {
	prep4Tests()

	id := time2id(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id, "")
	rp := p1.clone()
	md2 := []byte("Oh dear! This is a posting.")
	p2 := NewPosting(id, "").Set(md2)
	p3 := NewPosting(id, "")
	p3.Set(md2).Len()

	tests := []struct {
		name string
		post *TPosting
		want *TPosting
	}{
		{" 1", p1, rp},
		{" 2", p2, rp},
		{" 3", p3, rp},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.Clear()
			if (got.id != tt.want.id) ||
				!reflect.DeepEqual(got.markdown, tt.want.markdown) {
				t.Errorf("TPosting.Clear() = '%v',\n\t\t\twant '%v'", got, tt.want)
			}
		})
	}
} // Tes_tTPosting_Clear()

func Test_TPosting_clone(t *testing.T) {
	prep4Tests()

	id := time2id(time.Date(2019, 4, 14, 0, 0, 0, 0, time.Local))
	t1 := "Oh dear! This is a posting."
	p1 := NewPosting(id, t1)
	wp1 := NewPosting(id, t1)

	tests := []struct {
		name string
		post *TPosting
		want *TPosting
	}{
		{"1", p1, wp1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.clone()
			if (got.id != tt.want.id) ||
				// the `lastModified` field will be slightly different
				(string(got.markdown) != string(tt.want.markdown)) {
				t.Errorf("TPosting.clone() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_TPosting_clone()

func Test_TPosting_Delete(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1, "")

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "just a dummy")
	_, _ = p2.Store() // create a file

	tests := []struct {
		name    string
		post    *TPosting
		wantErr bool
	}{
		{"1", p1, false},
		{"2", p2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if err := p.Delete(); (err != nil) != tt.wantErr {
				t.Errorf("TPosting.Delete() error = %v, wantErr '%v'", err, tt.wantErr)
			}
		})
	}
} // Test_TPosting_Delete()

func Test_TPosting_Equal(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 1, 1, 0, 0, 0, -1, time.Local))
	p1 := NewPosting(id1, "")

	id2 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "")

	tests := []struct {
		name string
		post *TPosting
		id   uint64
		want bool
	}{
		{"1", p1, id2, false},
		{"2", p2, id2, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if got := p.Equal(tt.id); got != tt.want {
				t.Errorf("TPosting.Equal() = '%v', want '%v'", got, tt.want)
			}
		})
	}
} // TestTPosting_Equal()

func Test_TPosting_Exists(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 1, time.Local))
	p1 := NewPosting(id1, "")

	id2 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 2, time.Local))
	p2 := NewPosting(id2, "")

	id3 := time2id(time.Date(2019, 1, 1, 0, 0, 0, 3, time.Local))
	p3 := NewPosting(id3, "Hello World")
	_, _ = p3.Store()

	tests := []struct {
		name string
		post *TPosting
		want bool
	}{
		{"1", p1, false},
		{"2", p2, false},
		{"3", p3, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if got := p.Exists(); got != tt.want {
				t.Errorf("TPosting.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_TPosting_Exists()

func Test_TPosting_Load(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1, "")

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "Load: this is more nonsense")
	_, _ = p2.Store()

	p2.Clear()

	tests := []struct {
		name    string
		post    *TPosting
		wantErr bool
	}{
		{"1", p1, true},
		{" 2", p2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.post
			if err := p.Load(); (err != nil) != tt.wantErr {
				t.Errorf("TPosting.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
		tt.post.Delete() // clean up
	}
} // Test_TPosting_Load()

func Test_TPosting_Markdown(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 1, time.Local))
	md1 := "Markdown: this is a nonsensical posting"
	p1 := NewPosting(id1, md1)

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 2, time.Local))
	md2 := "Markdown: this is more nonsense"
	p2 := NewPosting(id2, md2)

	tests := []struct {
		name    string
		post    *TPosting
		want    []byte
		wantErr bool
	}{
		{"1", p1, []byte(md1), false},
		{"2", p2, []byte(md2), false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.Markdown()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPosting.Markdown() = [%s], want [%s]",
					got, tt.want)
			}
		})
	}
} // Test_TPosting_Markdown()

func Test_TPosting_pathFileName(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 0, time.Local))
	p1 := NewPosting(id1, "")
	rp1 := "/tmp/postings/2019158/158d2fcc0ff16000.md"

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 0, time.Local))
	p2 := NewPosting(id2, "")
	rp2 := "/tmp/postings/2019159/159b4b37fb6ac000.md"

	tests := []struct {
		name string
		post *TPosting
		want string
	}{
		{"1", p1, rp1},
		{"2", p2, rp2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.post.PathFileName(); got != tt.want {
				t.Errorf("TPosting.pathName() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPosting_pathFileName()

func Test_TPosting_Set(t *testing.T) {
	prep4Tests()

	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 1, time.Local))
	md1 := []byte("Set: this is obviously nonsense")
	p1 := NewPosting(id1, "")
	rp1 := p1.clone()
	rp1.markdown = md1

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 2, time.Local))
	md2 := []byte("Set: this is more nonsense")
	p2 := NewPosting(id2, "")
	rp2 := p2.clone()
	rp2.markdown = md2

	tests := []struct {
		name     string
		post     *TPosting
		markdown []byte
		want     *TPosting
	}{
		{"1", p1, md1, rp1},
		{"2", p2, md2, rp2},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.post.Set(tt.markdown)
			if (got.id != tt.want.id) ||
				(string(got.markdown) != string(tt.want.markdown)) {
				t.Errorf("TPosting.Set() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_TPosting_Set

func Test_TPosting_Store(t *testing.T) {
	prep4Tests()

	var len1 int
	id1 := time2id(time.Date(2019, 3, 19, 0, 0, 0, 1, time.Local))
	p1 := NewPosting(id1, "")

	id2 := time2id(time.Date(2019, 5, 4, 0, 0, 0, 2, time.Local))
	txt2 := "Store: this is more nonsense"
	p2 := NewPosting(id2, txt2)
	len2 := len(txt2)

	tests := []struct {
		name    string
		post    *TPosting
		want    int
		wantErr bool
	}{
		{"1", p1, len1, false},
		{"2", p2, len2, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.post.Store()
			if (err != nil) != tt.wantErr {
				t.Errorf("TPosting.Store() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TPosting.Store() = '%v',\nwant '%v'", got, tt.want)
			}
		})
		tt.post.Delete() // clean up
	}
} // Test_TPosting_Store()

func Test_TPosting_Time(t *testing.T) {
	prep4Tests()

	tm1 := time.Date(2019, 3, 19, 0, 0, 0, 1, time.Local)
	p1 := NewPosting(time2id(tm1), "")

	tm2 := time.Date(2019, 5, 4, 0, 0, 0, 2, time.Local)
	p2 := NewPosting(time2id(tm2), "")

	tm3 := time.Date(2000, 3, 2, 1, 2, 3, 4, time.UTC)
	p3 := NewPosting(time2id(tm3), "")

	tests := []struct {
		name string
		post *TPosting
		want time.Time
	}{
		{"1", p1, tm1},
		{"2", p2, tm2},
		{"3", p3, tm3},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.post.Time(); !got.Equal(tt.want) {
				t.Errorf("%q: TPosting.Time() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // Test_TPosting_Time()

/* _EoF_ */
