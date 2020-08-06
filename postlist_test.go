/*
   Copyright © 2019, 2020 M.Watermann, 10247 Berlin, Germany
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

func TestNewPostList(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
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
		{" 1", args{bd, "16"}, 24},
		{" 2", args{bd, "8"}, 50},
		{" 3", args{bd, "1\\d+"}, 72},
		{" 4", args{bd, "10\\d+"}, 0},
		{" 5", args{bd, "08\\s+08"}, 2},
		{" 6", args{bd, bd}, 72},
		{" 7", args{bd, "postings"}, 72},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SearchPostings(tt.args.aText); got.Len() != tt.want {
				t.Errorf("Search() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestSearchPostings()

func TestTPostList_Add(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	p1 := NewPosting("")
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

func TestTPostList_Article(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	pl1 := NewPostList()
	pl2 := NewPostList()
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		pl   *TPostList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", pl1, args{"156dfb3d4f4d7000"}, 1},
		{" 2", pl2, args{"1234567890123456"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Article(tt.args.aID); got.Len() != tt.want {
				t.Errorf("TPostList.Article() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Article()

func TestTPostList_in(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	p1 := NewPosting("").Set([]byte("# Hello World!"))
	p2 := NewPosting("").Set([]byte("I trust you're feeling good."))
	p3 := NewPosting("").Set([]byte("Goodbye!"))
	pl1 := NewPostList().Add(p1).Add(p2).Add(p3)
	wl1 := &TPostList{
		*p1,
		*p2,
		*p3,
	}
	tests := []struct {
		name string
		pl   *TPostList
		want *TPostList
	}{
		// TODO: Add test cases.
		{" 1", pl1, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.in(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPostList.in() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTPostList_in()

func TestTPostList_Len(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	p1 := NewPosting("").Set([]byte("11"))
	p2 := NewPosting("").Set([]byte("22"))
	p3 := NewPosting("").Set([]byte("33"))
	p4 := NewPosting("").Set([]byte("44"))
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

func TestTPostList_Sort(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	p1 := NewPosting("11").Set([]byte("11"))
	p2 := NewPosting("22").Set([]byte("22"))
	p3 := NewPosting("33").Set([]byte("33"))
	pl1 := NewPostList().Add(p2).Add(p3).Add(p1)
	wl1 := NewPostList().Add(p3).Add(p2).Add(p1)
	tests := []struct {
		name string
		pl   *TPostList
		want *TPostList
	}{
		// TODO: Add test cases.
		{" 1", pl1, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Sort(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TPostList.Sort() =\n%v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestTPostList_Sort()

func TestTPostList_IsSorted(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	p1 := NewPosting("11").Set([]byte("11"))
	p2 := NewPosting("22").Set([]byte("22"))
	p3 := NewPosting("33").Set([]byte("33"))
	pl1 := NewPostList().Add(p3).Add(p1).Add(p2)
	pl2 := NewPostList().Add(p3).Add(p2).Add(p1)
	pl3 := NewPostList().Add(p2).Add(p3).Add(p1).Sort()
	tests := []struct {
		name string
		pl   *TPostList
		want bool
	}{
		// TODO: Add test cases.
		{" 1", pl1, false},
		{" 2", pl2, true},
		{" 3", pl3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.IsSorted(); got != tt.want {
				t.Errorf("TPostList.IsSorted() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPostList_IsSorted()

func storeNewPost(aBaseDir string, aDay, aHour int) {
	t := time.Date(1970, 1, aDay, aHour, aHour, aHour, 0, time.Local)
	p := NewPosting(newID(t))
	p.Set([]byte(fmt.Sprintf("\n> %s\n\n%s\n\n@someone said%02d\n\n\t%02d\n#wewantitall%d", p.Date(), aBaseDir, aDay, aHour, aDay)))
	_, _ = p.Store()

	t = time.Date(2018, 12, aDay, aHour, aHour, aHour, 0, time.Local)
	p = NewPosting(newID(t))
	p.Set([]byte(fmt.Sprintf("\n> %s\n\n%s\n\n@someone said%02d\n\n\t%02d\n#wewantitall%d", p.Date(), aBaseDir, aDay, aHour, aDay)))
	_, _ = p.Store()
} // storeNewPost()

func prepareTestFiles() {
	bd, _ := filepath.Abs(PostingBaseDirectory())
	for i := 1; i < 13; i++ {
		storeNewPost(bd, i, 1)
		storeNewPost(bd, i, 8)
		storeNewPost(bd, i, 16)
	}
} // prepareTestFiles()

func TestTPostList_Month(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	prepareTestFiles()
	pl1 := NewPostList()
	pl2 := NewPostList()
	pl3 := NewPostList()
	pl4 := NewPostList()
	type args struct {
		aYear  int
		aMonth time.Month
	}
	tests := []struct {
		name string
		pl   *TPostList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", pl1, args{1970, 1}, 36},
		{" 2", pl2, args{2018, 12}, 36},
		{" 3", pl3, args{1970, 6}, 0},
		{" 4", pl4, args{2018, 12}, 36},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Month(tt.args.aYear, tt.args.aMonth); got.Len() != tt.want {
				t.Errorf("TPostList.Month() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Month()

func TestTPostList_Newest(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	prepareTestFiles()
	pl1 := NewPostList()
	type args struct {
		aNumber int
		aStart  int
	}
	tests := []struct {
		name    string
		pl      *TPostList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", pl1, args{10, 0}, false},
		{" 2", pl1, args{10, 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pl.Newest(tt.args.aNumber, tt.args.aStart); (err != nil) != tt.wantErr {
				t.Errorf("TPostList.Newest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTPostList_Newest()

func TestTPostList_Week(t *testing.T) {
	SetPostingBaseDirectory("/tmp/postings/")
	prepareTestFiles()
	pl1 := NewPostList()
	pl2 := NewPostList()
	pl3 := NewPostList()
	pl4 := NewPostList()
	// pl5 := NewPostList()
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
		{" 1", pl1, args{0, 0, 0}, 0},
		{" 2", pl2, args{1970, 1, 1}, 12},
		{" 3", pl3, args{2018, 12, 1}, 6},
		{" 4", pl4, args{2018, 12, 8}, 21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Week(tt.args.aYear, tt.args.aMonth, tt.args.aDay); got.Len() != tt.want {
				t.Errorf("TPostList.Week() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Week()
