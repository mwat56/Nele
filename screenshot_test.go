/*
Copyright © 2022, 2023 M.Watermann, 10247 Berlin, Germany

	    All rights reserved
	EMail : <support@mwat.de>
*/
package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/mwat56/screenshot"
)

func Test_checkScreenshotURLs(t *testing.T) {
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
			if gotRList := checkScreenshotURLs(tt.args.aTxt); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("checkForImgURL() = %v,\nwant %v", gotRList, tt.wantRList)
			}
		})
	}
} // Test_checkScreenshotURLs()

func Test_checkScreenshots(t *testing.T) {
	screenshot.SetImageDir("/tmp/")
	screenshot.SetImageAge(1)
	p1 := NewPosting("15d9c2334fce3991")
	p2 := NewPosting("15d9393f4f5f3bb4")

	type args struct {
		aPosting *TPosting
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{p1}},
		{" 2", args{p2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkScreenshots(tt.args.aPosting)
		})
	}
} // Test_checkScreenshots()

func Test_goUpdateAllLinkScreenshots(t *testing.T) {
	screenshot.SetImageDir("/tmp/")
	screenshot.SetImageAge(1)
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
			goUpdateAllLinkScreenshots(tt.args.aPostingBaseDir /*, tt.args.aImageURLdir*/)
			time.Sleep(time.Second)
		})
	}
	time.Sleep(time.Second)
} // Test_goUpdateAllLinkScreenshots()

func Test_goSetLinkScreenshots(t *testing.T) {
	screenshot.SetImageDir("/tmp/")
	screenshot.SetImageAge(1)
	imgURLdir := "/img/"
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
		{" 1", args{p1, imgURLdir}},
		{" 2", args{p2, imgURLdir}},
		{" 3", args{p3, imgURLdir}},
		{" 4", args{p4, imgURLdir}},
		{" 5", args{p5, imgURLdir}},
		{" 6", args{p6, imgURLdir}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goSetLinkScreenshots(tt.args.aPosting /*, tt.args.aImageDirectory*/)
		})
	}
} // Test_goSetLinkScreenshots()

func Test_prepPostText(t *testing.T) {
	imageURLdir := `/img/`
	p1 := []byte(`How the use of words changed over time in German parliament.

> [How has political language and the use of certain words changed over time?](https://mobile.twitter.com/MSchories/status/1235489948876320770)

*@Bundestag @Germany #Language #Parliament*`)
	i1 := `httpsmobiletwittercomMSchoriesstatus1235489948876320770`
	w1 := []byte(`How the use of words changed over time in German parliament.

> [![How has political language and the use of certain words changed over time?](` + imageURLdir + i1 + `)](https://mobile.twitter.com/MSchories/status/1235489948876320770)

*@Bundestag @Germany #Language #Parliament*`)

	p2 := []byte("bla \n> [link text two (a)](https://www.example.org/two/)\n bla\n > bla [„link text three“](https://www.example.org/three) bla")
	i2 := `httpswwwexampleorgtwo`
	w2 := []byte("bla \n> [![link text two (a)](/img/httpswwwexampleorgtwo)](https://www.example.org/two/)\n bla\n > bla [„link text three“](https://www.example.org/three) bla")

	type args struct {
		aPosting   []byte
		aImageName string
	}
	tests := []struct {
		name      string
		args      args
		wantRText []byte
	}{
		// TODO: Add test cases.
		{" 1", args{p1, i1}, w1},
		{" 2", args{p2, i2}, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linkMatches := ssLinkRE2.FindAllSubmatch(tt.args.aPosting, -1)
			link := &tLink{
				link:     string(linkMatches[0][1]),
				linkText: string(linkMatches[0][2]),
				linkURL:  string(linkMatches[0][3]),
			}
			if gotRText := prepPostText(tt.args.aPosting, link, tt.args.aImageName, imageURLdir); !reflect.DeepEqual(gotRText, tt.wantRText) {
				t.Errorf("prepPostText() = {%v},\nwant {%v}", string(gotRText), string(tt.wantRText))
			}
		})
	}
} // Test_prepPostText()

// R/O RegEx to extract link-text and link-URL from markup.
// Checking for the not-existence of the leading `!` should exclude
// embedded image links.
var ssLinkRE2 = regexp.MustCompile(
	`(?m)(?:^\s*\>[\t ]*)((?:[^\!\n\>][\t ]*)?\[([^\[]+?)\]\s*\(([^\]]+?)\))`)

//                                              122222222111111133333333311
// `[link-text](link-url)`
// 0 : complete RegEx match
// 1 : markdown link markup
// 2 : link text
// 3 : remote page URL

func Test_ssImageRE(t *testing.T) {
	var t1 string
	t2 := "bla \n> [„link text one“](https://www.example.org/one/)\n bla"
	t3 := "bla \n bla [„link text two“](https://www.example.org/two/)\n bla\n > [„link text three“](https://www.example.org/three) bla"
	t4 := `bla bla [link ext four](https://www.example.org/four/) bla.`
	t5 := `bla > bla [link ext five](https://www.example.org/five/) bla.`
	t6 := "bla \n> [![„alt text six“](/img/httpswwwexampleorgsix.png)](https://www.example.org/six/)\n bla"
	t7 := `bla \n> Hi there! [„link text seven](https://www.example.org/seven/)\n bla`
	t8 := `bla \n> ![„alt text eight](/img/httpswwwexampleorgeight.png)\n bla`
	t9 := "bla \n>\n[„link text nine“](https://www.example.org/nine/)\n bla"
	t10 := "\n>	[„link text ten“] (https://www.example.org/ten/) bla"
	t11 := "> [„link text eleven“](https://www.example.org/eleven/)\n bla"
	t12 := "> [„link text eleven“ (b)](https://www.example.org/eleven/)\n bla"

	tests := []struct {
		name     string
		aTxt     string
		matchNum int
	}{
		// TODO: Add test cases.
		{" 1", t1, 0},
		{" 2", t2, 1},
		{" 3", t3, 1},
		{" 4", t4, 0},
		{" 5", t5, 0},
		{" 6", t6, 0},
		{" 7", t7, 0},
		{" 8", t8, 0},
		{" 9", t9, 0},
		{"10", t10, 1},
		{"11", t11, 1},
		{"12", t12, 1},
	}
	var matchLen int

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches := ssLinkRE2.FindAllStringSubmatch(tt.aTxt, -1)
			if nil == gotMatches {
				matchLen = 0
			} else {
				matchLen = len(gotMatches)
			}
			if matchLen != tt.matchNum {
				t.Errorf("Test_ssLinkRE() =\n{%v},\nwant {%v},\n{%v}", matchLen, tt.matchNum, gotMatches)
			}
		})
	}
} // Test_ssImageRE()

func TestRemovePageScreenshots(t *testing.T) {
	screenshot.SetImageDir("/tmp/")
	t1 := NewPosting("")
	t2 := NewPosting("")
	type args struct {
		aPosting *TPosting
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{t1}},
		{" 2", args{t2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RemovePageScreenshots(tt.args.aPosting)
		})
	}
} // TestRemovePageScreenshots()
