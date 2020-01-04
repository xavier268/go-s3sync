package gosync

import (
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
)

// WalkObjects will push the s3 objects in a channel for further processing.
// It closes the object channel and call c.wait.Done() when finished.
func (c *Config) WalkObjects(wait *sync.WaitGroup) {

	li := new(s3.ListObjectsV2Input).SetBucket(c.bucket)
	c.s3.ListObjectsV2Pages(li, func(res *s3.ListObjectsV2Output, lastpage bool) bool {

		for _, o := range res.Contents {
			c.objects <- FromS3Object(o)
		}
		return !lastpage
	})

	close(c.objects)
	fmt.Println("Finished walking objects")
	wait.Done()
}

// ObjectWorker processes the objects.
// There are typically  multiple  workers running in parallel.
// It calls c.wait.Done() at the end.
func (c *Config) ObjectWorker(i int, wait *sync.WaitGroup) {
	fmt.Printf("Object worker %d started ....\n", i)
	for ob := range c.objects {
		fmt.Printf("%d)\tObject : %s\n", i, ob.String())
	}
	fmt.Printf("Object worker %d stopped ....\n", i)
	wait.Done()
}
