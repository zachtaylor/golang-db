# Package `db`

Package `db` provides database connection helpers based on `database/sql`

## Environment

Uses the following `taylz.io/env` keys

```
USER
PASSWORD
HOST
PORT
NAME
```

## Package `db/mysql`

Package `mysql` loads mysql driver using `"github.com/go-sql-driver/mysql"`

## `db-patch`

`go get taylz.io/db/cmd/db-patch`

Connect to a database using MySQL, and execute a series of patches

Patches are contained separately in files, known as patch files. These files
- contain SQL statements, which are executed as transactions (each patch will succeed or fail as a whole)
- begin with 4 numbers, identifying the patch number in sequence
- end with ".sql"
