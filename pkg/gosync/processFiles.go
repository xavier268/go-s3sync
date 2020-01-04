package gosync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// WalkFiles will walk and send files through the files channel.
// It will closes channel and calls c.wait.Done() at the end.
func (c *Config) WalkFiles(wait *sync.WaitGroup) {

	defer wait.Done()
	defer close(c.files)

	err := filepath.Walk(c.prefix,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}
			if info.IsDir() {
				// Just ignore dirs
				return nil
			}
			i := *new(SrcFile)
			i.absPath, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			i.updated = info.ModTime().UTC()
			i.size = info.Size()

			if len(i.absPath) >= c.maxKeyLength {
				return errors.New("file name exceeds allowed length : " + i.absPath)
			}
			// trigger file processing
			c.files <- i
			return nil
		})

	if err != nil {
		panic(err)
	}

	fmt.Println("FileWalker finished walking the files")

}

// FileWorker processes files from the channel.
// There are typically  multiple workers running in parallel.
// It calls c.wait.Done() at the end.
func (c *Config) FileWorker(i int, wait *sync.WaitGroup) {
	fmt.Printf("Fileworker %d started ..........\n", i)

	for sf := range c.files {
		fmt.Printf("%d)\t%s\n", i, sf.String())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Fileworker %d finished ..........\n", i)
	wait.Done()
}
