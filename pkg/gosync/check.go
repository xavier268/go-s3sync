package gosync

import "fmt"

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
	for i := 0; i < 10; i++ {
		c.wait.Add(1)
		go c.FileWorker(i)
	}
	for i := 0; i < 10; i++ {
		c.wait.Add(1)
		go c.ObjectWorker(i)
	}

	// Wait until all walkers and workers are finished.
	c.wait.Wait()
	fmt.Println("\nCheck finished")

}
