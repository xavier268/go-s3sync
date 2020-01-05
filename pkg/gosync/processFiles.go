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
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				c.DownloadObject(sf)
				fmt.Printf("\tDOWNLOADED OBJECT %s\n", c.mode.String())
			}
		case ModeRestoreMock:
			if err != nil { // S3 object not found ?
				fmt.Printf("\tDELETE FILE NEEDED %s\n", c.mode.String())
			}
			if out.LastModified.UTC().Before(sf.updated) || *out.ContentLength != sf.size {
				fmt.Printf("\tDOWNLOAD OBJECT NEEDED %s\n", c.mode.String())
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

// DownloadObject downloads a potentially large object from S3 to file,
// overwriting existing file.
func (c *Config) DownloadObject(sf SrcFile) {

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

// ProcessFile will check the file againts the S3 bucket,
// returning the Action needed.
func (c *Config) ProcessFile(sf SrcFile) Action {

	if sf.absPath == "" {
		return ActionNone
	}
	out, err := c.s3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(c.GetKey(sf)),
	})
	if err != nil {
		fmt.Println("Object was not found")
		return ActionUploadFile
	}
	if *out.ContentLength != sf.size {
		fmt.Println("Size of remote object does not match")
		return ActionUploadFile
	}
	if out.LastModified.Before(sf.updated) {
		fmt.Println("File is more recent than object")
		return ActionUploadFile
	}
	return ActionNone
}
