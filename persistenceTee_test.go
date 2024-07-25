/*
Copyright © 2024  M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"reflect"
	"testing"
)

func prepareTee() {
	AppArgs.persistence = "tee"
	prepareTestFiles()
} // prepareTee()

func TestTeePersistence_Count(t *testing.T) {
	prepareTee()

	tests := []struct {
		name string
		tp   *TeePersistence
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.Count(); got != tt.want {
				t.Errorf("%q: TTeePersistence.Count() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Count()

func TestTeePersistence_Create(t *testing.T) {
	prepareTee()

	tests := []struct {
		name    string
		tp      *TeePersistence
		post    *TPosting
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tp.Create(tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Create() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%q: TTeePersistence.Create() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Create()

func TestTeePersistence_Delete(t *testing.T) {
	prepareTee()

	tests := []struct {
		name    string
		tp      *TeePersistence
		id      uint64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.tp.Delete(tt.id); (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Delete() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
			}
		})
	}
} // TestTeePersistence_Delete()

func TestTeePersistence_Exists(t *testing.T) {
	prepareTee()

	tests := []struct {
		name string
		tp   *TeePersistence
		id   uint64
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.Exists(tt.id); got != tt.want {
				t.Errorf("%q: TeePersistence.Exists() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Exists()

func TestTeePersistence_PathFileName(t *testing.T) {
	prepareTee()

	tests := []struct {
		name string
		tp   *TeePersistence
		id   uint64
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tp.PathFileName(tt.id); got != tt.want {
				t.Errorf("%q: TeePersistence.PathFileName() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_PathFileName()

func TestTeePersistence_Read(t *testing.T) {
	prepareTee()

	tests := []struct {
		name    string
		tp      *TeePersistence
		id      uint64
		want    *TPosting
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tp.Read(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Read() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TeePersistence.Read() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Read()

func TestTeePersistence_Search(t *testing.T) {
	prepareTee()

	type tArgs struct {
		aText   string
		aOffset uint
		aLimit  uint
	}
	tests := []struct {
		name    string
		tp      *TeePersistence
		args    tArgs
		want    *TPostList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tp.Search(tt.args.aText, tt.args.aOffset, tt.args.aLimit)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Search() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%q: TeePersistence.Search() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Search()

func TestTeePersistence_Update(t *testing.T) {
	prepareTee()

	tests := []struct {
		name    string
		tp      *TeePersistence
		post    *TPosting
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tp.Update(tt.post)
			if (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Update() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("%q: TeePersistence.Update() = %v, want %v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestTeePersistence_Update()

func TestTeePersistence_Walk(t *testing.T) {
	prepareTee()

	tests := []struct {
		name    string
		tp      *TeePersistence
		wf      TWalkFunc
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tp.Walk(tt.wf); (err != nil) != tt.wantErr {
				t.Errorf("%q: TeePersistence.Walk() error = %v, wantErr %v",
					tt.name, err, tt.wantErr)
			}
		})
	}
} // TestTeePersistence_Walk()

/* _EoF_ */
