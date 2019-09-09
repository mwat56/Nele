/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package nele

import (
	"testing"
)

func Test_absolute(t *testing.T) {
	bd := "/var/tmp"
	type args struct {
		aBaseDir string
		aDir     string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{" 1", args{bd, ""}, ""},
		{" 2", args{bd, "/opt/"}, "/opt"},
		{" 3", args{bd, "./dir/"}, "/var/tmp/dir"},
		{" 4", args{bd, "../opt/file.txt"}, "/var/opt/file.txt"},
		{" 5", args{bd, "./bla.doc"}, "/var/tmp/bla.doc"},
		{" 6", args{"", "dir"}, "/home/matthias/devel/Go/src/github.com/mwat56/Nele/dir"},
		{" 7", args{"", "../../../dir"}, "/home/matthias/devel/Go/src/dir"},
		{" 8", args{"/", "../../../dir"}, "/dir"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := absolute(tt.args.aBaseDir, tt.args.aDir); got != tt.want {
				t.Errorf("absolute() = '%v',\nwant '%v'", got, tt.want)
			}
		})
	}
} // Test_absolute()

func Test_iniData(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{" 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readIniData()
		})
	}
}

/*
// to run this function the `init()` function must be commented out
func Test_initArgs(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initArguments()
		})
	}
} // Test_initArgs()
*/

func Test_kmg2Num(t *testing.T) {
	type args struct {
		aString string
	}
	tests := []struct {
		name     string
		args     args
		wantRInt int
	}{
		// TODO: Add test cases.
		{"0", args{""}, 0},
		{"1", args{"1"}, 1},
		{"2", args{"2kb"}, 2048},
		{"3", args{"3 MB"}, 3145728},
		{"4", args{"4 B"}, 4},
		{"5", args{"5gb"}, 5368709120},
		{"6", args{"10 Mb"}, 10485760},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRInt := kmg2Num(tt.args.aString); gotRInt != tt.wantRInt {
				t.Errorf("kmg2Num() = %v, want %v", gotRInt, tt.wantRInt)
			}
		})
	}
} // Test_kmg2Num()

func TestShowHelp(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{" 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ShowHelp()
		})
	}
} // TestShowHelp()
