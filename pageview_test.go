/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
	"time"

	"github.com/mwat56/pageview"
)

func Test_checkForImgURL(t *testing.T) {
	var t1 []byte
	var l1 tImgURLlist
	t2 := []byte(`bla \n> [![„Wir sind alle Opfer hier“](/img/httpswwwaddendumorgnewsopferstudium.png)](https://www.addendum.org/news/opferstudium/)\n bla`)
	l2 := tImgURLlist{
		tImgURL{
			`httpswwwaddendumorgnewsopferstudium.png`,
			`https://www.addendum.org/news/opferstudium/`,
		},
	}
	t3 := []byte(`bla \n> [![„Wir sind alle Opfer hier“](/img/httpswwwaddendumorgnewsopferstudium.png)](https://www.addendum.org/news/opferstudium/)\n bla\n> [![„Radikal den Kontakt abbrechen“](/img/httpswwwspiegeldepanoramahaeuslichegewaltwennmuttergeschlagenwirdwasmachtdasmitdenkinderna1291534html.png)](https://www.spiegel.de/panorama/haeusliche-gewalt-wenn-mutter-geschlagen-wird-was-macht-das-mit-den-kindern-a-1291534.html)`)
	l3 := tImgURLlist{
		tImgURL{
			`httpswwwaddendumorgnewsopferstudium.png`,
			`https://www.addendum.org/news/opferstudium/`,
		},
		tImgURL{
			`httpswwwspiegeldepanoramahaeuslichegewaltwennmuttergeschlagenwirdwasmachtdasmitdenkinderna1291534html.png`,
			`https://www.spiegel.de/panorama/haeusliche-gewalt-wenn-mutter-geschlagen-wird-was-macht-das-mit-den-kindern-a-1291534.html`,
		},
	}
	t4 := []byte(`> @Google holt sich [Anti-Gewerkschafts-Beratung](https://www.heise.de/newsticker/meldung/Google-holt-sich-Anti-Gewerkschafts-Beratung-4593692.html?view=print).`)

	type args struct {
		aTxt []byte
	}
	tests := []struct {
		name      string
		args      args
		wantRList tImgURLlist
	}{
		// TODO: Add test cases.
		{" 1", args{t1}, l1},
		{" 2", args{t2}, l2},
		{" 3", args{t3}, l3},
		{" 4", args{t4}, l1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := checkForImageURL(tt.args.aTxt); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("checkForImgURL() = %v,\nwant %v", gotRList, tt.wantRList)
			}
		})
	}
} // Test_checkForImgURL()

func Test_checkPageImages(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	pageview.SetMaxAge(1)
	var p1 TPosting
	p2 := NewPosting("15d9c2334fce3991")
	p3 := NewPosting("15d9393f4f5f3bb4")

	type args struct {
		aPosting *TPosting
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{&p1}},
		{" 2", args{p2}},
		{" 3", args{p3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPageImages(tt.args.aPosting)
		})
	}
} // Test_checkPageImages()

func Test_setPostingLinkViews(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	pageview.SetMaxAge(1)
	imgURLdir := "/img/"
	var p0 TPosting
	p1 := NewPosting("15d678172cfc527a")
	_ = p1.Load()
	p2 := NewPosting("15d9c2334fce3991")
	_ = p2.Load()
	p3 := NewPosting("15d9393f4f5f3bb4")
	_ = p3.Load()
	p4 := NewPosting("15d93196ab1b2899")
	_ = p4.Load()
	p5 := NewPosting("15d8b372f3186303")
	_ = p5.Load()
	p6 := NewPosting("15dbb86d6c2cdc2c")
	_ = p6.Load()
	type args struct {
		aPosting        *TPosting
		aImageDirectory string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 0", args{&p0, imgURLdir}},
		{" 1", args{p1, imgURLdir}},
		{" 2", args{p2, imgURLdir}},
		{" 3", args{p3, imgURLdir}},
		{" 4", args{p4, imgURLdir}},
		{" 5", args{p5, imgURLdir}},
		{" 6", args{p6, imgURLdir}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setPostingLinkViews(tt.args.aPosting, tt.args.aImageDirectory)
		})
	}
} // Test_setPostingLinkViews()

func Test_goUpdateAllLinkPreviews(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	pageview.SetMaxAge(1)
	type args struct {
		aPostingBaseDir string
		aImageURLdir    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{`/tmp/`, `/img/`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goUpdateAllLinkPreviews(tt.args.aPostingBaseDir, tt.args.aImageURLdir)
			time.Sleep(time.Second)
		})
	}
	time.Sleep(time.Second)
} // Test_goUpdateAllLinkPreviews()

func TestRemovePagePreviews(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	var t1 TPosting
	t2 := NewPosting("")
	t3 := NewPosting("")
	type args struct {
		aPosting *TPosting
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{&t1}},
		{" 2", args{t2}},
		{" 3", args{t3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemovePagePreviews(tt.args.aPosting)
		})
	}
} // TestRemovePagePreviews()
