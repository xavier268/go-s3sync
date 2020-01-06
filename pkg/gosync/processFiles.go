package gosync

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ProcessFiles performs a check on all files,
// checking what files or S3 objects should be changed.
func (c *Config) ProcessFiles() {

	// Set a new waitGroup
	wait := new(sync.WaitGroup)

	fmt.Println("CheckFiles started")

	// Start a couple of workers to process them
	// Each worker calls Done() when channel is closed.
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go c.FileWorker(i, wait)
	}

	// Start pushing in channel
	// Will close channel and call c.filesWait.Done() upon completion
	wait.Add(1)
	go c.WalkFiles(wait)

	// Wait until all walkers and workers are finished.
	wait.Wait()

	fmt.Println("\nCheckFiles finished")

}

// WalkFiles will walk and send files through the files channel.
// Directories are ignored, only the files inside are processed.
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

// FileWorker processes files from the channel.
// There are typically  multiple workers running in parallel.
// It calls c.wait.Done() at the end.
func (c *Config) FileWorker(i int, wait *sync.WaitGroup) {
	fmt.Printf("File worker %d started ..........\n", i)

	for sf := range c.files {
		fmt.Printf("File worker #%d\t%s\n", i, sf.String())
		if sf.absPath == "" {
			continue
		}
		out, err := c.s3.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(c.GetKey(sf)),
		})
		switch c.mode {
		case ModeBackup:
			if err != nil ||
				*out.ContentLength != sf.size ||
				out.LastModified.UTC().Before(sf.updated) {
				c.UploadFile(sf)
				fmt.Printf("\tUPLOADED %s\n", c.mode.String())
			}
		case ModeBackupMock:
			if err != nil ||
				*out.ContentLength != sf.size ||
				out.LastModified.UTC().Before(sf.updated) {
				fmt.Printf("\tUPLOAD %s\n", c.mode.String())
			}

		case ModeRestore:
			if err != nil { // S3 object not found ?
				c.DeleteFile(sf)
				fmt.Printf("\tDELETED FILE %s\n", c.mode.String())
				break
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				c.DownloadFile(sf)
				fmt.Printf("\tDOWNLOADED %s\n", c.mode.String())
			}
		case ModeRestoreMock:
			if err != nil { // S3 object not found ?
				fmt.Printf("\tDELETE FILE NEEDED %s\n", c.mode.String())
				break
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				fmt.Printf("\tDOWNLOAD NEEDED %s\n", c.mode.String())
			}

		default:
			fmt.Println("Mode code : ", c.mode)
			panic("Invalid mode in configuration ?! : ")

		}

	}
	fmt.Printf("File worker %d finished ..........\n", i)
	wait.Done()
}

// UploadFile upload a potentially large file to S3
func (c *Config) UploadFile(sf SrcFile) {

	file, err := os.Open(sf.absPath)
	if err != nil {
		fmt.Println(sf.String())
		panic(err)
	}
	defer file.Close()

	up := s3manager.NewUploader(c.sess)
	_, err = up.Upload(&s3manager.UploadInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.GetKey(sf)),
		Body:   file,
	})
	if err != nil {
		panic(err)
	}
}

// DeleteFile does just that ...
func (c *Config) DeleteFile(sf SrcFile) {
	err := os.Remove(sf.absPath)
	if err != nil {
		panic(err)
	}
}

// DownloadFile downloads a potentially large object from S3 to file,
// overwriting existing file.
func (c *Config) DownloadFile(sf SrcFile) {

	file, err := os.Create(sf.absPath)
	if err != nil {
		fmt.Println(sf.String())
		panic(err)
	}
	defer file.Close()

	down := s3manager.NewDownloader(c.sess)
	_, err = down.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(c.bucket),
			Key:    aws.String(c.GetKey(sf)),
		})

	if err != nil {
		fmt.Println(sf)
		panic(err)
	}
}
