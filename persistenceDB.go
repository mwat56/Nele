/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

			All rights reserved
		EMail : <support@mwat.de>
*/

package nele

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"sync"
	"time"
	"unsafe"

	_ "github.com/mattn/go-sqlite3"
	se "github.com/mwat56/sourceerror"
)

//lint:file-ignore ST1017 - I prefer Yoda conditions

/* Defined in `persistence.go`:
type (
	TPosting struct {
		id           uint64    // integer representation of date/time
		lastModified time.Time // file modification time
		markdown     []byte    // article contents in Markdown markup
	}

	TPostList []TPosting

	TWalkFunc func(aID uint64) error

	IPersistence interface {
		Create(aPost *TPosting) (int, error)
		Read(aID uint64) (*TPosting, error)
		Update(aPost *TPosting) (int, error)
		Delete(aID uint64) error

		Count() int
		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		Rename(aOldID, aNewID uint64) error
		Search(aText string, aOffset, aLimit uint) (*TPostList, error)
		Walk(aWalkFunc TWalkFunc) error
	}
)
*/

const (
	maxSqliteInt int64 = math.MaxInt64

	u2iOffset uint64 = 1 << 63
)

type (
	// `TDBpersistence` is a database-based `IPersistence` implementation.
	TDBpersistence struct {
		_    struct{}
		db   *sql.DB       // the database to use
		mtx  *sync.RWMutex // pointer to avoid copying warnings
		fts5 bool          // whether SQLite supports full-text search
	}
)

// --------------------------------------------------------------------------
// constructor function

// `NewDBpersistence()` creates a new instance of `TDBpersistence`.
//
// In case of errors initialising the database connection, the function
// returns a `nil` value.
//
// Parameters:
//   - `aName`: The name of the database file to use.
//
// Returns:
//   - `*TDBpersistence`: A persistence instance instance.
func NewDBpersistence(aName string) *TDBpersistence {
	fName := filepath.Join(poPostingBaseDirectory, aName)
	dbInstance, hasFTS, err := initDatabase(fName)
	if nil != err {
		return nil
	}

	return &TDBpersistence{
		db:   dbInstance,
		mtx:  new(sync.RWMutex),
		fts5: hasFTS,
	}
} // NewDBpersistence()

// --------------------------------------------------------------------------
// private helper functions:

// `dbInt2id()` converts a 64-bit integer to a uint64, handling
// potential overflow.
//
// Parameters:
//   - `aInt`: The 64-bit integer to convert.
//
// Returns:
//   - `uint64`: The input integer, potentially with an offset applied.
func dbInt2id(aInt int64) uint64 {
	if 0 > aInt {
		offset := u2iOffset
		return uint64(aInt + int64(offset))
	}

	return uint64(aInt)
} //dbInt2id ()

// `dbInt2time()` converts a signed 64-bit integer from SQLite to
// a Go `time.Time` value.
//
// Parameters:
//   - `aInt`: A 64-bit integer ID to be converted.
//
// Returns:
//   - `time.Time`: The UnixNano value of the provided time.Time.
func dbInt2time(aInt int64) time.Time {
	return time.Unix(0, aInt)
} //dbInt2time ()

// `id2dbInt()` converts a 64-bit integer to a database-compatible integer.
//
// This function is used to convert the unsigned 64-bit integer IDs used
// in the application to the 64-bit integer IDs used in the SQLite database.
// If the provided `aID` is larger than the maximum SQLite integer
// (9223372036854775807), the function subtracts the offset `u2iOffset`
// (1 << 63) to ensure that the resulting integer is within the valid
// range for SQLite integers.
//
// Parameters:
//   - `aID`: The unsigned 64-bit integer ID to be converted.
//
// Returns:
//   - `int64`: The converted integer.
func id2dbInt(aID uint64) int64 {
	if aID > uint64(maxSqliteInt) {
		return int64(aID - u2iOffset)
	}

	return int64(aID)
} // id2dbInt()

// `time2dbInt()` converts a `time.Time` value to a signed 64-bit integer
// value suitable for SQLite's INTEGER field.
//
// Parameters:
//   - `aTime`: A time.Time value to be converted to a 64-bit integer.
//
// Returns:
//   - `int64`: The converted integer.
func time2dbInt(aTime time.Time) int64 {
	return aTime.UnixNano()
} // time2dbInt()

// --------------------------------------------------------------------------

// `init()` ensures proper interface implementation.
func init() {
	var (
		_ IPersistence = TDBpersistence{}
		_ IPersistence = (*TDBpersistence)(nil)
	)
} // init()

// --------------------------------------------------------------------------

// The SQL statement to create the database table
const dbInitTable = `
	CREATE TABLE IF NOT EXISTS "postings" (
		"id" INTEGER PRIMARY KEY,
		"lastModified" INTEGER NOT NULL,
		"markdown" TEXT NOT NULL
	);
`

// `initDatabase()` initialises a new SQLite database connection and
// checks whether it supports full-text search (FTS5).
//
// The function takes a path to the SQLite database file as a parameter
// and returns a pointer to a new SQLite database connection, a boolean
// indicating whether the database supports FTS5, and an error if any
// occurs during the initialisation process.

// The database connection is opened using the provided path and the
// "sqlite3" driver. If the database creation SQL statement fails to
// execute, the function returns `nil`, `false`, and an `error`.

// After the database connection is established, the function checks
// whether the SQLite database supports FTS5. If it does, the function
// creates an FTS5 virtual table referencing the regular table and adds
// triggers to keep the FTS table in sync with the regular table.

// The function returns the SQLite database connection, a boolean
// indicating whether the database supports FTS5, and nil if no errors occur.
//
// Parameters:
//   - `aPathFile`: The path to the SQLite database file.
//   - `*sql.DB`: A pointer to a new SQLite database connection.
//
// Returns:
//   - `bool`: An indicator for whether the database supports FTS5.
//   - `error`: An error if any occurs during the initialisation process.
func initDatabase(aPathFile string) (*sql.DB, bool, error) {
	// `cache=shared` is essential to avoid running out of file
	// handles since each query seems to hold its own file handle.
	// `loc=auto` gets `time.Time` with current locale.
	dsn := `file:` + aPathFile + `?cache=shared&loc=auto`

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		// failed to open database
		return nil, false, se.Wrap(err, 3)
	}

	// Create the table
	if _, err = db.Exec(dbInitTable); err != nil {
		db.Close()
		return nil, false, se.Wrap(err, 2)
	}

	// Check and add FTS5 database
	hasFTS, err := initFTS5(db)
	if err != nil {
		// db.Close()
		return db, false, nil // err
	}

	return db, hasFTS, nil
} // initDatabase()

// For the full-text search to work, we need to use the following build tag:
//
//	go build -tags "sqlite_fts5"

// `initFTS5()` checks the SQLite database for the FTS5 full-text
// search engine.
//
// It returns a boolean value indicating whether the SQLite database
// supports FTS5, and a possible error.
//
// Parameters:
//   - `aDB`: The SQLite database connection.
//
// Returns:
//   - `bool`: `true` if the SQLite database supports FTS5, `false` otherwise.
//   - `error`: A possible error, or `nil` on success.
func initFTS5(aDB *sql.DB) (bool, error) {
	const check4FTS5 = `SELECT sqlite_compileoption_used('ENABLE_FTS5')`
	var fts string

	if err := aDB.QueryRow(check4FTS5).Scan(&fts); nil != err {
		return false, se.Wrap(err, 1)
	}
	if "1" != fts {
		return false, nil
	}

	const dbAddFTS5 = `
	-- FTS virtual table referencing the regular table
	CREATE VIRTUAL TABLE IF NOT EXISTS postings_FTS USING FTS5(
		markdown,
		content='postings',
		content_rowid='id'
	);

	-- Trigger to keep FTS table in sync
	CREATE TRIGGER postings_ai AFTER INSERT ON postings BEGIN
		INSERT INTO postings_FTS(rowid, markdown) VALUES (new.id, new.markdown);
	END;

	CREATE TRIGGER postings_ad AFTER DELETE ON postings BEGIN
		INSERT INTO postings_FTS(postings_FTS, rowid, markdown) VALUES('delete', old.id, old.markdown);
	END;

	CREATE TRIGGER postings_au AFTER UPDATE ON postings BEGIN
		INSERT INTO postings_FTS(postings_FTS, rowid, markdown) VALUES('delete', old.id, old.markdown);
		INSERT INTO postings_FTS(rowid, markdown) VALUES (new.id, new.markdown);
	END;
`
	if _, err := aDB.Exec(dbAddFTS5); nil != err {
		return false, se.Wrap(err, 1)
	}

	return true, nil
} // initFTS5()

// --------------------------------------------------------------------------
// TDBpersistence methods

const dbGetCount = `SELECT COUNT(*) FROM postings`

// `Count()` returns the number of postings currently available.
//
// NOTE: This method is very resource intensive as it has to count all the
// posts stored in the filesystem.
//
// Returns:
//   - `int`: The number of available postings, or `0` in case of errors.
func (dbp TDBpersistence) Count() int {
	var result int

	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<1)
	defer cancel()

	if err := dbp.db.QueryRowContext(ctx, dbGetCount).Scan(&result); err != nil {
		return 0 //, fmt.Errorf("error counting rows: %v", err)
	}

	return result
} // Count()

const dbCreateRow = `INSERT INTO postings(id, lastModifies, markdown) VALUES(?, ?, ?)`

// `Create()` creates a new posting in the filesystem.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
//   - `aPost`: The `TPosting` instance containing the article's data.
//
// Returns:
//   - `int`: The number of bytes stored.
//   - 'error`:` A possible error, or `nil` on success.
func (dbp TDBpersistence) Create(aPost *TPosting) (int, error) {
	if nil == aPost {
		return 0, se.Wrap(ErrEmptyPosting, 1)
	}
	dbp.mtx.Lock()
	defer dbp.mtx.Unlock()

	dbID := id2dbInt(aPost.id)
	dbLM := time2dbInt(aPost.lastModified)
	dbText := string(aPost.markdown)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<2)
	defer cancel()

	result, err := dbp.db.ExecContext(ctx, dbCreateRow, dbID, dbLM, dbText)
	if err != nil {
		return 0, se.Wrap(err, 3)
	}

	if _, err = result.LastInsertId(); err != nil {
		return 0, se.Wrap(err, 1)
	}

	return int(unsafe.Sizeof(aPost.id)) +
		int(unsafe.Sizeof(aPost.lastModified)) +
		aPost.Len(), nil
} // Create()

const dbDeleteRow = `DELETE FROM postings WHERE id = ?`

// `Delete()` removes the posting/article from the filesystem
// and returns a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to delete.
//
// Returns:
//   - 'error`: A possible I/O error, or `nil` on success.
//
// Side Effects:
//   - Invalidates the internal count cache.
func (dbp TDBpersistence) Delete(aID uint64) error {
	dbp.mtx.Lock()
	defer dbp.mtx.Unlock()

	dbID := id2dbInt(aID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<2)
	defer cancel()

	res, err := dbp.db.ExecContext(ctx, dbDeleteRow, dbID)
	if err != nil {
		return se.Wrap(err, 2)
	}

	// Check the number of rows affected
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return se.Wrap(err, 2)
	}

	if 0 == rowsAffected {
		return fmt.Errorf("no rows deleted")
	}

	return nil
} // Delete()

const dbExistRow = `SELECT EXISTS(SELECT 1 FROM postings WHERE id = ?)`

// `Exists()` checks if a file with the given ID exists in the filesystem.
//
// It returns a boolean value indicating whether the file exists.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to check.
//
// Returns:
//   - `bool`: `true` if the file exists, `false` otherwise.
func (dbp TDBpersistence) Exists(aID uint64) bool {
	dbp.mtx.RLock()
	defer dbp.mtx.RUnlock()

	var result bool
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<1)
	defer cancel()

	if err := dbp.db.QueryRowContext(ctx, dbExistRow, aID).Scan(&result); err != nil {
		return false
	}

	return result
} // Exists()

// `PathFileName()` returns the posting's complete path-/filename.
//
// The returned path-/filename is in the format:
//
//	<base_directory>/<posting_id>.md
//
// Parameters:
//   - `aID`: The unique identifier of the posting to handle.
//
// Returns:
//   - `*string`: The path-/filename associated with `aID`.
func (dbp TDBpersistence) PathFileName(aID uint64) string {
	return poPostingBaseDirectory
} // PathFileName()

const dbReadRow = `SELECT id, lastModified, markdown FROM postings WHERE id = ?`

// `Read()` reads the posting from disk, returning a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to be read.
//
// Returns:
//   - `*TPosting`: The `TPosting` instance containing the article's data, or `nil` if the record doesn't exist.
//   - 'error`: A possible I/O error, or `nil` on success.
func (dbp TDBpersistence) Read(aID uint64) (*TPosting, error) {
	dbp.mtx.RLock()
	defer dbp.mtx.RUnlock()

	var (
		dbID, dbLM int64
		dbText     string
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<1)
	defer cancel()

	err := dbp.db.QueryRowContext(ctx, dbReadRow, id2dbInt(aID)).
		Scan(&dbID, &dbLM, &dbText)
	if err != nil {
		return nil, se.Wrap(err, 3)
	}

	post := &TPosting{
		id:           dbInt2id(dbID),
		lastModified: dbInt2time(dbLM),
		markdown:     []byte(dbText),
	}
	return post, nil
} // Read()

const dbRenameRow = `UPDATE postings SET id = ? WHERE id = ?"`

// `Rename()` renames a posting from its old ID to a new ID.
//
// Parameters:
//   - aOldID: The unique identifier of the posting to be renamed.
//   - aNewID: The new unique identifier for the new posting.
//
// Returns:
//   - `error`: An error if the operation fails, or `nil` on success.
func (dbp TDBpersistence) Rename(aOldID, aNewID uint64) error {
	dbp.mtx.Lock()
	defer dbp.mtx.Unlock()

	dbOldID, dbNewID := id2dbInt(aOldID), id2dbInt(aNewID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<1)
	defer cancel()

	result, err := dbp.db.ExecContext(ctx, dbRenameRow, dbNewID, dbOldID)
	if err != nil {
		return se.Wrap(err, 2)
	}

	// Get the number of affected rows
	if _, err = result.RowsAffected(); err != nil {
		return se.Wrap(err, 1)
	}

	return nil
} // Rename()

const (
	dbSearchLIKE = `SELECT id, lastModified, markup FROM postings WHERE markup LIKE ? LIMIT ? OFFSET ? ORDER BY id DESC`

	dbSearchMATCH = `SELECT id, lastModified, markup FROM postings WHERE markup MATCH ? LIMIT ? OFFSET ? ORDER BY id DESC`
)

// `Search()` retrieves a list of postings based on a search term.
//
// The method uses SQLite's FTS5 (Full-Text Search) feature to perform
// the search. If the underlying database does not support FTS5, the
// method falls back to a LIKE-based search.
//
// A zero value of `aLimit` means: no limit alt all.
//
// The returned `TPostList` type is a slice of `TPosting` instances, where
// `TPosting` is a struct representing a single posting. If the returned
// slice is an empty list then no matching postings were found; if it is
// `nil` it means there was an error retrieving the matches.
//
// Parameters:
//   - `aText`: The search query string.
//   - `aOffset`: An offset in the database result set of the search results.
//   - `aLimit`: The maximum number of search results to return.
//
// Returns:
//   - `*TPostList`: The list of search results, or `nil` in case of errors.
//   - `error`: If the search operation fails, or `nil` on success.
func (dbp TDBpersistence) Search(aText string, aOffset, aLimit uint) (*TPostList, error) {
	dbp.mtx.RLock()
	defer dbp.mtx.RUnlock()

	if 0 == aLimit {
		aLimit = 1 << 15 // 64K
	}

	var (
		err    error
		rows   *sql.Rows
		search string
	)
	if dbp.fts5 {
		search = dbSearchMATCH
	} else {
		aText = fmt.Sprintf("%%%s%%", aText)
		search = dbSearchLIKE
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<2)
	defer cancel()

	if rows, err = dbp.db.QueryContext(ctx, search, aText, aLimit, aOffset); err != nil {
		return nil, se.Wrap(err, 1)
	}
	defer rows.Close()

	postlist := NewPostList()
	for rows.Next() {
		var (
			dbID, dbLM int64
			dbText     string
		)
		if err = rows.Scan(&dbID, &dbLM, &dbText); err != nil {
			return nil, se.Wrap(err, 1)
		}
		post := &TPosting{
			id:           dbInt2id(dbID),
			lastModified: dbInt2time(dbLM),
			markdown:     []byte(dbText),
		}
		postlist.insert(post)
	}

	if err = rows.Err(); err != nil {
		return nil, se.Wrap(err, 1)
	}

	return postlist, nil
} // Search()

const dbUpdateRow = `UPDATE postings SET lastModified = ?, markdown = ? WHERE id = ?"`

// `Update()` updates the article's Markdown in the database.
//
// It returns the number of bytes stored and a possible I/O error.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
//   - `aPost`: A `TPosting` instance containing the article's data.
//
// Returns:
//   - `int`: The number of bytes written to the file.
//   - 'error`:` A possible I/O error, or `nil` on success.
//
// Side Effects:
//   - Invalidates the internal count cache.
func (dbp TDBpersistence) Update(aPost *TPosting) (int, error) {
	dbp.mtx.Lock()
	defer dbp.mtx.Unlock()

	dbID := id2dbInt(aPost.id)
	dbLM := time2dbInt(aPost.lastModified)
	dbText := string(aPost.markdown)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<2)
	defer cancel()

	result, err := dbp.db.ExecContext(ctx, dbUpdateRow, dbLM, dbText, dbID)
	if err != nil {
		return 0, se.Wrap(err, 2)
	}

	// Get the number of affected rows
	if _, err = result.RowsAffected(); err != nil {
		return 0, se.Wrap(err, 1)
	}

	return int(unsafe.Sizeof(aPost.id)) +
		int(unsafe.Sizeof(aPost.lastModified)) +
		aPost.Len(), nil
} // Update()

const dbWalkRows = `SELECT id FROM postings ORDER BY id DESC;`

// `Walk()` visits all existing postings, calling `aWalkFunc`
// for each posting.
//
// Parameters:
//   - `aWalkFunc`: The function to call for each posting.
//
// Returns:
//   - `error`: a possible error occurring the traversal process.
func (dbp TDBpersistence) Walk(aWalkFunc TWalkFunc) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second<<3)
	defer cancel()

	rows, err := dbp.db.QueryContext(ctx, dbWalkRows)
	if err != nil {
		return se.Wrap(err, 2)
	}
	defer rows.Close()

dirLoop: // Iterate over rows
	for rows.Next() {

		// Call the callback function with the row data
		var dbID int64

		if err := rows.Scan(&dbID); err != nil {
			continue
		}
		id := dbInt2id(dbID)

		if err := aWalkFunc(id); nil != err {
			if errors.Is(err, ErrSkipAll) {
				break dirLoop
			}
			return se.Wrap(err, 4)
		}
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return se.Wrap(err, 1)
	}

	return nil
} // Walk()

/* _EoF_ */
