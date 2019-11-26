/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package nele

//lint:file-ignore ST1017 - I prefer Yoda conditions

import (
	"testing"

	"github.com/mwat56/pageview"
)

func Test_goSetPostingLinkViews(t *testing.T) {
	p1 := ``
	p2 := `/home/matthias/devel/Go/src/github.com/mwat56/Nele/postings/201915d/15da6104b009723d.md`
	imgDir := "./img/"
	pageview.SetCacheDirectory(imgDir)
	pageview.SetMaxAge(1)
	type args struct {
		aPosting, aImageDirectory string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{" 1", args{p1, imgDir}},
		{" 2", args{p2, imgDir}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goSetPostingLinkViews(tt.args.aPosting, tt.args.aImageDirectory)
		})
	}
} // Test_goSetPostingLinkViews()
