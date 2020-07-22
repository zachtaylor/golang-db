package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // load mysql driver
	"taylz.io/db"
)

// Open creates a db connection using mysql
func Open(dataSourceName string) (*db.DB, error) {
	if db, err := sql.Open("mysql", dataSourceName); err != nil {
		return nil, err
	} else if err = db.Ping(); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

// CreatePatchTable creates the patch table
func CreatePatchTable(x *db.DB) error {
	return db.ExecTx(x, `CREATE TABLE patch(patch INTEGER) ENGINE=InnoDB; INSERT INTO patch (patch) VALUES (0);`)
}
