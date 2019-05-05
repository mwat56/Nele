package blog

import (
	"testing"
)

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
