package synchro

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// newS3 returns a new, configured client
func newS3(c *Config) *s3.S3 {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(c.region),
	})
	if err != nil {
		fmt.Printf("Configuration : %v\n", c)
		panic(err)
	}
	return s3.New(sess)
}

// Check performs a sync check, checking what S3 changes will be needed.
// Used for debugging.
func Check(c *Config) {

	// Start pushing files into the files channel
	// Will close channel when all files are sent
	go c.WalkFiles()

	// Start a couple of workers to process them
	// Each worker call c.fileWait.Done() upon completion.
	for i := 0; i < c.filesConcur; i++ {
		c.filesWait.Add(1)
		go c.FileWorker(i)
	}

	// Wait until all workers are finished.
	c.filesWait.Wait()
	fmt.Println("File check finished")

}
