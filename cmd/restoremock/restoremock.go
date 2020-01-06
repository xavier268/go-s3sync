package main

import (
	"fmt"

	"github.com/go-s3sync/pkg/gosync"
)

func main() {
	fmt.Println("Simulating restoring files from s3")

	c := gosync.NewConfig().SetMode(gosync.ModeRestoreMock)
	fmt.Println(c)

	c.ProcessObjects()
	c.ProcessFiles()
}
