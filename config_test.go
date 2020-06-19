/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
               EMail : <support@mwat.de>
*/

package nele

import (
	"flag"
	"reflect"
	"testing"

	"github.com/mwat56/ini"
)

// `parseAppArgsDebug()` calls `parseAppArgs()` and returns `AppArgs`.
//
// This function is meant for unit testing only.
func parseAppArgsDebug() *TAppArgs {
	flag.CommandLine = flag.NewFlagSet(`Nele`, flag.ExitOnError)

	// Define some flags used by `testing` to avoid
	// bailing out during the test.
	var coverprofile, run, testlogfile, timeout string
	flag.CommandLine.StringVar(&coverprofile, `test.coverprofile`, coverprofile,
		"coverprofile for tests")
	flag.CommandLine.StringVar(&run, `test.run`, run,
		"run for tests")
	flag.CommandLine.StringVar(&testlogfile, `test.testlogfile`, testlogfile,
		"testlogfile for tests")
	flag.CommandLine.StringVar(&timeout, `test.timeout`, timeout,
		"timeout for tests")

	parseAppArgs()

	return &AppArgs
} // parseAppArgsDebug()

// `readAppArgDebug()` calls `readAppArgs()` and returns `AppArgs`.
//
// This function is meant for unit testing only.
func readAppArgsDebug() *TAppArgs {
	flag.CommandLine = flag.NewFlagSet(`Nele`, flag.ExitOnError)
	AppArgs = TAppArgs{}

	setAppArgs()
	readAppArgs()

	return &AppArgs
} // readAppArgsDebug()

// `setAppArgsDebug()` calls `setAppArgs()` and returns `AppArgs`.
//
// This function is meant for unit testing only.
func setAppArgsDebug() *TAppArgs {
	flag.CommandLine = flag.NewFlagSet(`Nele`, flag.ExitOnError)
	AppArgs = TAppArgs{}

	var ini1 ini.TIniList
	// Clear/reset the INI values to simulate missing INI file(s):
	iniValues = tArguments{*ini1.GetSection(``)}

	setAppArgs()

	return &AppArgs
} // setAppArgsDebug()

/* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * */

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
		{" 1", args{bd, ``}, `/home/matthias/devel/Go/src/github.com/mwat56/nele`},
		{" 2", args{bd, "/opt/"}, "/opt"},
		{" 3", args{bd, "./dir/"}, "/var/tmp/dir"},
		{" 4", args{bd, "../opt/file.txt"}, "/var/opt/file.txt"},
		{" 5", args{bd, "./bla.doc"}, "/var/tmp/bla.doc"},
		{" 6", args{"", "dir"}, "/home/matthias/devel/Go/src/github.com/mwat56/nele/dir"},
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

func Test_kmg2Num(t *testing.T) {
	type args struct {
		aString string
	}
	tests := []struct {
		name     string
		args     args
		wantRInt int64
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

func Test_parseAppArgsDebug(t *testing.T) {
	expected := &TAppArgs{}
	tests := []struct {
		name string
		want *TAppArgs
	}{
		// TODO: Add test cases.
		{" 1", expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAppArgsDebug(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAppArgsDebug() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_parseAppArgsDebug()

func Test_readAppArgsDebug(t *testing.T) {
	expected := &TAppArgs{
		Addr:          `127.0.0.1:8181`,
		BlogName:      `<! BlogName not configured !>`,
		DataDir:       `/home/matthias/devel/Go/src/github.com/mwat56/nele`,
		delWhitespace: true,
		GZip:          true,
		HashFile:      "/home/matthias/devel/Go/src/github.com/mwat56/nele/HashFile.db",
		Lang:          `en`,
		listen:        `127.0.0.1`,
		MaxFileSize:   10485760,
		mfs:           `10485760`,
		port:          8181,
		Realm:         `My Blog`,
		Theme:         `dark`,
	}
	tests := []struct {
		name string
		want *TAppArgs
	}{
		// TODO: Add test cases.
		{" 1", expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readAppArgsDebug(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readAppArgsDebug() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_readAppArgsDebug()

func Test_setAppArgsDebug(t *testing.T) {
	expected := &TAppArgs{
		BlogName:      `<! BlogName not configured !>`,
		DataDir:       `/home/matthias/devel/Go/src/github.com/mwat56/nele`,
		delWhitespace: true,
		GZip:          true,
		HashFile:      "/home/matthias/devel/Go/src/github.com/mwat56/nele/HashFile.db",
		Lang:          `en`,
		listen:        `127.0.0.1`,
		mfs:           `10485760`,
		port:          8181,
		Realm:         `My Blog`,
		Theme:         `dark`,
	}
	tests := []struct {
		name string
		want *TAppArgs
	}{
		// TODO: Add test cases.
		{" 1", expected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setAppArgsDebug(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setAppArgsDebug() = %s,\nwant %s", got, tt.want)
			}
		})
	}
} // Test_setAppArgsDebug()

func TestShowHelp(t *testing.T) {
	_ = setAppArgsDebug()
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
