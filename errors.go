package db

import "errors"

// ErrPatchTable is returned by Patch when the patch table doesn't exist
var ErrPatchTable = errors.New(`table "patch" does not exist`)

// ErrSQLPanic is returned by ExecTx when it encounters a panic
var ErrSQLPanic = errors.New(`sql panic`)

// ErrTxEmpty is returned by ExecTx when tx has no statements
var ErrTxEmpty = errors.New(`patch file contains no statements`)
