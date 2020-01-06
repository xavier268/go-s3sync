package gosync

import (
	"fmt"
	"io"
	"os"
	"path"
)

// RemoveAllEmptyDirs remove all empty directtories from the source file.
// Mostly useful after a restore, if file folder was not initially empty.
// Never called implicitely on backup/restore.
// Special care is taken to ensure nested empty dirs are also removed.
func (c *Config) RemoveAllEmptyDirs() {
	c.clean(c.prefix)
}

// Recursively clen empty dirs from the root path.
// root must be a dir, expressed in absolute path.
// Cleaning will happen leaves first,
// to ensure nested empty dires are correctly cleaned.
func (c *Config) clean(root string) {

	fmt.Printf("Clean is looking at dir %s\n", root)

	info, err := os.Stat(root)
	if err != nil {
		fmt.Println(err)
		return
	}

	if !info.IsDir() {
		// Ingnore non dirs ..
		return
	}

	// First, canvass the non empty dir recursiveley.
	if info.IsDir() && !c.isEmptyDir(root) {

		f, err := os.Open(root)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		dirs, err := f.Readdirnames(-1)
		if err != nil {
			panic(err)
		}
		for _, d := range dirs {
			ad := path.Join(c.prefix, d)
			if err != nil {
				panic(err)
			}
			// recurse ...
			c.clean(ad)
		}

	}

	// Then, and only then, reevaluate if dir is empty.
	// Nested sub-dirs might have been already removed...
	// Whatever, but never remove the "prefix" top level dir !
	if c.isEmptyDir(root) && root != c.prefix {
		fmt.Printf("Removing empty dir : %s\n", root)
		err := os.Remove(root)
		if err != nil {
			panic(err)
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
