/*
Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func prepareTestFiles() {
	prep4Tests()

	bd, _ := filepath.Abs(PostingBaseDirectory())
	for i := 1; i < 13; i++ {
		storeNewPost(bd, i, 1)
		storeNewPost(bd, i, 8)
		storeNewPost(bd, i, 16)
	}
} // prepareTestFiles()

func storeNewPost(aBaseDir string, aDay, aHour int) {
	t := time.Date(1970, 1, aDay, aHour, aHour, aHour, 0, time.Local)
	p := NewPosting(time2id(t), "")
	p.Set([]byte(fmt.Sprintf("\n> %s\n\n%s\n\n@someone said%02d\n\n\t%02d\n#wewantitall%d", p.Date(), aBaseDir, aDay, aHour, aDay)))
	_, _ = p.Store()

	t = time.Date(2018, 12, aDay, aHour, aHour, aHour, 0, time.Local)
	p = NewPosting(time2id(t), "")
	p.Set([]byte(fmt.Sprintf("\n> %s\n\n%s\n\n@someone said%02d\n\n\t%02d\n#wewantitall%d", p.Date(), aBaseDir, aDay, aHour, aDay)))
	_, _ = p.Store()
} // storeNewPost()

func TestNewPostList(t *testing.T) {
	prepareTestFiles()

	wl1 := &TPostList{}
	tests := []struct {
		name string
		want *TPostList
	}{
		// TODO: Add test cases.
		{" 1", wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostList() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestNewPostList()

func TestSearchPostings(t *testing.T) {
	prepareTestFiles()

	tests := []struct {
		name string
		text string
		want int
	}{
		// TODO: Add test cases.
		{"1", "16", 24},
		{"2", "8", 50},
		{"3", "1\\d+", 72},
		{"4", "10\\d+", 0},
		{"5", "08\\s+08", 2},
		{"6", "postings", 72},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SearchPostings(tt.text); got.Len() != tt.want {
				t.Errorf("SearchPostings() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestSearchPostings()

func TestTPostList_Add(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(0, "")
	pl1 := NewPostList()
	wl1 := &TPostList{
		*p1,
	}
	type args struct {
		aPosting *TPosting
	}
	tests := []struct {
		name string
		pl   *TPostList
		args args
		want *TPostList
	}{
		// TODO: Add test cases.
		{" 1", pl1, args{p1}, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Add(tt.args.aPosting); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPostList.Add() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPostList_Add()

func TestTPostList_Delete(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(0, "")
	pl1 := NewPostList()
	wl1 := NewPostList()
	wb1 := false
	// ---
	p2 := NewPosting(0, "")
	pl2 := NewPostList().Add(p2)
	wl2 := NewPostList()
	wb2 := true
	// ---
	p3 := NewPosting(0, "")
	pl3 := NewPostList().Add(p1).Add(p2).Add(p3)
	wl3 := NewPostList().Add(p1).Add(p2)
	wb3 := true
	// ---
	p4 := NewPosting(0, "")
	pl4 := NewPostList().Add(p1).Add(p2).Add(p4).Add(p3)
	wl4 := NewPostList().Add(p1).Add(p2).Add(p3)
	wb4 := true
	// ---
	p5 := NewPosting(0, "")
	pl5 := NewPostList().Add(p1).Add(p2).Add(p3).Add(p4)
	wl5 := NewPostList().Add(p1).Add(p2).Add(p3).Add(p4)
	wb5 := false
	// ---

	tests := []struct {
		name     string
		pl       *TPostList
		aPosting *TPosting
		want     *TPostList
		want1    bool
	}{
		// TODO: Add test cases.
		{" 1", pl1, p1, wl1, wb1},
		{" 2", pl2, p2, wl2, wb2},
		{" 3", pl3, p3, wl3, wb3},
		{" 4", pl4, p4, wl4, wb4},
		{" 5", pl5, p5, wl5, wb5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.pl.Delete(tt.aPosting)
			if got1 != tt.want1 {
				t.Errorf("TPostList.Delete() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPostList.Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPostList_Delete()

func TestTPostList_insert(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(111, "> 111")
	p2 := NewPosting(222, "> 222")
	p3 := NewPosting(333, "> 333")
	p4 := NewPosting(444, "> 444")
	p5 := NewPosting(188, "> 188")

	pl1 := NewPostList()
	pl2 := NewPostList().Add(p3).Add(p2).Add(p1)
	pl3 := NewPostList().Add(p1).Add(p2).Add(p3) //.Sort()

	tests := []struct {
		name string
		pl   *TPostList
		post *TPosting
		want bool
	}{
		{"1", pl1, p1, true},  // first entry
		{"2", pl2, p4, true},  // after last entry
		{"3", pl3, p2, false}, // middle existing
		{"4", pl2, p5, true},  // middle/new
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.insert(tt.post); got != tt.want {
				t.Errorf("%q: TPostList.insert() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTPostList_insert()

func TestTPostList_IsSorted(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(11, "11")
	p2 := NewPosting(22, "22")
	p3 := NewPosting(33, "33")
	pl1 := NewPostList().Add(p3).Add(p1).Add(p2)
	pl2 := NewPostList().Add(p3).Add(p2).Add(p1)
	pl3 := NewPostList().Add(p2).Add(p3).Add(p1) // .Sort()

	tests := []struct {
		name string
		pl   *TPostList
		want bool
	}{
		{"1", pl1, true},
		{"2", pl2, true},
		{"3", pl3, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.IsSorted(); got != tt.want {
				t.Errorf("%q: TPostList.IsSorted() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTPostList_IsSorted()

func TestTPostList_Len(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(0, "").Set([]byte("11"))
	p2 := NewPosting(0, "").Set([]byte("22"))
	p3 := NewPosting(0, "").Set([]byte("33"))
	p4 := NewPosting(0, "").Set([]byte("44"))
	pl1 := NewPostList().Add(p3).Add(p1).Add(p2)
	pl2 := NewPostList().Add(p1).Add(p2).Add(p3).Add(p4)
	tests := []struct {
		name string
		pl   *TPostList
		want int
	}{
		// TODO: Add test cases.
		{" 1", pl1, 3},
		{" 2", pl2, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Len(); got != tt.want {
				t.Errorf("TPostList.Len() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPostList_Len()

func TestTPostList_Month(t *testing.T) {
	prepareTestFiles()

	pl1 := NewPostList()
	pl2 := NewPostList()
	pl3 := NewPostList()
	pl4 := NewPostList()

	type tArgs struct {
		aYear  int
		aMonth time.Month
	}
	tests := []struct {
		name string
		pl   *TPostList
		args tArgs
		want int
	}{
		// TODO: Add test cases.
		{" 1", pl1, tArgs{1970, 1}, 36},
		{" 2", pl2, tArgs{2018, 12}, 36},
		{" 3", pl3, tArgs{1970, 6}, 0},
		{" 4", pl4, tArgs{2018, 12}, 36},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Month(tt.args.aYear, tt.args.aMonth); got.Len() != tt.want {
				t.Errorf("%q: TPostList.Month() = %v, want %v", tt.name, got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Month()

func TestTPostList_Newest(t *testing.T) {
	prepareTestFiles()

	pl1 := NewPostList()
	type tArgs struct {
		aNumber int
		aStart  int
	}
	tests := []struct {
		name    string
		pl      *TPostList
		args    tArgs
		wantErr bool
	}{
		// TODO: Add test cases.
		{"1", pl1, tArgs{10, 0}, false},
		{"2", pl1, tArgs{10, 10}, false},
		{"3", pl1, tArgs{5, 15}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pl.Newest(tt.args.aNumber, tt.args.aStart); (err != nil) != tt.wantErr {
				t.Errorf("%q: TPostList.Newest() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if pLen := tt.pl.Len() - 1; tt.args.aNumber < pLen {
				t.Errorf("%q: TPostList.Newest() number = %d, wanted %d",
					tt.name, pLen, tt.args.aNumber)
			}
		})
	}
} // TestTPostList_Newest()

func TestTPostList_Sort(t *testing.T) {
	prepareTestFiles()

	p1 := NewPosting(11, "> 11")
	p2 := NewPosting(22, "> 22")
	p3 := NewPosting(33, "> 33")
	pl1 := NewPostList().Add(p2).Add(p3).Add(p1)
	wl1 := NewPostList().Add(p3).Add(p2).Add(p1)

	tests := []struct {
		name string
		pl   *TPostList
		want *TPostList
	}{
		{"1", pl1, wl1},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TPostList.Sort() =\n%v\n>>> want >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTPostList_Sort()

func TestTPostList_Week(t *testing.T) {
	prepareTestFiles()

	pl1 := NewPostList()
	pl2 := NewPostList()
	pl3 := NewPostList()
	pl4 := NewPostList()

	type args struct {
		aYear  int
		aMonth time.Month
		aDay   int
	}
	tests := []struct {
		name string
		pl   *TPostList
		args args
		want int
	}{
		// TODO: Add test cases.
		{"1", pl1, args{1, 1, 1}, 0},
		{"2", pl2, args{1970, 1, 1}, 12},
		{"3", pl3, args{2018, 12, 1}, 6},
		{"4", pl4, args{2018, 12, 8}, 21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Week(tt.args.aYear, tt.args.aMonth, tt.args.aDay); got.Len() != tt.want {
				t.Errorf("TPostList.Week() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Week()

/* _EoF_ */
