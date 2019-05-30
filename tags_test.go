/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
)

func Test_markupTags(t *testing.T) {
	p1 := []byte(`bla #hash1 bla _@mention1_ bla&#39; <a href="page#fragment">bla</a> [link text](http://host.com/page#frag2) #hash2`)
	w1 := []byte(`bla <a href="/hl/hash1" class="smaller">#hash1</a> bla _<a href="/ml/mention1" class="smaller">@mention1</a>_ bla&#39; <a href="page#fragment">bla</a> [link text](http://host.com/page<a href="/hl/frag2" class="smaller">#frag2</a>) <a href="/hl/hash2" class="smaller">#hash2</a>`)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := markupTags(tt.args.aPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("markupTags() = %s,\nwant %s", got, tt.want)
			}
		})
	}
}
