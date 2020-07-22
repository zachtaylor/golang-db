package patch

import (
	"io/ioutil"
	"strconv"
)

// GetFiles scans the dir for files named as "\d\d\d\d.*\.sql", returns map(patchid->filename)
func GetFiles(dir string) map[int]string {
	patches := make(map[int]string)
	if files, err := ioutil.ReadDir(dir); err == nil {
		for _, f := range files {
			if name := f.Name(); len(name) < 8 {
				// file name too short
			} else if ext := name[len(name)-4:]; ext != ".sql" {
				// file name does not end with ".sql"
			} else if id, err := strconv.ParseInt(name[:4], 10, 0); err != nil {
				// file name does not start with 4 numbers
			} else {
				patches[int(id)] = dir + name
			}
		}
	}
	return patches
}
