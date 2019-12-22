/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"

	"github.com/mwat56/hashtags"
)

func Test_MarkupTags(t *testing.T) {
	p1 := []byte(`bla #hash1 bla _@mention1_ bla&#39; <a href="page#fragment">bla</a>
	[link text](http://host.com/page#frag2) #hash2`)
	w1 := []byte(`bla <a href="/hl/hash1" class="smaller">#hash1</a> bla _<a href="/ml/mention1" class="smaller">@mention1</a>_ bla&#39; <a href="page#fragment">bla</a>
	[link text](http://host.com/page<a href="/hl/frag2" class="smaller">#frag2</a>) <a href="/hl/hash2" class="smaller">#hash2</a>`)
	p2 := []byte(`
*@Antoni_Comín @Carles_Puigdemont @Catalonia @EU @Immunity @Oriol_Junqueras @Spain*`)
	w2 := []byte(`
*<a href="/ml/Antoni_Comín" class="smaller">@Antoni_Comín</a> <a href="/ml/Carles_Puigdemont" class="smaller">@Carles_Puigdemont</a> <a href="/ml/Catalonia" class="smaller">@Catalonia</a> <a href="/ml/EU" class="smaller">@EU</a> <a href="/ml/Immunity" class="smaller">@Immunity</a> <a href="/ml/Oriol_Junqueras" class="smaller">@Oriol_Junqueras</a> <a href="/ml/Spain" class="smaller">@Spain</a>*`)
	type args struct {
		aPage []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 1", args{p1}, w1},
		{" 2", args{p2}, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarkupTags(tt.args.aPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarkupTags() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_MarkupTags()

func TestReplaceTag(t *testing.T) {
	hashtags.UseBinaryStorage = false
	l1, _ := hashtags.New(`./TestReplaceTag.db`)
	type args struct {
		aList       *hashtags.THashList
		aSearchTag  string
		aReplaceTag string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{l1, ``, ``}},
		{" 2", args{l1, `@chelseamanning`, `Chelsea_Manning`}},
		{" 3", args{l1, `chelseamanning`, `@Chelsea_Manning`}},
		{" 4", args{l1, `@chelseamanning`, `@Chelsea_Manning`}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReplaceTag(tt.args.aList, tt.args.aSearchTag, tt.args.aReplaceTag)
		})
	}
} // TestReplaceTag()
