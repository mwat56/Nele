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

func TestRemoveWhiteSpace(t *testing.T) {
	txtIn1 := []byte(`<hr />
	<p>Here is an example of AppleScript:</p>
	<pre>
	tell application &quot;Foo&quot;
	  beep
	end tell
	</pre>
	<hr />`)
	txtOut1 := []byte(`<hr /><p>Here is an example of AppleScript:</p><pre>
	tell application &quot;Foo&quot;
	  beep
	end tell
	</pre><hr />`)
	type args struct {
		aPage []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
		{" 1", args{txtIn1}, txtOut1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveWhiteSpace(tt.args.aPage); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveWhiteSpace() = [%s],\nwant [%s]", got, tt.want)
			}
		})
	}
} // TestRemoveWhiteSpace()
