package main

import (
	"testing"
)

func TestMd2Html(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(); (err != nil) != tt.wantErr {
				t.Errorf("Md2Html() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} // TestMd2Html()
