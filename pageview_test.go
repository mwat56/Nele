/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"

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
			if gotRList := checkForImgURL(tt.args.aTxt); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("checkForImgURL() = %v,\nwant %v", gotRList, tt.wantRList)
			}
		})
	}
} // Test_checkForImgURL()

func Test_goCheckPageImages(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	pageview.SetMaxAge(1)
	pageImgDir := pageview.ImageDirectory()
	imgURLdir := "/img/"
	var p1 TPosting
	p2 := NewPosting("15d9c2334fce3991")
	p3 := NewPosting("15d9393f4f5f3bb4")

	type args struct {
		aPosting     *TPosting
		aImageURLdir string
		aImageDir    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{&p1, imgURLdir, pageImgDir}},
		{" 2", args{p2, imgURLdir, pageImgDir}},
		{" 3", args{p3, imgURLdir, pageImgDir}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goCheckPageImages(tt.args.aPosting, tt.args.aImageURLdir, tt.args.aImageDir)
		})
	}
} // Test_goCheckPageImages()

func Test_goSetPostingLinkViews(t *testing.T) {
	pageview.SetImageDirectory("/tmp/")
	pageview.SetMaxAge(1)
	pageImgDir := pageview.ImageDirectory()
	imgURLdir := "/img/"
	var p1 TPosting
	p2 := NewPosting("15d9c2334fce3991")
	type args struct {
		aPosting        *TPosting
		aImageDirectory string
		aCacheDirectory string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{&p1, imgURLdir, pageImgDir}},
		{" 2", args{p2, imgURLdir, pageImgDir}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goSetPostingLinkViews(tt.args.aPosting, tt.args.aImageDirectory, tt.args.aCacheDirectory)
		})
	}
} // Test_goSetPostingLinkViews()