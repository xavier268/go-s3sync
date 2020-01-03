package synchro

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Config is the sync configuration object
// It is safe for parrallel use.
type Config struct {
	targetBucket,
	region,
	fileRoot string
	maxLen int

	s3 *s3.S3

	files       chan SrcFile
	filesWait   sync.WaitGroup
	filesConcur int
}

// NewTestConfig provides a test configuration
func NewTestConfig() *Config {
	c := new(Config)
	c.targetBucket = "bup3.photos.gandillot.com"
	c.region = "eu-west-1"
	c.fileRoot = "/home/xavier/Desktop/go-s3sync"
	c.maxLen = 1000 // Name limit - real is 1024

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.region),
	})
	if err != nil {
		panic(err)
	}
	c.s3 = s3.New(sess)

	c.files = make(chan SrcFile, 20)
	c.filesConcur = 10 // 10 files processed in parallel
	return c
}

// WalkFiles will walk and send files through the files channel.
func (c *Config) WalkFiles() {

	err := filepath.Walk(c.fileRoot,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}
			if info.IsDir() {
				// Just ignore dirs
				return nil
			}
			i := *new(SrcFile)
			i.abspath, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			i.updated = info.ModTime().UTC()

			if len(i.abspath) >= c.maxLen {
				return errors.New("file name exceeds allowed length : " + i.abspath)
			}

			c.files <- i
			return nil

		})

	// Close channel,
	// all files have been sent.
	close(c.files)

	if err != nil {
		panic(err)
	}
}

// FileWorker processes files from the channel.
// There would typically be multiple file workers.
func (c *Config) FileWorker(i int) {
	fmt.Printf("Fileworker %d started ..........\n", i)

	for sf := range c.files {
		fmt.Printf("%d)\t%s\n", i, sf.String())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Fileworker %d finished ..........\n", i)
	c.filesWait.Done()
}
