package gosync

import (
	"io"
	"os"
	"path/filepath"
)

// RemoveAllEmptyDirs remove all empty directories from the source file.
// Mostly useful after a restore, if file folder was not initially empty.
// Never called implicitely on backup/restore.
// Special care is taken to ensure nested empty dirs are also removed.
func (c *Config) RemoveAllEmptyDirs() {

	var touched bool = true

	// Iterated as long as we are touching something ...
	for touched {

		touched = false

		e2 := filepath.Walk(c.prefix,
			func(path string, info os.FileInfo, e1 error) error {
				if e1 != nil {
					panic(e1)
				}
				if path == c.prefix {
					// don't remove root prefix dir !
					return nil
				}
				if info.IsDir() && c.isEmptyDir(path) {
					e3 := os.Remove(path)
					if e3 != nil {
						panic(e3)
					}
					touched = true
				}
				return nil
			})

		if e2 != nil {
			panic(e2)
		}

	}
}

// isEmptyDir test for an empty dir
// Assuming we already know it is a dir...
func (c *Config) isEmptyDir(dirname string) bool {

	f, err := os.Open(dirname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Minimal load on the system ...
	if err == io.EOF {
		return true
	}
	return false
}
