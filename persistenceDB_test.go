/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

func TestTDBpersistence_init(t *testing.T) {
	tests := []struct {
		name string
		dbp  *TDBpersistence
	}{
		{
			name: "NilPointerAfterSettingToAnotherType",
			dbp: func() *TDBpersistence {
				dbp := NewDBpersistence("test.db")
				var _ IPersistence = dbp
				return nil
			}(),
		},
		{
			name: "NilPointerAfterSettingToNil",
			dbp: func() *TDBpersistence {
				var dbp *TDBpersistence
				var _ IPersistence = dbp
				return nil
			}(),
		},
		{
			name: "NilPointerAfterSettingToNilAndThenSettingToAnotherType",
			dbp: func() *TDBpersistence {
				var dbp *TDBpersistence
				var _ IPersistence = dbp
				dbp = NewDBpersistence("test.db")
				if nil == dbp {
					return nil
				}
				return nil
			}(),
		},
		{
			name: "NilPointerAfterSettingToAnotherTypeAndThenSettingToNil",
			dbp: func() *TDBpersistence {
				dbp := NewDBpersistence("test.db")
				var _ IPersistence = dbp
				dbp = nil
				return nil
			}(),
		},
		{
			name: "NilPointerAfterSettingToNilAndThenSettingToNilAgain",
			dbp: func() *TDBpersistence {
				var dbp *TDBpersistence
				var _ IPersistence = dbp
				dbp = nil
				return nil
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dbp != nil {
				t.Errorf("Expected nil, but got %v", tt.dbp)
			}
		})
	}
} // TestTDBpersistence_init()

func TestNewDBpersistence(t *testing.T) {
	fn1 := "tstDB1.db"
	fName := filepath.Join(poPostingBaseDirectory, fn1)
	defer func() {
		os.Remove(fName)
	}()

	var (
		dbInstance *sql.DB
		hasFTS     bool
	)
	dbInstance, hasFTS, _ = initDatabase(fName)
	wd1 := &TDBpersistence{
		db:   dbInstance,
		fts5: hasFTS,
	}

	tests := []struct {
		name    string
		fName   string
		want    *TDBpersistence
		wantNIL bool
	}{
		{"1", fn1, wd1, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDBpersistence(tt.fName)
			if (tt.wantNIL && (nil != got)) || ((nil == got) && !tt.wantNIL) {
				t.Errorf("%q: NewDBpersistence() = \n%v\n>>> want >>>\n%v",
					tt.name, got, tt.want)
			}
		})
	}
} // TestNewDBpersistence()
