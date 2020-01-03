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
	// Will close channel and call c.filesWait.Done() upon completion
	c.wait.Add(2)
	go c.WalkFiles()
	go c.WalkObjects()

	// Start a couple of workers to process them
	// Each worker calls Done()
	for i := 0; i < c.filesConcur; i++ {
		c.wait.Add(1)
		go c.FileWorker(i)
	}
	for i := 0; i < c.objectsConcur; i++ {
		c.wait.Add(1)
		go c.ObjectWorker(i)
	}

	// Wait until all walkers and workers are finished.
	c.wait.Wait()
	fmt.Println("\nCheck finished")

}
