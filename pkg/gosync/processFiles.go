package gosync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ProcessFiles performs a check on all files,
// checking what files or S3 objects should be changed.
// If we re not in the xxxmock mode, changes will be made asynchroneously.
func (c *Config) ProcessFiles() {

	// Set a new waitGroup
	wait := new(sync.WaitGroup)

	fmt.Println("\nCheckFiles started")

	// Start a couple of workers to process them
	// Each worker calls Done() when channel is closed.
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go c.fileWorker(i, wait)
	}

	// Start pushing in channel
	// Will close channel and call c.filesWait.Done() upon completion
	wait.Add(1)
	go c.walkFiles(wait)

	// Wait until all walkers and workers are finished.
	wait.Wait()

	fmt.Println("\nCheckFiles finished")

}

// walkFiles will walk and send files through the files channel.
// Directories are ignored, only the files inside are processed.
// It will closes channel and calls c.wait.Done() at the end.
func (c *Config) walkFiles(wait *sync.WaitGroup) {

	defer wait.Done()
	defer close(c.files)

	err := filepath.Walk(c.prefix,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}
			if info.IsDir() {
				// Just ignore dirs, do nothing
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

// fileWorker processes files from the channel.
// There are typically  multiple workers running in parallel.
// It calls c.wait.Done() at the end.
func (c *Config) fileWorker(i int, wait *sync.WaitGroup) {
	fmt.Printf("File worker %d started ..........\n", i)

	for sf := range c.files {

		out, err := c.s3.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(c.getKey(sf)),
		})

		switch c.mode {
		case ModeBackup:
			if err != nil ||
				*out.ContentLength != sf.size ||
				out.LastModified.UTC().Before(sf.updated) {
				c.uploadFile(sf)
				fmt.Printf("UPLOADED %s\t%s\n", c.mode.String(), sf.String())
			}
		case ModeBackupMock:
			if err != nil ||
				*out.ContentLength != sf.size ||
				out.LastModified.UTC().Before(sf.updated) {
				fmt.Printf("UPLOADED %s\t%s\n", c.mode.String(), sf.String())
			}

		case ModeRestore:
			if err != nil { // S3 object not found ?
				c.deleteFile(sf)
				fmt.Printf("\tDELETED FILE %s\t%s\n", c.mode.String(), sf.String())
				break
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				c.downloadFile(sf)
				fmt.Printf("\tDOWNLOADED %s\t%s\n", c.mode.String(), sf.String())
			}
		case ModeRestoreMock:
			if err != nil { // S3 object not found ?
				fmt.Printf("\tDELETED FILE %s\t%s\n", c.mode.String(), sf.String())
				break
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				fmt.Printf("\tDOWNLOADED %s\t%s\n", c.mode.String(), sf.String())
			}

		default:
			fmt.Println("Mode code : ", c.mode)
			panic("Invalid mode in configuration ?! : ")
		}

	}
	fmt.Printf("File worker %d finished ..........\n", i)
	wait.Done()
}
