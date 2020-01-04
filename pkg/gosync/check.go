package gosync

import (
	"fmt"
	"sync"
)

// CheckObjects performs a sync check, checking what S3 changes will be needed.
// Used for debugging,  thread safe.
func (c *Config) CheckObjects() {

	// Set a new waitGroup
	wait := new(sync.WaitGroup)

	fmt.Println("CheckObjects started")

	// Start pushing in channel
	// Will close channel and call c.filesWait.Done() upon completion
	wait.Add(1)
	go c.WalkObjects(wait)

	// Start a couple of workers to process them
	// Each worker calls Done()
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go c.ObjectWorker(i, wait)
	}

	// Wait until all walkers and workers are finished.
	wait.Wait()

	fmt.Println("\nCheckObjects finished")

}

// CheckFiles performs a sync check, checking what files are not in sync.
// Used for debugging,  thread safe.
func (c *Config) CheckFiles() {

	// Set a new waitGroup
	wait := new(sync.WaitGroup)

	fmt.Println("CheckFiles started")

	// Start pushing in channel
	// Will close channel and call c.filesWait.Done() upon completion
	wait.Add(1)
	go c.WalkFiles(wait)

	// Start a couple of workers to process them
	// Each worker calls Done()
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go c.FileWorker(i, wait)
	}

	// Wait until all walkers and workers are finished.
	wait.Wait()

	fmt.Println("\nCheckFiles finished")

}
