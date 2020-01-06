package gosync

import (
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestRemoveAllContent1(t *testing.T) {

	c := tConfig{NewDefaultConfig()}
	c.createDummyFoldersAndFiles()
	c.removeAllFileContent()
	//We should be empty ...
	if !c.isEmptyDir(c.prefix) {
		t.Fatal("The test directory content is not empty.")
	}
	// but root folder should still be there ...
	_, err := os.Stat(c.prefix)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveEmptyDirs(t *testing.T) {

	c := tConfig{NewDefaultConfig()}
	c.removeAllFileContent()

	// Now, clean start ...

	c.createDummyFolders()
	c.RemoveAllEmptyDirs()
	//We should be empty ...
	if !c.isEmptyDir(c.prefix) {
		t.Fatal("The test directory content is not empty.")
	}
	// but root folder should still be there ...
	_, err := os.Stat(c.prefix)
	if err != nil {
		t.Fatal(err)
	}
}

// ************* utilities ******************

// testable Configuration.
type tConfig struct {
	*Config
}

// recursively remove all content
// inside the root dir (prefix).
// Used for testing, and resetting to a known initial state
// the root directory.
func (c tConfig) removeAllFileContent() {
	d, err := os.Open(c.prefix)
	if err != nil {
		panic(err)
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(c.prefix, name))
		if err != nil {
			panic(err)
		}
	}
}

// create a dummy test folder structure
func (c tConfig) createDummyFolders() {

	os.MkdirAll(path.Join(c.prefix, "fone", "ftwo", "ftree1"), c.dirPerm)
	os.MkdirAll(path.Join(c.prefix, "fone", "ftwo", "ftree2"), c.dirPerm)
	os.MkdirAll(path.Join(c.prefix, "fone", "ftwo2", "ftree4"), c.dirPerm)
	os.MkdirAll(path.Join(c.prefix, "fone", "ftwo3"), c.dirPerm)
	os.MkdirAll(path.Join(c.prefix, "fone", "ftwo4"), c.dirPerm)
	os.MkdirAll(path.Join(c.prefix, "fone2"), c.dirPerm)

}

// Create 3 empty files
func (c tConfig) createDummyFiles() {
	os.Create(path.Join(c.prefix, "one"))
	os.Create(path.Join(c.prefix, "two"))
	os.Create(path.Join(c.prefix, "three"))
}

// create dummy folders, some with files inside
func (c tConfig) createDummyFoldersAndFiles() {

	c.createDummyFolders()
	c.createDummyFiles()

	os.Create(path.Join(c.prefix, "fone/one"))
	os.Create(path.Join(c.prefix, "fone/ftwo/two"))

}
