/*
Copyright © 2024 M.Watermann, 10247 Berlin, Germany

		All rights reserved
	EMail : <support@mwat.de>
*/
package nele

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"path/filepath"
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

		Count() uint32
		Exists(aID uint64) bool
		PathFileName(aID uint64) string
		Rename(aOldID, aNewID uint64) error
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
		_  struct{}
		db *sql.DB // the database to use
		// mtx *sync.RWMutex // pointer to avoid copying warnings
		fts5 bool // whether SQLite supports full-text search
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
		db: dbInstance,
		// mtx: new(sync.RWMutex),
		fts5: hasFTS,
	}
} // NewDBpersistence()

// --------------------------------------------------------------------------
// private helper functions:

// `dbInt2time()` returns a date/time represented by `aID`.
//
// Parameters:
// - `aID`: A posting's ID to be converted to a `time.Time`.
//
// Returns:
// - `time.Time`: The UnixNano value of the provided time.Time.
func dbInt2time(aInt int64) time.Time {
	return time.Unix(0, aInt)
} //dbInt2time ()

func time2dbInt(aTime time.Time) int64 {
	return aTime.UnixNano()
} // time2dbInt()

func id2dbInt(aID uint64) int64 {
	if aID > uint64(maxSqliteInt) {
		return int64(aID - u2iOffset)
	}

	return int64(aID)
} // id2dbInt()

func dbInt2id(aInt int64) uint64 {
	if 0 > aInt {
		offset := u2iOffset
		return uint64(aInt + int64(offset))
	}

	return uint64(aInt)
} //dbInt2id ()

// --------------------------------------------------------------------------

// `init()` ensures proper interface implementation.
func init() {
	var (
		_ IPersistence = TDBpersistence{}
		_ IPersistence = (*TDBpersistence)(nil)
	)
} // init()

// --------------------------------------------------------------------------

// The SQL statements to create the database table
const dbCreationSQL = `
	CREATE TABLE IF NOT EXISTS postings (
		id INTEGER PRIMARY KEY,
		lastModified INTEGER NOT NULL,
		markdown TEXT NOT NULL
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
	if _, err = db.Exec(dbCreationSQL); err != nil {
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

	if "1" == fts {
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
	}

	return false, nil
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

	if err := dbp.db.QueryRow(dbGetCount).Scan(&result); err != nil {
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
	dbID := id2dbInt(aPost.id)
	dbLM := time2dbInt(aPost.lastModified)
	dbText := string(aPost.markdown)

	result, err := dbp.db.Exec(dbCreateRow, dbID, dbLM, dbText)
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

const dbDelRecord = `DELETE FROM postings WHERE id = ?`

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
	dbID := id2dbInt(aID)
	res, err := dbp.db.Exec(dbDelRecord, dbID)
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

const dbExistRecord = `SELECT EXISTS(SELECT 1 FROM postings WHERE id = ?)`

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
	var result bool

	err := dbp.db.QueryRow(dbExistRecord, aID).Scan(&result)
	if err != nil {
		return false //, err
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

const dbReadRecord = `SELECT id, lastModified, markdown FROM postings WHERE id = ?`

// `Read()` reads the posting from disk, returning a possible I/O error.
//
// Parameters:
//   - `aID`: The unique identifier of the posting to be read.
//
// Returns:
//   - `*TPosting`: The `TPosting` instance containing the article's data, or `nil` if the record doesn't exist.
//   - 'error`: A possible I/O error, or `nil` on success.
func (dbp TDBpersistence) Read(aID uint64) (*TPosting, error) {
	var (
		dbID, dbLM int64
		dbText     string
	)

	err := dbp.db.QueryRow(dbReadRecord, id2dbInt(aID)).
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

const dbRenRecord = `UPDATE postings SET id = ? WHERE id = ?"`

// `Rename()` renames a posting from its old ID to a new ID.
//
// Parameters:
//   - aOldID: The unique identifier of the posting to be renamed.
//   - aNewID: The new unique identifier for the new posting.
//
// Returns:
//   - `error`: An error if the operation fails, or `nil` on success.
func (dbp TDBpersistence) Rename(aOldID, aNewID uint64) error {
	dbOldID, dbNewID := id2dbInt(aOldID), id2dbInt(aNewID)

	result, err := dbp.db.Exec(dbRenRecord, dbNewID, dbOldID)
	if err != nil {
		return se.Wrap(err, 2)
	}

	// Get the number of affected rows
	if _, err = result.RowsAffected(); err != nil {
		return se.Wrap(err, 1)
	}

	return nil
} // Rename()

const dbSearchText = `SELECT id, lastModified, markup FROM postings WHERE postings MATCH ? ORDER BY id DESC`

func (dbp TDBpersistence) Search(aText string) (*TPostList, error) {
	rows, err := dbp.db.Query(dbSearchText, aText)
	if err != nil {
		return nil, se.Wrap(err, 2)
	}
	defer rows.Close()

	postlist := NewPostList()
	for rows.Next() {
		var (
			dbID, dbLM int64
			dbText     string
		)
		if err := rows.Scan(&dbID, &dbLM, &dbText); err != nil {
			return nil, se.Wrap(err, 1)
		}
		post := &TPosting{
			id:           dbInt2id(dbID),
			lastModified: dbInt2time(dbLM),
			markdown:     []byte(dbText),
		}
		postlist.insert(post)
	}

	if err := rows.Err(); err != nil {
		return nil, se.Wrap(err, 1)
	}

	return postlist, nil
} // Search()

const dbUpdRecord = `UPDATE postings SET lastModified = ?, markdown = ? WHERE id = ?"`

// `Update()` updates the article's Markdown on disk.
//
// It returns the number of bytes written to the file and a possible I/O error.
//
// If the provided `aPost` is `nil`, an `ErrEmptyPosting` error
// is returned.
//
// Parameters:
// - `aPost`: A `TPosting` instance containing the article's data.
//
// Returns:
// - `int`: The number of bytes written to the file.
// - 'error`:` A possible I/O error, or `nil` on success.
//
// Side Effects:
// - Invalidates the internal count cache.
func (dbp TDBpersistence) Update(aPost *TPosting) (int, error) {
	dbID := id2dbInt(aPost.id)
	dbLM := time2dbInt(aPost.lastModified)
	dbText := string(aPost.markdown)

	result, err := dbp.db.Exec(dbUpdRecord, dbLM, dbText, dbID)
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

const dbWalkRecords = `SELECT id FROM postings ORDER BY id DESC;`

// `Walk()` visits all existing postings, calling `aWalkFunc`
// for each posting.
//
// Parameters:
//   - `aWalkFunc`: The function to call for each posting.
//
// Returns:
//   - `error`: a possible error occurring the traversal process.
func (dbp TDBpersistence) Walk(aWalkFunc TWalkFunc) error {
	rows, err := dbp.db.Query(dbWalkRecords)
	if err != nil {
		return se.Wrap(err, 2)
	}
	defer rows.Close()

dirLoop:
	// Iterate over rows
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
			return se.Wrap(err, 7)
		}
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return se.Wrap(err, 1)
	}

	return nil
} // Walk()

/*
func iterateTable(db *sql.DB, callback func(map[string]interface{}) error) error {
	// Query to get all rows from the table

	// Execute the query
	rows, err := db.Query(dbWalkRecords)
	if err != nil {
		return se.Wrap(err, 2)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// Prepare a slice to hold the column values
	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	// Iterate over rows
	for rows.Next() {
		// Scan the row into the values slice
		err := rows.Scan(values...)
		if err != nil {
			return err
		}

		// Create a map to hold the row data
		rowData := make(map[string]interface{})
		for i, col := range columns {
			rowData[col] = *(values[i].(*interface{}))
		}

		// Call the callback function with the row data
		err = callback(rowData)
		if err != nil {
			return err
		}
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}
*/

/*
func main() {
	db, err := sql.Open("sqlite3", "path/to/your/database.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	results, err := SearchDocuments(db, "your search query")
	if err != nil {
		log.Fatalf("Search failed: %v", err)
	}

	for _, doc := range results {
		fmt.Printf("Title: %s\nBody: %s\n\n", doc.Title, doc.Body)
	}
}
*/

/* _EoF_ */
