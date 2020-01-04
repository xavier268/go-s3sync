package gosync

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/service/s3"
)

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
	fmt.Printf("Object worker %d started ....\n", i)
	for ob := range c.objects {
		fmt.Printf("%d)\tAction : %s\n\tObject : %s\n",
			i,
			c.CheckObject(ob),
			ob.String())
	}
	fmt.Printf("Object worker %d stopped ....\n", i)
	wait.Done()
}

// CheckObject verify if the s3 object is in sync with the local file system.
// It returns the action needed.
func (c *Config) CheckObject(o DstObject) Action {
	ap := o.GetAbsPath(c)
	fi, err := os.Stat(ap)
	if err != nil || fi.IsDir() {
		return ActionDeleteObject
	}
	if fi.ModTime().Before(o.updated) && fi.Size() == o.size {
		return ActionNone
	}
	return ActionUploadFile
}
