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
	filesConcur int

	objects       chan s3.Object
	objectsConcur int

	wait sync.WaitGroup
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

	c.files = make(chan SrcFile, 2000)
	c.filesConcur = 50 // nbr of files processed in parallel

	c.objects = make(chan s3.Object, 2000)
	c.objectsConcur = 50

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

	fmt.Println("FileWalker finished walking the files")
	c.wait.Done()
}

// FileWorker processes files from the channel.
// There are typically  multiple workers running in parallel.
func (c *Config) FileWorker(i int) {
	fmt.Printf("Fileworker %d started ..........\n", i)

	for sf := range c.files {
		fmt.Printf("%d)\t%s\n", i, sf.String())
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Printf("Fileworker %d finished ..........\n", i)
	c.wait.Done()
}

// WalkObjects will push the s3 objects in a channel for further processing.
// It closes the object channel and call c.wait.Done() when finished.
func (c *Config) WalkObjects() {

	li := new(s3.ListObjectsV2Input).SetBucket(c.targetBucket)
	c.s3.ListObjectsV2Pages(li, func(res *s3.ListObjectsV2Output, lastpage bool) bool {

		for _, o := range res.Contents {
			c.objects <- *o
		}
		return !lastpage
	})

	close(c.objects)
	fmt.Println("Finished walking objects")
	c.wait.Done()
}

// ObjectWorker processes the objects.
// There are typically  multiple  workers running in parallel.
func (c *Config) ObjectWorker(i int) {
	fmt.Printf("Object worker %d started ....\n", i)
	for ob := range c.objects {
		fmt.Printf("%d)\tObject : %s\n", i, *ob.Key)
	}
	fmt.Printf("Object worker %d stopped ....\n", i)
	c.wait.Done()
}
