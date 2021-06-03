package db // import "ztaylor.me/db"

import (
	"database/sql"
	"errors"
	"strings"
)

// ErrPatchTable is returned by Patch when the patch table doesn't exist
var ErrPatchTable = errors.New(`table "patch" does not exist`)

// ErrSQLPanic is returned by ExecTx when it encounters a panic
var ErrSQLPanic = errors.New(`sql panic`)

// ErrTxEmpty is returned by ExecTx when tx has no statements
var ErrTxEmpty = errors.New(`transaction contains no statements`)

// DB = sql.DB
type DB = sql.DB

// Result = sql.Result
type Result = sql.Result

// Scanner provides a header for generic SQL data set
type Scanner interface {
	Scan(...interface{}) error
}

// DSN returns a formatted DSN string
func DSN(user, password, host, port, name string) string {
	return user + `:` + password + `@tcp(` + host + `:` + port + `)/` + name
}

// ExecTx uses database transaction to apply SQL statements
func ExecTx(db *DB, sql string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	isEmpty := true
	defer func() {
		if p := recover(); p != nil {
			err = ErrSQLPanic
		}
		if isEmpty {
			err = ErrTxEmpty
		}
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	for _, stmt := range strings.Split(sql, `;`) {
		if stmt = strings.Trim(stmt, "\n\r \t"); stmt != "" {
			isEmpty = false
			if _, err = tx.Exec(stmt); err != nil {
				break
			}
		}
	}
	return err
}
