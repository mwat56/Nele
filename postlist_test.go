package blog

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestNewPostList(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	wl1 := &TPostList{
		// TPosting{
		// 	/* basedir: bd, */
		// 	id: "~~~~~~~~~~~~~~~~",
		// },
	}
	// type args struct {
	// 	aBaseDir string
	// }
	tests := []struct {
		name string
		// args args
		want *TPostList
	}{
		// TODO: Add test cases.
		{" 1" /* args{bd}, */, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostList( /* tt.args.aBaseDir */ ); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostList() = %v,\nwant %v", got, tt.want)
			}
		})
	}
} // TestNewPostList()

func TestTPostList_Add(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	p1 := NewPosting( /* bd */ )
	pl1 := NewPostList( /* bd */ )
	wl1 := &TPostList{
		// TPosting{
		// 	basedir: bd,
		// 	id:      "~~~~~~~~~~~~~~~~",
		// },
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
	PostingBaseDirectory = "/tmp/postings/"
	pl1 := NewPostList( /* bd */ )
	pl2 := NewPostList( /* bd */ )
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
		{" 1", pl1, args{"1580002c0c472200"}, 1},
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
	PostingBaseDirectory = "/tmp/postings/"
	p1 := NewPosting( /* bd */ ).Set([]byte("# Hello World!"))
	p2 := NewPosting( /* bd */ ).Set([]byte("I trust you're feeling good."))
	p3 := NewPosting( /* bd */ ).Set([]byte("Goodbye!"))
	pl1 := NewPostList( /* bd */ ).Add(p1).Add(p2).Add(p3)
	wl1 := &TPostList{
		// TPosting{
		// 	basedir: bd,
		// 	id:      "~~~~~~~~~~~~~~~~",
		// },
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
				t.Errorf("TPostList.in() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTPostList_in()

func TestTPostList_Len(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	p1 := NewPosting( /* bd */ ).Set([]byte("11"))
	p2 := NewPosting( /* bd */ ).Set([]byte("22"))
	p3 := NewPosting( /* bd */ ).Set([]byte("33"))
	p4 := NewPosting( /* bd */ ).Set([]byte("44"))
	pl1 := NewPostList( /* bd */ ).Add(p3).Add(p1).Add(p2)
	pl2 := NewPostList( /* bd */ ).Add(p1).Add(p2).Add(p3).Add(p4)
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
	PostingBaseDirectory = "/tmp/postings/"
	p1 := newPosting( /* bd,  */ "11").Set([]byte("11"))
	p2 := newPosting( /* bd,  */ "22").Set([]byte("22"))
	p3 := newPosting( /* bd,  */ "33").Set([]byte("33"))
	pl1 := NewPostList( /* bd */ ).Add(p2).Add(p3).Add(p1)
	wl1 := NewPostList( /* bd */ ).Add(p3).Add(p2).Add(p1)
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
	PostingBaseDirectory = "/tmp/postings/"
	p1 := newPosting( /* bd,  */ "11").Set([]byte("11"))
	p2 := newPosting( /* bd,  */ "22").Set([]byte("22"))
	p3 := newPosting( /* bd,  */ "33").Set([]byte("33"))
	pl1 := NewPostList( /* bd */ ).Add(p3).Add(p1).Add(p2)
	pl2 := NewPostList( /* bd */ ).Add(p3).Add(p2).Add(p1)
	pl3 := NewPostList( /* bd */ ).Add(p2).Add(p3).Add(p1).Sort()
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
	// n := time.Now()
	// y, m := n.Year(), n.Month()
	y, m := 2018, time.December
	t := time.Date(y, m, aDay, aHour, aHour, aHour, 0, time.Local)
	p := newPosting( /* aBaseDir,  */ newID(t)).
		Set([]byte(fmt.Sprintf("\n> %s\n\n%02d\n\n\t%02d\n", aBaseDir, aDay, aHour)))
	p.Store()
} // storeNewPost()

func prepareTestFiles( /* aBaseDir string */ ) {
	bd, _ := filepath.Abs(PostingBaseDirectory /* aBaseDir */)
	for i := 0; i < 111; i++ {
		storeNewPost(bd, i, 1)
		storeNewPost(bd, i, 8)
		storeNewPost(bd, i, 16)
	}
} // prepareTestFiles()

func TestTPostList_Month(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	prepareTestFiles( /* bd */ )
	pl1 := NewPostList( /* bd */ )
	pl2 := NewPostList( /* bd */ )
	pl3 := NewPostList( /* bd */ )
	pl4 := NewPostList( /* bd */ )
	// pl5 := NewPostList(bd)
	// pl6 := NewPostList(bd)
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
		{" 1", pl1, args{2019, 1}, 93},
		{" 2", pl2, args{2019, 2}, 84},
		{" 3", pl3, args{2019, 3}, 60},
		{" 4", pl4, args{2019, 4}, 0},
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
	PostingBaseDirectory = "/tmp/postings/"
	prepareTestFiles( /* bd */ )
	pl1 := NewPostList( /* bd */ )
	type args struct {
		aNumber int
	}
	tests := []struct {
		name    string
		pl      *TPostList
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", pl1, args{10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.pl.Newest(tt.args.aNumber); (err != nil) != tt.wantErr {
				t.Errorf("TPostList.Newest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestTPostList_Newest()

func TestTPostList_Week(t *testing.T) {
	PostingBaseDirectory = "/tmp/postings/"
	prepareTestFiles( /* bd */ )
	pl1 := NewPostList( /* bd */ )
	pl2 := NewPostList( /* bd */ )
	pl3 := NewPostList( /* bd */ )
	pl4 := NewPostList( /* bd */ )
	// pl5 := NewPostList(/* bd */)
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
		{" 2", pl2, args{2019, 1, 1}, 21},
		{" 3", pl3, args{2019, 2, 2}, 21},
		{" 4", pl4, args{2019, 3, 3}, 21},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pl.Week(tt.args.aYear, tt.args.aMonth, tt.args.aDay); got.Len() != tt.want {
				t.Errorf("TPostList.Week() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTPostList_Week()
