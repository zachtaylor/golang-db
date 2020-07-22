package patch

import "taylz.io/db"

// Get returns the current patch number for the database
//
// returns -1, ErrPatchTable if the table doesn't exist
func Get(conn *db.DB) (int, error) {
	if patch, err := scanPatch(conn); err == nil {
		return patch, nil
	} else if e := err.Error(); len(e) > 10 && e[:10] == "Error 1146" {
		return -1, db.ErrPatchTable
	} else {
		return -1, err
	}
}
func scanPatch(conn *db.DB) (patch int, err error) {
	err = conn.QueryRow("SELECT * FROM patch").Scan(&patch)
	return
}
