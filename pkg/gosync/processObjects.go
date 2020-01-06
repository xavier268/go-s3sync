package gosync

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ProcessObjects performs a check on all s3 objects,
// checking what S3 or files changes are needed.
func (c *Config) ProcessObjects() {

	// Set a new waitGroup
	wait := new(sync.WaitGroup)

	fmt.Println("CheckObjects started")

	// Start a couple of workers to process them
	// Each worker calls Done() when channel is closed
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go c.ObjectWorker(i, wait)
	}

	// Start pushing in channel
	// Will close channel and call c.filesWait.Done() upon completion
	wait.Add(1)
	go c.WalkObjects(wait)

	// Wait until all walkers and workers are finished.
	wait.Wait()

	fmt.Println("\nCheckObjects finished")

}

// WalkObjects will push the s3 objects in a channel for further processing.
// It closes the object channel and call c.wait.Done() when finished.
func (c *Config) WalkObjects(wait *sync.WaitGroup) {

	defer wait.Done()
	defer close(c.objects)

	li := new(s3.ListObjectsV2Input).SetBucket(c.bucket)
	err := c.s3.ListObjectsV2Pages(li, func(res *s3.ListObjectsV2Output, lastpage bool) bool {

		for _, o := range res.Contents {
			c.objects <- c.DstObjectFromS3Object(o)
		}
		return !lastpage
	})

	if err != nil {
		panic(err)
	}
	fmt.Println("Finished walking objects")

}

// ObjectWorker processes the objects.
// There are typically  multiple  workers running in parallel.
// It calls c.wait.Done() at the end.
func (c *Config) ObjectWorker(i int, wait *sync.WaitGroup) {

	defer wait.Done()

	fmt.Printf("Object worker %d started ....\n", i)
	for ob := range c.objects {
		// look for corresponding file info
		fi, err := os.Stat(ob.GetAbsPath(c))

		switch c.mode {

		case ModeBackup:
			if err != nil || fi.IsDir() {
				// no file, delete the corresponding s3 object
				c.DeleteObject(ob)
				fmt.Printf("\tDELETED\t%s\n", c.mode.String())
				break
			}
			if fi.ModTime().UTC().After(ob.updated) || fi.Size() != ob.size {
				// refresh needed
				c.UploadObject(ob)
				fmt.Printf("\tUPLOADED\t%s\n", c.mode.String())
			}

		case ModeBackupMock:
			if err != nil || fi.IsDir() {
				// no file, delete the corresponding s3 object
				fmt.Printf("\tDELETE NEEDED\t%s\n", c.mode.String())
				break
			}
			if fi.ModTime().UTC().After(ob.updated) || fi.Size() != ob.size {
				// refresh needed
				fmt.Printf("\tUPLOAD NEEDED\t%s\n", c.mode.String())
			}
		case ModeRestore:
			if err != nil ||
				fi.IsDir() ||
				fi.Size() != ob.size ||
				fi.ModTime().UTC().After(ob.updated) {
				// need to download from s3
				c.DownloadObject(ob)
				fmt.Printf("\tDOWNLOADED\t%s\n", c.mode.String())

			}

		case ModeRestoreMock:
			if err != nil ||
				fi.IsDir() ||
				fi.Size() != ob.size ||
				fi.ModTime().UTC().After(ob.updated) {
				// need to download from s3
				fmt.Printf("\tDOWNLOAD NEEDED\t%s\n", c.mode.String())
			}

		default:
			panic("invalid mode specified in configuration")
		}

	}
	fmt.Printf("Object worker %d stopped ....\n", i)
}

// // CheckObject verify if the s3 object is in sync with the local file system.
// // It returns the action needed.
// func (c *Config) CheckObject(o DstObject) Action {
// 	ap := o.GetAbsPath(c)
// 	fi, err := os.Stat(ap)
// 	if err != nil || fi.IsDir() {
// 		fmt.Println("No corresponding file")
// 		return ActionDeleteObject
// 	}
// 	if fi.Size() != o.size {
// 		fmt.Println("File and object size do not match")
// 		return ActionUploadFile
// 	}
// 	if fi.ModTime().After(o.updated) {
// 		fmt.Println("File was modified after being saved")
// 		return ActionUploadFile
// 	}
// 	return ActionNone
// }

// DeleteObject delete the provided object from s3
func (c *Config) DeleteObject(ob DstObject) {

	_, err := c.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(ob.key),
	})

	if err != nil {
		panic(err)
	}

}

// UploadObject refresh the s3 object from corresponding file
func (c *Config) UploadObject(ob DstObject) {
	c.UploadFile(
		SrcFile{
			absPath: ob.GetAbsPath(c),
		})
}

// DownloadObject downloads an S3 object to the local file system.
func (c *Config) DownloadObject(ob DstObject) {
	c.DownloadFile(
		SrcFile{
			absPath: ob.GetAbsPath(c),
		})

}
