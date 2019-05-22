package blog

import (
	"testing"
)

func Test_absolute(t *testing.T) {
	bd := "/var/tmp"
	d1, w1 := "", ""
	d2, w2 := "/opt/", "/opt"
	d3, w3 := "./dir/", "/var/tmp/dir"
	d4, w4 := "../opt/file.txt", "/var/opt/file.txt"
	d5, w5 := "./bla.doc", "/var/tmp/bla.doc"
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
		{" 1", args{bd, d1}, w1},
		{" 2", args{bd, d2}, w2},
		{" 3", args{bd, d3}, w3},
		{" 4", args{bd, d4}, w4},
		{" 5", args{bd, d5}, w5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := absolute(tt.args.aBaseDir, tt.args.aDir); got != tt.want {
				t.Errorf("absolute() = '%v',\nwant '%v'", got, tt.want)
			}
		})
	}
} // Test_absolute()

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
